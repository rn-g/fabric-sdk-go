/*
Copyright SecureKey Technologies Inc. All Rights Reserved.


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at


      http://www.apache.org/licenses/LICENSE-2.0


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package events

import (
	"fmt"

	consumer "github.com/hyperledger/fabric-sdk-go/events/consumer"
	common "github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

// EventHub ...
type EventHub struct {
	// Map of clients registered for chaincode events
	chaincodeRegistrants map[string][]*ChainCodeCBE
	// Map of clients registered for block events
	blockRegistrants []func(*common.Block, string, string)
	// Map of clients registered for transactional events
	txRegistrants map[string]func(string, error)
	// peer addr to connect to
	peerAddr string
	// grpc event client interface
	client *consumer.EventsClient
	// fabric connection state of this eventhub
	connected bool
}

// ChainCodeCBE ...
/**
 * The ChainCodeCBE is used internal to the EventHub to hold chaincode
 * event registration callbacks.
 */
type ChainCodeCBE struct {
	// chaincode id
	CCID string
	// event name regex filter
	EventNameFilter string
	// callback function to invoke on successful filter match
	CallbackFunc func(*pb.ChaincodeEvent)
}

// NewEventHub ...
func NewEventHub() *EventHub {
	chaincodeRegistrants := make(map[string][]*ChainCodeCBE)
	blockRegistrants := make([]func(*common.Block, string, string), 0)
	txRegistrants := make(map[string]func(string, error))

	eventHub := &EventHub{chaincodeRegistrants: chaincodeRegistrants, blockRegistrants: blockRegistrants, txRegistrants: txRegistrants}

	return eventHub
}

// SetPeerAddr ...
/**
 * Set peer url for event source<p>
 * Note: Only use this if creating your own EventHub. The chain
 * creates a default eventHub that most Node clients can
 * use (see eventHubConnect, eventHubDisconnect and getEventHub).
 * @param {string} peeraddr peer url
 */
func (eventHub *EventHub) SetPeerAddr(peerURL string) {
	eventHub.peerAddr = peerURL
}

// Isconnected ...
/**
 * Get connected state of eventhub
 * @returns true if connected to event source, false otherwise
 */
func (eventHub *EventHub) Isconnected() bool {
	return eventHub.connected
}

// Connect ...
/**
 * Establishes connection with peer event source<p>
 */
func (eventHub *EventHub) Connect() error {
	if eventHub.peerAddr == "" {
		return fmt.Errorf("eventHub.peerAddr is empty")
	}
	eventHub.blockRegistrants = make([]func(*common.Block, string, string), 0)
	eventHub.blockRegistrants = append(eventHub.blockRegistrants, eventHub.txCallback)

	eventsClient, _ := consumer.NewEventsClient(eventHub.peerAddr, 5, eventHub)
	if err := eventsClient.Start(); err != nil {
		eventsClient.Stop()
		return fmt.Errorf("Error from eventsClient.Start (%s)", err.Error())

	}
	eventHub.connected = true
	eventHub.client = eventsClient
	return nil
}

//GetInterestedEvents implements consumer.EventAdapter interface for registering interested events
func (eventHub *EventHub) GetInterestedEvents() ([]*pb.Interest, error) {
	return []*pb.Interest{{EventType: pb.EventType_BLOCK}, {EventType: pb.EventType_REJECTION}}, nil
}

//Recv implements consumer.EventAdapter interface for receiving events
func (eventHub *EventHub) Recv(msg *pb.Event) (bool, error) {
	switch msg.Event.(type) {
	case *pb.Event_Block:
		blockEvent := msg.Event.(*pb.Event_Block)
		logger.Debugf("Recv blockEvent:%v\n", blockEvent)
		for _, v := range eventHub.blockRegistrants {
			v(blockEvent.Block, "", "")
		}
		return true, nil
	case *pb.Event_Rejection:
		rejectionEvent := msg.Event.(*pb.Event_Rejection)
		logger.Debugf("Recv rejectionEvent:%v\n", rejectionEvent)
		for _, v := range eventHub.blockRegistrants {
			v(nil, "", rejectionEvent.Rejection.ErrorMsg)
		}
		return true, nil
	default:
		return true, nil
	}
}

