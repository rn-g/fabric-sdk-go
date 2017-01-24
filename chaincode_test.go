package fabric_sdk_go

import (
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	config "github.com/hyperledger/fabric-sdk-go/config"
	cop "github.com/hyperledger/fabric-sdk-go/cop"
	"github.com/hyperledger/fabric/bccsp"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"

	"github.com/hyperledger/fabric/bccsp/sw"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func TestChainCodeInvoke(t *testing.T) {
	// Configuration
	client := NewClient()

	ks := &sw.FileBasedKeyStore{}
	if err := ks.Init(nil, config.GetKeyStorePath(), false); err != nil {
		t.Fatalf("Failed initializing key store [%s]", err)
	}

	cryptoSuite, err := bccspFactory.GetBCCSP(&bccspFactory.SwOpts{Ephemeral_: true, SecLevel: config.GetSecurityLevel(),
		HashFamily: config.GetSecurityAlgorithm(), KeyStore: ks})
	if err != nil {
		t.Fatalf("Failed getting ephemeral software-based BCCSP [%s]", err)
	}

	client.SetCryptoSuite(cryptoSuite)
	if client.GetUserContext("admin") == nil {
		fcs, err := cop.NewFabricCOPServices(config.GetMspUrl(), config.GetMspClientPath())
		if err != nil {
			t.Fatalf("NewFabricCOPServices return error: %v", err)
		}
		key, cert, err := fcs.Enroll("admin", "adminpw")
		block, _ := pem.Decode(key)
		if err != nil {
			t.Fatalf("Enroll return error: %v", err)
		}
		user := NewUser("admin")
		k, err := client.GetCryptoSuite().KeyImport(block.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: false})
		if err != nil {
			t.Fatalf("KeyImport return error: %v", err)
		}
		user.SetPrivateKey(k)
		user.SetEnrollmentCertificate(cert)
		client.SetUserContext(user)
	}

	querychain, err := client.NewChain("querychain")
	if err != nil {
		t.Fatalf("NewChain return error: %v", err)
	}

	for _, p := range config.GetPeersConfig() {
		endorser := CreateNewPeer(fmt.Sprintf("%s:%s", p.Host, p.Port))
		querychain.AddPeer(endorser)
		break
	}

	invokechain, err := client.NewChain("invokechain")
	if err != nil {
		t.Fatalf("NewChain return error: %v", err)
	}
	orderer := CreateNewOrderer(fmt.Sprintf("%s:%s", config.GetOrdererHost(), config.GetOrdererPort()))
	invokechain.AddOrderer(orderer)

	for _, p := range config.GetPeersConfig() {
		endorser := CreateNewPeer(fmt.Sprintf("%s:%s", p.Host, p.Port))
		invokechain.AddPeer(endorser)
	}

	// Get Query value before invoke
	value, err := getQueryValue(t, querychain)
	if err != nil {
		t.Fatalf("getQueryValue return error: %v", err)
	}
	fmt.Printf("*** QueryValue before invoke %s\n", value)

	err = invoke(t, invokechain)
	if err != nil {
		t.Fatalf("invoke return error: %v", err)
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
		t.Fatalf("SendTransaction didn't change the QueryValue")

	}

}

func getQueryValue(t *testing.T, chain *Chain) (string, error) {

	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "b")
	signedProposal, _, err := chain.CreateTransactionProposal("mycc", "test_chainid", args, true)
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
	signedProposal, proposal, err := chain.CreateTransactionProposal("mycc", "test_chainid", args, true)
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

func TestMain(m *testing.M) {
	err := config.InitConfig("./test_resources/config_test.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(m.Run())
}
