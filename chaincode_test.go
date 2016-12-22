package fabric_sdk_go

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	pb "github.com/hyperledger/fabric/protos/peer"
	config "sk-git.securekey.com/vme/fabric-sdk-go/config"
	crypto "sk-git.securekey.com/vme/fabric-sdk-go/crypto"
)

func TestChainCodeInvoke(t *testing.T) {
	// Configuration
	privateKey, err := loadEnrollmentPrivateKey()
	if err != nil {
		t.Errorf("loadEnrollmentPrivateKey return error: %v", err)
	}
	publicKey, err := loadEnrollmentPublicKey()
	if err != nil {
		t.Errorf("loadEnrollmentPublicKey return error: %v", err)
	}
	client := NewClient()
	cs_ecdsa_aes := crypto.CryptoSuite_ECDSA_AES{}
	client.SetCryptoSuite(cs_ecdsa_aes)
	user := NewUser("testuser")
	user.SetEnrollment(privateKey, publicKey)
	client.SetUserContext(user)

	querychain, err := client.NewChain("querychain")
	if err != nil {
		t.Errorf("NewChain return error: %v", err)
		return
	}

	for _, p := range config.GetPeersConfig() {
		endorser := CreateNewPeer(fmt.Sprintf("%s:%s", p.Host, p.Port))
		querychain.AddPeer(endorser)
		break
	}

	invokechain, err := client.NewChain("invokechain")
	if err != nil {
		t.Errorf("NewChain return error: %v", err)
		return
	}
	orderer := CreateNewOrderer(fmt.Sprintf("%s:%s", config.GetOrdererHost(), config.GetOrdererPort()))
	invokechain.AddOrderer(orderer)

	for _, p := range config.GetPeersConfig() {
		endorser := CreateNewPeer(fmt.Sprintf("%s:%s", p.Host, p.Port))
		invokechain.AddPeer(endorser)
	}
	endorser := CreateNewPeer(fmt.Sprintf("%s:%s", "localhost", "7051"))
	invokechain.AddPeer(endorser)
	endorser = CreateNewPeer(fmt.Sprintf("%s:%s", "localhost", "7056"))
	invokechain.AddPeer(endorser)

	// Get Query value before invoke
	value, err := getQueryValue(t, querychain)
	if err != nil {
		t.Errorf("getQueryValue return error: %v", err)
		return
	}
	fmt.Printf("*** QueryValue before invoke %s\n", value)

	err = invoke(t, invokechain)
	if err != nil {
		t.Errorf("invoke return error: %v", err)
		return
	}

	fmt.Println("need to wait now for the committer to catch up")
	time.Sleep(time.Second * 20)
	valueAfterInvoke, err := getQueryValue(t, querychain)
	if err != nil {
		t.Errorf("getQueryValue return error: %v", err)
		return
	}
	fmt.Printf("*** QueryValue after invoke %s\n", valueAfterInvoke)

	valueInt, _ := strconv.Atoi(value)
	valueInt = valueInt + 1
	valueAfterInvokeInt, _ := strconv.Atoi(valueAfterInvoke)
	if valueInt != valueAfterInvokeInt {
		t.Errorf("SendTransaction didn't change the QueryValue")
		return

	}

}

func getQueryValue(t *testing.T, chain *Chain) (string, error) {

	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "b")
	signedProposal, _, err := chain.CreateTransactionProposal("mycc2", "**TEST_CHAINID**", args, true)
	if err != nil {
		return "", fmt.Errorf("SendTransactionProposal return error: %v", err)
	}
	transactionProposalResponse, err := chain.SendTransactionProposal(signedProposal, 0)
	if err != nil {
		return "", fmt.Errorf("SendTransactionProposal return error: %v", err)
	}

	for _, v := range transactionProposalResponse {
		if v.Err != nil {
			return "", fmt.Errorf("Endorser %s return error: %v", v.Endorser, v.Err)
		}
		return string(v.ProposalResponse.GetResponse().Payload), nil
	}
	return "", nil
}

func invoke(t *testing.T, chain *Chain) error {

	var args []string
	args = append(args, "invoke")
	args = append(args, "move")
	args = append(args, "a")
	args = append(args, "b")
	args = append(args, "1")
	signedProposal, proposal, err := chain.CreateTransactionProposal("mycc2", "**TEST_CHAINID**", args, true)
	if err != nil {
		return fmt.Errorf("SendTransactionProposal return error: %v", err)
	}
	transactionProposalResponse, err := chain.SendTransactionProposal(signedProposal, 0)
	if err != nil {
		return fmt.Errorf("SendTransactionProposal return error: %v", err)
	}

	var proposalResponses []*pb.ProposalResponse
	for _, v := range transactionProposalResponse {
		if v.Err != nil {
			return fmt.Errorf("Endorser %s return error: %v", v.Endorser, v.Err)
		}
		proposalResponses = append(proposalResponses, v.ProposalResponse)
		fmt.Printf("Endorser '%s' return ProposalResponse:%v\n", v.Endorser, v.ProposalResponse.GetResponse())
	}

	tx, err := chain.CreateTransaction(proposal, proposalResponses)
	if err != nil {
		return fmt.Errorf("CreateTransaction return error: %v", err)

	}
	transactionResponse, err := chain.SendTransaction(proposal, tx)
	if err != nil {
		return fmt.Errorf("SendTransaction return error: %v", err)

	}
	for _, v := range transactionResponse {
		if v.Err != nil {
			return fmt.Errorf("Orderer %s return error: %v", v.Orderer, v.Err)
		}
	}
	return nil

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
