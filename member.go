package fabric_sdk_go

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric/protos/peer"
	protos_utils "github.com/hyperledger/fabric/protos/utils"
	"github.com/op/go-logging"
	config "sk-git.securekey.com/vme/fabric-sdk-go/config"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

type Member struct {
	Name           string
	Roles          []string
	Affiliation    string
	EnrollmentKeys *Enrollment
	Chain          *Chain //The Chain object associated with this member.
	TcertBatchSize int    // The number of tcerts to get in each batch

}

type TransactionProposalRequest struct {
	Targets      []*Peer //An array or single Endorsing objects as the targets of the request
	ChaincodeId  string  //The id of the chaincode to perform the transaction proposal
	ChainId      string  //required - String of the name of the chain
	FunctionName string
	TxId         string   //required - String of the transaction id
	Args         []string //an array of arguments specific to the chaincode 'invoke'
}

type TransactionProposalResponse struct {
	Endorser         string
	ProposalResponse *pb.ProposalResponse
	Err              error
}

type Enrollment struct {
	PrivateKey      []byte
	PublicKey       []byte
	EcdsaPrivateKey *ecdsa.PrivateKey
}

/**
 * Constructor for a member.
 *
 * @param {string} memberName - The member name.
 * @param {Chain} chain - The Chain object associated with this member.
 */
func CreateNewMember(memberName string, chain *Chain) *Member {
	return &Member{Name: memberName, Chain: chain, TcertBatchSize: chain.TcertBatchSize}
}

func (m *Member) SetEnrollment(privateKey []byte, publicKey []byte) error {
	pemkey, _ := pem.Decode(privateKey)
	enrollmentPrivateKey, err := x509.ParsePKCS8PrivateKey(pemkey.Bytes)
	if err != nil {
		return err
	}
	ecPrivateKey, ok := enrollmentPrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return fmt.Errorf("key not EC")
	}

	m.EnrollmentKeys = &Enrollment{PrivateKey: privateKey, PublicKey: publicKey, EcdsaPrivateKey: ecPrivateKey}
	return nil
}

func (m *Member) SendTransactionProposal(transactionProposalRequest TransactionProposalRequest) (map[string]*TransactionProposalResponse, *pb.Proposal, error) {
	logger.Debugf("Member.sendTransactionProposal - request:%v\n", transactionProposalRequest)

	err := checkProposalRequest(transactionProposalRequest)
	if err != nil {
		return nil, nil, err
	}
	// create a proposal from a ChaincodeInvocationSpec
	prop, err := protos_utils.CreateChaincodeProposal(transactionProposalRequest.TxId, transactionProposalRequest.ChainId, createCIS(transactionProposalRequest),
		m.EnrollmentKeys.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not create chaincode proposal, err %s\n", err)
	}

	signedProposal, err := signProposal(prop, m.EnrollmentKeys.EcdsaPrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("signProposal return error, err %s\n", err)
	}
	transactionProposalResponseMap, err := SendPeersProposal(transactionProposalRequest.Targets, signedProposal)
	if err != nil {
		return nil, nil, fmt.Errorf("SendPeersProposal return error, err %s\n", err)
	}
	return transactionProposalResponseMap, prop, nil
}

func (m *Member) SendTransaction(proposal *pb.Proposal, resps []*pb.ProposalResponse) error {
	logger.Debugf("Member.SendTransaction - proposal:%v\n", proposal)
	if len(resps) < 1 {
		return fmt.Errorf("Missing 'ProposalResponse'")
	}
	if proposal == nil {
		return fmt.Errorf("Missing 'Proposal' object")
	}
	if m.Chain.Orderer == nil {
		return fmt.Errorf("Member.sendTransaction - no orderer defined")
	}
	envelope, err := createSignedTx(proposal, m.EnrollmentKeys.EcdsaPrivateKey, resps)
	if err != nil {
		return err
	}
	err = m.Chain.Orderer.SendBroadcast(envelope)
	if err != nil {
		return err
	}
	return nil
}

func createCIS(transactionProposalRequest TransactionProposalRequest) *pb.ChaincodeInvocationSpec {
	arry := make([][]byte, len(transactionProposalRequest.Args)+1)
	arry[0] = []byte(transactionProposalRequest.FunctionName)
	for i, arg := range transactionProposalRequest.Args {
		arry[i+1] = []byte(arg)
	}

	return &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        pb.ChaincodeSpec_GOLANG,
			ChaincodeID: &pb.ChaincodeID{Name: transactionProposalRequest.ChaincodeId},
			CtorMsg:     &pb.ChaincodeInput{Args: arry}}}
}

