package keysstore

/**
 * Abstract class for a Key-Value store. The Chain class uses this store
 * to save sensitive information such as authenticated user's private keys,
 * certificates, etc.
 *
 */
type KeyValueStore interface {
	/**
	 * Get the value associated with name.
	 *
	 * @param {string} name of the key
	 * @returns {[]byte}
	 */
	GetValue(name string) ([]byte, error)

	/**
	 * Set the value associated with name.
	 * @param {string} name of the key to save
	 * @param {[]byte} value to save
	 */
	SetValue(name string, value []byte) error
}
