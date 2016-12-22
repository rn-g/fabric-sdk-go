package fabric_sdk_go

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

/**
 * The User struct represents users that have been enrolled and represented by
 * an enrollment certificate (ECert) and a signing key. The ECert must have
 * been signed by one of the CAs the blockchain network has been configured to trust.
 * An enrolled user (having a signing key and ECert) can conduct chaincode deployments,
 * transactions and queries with the Chain.
 *
 * User ECerts can be obtained from a CA beforehand as part of deploying the application,
 * or it can be obtained from the optional Fabric COP service via its enrollment process.
 *
 * Sometimes User identities are confused with Peer identities. User identities represent
 * signing capability because it has access to the private key, while Peer identities in
 * the context of the application/SDK only has the certificate for verifying signatures.
 * An application cannot use the Peer identity to sign things because the application doesn’t
 * have access to the Peer identity’s private key.
 *
 */
type User struct {
	name                  string
	roles                 []string
	enrollmentKeys        *Enrollment
	enrollmentCertificate *pem.Block
}

/**
 * This structure temporary until we have tcerts
 */
type Enrollment struct {
	PrivateKey      []byte
	PublicKey       []byte
	EcdsaPrivateKey *ecdsa.PrivateKey
}

/**
 * Constructor for a user.
 *
 * @param {string} name - The user name
 */
func NewUser(name string) *User {
	return &User{name: name}
}

/**
 * Get the user name.
 * @returns {string} The user name.
 */
func (u *User) GetName() string {
	return u.name
}

/**
 * Get the roles.
 * @returns {[]string} The roles.
 */
func (u *User) GetRoles() []string {
	return u.roles
}

/**
 * Set the roles.
 * @param roles {[]string} The roles.
 */
func (u *User) SetRoles(roles []string) {
	u.roles = roles
}

/**
 * Returns the underlying ECert representing this user’s identity.
 */
func (u *User) GetEnrollmentCertificate() *pem.Block {
	return u.enrollmentCertificate
}

/**
 * Set the user’s Enrollment Certificate.
 */
func (u *User) SetEnrollmentCertificate(pem *pem.Block) {
	u.enrollmentCertificate = pem
}

/**
 * deprecated.
 */
func (u *User) SetEnrollment(privateKey []byte, publicKey []byte) error {
	pemkey, _ := pem.Decode(privateKey)
	enrollmentPrivateKey, err := x509.ParsePKCS8PrivateKey(pemkey.Bytes)
	if err != nil {
		return err
	}
	ecPrivateKey, ok := enrollmentPrivateKey.(*ecdsa.PrivateKey)
	if !ok {
		return fmt.Errorf("key not EC")
	}

	u.enrollmentKeys = &Enrollment{PrivateKey: privateKey, PublicKey: publicKey, EcdsaPrivateKey: ecPrivateKey}
	return nil
}

/**
 * deprecated.
 */
func (u *User) GetEnrollment() *Enrollment {
	return u.enrollmentKeys
}

/**
 * Gets a batch of TCerts to use for transaction. there is a 1-to-1 relationship between
 * TCert and Transaction. The TCert can be generated locally by the SDK using the user’s crypto materials.
 * @param {int} count how many in the batch to obtain
 * @param {[]string} attributes  list of attributes to include in the TCert
 * @return {[]tcert} An array of TCerts
 */
func (u *User) GenerateTcerts(count int, attributes []string) {

}