func checkProposalRequest(transactionProposalRequest TransactionProposalRequest) error {
	if transactionProposalRequest.ChaincodeId == "" {
		return fmt.Errorf("Missing 'ChaincodeId' parameter in the proposal request")
	} else if transactionProposalRequest.ChainId == "" {
		return fmt.Errorf("Missing 'ChainId' parameter in the proposal request")
	} else if len(transactionProposalRequest.Targets) < 1 {
		return fmt.Errorf("Missing 'Targets' parameter in the proposal request")
	} else if transactionProposalRequest.TxId == "" {
		return fmt.Errorf("Missing 'TxId' parameter in the proposal request")
	} else if transactionProposalRequest.ChainId == "" {
		return fmt.Errorf("Missing 'ChainId' parameter in the proposal request")
	} else if len(transactionProposalRequest.Args) < 1 {
		return fmt.Errorf("Missing 'Args' parameter in the proposal request")
	}
	return nil
}

func signProposal(proposal *pb.Proposal, enrollmentPrivateKey *ecdsa.PrivateKey) (*pb.SignedProposal, error) {
	proposalBytes, err := protos_utils.GetBytesProposal(proposal)
	if err != nil {
		return nil, err
	}
	err = primitives.SetSecurityLevel(config.GetSecurityAlgorithm(), config.GetSecurityLevel())
	if err != nil {
		return nil, err
	}
	signature, err := primitives.ECDSASign(enrollmentPrivateKey, proposalBytes)
	if err != nil {
		return nil, err
	}
	signedProposal := &pb.SignedProposal{ProposalBytes: proposalBytes, Signature: signature}
	return signedProposal, nil

}

func SendPeersProposal(peers []*Peer, signedProposal *pb.SignedProposal) (map[string]*TransactionProposalResponse, error) {
	transactionProposalResponseMap := make(map[string]*TransactionProposalResponse)
	var wg sync.WaitGroup
	for _, p := range peers {
		wg.Add(1)
		go func(peer *Peer, wg *sync.WaitGroup, tprm map[string]*TransactionProposalResponse) {
			defer wg.Done()
			var err error
			var proposalResponse *pb.ProposalResponse
			var transactionProposalResponse *TransactionProposalResponse
			logger.Debugf("Send ProposalRequest to peer :%s\n", peer.Url)
			if proposalResponse, err = peer.sendProposal(signedProposal); err != nil {
				logger.Debugf("Receive Error Response :%v\n", proposalResponse)
				transactionProposalResponse = &TransactionProposalResponse{peer.Url, nil, fmt.Errorf("Error calling endorser '%s':  %s", peer.Url, err)}
			} else {
				logger.Debugf("Receive Proposal Response :%v\n", proposalResponse)
				transactionProposalResponse = &TransactionProposalResponse{peer.Url, proposalResponse, nil}
			}
			tprm[transactionProposalResponse.Endorser] = transactionProposalResponse
		}(p, &wg, transactionProposalResponseMap)
	}
	wg.Wait()
	return transactionProposalResponseMap, nil
}

// assemble an Envelope message from proposal, endorsements and a signer.
// This function should be called by a client when it has collected enough endorsements
// for a proposal to create a transaction and submit it to peers for ordering
func createSignedTx(proposal *pb.Proposal, enrollmentPrivateKey *ecdsa.PrivateKey, resps []*pb.ProposalResponse) (*common.Envelope, error) {
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

	// ensure that all actions are bitwise equal and that they are successful
	var a1 []byte
	for n, r := range resps {
		if n == 0 {
			a1 = r.Payload
			if r.Response.Status != 200 {
				return nil, fmt.Errorf("Proposal response was not successful, error code %d, msg %s", r.Response.Status, r.Response.Message)
			}
			continue
		}

		if bytes.Compare(a1, r.Payload) != 0 {
			return nil, fmt.Errorf("ProposalResponsePayloads do not match")
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

	// sign the payload
	err = primitives.SetSecurityLevel(config.GetSecurityAlgorithm(), config.GetSecurityLevel())
	if err != nil {
		return nil, err
	}
	signature, err := primitives.ECDSASign(enrollmentPrivateKey, paylBytes)
	if err != nil {
		return nil, err
	}
	// here's the envelope
	return &common.Envelope{Payload: paylBytes, Signature: signature}, nil
}
