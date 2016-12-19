package fabric_sdk_go

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/util"
	pb "github.com/hyperledger/fabric/protos/peer"
	config "sk-git.securekey.com/vme/fabric-sdk-go/config"
)

func TestChainCodeInvoke(t *testing.T) {
	privateKey, err := loadEnrollmentPrivateKey()
	if err != nil {
		t.Errorf("loadEnrollmentPrivateKey return error: %v", err)
	}
	publicKey, err := loadEnrollmentPublicKey()
	if err != nil {
		t.Errorf("loadEnrollmentPublicKey return error: %v", err)
	}
	chain := CreateNewChain("testchain")
	chain.Orderer = CreateNewOrderer(fmt.Sprintf("%s:%s", config.GetOrdererHost(), config.GetOrdererPort()))
	user := chain.GetMember("admin")
	user.SetEnrollment(privateKey, publicKey)
	value := getQueryValue(t, user)
	fmt.Printf("*** QueryValue before invoke %s\n", value)

	var endorsers []*Peer
	for _, peer := range config.GetPeersConfig() {
		endorsers = append(endorsers, CreateNewPeer(fmt.Sprintf("%s:%s", peer.Host, peer.Port)))
	}
	var args []string
	args = append(args, "move")
	args = append(args, "a")
	args = append(args, "b")
	args = append(args, "1")
	transactionProposalRequest := TransactionProposalRequest{Targets: endorsers, ChaincodeId: "mycc2", FunctionName: "invoke", Args: args,
		ChainId: "**TEST_CHAINID**", TxId: util.GenerateUUID()}
	transactionProposalResponse, proposal, err := user.SendTransactionProposal(transactionProposalRequest)
	if err != nil {
		t.Errorf("SendTransactionProposal return error: %v", err)
		return
	}

	var proposalResponses []*pb.ProposalResponse
	for _, v := range transactionProposalResponse {
		if v.Err != nil {
			t.Errorf("Endorser %s return error: %v", v.Endorser, v.Err)
		}
		proposalResponses = append(proposalResponses, v.ProposalResponse)
		fmt.Printf("Endorser '%s' return ProposalResponse:%v\n", v.Endorser, v.ProposalResponse.GetResponse())
	}
	err = user.SendTransaction(proposal, proposalResponses)
	if err != nil {
		t.Errorf("SendTransaction return error: %v", err)
		return

	}
	fmt.Println("need to wait now for the committer to catch up")
	time.Sleep(time.Second * 20)
	valueAfterInvoke := getQueryValue(t, user)
	fmt.Printf("*** QueryValue after invoke %s\n", valueAfterInvoke)

	valueInt, _ := strconv.Atoi(value)
	valueInt = valueInt + 1
	valueAfterInvokeInt, _ := strconv.Atoi(valueAfterInvoke)
	if valueInt != valueAfterInvokeInt {
		t.Errorf("SendTransaction didn't change the QueryValue")
		return

	}

}

func getQueryValue(t *testing.T, user *Member) string {
	var endorsers []*Peer
	for _, peer := range config.GetPeersConfig() {
		endorsers = append(endorsers, CreateNewPeer(fmt.Sprintf("%s:%s", peer.Host, peer.Port)))
		break
	}
	var args []string
	args = append(args, "query")
	args = append(args, "b")
	transactionProposalRequest := TransactionProposalRequest{Targets: endorsers, ChaincodeId: "mycc2", FunctionName: "invoke", Args: args,
		ChainId: "**TEST_CHAINID**", TxId: util.GenerateUUID()}
	transactionProposalResponse, _, err := user.SendTransactionProposal(transactionProposalRequest)
	if err != nil {
		t.Errorf("SendTransactionProposal return error: %v", err)
	}

	for _, v := range transactionProposalResponse {
		if v.Err != nil {
			t.Errorf("Endorser %s return error: %v", v.Endorser, v.Err)
		}
		return string(v.ProposalResponse.GetResponse().Payload)
	}
	return ""
}

func loadEnrollmentPrivateKey() ([]byte, error) {
	raw, err := ioutil.ReadFile("./test_resources/private.pem")
	if err != nil {
		return nil, err
	}
	return raw, nil

}

func loadEnrollmentPublicKey() ([]byte, error) {
	raw, err := ioutil.ReadFile("./test_resources/public.pem")
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func TestMain(m *testing.M) {
	err := config.InitConfig("./test_resources/config_test.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(m.Run())
}
