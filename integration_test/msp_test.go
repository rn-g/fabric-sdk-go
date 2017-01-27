package integration_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	fabric_sdk "github.com/hyperledger/fabric-sdk-go"
	config "github.com/hyperledger/fabric-sdk-go/config"
	kvs "github.com/hyperledger/fabric-sdk-go/keyvaluestore"
	"github.com/hyperledger/fabric/bccsp"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/bccsp/sw"

	msp "github.com/hyperledger/fabric-sdk-go/msp"
)

// this test uses the MSPServices to enroll a user, and
// saves the enrollment materials into a key value store.
// then uses the Client class to load the member from the
// key value store
func TestEnroll(t *testing.T) {
	InitConfigForMsp()
	client := fabric_sdk.NewClient()
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
	stateStore, err := kvs.CreateNewFileKeyValueStore("/tmp/enroll_user")
	if err != nil {
		t.Fatalf("CreateNewFileKeyValueStore return error[%s]", err)
	}
	client.SetStateStore(stateStore)

	msps, err := msp.NewMSPServices(config.GetMspUrl(), config.GetMspClientPath())
	if err != nil {
		t.Fatalf("NewMSPServices return error: %v", err)
	}
	key, cert, err := msps.Enroll("testUser2", "user2")
	if err != nil {
		t.Fatalf("Enroll return error: %v", err)
	}
	if key == nil {
		t.Fatalf("private key return from Enroll is nil")
	}
	if cert == nil {
		t.Fatalf("cert return from Enroll is nil")
	}

	certPem, _ := pem.Decode(cert)
	if err != nil {
		t.Fatalf("pem Decode return error: %v", err)
	}

	cert509, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		t.Fatalf("x509 ParseCertificate return error: %v", err)
	}
	if cert509.Subject.CommonName != "testUser2" {
		t.Fatalf("CommonName in x509 cert is not the enrollmentID")
	}

	keyPem, _ := pem.Decode(key)
	if err != nil {
		t.Fatalf("pem Decode return error: %v", err)
	}
	user := fabric_sdk.NewUser("testUser2")
	k, err := client.GetCryptoSuite().KeyImport(keyPem.Bytes, &bccsp.ECDSAPrivateKeyImportOpts{Temporary: false})
	if err != nil {
		t.Fatalf("KeyImport return error: %v", err)
	}
	user.SetPrivateKey(k)
	user.SetEnrollmentCertificate(cert)
	err = client.SetUserContext(user, false)
	if err != nil {
		t.Fatalf("client.SetUserContext return error: %v", err)
	}
	user, err = client.GetUserContext("testUser2")
	if err != nil {
		t.Fatalf("client.GetUserContext return error: %v", err)
	}
	if user == nil {
		t.Fatalf("client.GetUserContext return nil")
	}

}

func InitConfigForMsp() {
	config.InitConfig("./test_resources/config/config_test.yaml")
}
