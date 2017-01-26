package fabric_sdk_go

import (
	_ "bytes"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/common/util"
	msp "github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"

	protos_utils "github.com/hyperledger/fabric/protos/utils"
	"github.com/op/go-logging"

	config "github.com/hyperledger/fabric-sdk-go/config"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

/**
 * The “Chain” object captures settings for a channel, which is created by
 * the orderers to isolate transactions delivery to peers participating on channel.
 * A chain must be initialized after it has been configured with the list of peers
 * and orderers. The initialization sends a CONFIGURATION transaction to the orderers
 * to create the specified channel and asks the peers to join that channel.
 *
 */
type Chain struct {
	name            string // Name of the chain is only meaningful to the client
	securityEnabled bool   // Security enabled flag
	peers           map[string]*Peer
	tcertBatchSize  int // The number of tcerts to get in each batch
	orderers        map[string]*Orderer
	clientContext   *Client
}

/**
 * The TransactionProposalResponse object result return from endorsers.
 */
type TransactionProposalResponse struct {
	Endorser         string
	ProposalResponse *pb.ProposalResponse
	Err              error
}

/**
 * The TransactionProposalResponse object result return from orderers.
 */
type TransactionResponse struct {
	Orderer string
	Err     error
}

/**
 * @param {string} name to identify different chain instances. The naming of chain instances
 * is enforced by the ordering service and must be unique within the blockchain network
 * @param {Client} clientContext An instance of {@link Client} that provides operational context
 * such as submitting User etc.
 */
func NewChain(name string, client *Client) (*Chain, error) {
	if name == "" {
		return nil, fmt.Errorf("Failed to create Chain. Missing requirement 'name' parameter.")
	}
	if client == nil {
		return nil, fmt.Errorf("Failed to create Chain. Missing requirement 'clientContext' parameter.")
	}
	p := make(map[string]*Peer)
	o := make(map[string]*Orderer)
	c := &Chain{name: name, securityEnabled: config.IsSecurityEnabled(), peers: p,
		tcertBatchSize: config.TcertBatchSize(), orderers: o, clientContext: client}
	logger.Infof("Constructed Chain instance: %v", c)

	return c, nil
}

/**
 * Get the chain name.
 * @returns {string} The name of the chain.
 */
func (c *Chain) GetName() string {
	return c.name
}

/**
 * Determine if security is enabled.
 */
func (c *Chain) IsSecurityEnabled() bool {
	return c.securityEnabled
}

/**
 * Get the tcert batch size.
 */
func (c *Chain) GetTCertBatchSize() int {
	return c.tcertBatchSize
}

/**
 * Set the tcert batch size.
 */
func (c *Chain) SetTCertBatchSize(batchSize int) {
	c.tcertBatchSize = batchSize
}

/**
 * Add peer endpoint to chain.
 * @param {Peer} peer An instance of the Peer that has been initialized with URL,
 * TLC certificate, and enrollment certificate.
 */
func (c *Chain) AddPeer(peer *Peer) {
	c.peers[peer.GetUrl()] = peer
}

/**
 * Remove peer endpoint from chain.
 * @param {Peer} peer An instance of the Peer.
 */
func (c *Chain) RemovePeer(peer Peer) {
	delete(c.peers, peer.GetUrl())
}

/**
 * Get peers of a chain from local information.
 * @returns {[]Peer} The peer list on the chain.
 */
func (c *Chain) GetPeers() []*Peer {
	var peersArray []*Peer
	for _, v := range c.peers {
		peersArray = append(peersArray, v)
	}
	return peersArray
}

/**
 * Add orderer endpoint to a chain object, this is a local-only operation.
 * A chain instance may choose to use a single orderer node, which will broadcast
 * requests to the rest of the orderer network. Or if the application does not trust
 * the orderer nodes, it can choose to use more than one by adding them to the chain instance.
 * All APIs concerning the orderer will broadcast to all orderers simultaneously.
 * @param {Orderer} orderer An instance of the Orderer class.
 */
func (c *Chain) AddOrderer(orderer *Orderer) {
	c.orderers[orderer.Url] = orderer
}

/**
 * Remove orderer endpoint from a chain object, this is a local-only operation.
 * @param {Orderer} orderer An instance of the Orderer class.
 */
func (c *Chain) RemoveOrderer(orderer *Orderer) {
	delete(c.orderers, orderer.Url)

}

/**
 * Get orderers of a chain.
 */
func (c *Chain) GetOrderers() []*Orderer {
	var orderersArray []*Orderer
	for _, v := range c.orderers {
		orderersArray = append(orderersArray, v)
	}
	return orderersArray
}

/**
 * Calls the orderer(s) to start building the new chain, which is a combination
 * of opening new message stream and connecting the list of participating peers.
 * This is a long-running process. Only one of the application instances needs
 * to call this method. Once the chain is successfully created, other application
 * instances only need to call getChain() to obtain the information about this chain.
 * @returns {bool} Whether the chain initialization process was successful.
 */
func (c *Chain) InitializeChain() bool {
	return false
}

/**
 * Calls the orderer(s) to update an existing chain. This allows the addition and
 * deletion of Peer nodes to an existing chain, as well as the update of Peer
 * certificate information upon certificate renewals.
 * @returns {bool} Whether the chain update process was successful.
 */
func (c *Chain) UpdateChain() bool {
	return false
}

/**
 * Get chain status to see if the underlying channel has been terminated,
 * making it a read-only chain, where information (transactions and states)
 * can be queried but no new transactions can be submitted.
 * @returns {bool} Is read-only, true or not.
 */
func (c *Chain) IsReadonly() bool {
	return false //to do
}

/**
 * Queries for various useful information on the state of the Chain
 * (height, known peers).
 * @returns {object} With height, currently the only useful info.
 */
func (c *Chain) QueryInfo() {
	//to do
}

/**
 * Queries the ledger for Block by block number.
 * @param {int} blockNumber The number which is the ID of the Block.
 * @returns {object} Object containing the block.
 */
func (c *Chain) QueryBlock(blockNumber int) {
	//to do
}

/**
 * Queries the ledger for Transaction by number.
 * @param {int} transactionID
 * @returns {object} Transaction information containing the transaction.
 */
func (c *Chain) QueryTransaction(transactionID int) {
	//to do
}

/**
 * Create  a proposal for transaction. This involves assembling the proposal
 * with the data (chaincodeName, function to call, arguments, etc.) and signing it using the private key corresponding to the
 * ECert to sign.
 */
func (c *Chain) CreateTransactionProposal(chaincodeName string, chainId string, args []string, sign bool) (*pb.SignedProposal, *pb.Proposal, error) {

	arry := make([][]byte, len(args))
	for i, arg := range args {
		arry[i] = []byte(arg)
	}
	ccis := &pb.ChaincodeInvocationSpec{ChaincodeSpec: &pb.ChaincodeSpec{
		Type: pb.ChaincodeSpec_GOLANG, ChaincodeID: &pb.ChaincodeID{Name: chaincodeName},
		Input: &pb.ChaincodeInput{Args: arry}}}

	txid := util.GenerateUUID()

	user, err := c.clientContext.GetUserContext("")
	if err != nil {
		return nil, nil, fmt.Errorf("GetUserContext return error: %s\n", err)
	}
	serializedIdentity := &msp.SerializedIdentity{Mspid: config.GetMspId(), IdBytes: user.GetEnrollmentCertificate()}
	creatorId, err := proto.Marshal(serializedIdentity)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not Marshal serializedIdentity, err %s\n", err)
	}
	// create a proposal from a ChaincodeInvocationSpec
	proposal, err := protos_utils.CreateChaincodeProposal(txid, common.HeaderType_ENDORSER_TRANSACTION, chainId, ccis, creatorId)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not create chaincode proposal, err %s\n", err)
	}

	proposalBytes, err := protos_utils.GetBytesProposal(proposal)
	if err != nil {
		return nil, nil, err
	}
	cryptoSuite := c.clientContext.GetCryptoSuite()
	digest, err := cryptoSuite.Hash(proposalBytes, &bccsp.SHAOpts{})
	if err != nil {
		return nil, nil, err
	}
	signature, err := cryptoSuite.Sign(user.GetPrivateKey(),
		digest, nil)
	if err != nil {
		return nil, nil, err
	}
	signedProposal := &pb.SignedProposal{ProposalBytes: proposalBytes, Signature: signature}
	return signedProposal, proposal, nil
}

