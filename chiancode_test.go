package fabric_sdk_go

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hyperledger/fabric/core/util"
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
	user := chain.GetMember("admin")
	user.SetEnrollment(privateKey, publicKey)
	var endorsers []*Peer
	for _, peer := range config.GetPeersConfig() {
		endorsers = append(endorsers, CreateNewPeer(fmt.Sprintf("%s:%s", peer.Host, peer.Port)))

	}
	var args []string
	args = append(args, "move")
	args = append(args, "a")
	args = append(args, "b")
	args = append(args, "100")
	transactionProposalRequest := TransactionProposalRequest{Targets: endorsers, ChaincodeId: "mycc2", FunctionName: "invoke", Args: args,
		ChainId: "**TEST_CHAINID**", TxId: util.GenerateUUID()}
	transactionProposalResponse, err := user.SendTransactionProposal(transactionProposalRequest)
	if err != nil {
		t.Errorf("SendTransactionProposal return error: %v", err)
	}

	for _, v := range transactionProposalResponse {
		if v.Err != nil {
			t.Errorf("Endorser %s return error: %v", v.Endorser, v.Err)
		}
		fmt.Printf("Endorser '%s' return ProposalResponse: %v\n", v.Endorser, v.Proposal)
	}

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