// Disconnected ...
/**
 * Disconnects peer event source<p>
 * Note: Only use this if creating your own EventHub. The chain
 * class creates a default eventHub that most Node clients can
 * use (see eventHubConnect, eventHubDisconnect and getEventHub).
 */
func (eventHub *EventHub) Disconnected(err error) {
	if !eventHub.connected {
		return
	}
	eventHub.client.Stop()
	eventHub.connected = false

}

// RegisterChaincodeEvent ...
/**
 * Register a callback function to receive chaincode events.
 * @param {string} ccid string chaincode id
 * @param {string} eventname string The regex string used to filter events
 * @param {function} callback Function Callback function for filter matches
 * that takes a single parameter which is a json object representation
 * of type "message ChaincodeEvent"
 * @returns {object} ChainCodeCBE object that should be treated as an opaque
 * handle used to unregister (see unregisterChaincodeEvent)
 */
func (eventHub *EventHub) RegisterChaincodeEvent(ccid string, eventname string, callback func(*pb.ChaincodeEvent)) *ChainCodeCBE {
	if !eventHub.connected {
		return nil
	}
	cbe := ChainCodeCBE{CCID: ccid, EventNameFilter: eventname, CallbackFunc: callback}
	cbeArray := eventHub.chaincodeRegistrants[ccid]
	if cbeArray == nil && len(cbeArray) <= 0 {
		cbeArray = make([]*ChainCodeCBE, 0)
		cbeArray = append(cbeArray, &cbe)
		eventHub.chaincodeRegistrants[ccid] = cbeArray
	} else {
		cbeArray = append(cbeArray, &cbe)
		eventHub.chaincodeRegistrants[ccid] = cbeArray
	}
	return &cbe
}

// UnregisterChaincodeEvent ...
/**
 * Unregister chaincode event registration
 * @param {object} ChainCodeCBE handle returned from call to
 * registerChaincodeEvent.
 */
func (eventHub *EventHub) UnregisterChaincodeEvent(cbe *ChainCodeCBE) {
	if !eventHub.connected {
		return
	}
	cbeArray := eventHub.chaincodeRegistrants[cbe.CCID]
	if len(cbeArray) <= 0 {
		logger.Debugf("No event registration for ccid %s \n", cbe.CCID)
		return
	}
	for i, v := range cbeArray {
		if v.EventNameFilter == cbe.EventNameFilter {

			cbeArray = append(cbeArray[:i], cbeArray[i+1:]...)

		}
	}
	if len(cbeArray) <= 0 {
		delete(eventHub.chaincodeRegistrants, cbe.CCID)
	}

}

// RegisterTxEvent ...
/**
 * Register a callback function to receive transactional events.<p>
 * Note: transactional event registration is primarily used by
 * the sdk to track deploy and invoke completion events. Nodejs
 * clients generally should not need to call directly.
 * @param {string} txid string transaction id
 * @param {function} callback Function that takes a single parameter which
 * is a json object representation of type "message Transaction"
 */
func (eventHub *EventHub) RegisterTxEvent(txID string, callback func(string, error)) {
	logger.Debugf("reg txid %s\n", txID)
	eventHub.txRegistrants[txID] = callback
}

// UnregisterTxEvent ...
/**
 * Unregister transactional event registration.
 * @param txid string transaction id
 */
func (eventHub *EventHub) UnregisterTxEvent(txID string) {
	delete(eventHub.txRegistrants, txID)
}

/**
 * private internal callback for processing tx events
 * @param {object} block json object representing block of tx
 * from the fabric
 */
func (eventHub *EventHub) txCallback(block *common.Block, txID string, errMsg string) {
	logger.Debugf("txCallback block=%v\n", block)

	for _, v := range block.Data.Data {
		if env, err := utils.GetEnvelopeFromBlock(v); err != nil {
			return
		} else if env != nil {
			// get the payload from the envelope
			payload, err := utils.GetPayload(env)
			if err != nil {
				return
			}

			callback := eventHub.txRegistrants[payload.Header.ChainHeader.TxID]
			if callback != nil {
				callback(payload.Header.ChainHeader.TxID, nil)
			}
		}

	}

}