func (c *Chain) SendTransactionProposal(signedProposal *pb.SignedProposal, retry int) (map[string]*TransactionProposalResponse, error) {
	transactionProposalResponseMap := make(map[string]*TransactionProposalResponse)
	var wg sync.WaitGroup
	for _, p := range c.peers {
		wg.Add(1)
		go func(peer *Peer, wg *sync.WaitGroup, tprm map[string]*TransactionProposalResponse) {
			defer wg.Done()
			var err error
			var proposalResponse *pb.ProposalResponse
			var transactionProposalResponse *TransactionProposalResponse
			logger.Debugf("Send ProposalRequest to peer :%s\n", peer.GetUrl())
			if proposalResponse, err = peer.SendProposal(signedProposal); err != nil {
				logger.Debugf("Receive Error Response :%v\n", proposalResponse)
				transactionProposalResponse = &TransactionProposalResponse{peer.GetUrl(), nil, fmt.Errorf("Error calling endorser '%s':  %s", peer.GetUrl(), err)}
			} else {
				prp1, _ := protos_utils.GetProposalResponsePayload(proposalResponse.Payload)
				act1, _ := protos_utils.GetChaincodeAction(prp1.Extension)
				logger.Debugf("%s ProposalResponsePayload Extension ChaincodeAction Results\n%s\n", peer.GetUrl(), string(act1.Results))

				logger.Debugf("Receive Proposal ChaincodeActionResponse :%v\n", proposalResponse)
				transactionProposalResponse = &TransactionProposalResponse{peer.GetUrl(), proposalResponse, nil}
			}
			tprm[transactionProposalResponse.Endorser] = transactionProposalResponse
		}(p, &wg, transactionProposalResponseMap)
	}
	wg.Wait()
	return transactionProposalResponseMap, nil
}

