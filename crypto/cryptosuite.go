package cryptosuite

/**
 * Abstract class for a suite of crypto algorithms used by the SDK to perform encryption,
 * decryption and secure hashing. A complete suite includes libraries for asymmetric
 * keys (such as ECDSA or RSA), symmetric keys (such as AES) and secure hash (such as
 * SHA2/3).
 */
type CryptoSuite interface {
	/**
	 * Generate a key
	 *
	 * @param {ephemeral} ephemeral: true if the key to generate has to be ephemeral
	 * @returns {Key} Key instance
	 */
	GenerateKey(ephemeral bool) (interface{}, error)

	/**
	 * Derives a key from k using opts.
	 * @param {Key} key the source key
	 * @returns {Key} derived key
	 */
	DeriveKey(key interface{}) (interface{}, error)

	/**
	 * Imports a key from its raw representation using opts.
	 * @param {[]byte} raw Raw bytes of the key to import
	 * @param {string} algorithm:an identifier for the algorithm to be used
	 * @returns {Key} An instance of the Key wrapping the raw key bytes
	 */
	ImportKey(raw []byte, algorithm string) (interface{}, error)

	/**
	 * Returns the key this CSP associates to the Subject Key Identifier ski.
	 *
	 * @param {[]byte} ski Subject Key Identifier specific to a Crypto Suite implementation
	 * @returns {Key} Key instance  corresponding to the ski
	 */
	GetKey(ski []byte) (interface{}, error)

	/**
	 * Hashes messages msg using options opts.
	 *
	 * @param {[]byte} msg Source message to be hashed
	 * @param {string} algorithm: an identifier for the algorithm to be used, such as "SHA3"
	 * @returns {[]byte} The hashed digest
	 */
	Hash(msg []byte, algorithm string) ([]byte, error)

	/**
	 * Signs digest using key k.
	 * The opts argument should be appropriate for the algorithm used.
	 *
	 * @param {Key} key Signing key (private key)
	 * @param {[]byte} digest The message digest to be signed
	 * @param {string} algorithm: the function to use to hash
	 * @param {int} securityLevel: the security Level
	 * @returns {[]byte} the resulting signature
	 */
	Sign(key interface{}, digest []byte, algorithm string, securityLevel int) ([]byte, error)

	/**
	 * Verifies signature against key k and digest
	 * The opts argument should be appropriate for the algorithm used.
	 *
	 * @param {Key} key Signing verification key (public key)
	 * @param {[]byte} signature The signature to verify
	 * @param {[]byte} digest The digest that the signature was created for
	 * @returns {bool} true if the signature verifies successfully
	 */
	Verify(key interface{}, signature []byte, digest []byte) (bool, error)

	/**
	 * Encrypts plaintext using key k.
	 * The opts argument should be appropriate for the algorithm used.
	 *
	 * @param {Key} key Encryption key (public key)
	 * @param {[]byte} plainText Plain text to encrypt
	 * @param {Object} opts Encryption options
	 * @returns {[]byte} Cipher text after encryption
	 */
	Encrypt(key interface{}, plaintext []byte) ([]byte, error)

	/**
	 * Decrypts ciphertext using key k.
	 * The opts argument should be appropriate for the algorithm used.
	 *
	 * @param {Key} key Decryption key (private key)
	 * @param {[]byte} cipherText Cipher text to decrypt
	 * @returns {[]byte} Plain text after decryption
	 */
	Decrypt(key interface{}, cipherText []byte) ([]byte, error)
}
