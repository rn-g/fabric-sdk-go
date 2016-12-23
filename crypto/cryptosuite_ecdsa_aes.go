package cryptosuite

import (
	"fmt"

	"github.com/hyperledger/fabric/core/crypto/primitives"
)

type CryptoSuite_ECDSA_AES struct {
}

func (cs CryptoSuite_ECDSA_AES) GenerateKey(ephemeral bool) (interface{}, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (cs CryptoSuite_ECDSA_AES) DeriveKey(key interface{}) (interface{}, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (cs CryptoSuite_ECDSA_AES) ImportKey(raw []byte, algorithm string) (interface{}, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (cs CryptoSuite_ECDSA_AES) GetKey(ski []byte) (interface{}, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (cs CryptoSuite_ECDSA_AES) Hash(msg []byte, algorithm string) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (cs CryptoSuite_ECDSA_AES) Sign(key interface{}, digest []byte, algorithm string, securityLevel int) ([]byte, error) {
	err := primitives.SetSecurityLevel(algorithm, securityLevel)
	if err != nil {
		return nil, err
	}
	signature, err := primitives.ECDSASign(key, digest)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (cs CryptoSuite_ECDSA_AES) Verify(key interface{}, signature []byte, digest []byte) (bool, error) {
	return false, fmt.Errorf("Not implemented yet")
}

func (cs CryptoSuite_ECDSA_AES) Encrypt(key interface{}, plaintext []byte) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

func (cs CryptoSuite_ECDSA_AES) Decrypt(key interface{}, cipherText []byte) ([]byte, error) {
	return nil, fmt.Errorf("Not implemented yet")
}