/**
 * Create a transaction with proposal response, following the endorsement policy.
 */
func (c *Chain) CreateTransaction(proposal *pb.Proposal, resps []*pb.ProposalResponse) (*pb.Transaction, error) {
	if len(resps) == 0 {
		return nil, fmt.Errorf("At least one proposal response is necessary")
	}

	// the original header
	hdr, err := protos_utils.GetHeader(proposal.Header)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proposal header")
	}

	// the original payload
	pPayl, err := protos_utils.GetChaincodeProposalPayload(proposal.Payload)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proposal payload")
	}

	// get header extensions so we have the visibility field
	hdrExt, err := protos_utils.GetChaincodeHeaderExtension(hdr)
	if err != nil {
		return nil, err
	}

	// This code is commented out because the ProposalResponsePayload Extension ChaincodeAction Results
	// return from endorsements is different so the compare will fail

	//	var a1 []byte
	//	for n, r := range resps {
	//		if n == 0 {
	//			a1 = r.Payload
	//			if r.Response.Status != 200 {
	//				return nil, fmt.Errorf("Proposal response was not successful, error code %d, msg %s", r.Response.Status, r.Response.Message)
	//			}
	//			continue
	//		}

	//		if bytes.Compare(a1, r.Payload) != 0 {
	//			return nil, fmt.Errorf("ProposalResponsePayloads do not match")
	//		}
	//	}

	for _, r := range resps {
		if r.Response.Status != 200 {
			return nil, fmt.Errorf("Proposal response was not successful, error code %d, msg %s", r.Response.Status, r.Response.Message)
		}
	}

	// fill endorsements
	endorsements := make([]*pb.Endorsement, len(resps))
	for n, r := range resps {
		endorsements[n] = r.Endorsement
	}
	// create ChaincodeEndorsedAction
	cea := &pb.ChaincodeEndorsedAction{ProposalResponsePayload: resps[0].Payload, Endorsements: endorsements}

	// obtain the bytes of the proposal payload that will go to the transaction
	propPayloadBytes, err := protos_utils.GetBytesProposalPayloadForTx(pPayl, hdrExt.PayloadVisibility)
	if err != nil {
		return nil, err
	}

	// get the bytes of the signature header, that will be the header of the TransactionAction
	sHdrBytes, err := protos_utils.GetBytesSignatureHeader(hdr.SignatureHeader)
	if err != nil {
		return nil, err
	}

	// serialize the chaincode action payload
	cap := &pb.ChaincodeActionPayload{ChaincodeProposalPayload: propPayloadBytes, Action: cea}
	capBytes, err := protos_utils.GetBytesChaincodeActionPayload(cap)
	if err != nil {
		return nil, err
	}

	// create a transaction
	taa := &pb.TransactionAction{Header: sHdrBytes, Payload: capBytes}
	taas := make([]*pb.TransactionAction, 1)
	taas[0] = taa
	tx := &pb.Transaction{Actions: taas}

	return tx, nil

}

