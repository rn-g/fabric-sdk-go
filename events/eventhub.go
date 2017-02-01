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
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

// EventHub ...
type EventHub struct {
	// Map of clients registered for chaincode events
	chaincodeRegistrants map[string]*ChainCodeCBE
	// Map of clients registered for block events
	blockRegistrants map[string]func(*common.Block, error)
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
	ccid string
	// event name regex filter
	eventNameFilter string
	// callback function to invoke on successful filter match
	callbackFunc func(*pb.ChaincodeEvent)
}

// SetPeerAddr ...
/**
 * Set peer url for event source<p>
 * Note: Only use this if creating your own EventHub. The chain
 * creates a default eventHub that most Node clients can
 * use (see eventHubConnect, eventHubDisconnect and getEventHub).
 * @param {string} peeraddr peer url
 */
func (eventHub EventHub) SetPeerAddr(peerAddr string) {
	eventHub.peerAddr = peerAddr
}

// Isconnected ...
/**
 * Get connected state of eventhub
 * @returns true if connected to event source, false otherwise
 */
func (eventHub EventHub) Isconnected() bool {
	return eventHub.connected
}

// Connect ...
/**
 * Establishes connection with peer event source<p>
 */
func (eventHub *EventHub) Connect() error {

	eventsClient, err := consumer.NewEventsClient(eventHub.peerAddr, 5, eventHub)
	if err != nil {
		return fmt.Errorf("Error from consumer.NewEventsClient (%s)", err.Error())
	}
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
			v(blockEvent.Block, nil)
		}
		return true, nil
	case *pb.Event_Rejection:
		rejectionEvent := msg.Event.(*pb.Event_Rejection)
		logger.Debugf("Recv rejectionEvent:%v\n", rejectionEvent)
		for _, v := range eventHub.blockRegistrants {
			v(nil, fmt.Errorf(rejectionEvent.Rejection.ErrorMsg))
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
	//unregisterBlockEvent()
	eventHub.client.Stop()
	eventHub.connected = false

}
