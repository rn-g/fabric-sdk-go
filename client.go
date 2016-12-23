package fabric_sdk_go

import (
	"fmt"

	crypto "github.com/hyperledger/fabric-sdk-go/crypto"
	keysstore "github.com/hyperledger/fabric-sdk-go/keysstore"
)

/**
 * Main interaction handler with end user. A client instance provides a handler to interact
 * with a network of peers, orderers and optionally member services. An application using the
 * SDK may need to interact with multiple networks, each through a separate instance of the Client.
 *
 * Each client when initially created should be initialized with configuration data from the
 * consensus service, which includes a list of trusted roots, orderer certificates and IP addresses,
 * and a list of peer certificates and IP addresses that it can access. This must be done out of band
 * as part of bootstrapping the application environment. It is also the responsibility of the application
 * to maintain the configuration of a client as the SDK does not persist this object.
 *
 * Each Client instance can maintain several {@link Chain} instances representing channels and the associated
 * private ledgers.
 *
 *
 */
type Client struct {
	chains      map[string]*Chain
	cryptoSuite crypto.CryptoSuite
	stateStore  *keysstore.KeyValueStore
	userContext *User
}

/**
 * Returns a Client instance
 */
func NewClient() *Client {
	chains := make(map[string]*Chain)
	c := &Client{chains: chains, cryptoSuite: nil, stateStore: nil, userContext: nil}
	return c
}

/**
 * Returns a chain instance with the given name. This represents a channel and its associated ledger
 * (as explained above), and this call returns an empty object. To initialize the chain in the blockchain network,
 * a list of participating endorsers and orderer peers must be configured first on the returned object.
 * @param {string} name The name of the chain.  Recommend using namespaces to avoid collision.
 * @returns {Chain} The uninitialized chain instance.
 * @returns {Error} if the chain by that name already exists in the application's state store
 */
func (c *Client) NewChain(name string) (*Chain, error) {
	if _, ok := c.chains[name]; ok {
		return nil, fmt.Errorf("Chain %s already exists", name)
	}
	var err error
	c.chains[name], err = NewChain(name, c)
	if err != nil {
		return nil, err
	}
	return c.chains[name], nil

}

/**
 * Get a {@link Chain} instance from the state storage. This allows existing chain instances to be saved
 * for retrieval later and to be shared among instances of the application. Note that it’s the
 * application/SDK’s responsibility to record the chain information. If an application is not able
 * to look up the chain information from storage, it may call another API that queries one or more
 * Peers for that information.
 * @param {string} name The name of the chain.
 * @returns {Chain} The chain instance
 */
func (c *Client) GetChain(name string) *Chain {
	return c.chains[name]
}

/**
 * This is a network call to the designated Peer(s) to discover the chain information.
 * The target Peer(s) must be part of the chain to be able to return the requested information.
 * @param {string} name The name of the chain.
 * @param {[]Peer} peers Array of target Peers to query.
 * @returns {Chain} The chain instance for the name or error if the target Peer(s) does not know
 * anything about the chain.
 */
func (c *Client) QueryChainInfo(name string, peers []*Peer) (*Chain, error) {
	return nil, fmt.Errorf("Not implemented yet")
}

/**
 * The SDK should have a built-in key value store implementation (suggest a file-based implementation to allow easy setup during
 * development). But production systems would want a store backed by database for more robust storage and clustering,
 * so that multiple app instances can share app state via the database (note that this doesn’t necessarily make the app stateful).
 * This API makes this pluggable so that different store implementations can be selected by the application.
 */
func (c *Client) SetStateStore(stateStore *keysstore.KeyValueStore) {
	c.stateStore = stateStore
}

/**
 * A convenience method for obtaining the state store object in use for this client.
 */
func (c *Client) GetStateStore() *keysstore.KeyValueStore {
	return c.GetStateStore()
}

/**
 * A convenience method for obtaining the state store object in use for this client.
 */
func (c *Client) SetCryptoSuite(cryptoSuite crypto.CryptoSuite) {
	c.cryptoSuite = cryptoSuite
}

/**
 * A convenience method for obtaining the CryptoSuite object in use for this client.
 */
func (c *Client) GetCryptoSuite() crypto.CryptoSuite {
	return c.cryptoSuite
}

/**
 * Sets an instance of the User class as the security context of this client instance. This user’s credentials (ECert) will be
 * used to conduct transactions and queries with the blockchain network. Upon setting the user context, the SDK saves the object
 * in a persistence cache if the “state store” has been set on the Client instance. If no state store has been set,
 * this cache will not be established and the application is responsible for setting the user context again when the application
 * crashed and is recovered.
 */
func (c *Client) SetUserContext(user *User) {
	c.userContext = user
}

/**
 * The client instance can have an optional state store. The SDK saves enrolled users in the storage which can be accessed by
 * authorized users of the application (authentication is done by the application outside of the SDK).
 * This function attempts to load the user by name from the local storage (via the KeyValueStore interface).
 * The loaded user object must represent an enrolled user with a valid enrollment certificate signed by a trusted CA
 * (such as the COP server).
 */
func (c *Client) GetUserContext() *User {
	return c.userContext
}