/**
 * Send a transaction to the chain’s orderer service (one or more orderer endpoints) for consensus and committing to the ledger.
 * This call is asynchronous and the successful transaction commit is notified via a BLOCK or CHAINCODE event. This method must provide a mechanism for applications to attach event listeners to handle “transaction submitted”, “transaction complete” and “error” events.
 * Note that under the cover there are two different kinds of communications with the fabric backend that trigger different events to
 * be emitted back to the application’s handlers:
 * 1-)The grpc client with the orderer service uses a “regular” stateless HTTP connection in a request/response fashion with the “broadcast” call.
 * The method implementation should emit “transaction submitted” when a successful acknowledgement is received in the response,
 * or “error” when an error is received
 * 2-)The method implementation should also maintain a persistent connection with the Chain’s event source Peer as part of the
 * internal event hub mechanism in order to support the fabric events “BLOCK”, “CHAINCODE” and “TRANSACTION”.
 * These events should cause the method to emit “complete” or “error” events to the application.
 */
func (c *Chain) SendTransaction(proposal *pb.Proposal, tx *pb.Transaction) (map[string]*TransactionResponse, error) {
	if c.orderers == nil || len(c.orderers) == 0 {
		return nil, fmt.Errorf("orderers is nil")
	}
	if proposal == nil {
		return nil, fmt.Errorf("proposal is nil")
	}
	if tx == nil {
		return nil, fmt.Errorf("Transaction is nil")
	}
	// the original header
	hdr, err := protos_utils.GetHeader(proposal.Header)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the proposal header")
	}
	// serialize the tx
	txBytes, err := protos_utils.GetBytesTransaction(tx)
	if err != nil {
		return nil, err
	}

	// create the payload
	payl := &common.Payload{Header: hdr, Data: txBytes}
	paylBytes, err := protos_utils.GetBytesPayload(payl)
	if err != nil {
		return nil, err
	}

	cryptoSuite := c.clientContext.GetCryptoSuite()
	digest, err := cryptoSuite.Hash(paylBytes, &bccsp.SHAOpts{})
	if err != nil {
		return nil, err
	}
	user, err := c.clientContext.GetUserContext("")
	if err != nil {
		return nil, fmt.Errorf("GetUserContext return error: %s\n", err)
	}
	signature, err := cryptoSuite.Sign(user.GetPrivateKey(),
		digest, nil)
	if err != nil {
		return nil, err
	}
	// here's the envelope
	envelope := &common.Envelope{Payload: paylBytes, Signature: signature}

	transactionResponseMap := make(map[string]*TransactionResponse)
	var wg sync.WaitGroup
	for _, o := range c.orderers {
		wg.Add(1)
		go func(orderer *Orderer, wg *sync.WaitGroup, trm map[string]*TransactionResponse) {
			defer wg.Done()
			var err error
			var transactionResponse *TransactionResponse

			logger.Debugf("Send TransactionRequest to orderer :%s\n", orderer.Url)
			if err = orderer.sendBroadcast(envelope); err != nil {
				logger.Debugf("Receive Error Response from orderer :%v\n", err)
				transactionResponse = &TransactionResponse{orderer.Url, fmt.Errorf("Error calling endorser '%s':  %s", orderer.Url, err)}
			} else {
				logger.Debugf("Receive Success Response from orderer\n")
				transactionResponse = &TransactionResponse{orderer.Url, nil}
			}
			trm[transactionResponse.Orderer] = transactionResponse
		}(o, &wg, transactionResponseMap)
	}
	wg.Wait()
	return transactionResponseMap, nil

}
