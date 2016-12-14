package fabric_sdk_go

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"

	"github.com/hyperledger/fabric/core/crypto/primitives"
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
	Endorser string
	Proposal *pb.ProposalResponse
	Err      error
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
func CreateNewMember(memberName string, chain Chain) *Member {
	return &Member{Name: memberName, TcertBatchSize: chain.TcertBatchSize}
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

func (m *Member) SendTransactionProposal(transactionProposalRequest TransactionProposalRequest) (map[string]*TransactionProposalResponse, error) {
	logger.Debugf("Member.sendTransactionProposal - request:%v\n", transactionProposalRequest)

	err := checkProposalRequest(transactionProposalRequest)
	if err != nil {
		return nil, err
	}
	// create a proposal from a ChaincodeInvocationSpec
	prop, err := protos_utils.CreateChaincodeProposal(transactionProposalRequest.TxId, transactionProposalRequest.ChainId, createCIS(transactionProposalRequest),
		m.EnrollmentKeys.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("Could not create chaincode proposal, err %s\n", err)
	}

	signedProposal, err := signProposal(prop, m.EnrollmentKeys.EcdsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("signProposal return error, err %s\n", err)
	}
	transactionProposalResponseMap, err := SendPeersProposal(transactionProposalRequest.Targets, signedProposal)
	if err != nil {
		return nil, fmt.Errorf("SendPeersProposal return error, err %s\n", err)
	}
	return transactionProposalResponseMap, nil
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
