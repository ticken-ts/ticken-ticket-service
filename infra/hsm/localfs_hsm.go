package hsm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"path"
)

const KeyStoreFolder = ".keystore"

type LocalFileSystemHSM struct {
	gcm       cipher.AEAD
	storePath string
}

func NewLocalFileSystemHSM(encryptionKey string, storePath string) (*LocalFileSystemHSM, error) {
	hsm := new(LocalFileSystemHSM)

	keyStorePath := path.Join(storePath, KeyStoreFolder)

	err := os.MkdirAll(keyStorePath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("unable to init keystore folder: %s", err.Error())
	}

	// generate a new aes cipher using our 32 byte long key
	cypher, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return nil, fmt.Errorf("unable to init hsm cypher: %s", err.Error())
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(cypher)
	if err != nil {
		return nil, fmt.Errorf("unable to init hsm gmc: %s", err.Error())
	}

	hsm.gcm = gcm
	hsm.storePath = keyStorePath

	return hsm, nil
}

func (hsm *LocalFileSystemHSM) Store(data []byte) (string, error) {
	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, hsm.gcm.NonceSize())

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error reading nonce: %s", err.Error())
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	encryptedData := hsm.gcm.Seal(nonce, nonce, data, nil)

	key := genKey()

	err := os.WriteFile(path.Join(hsm.storePath, key), encryptedData, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error during storing encrypted data: %s", err.Error())
	}

	return key, nil
}

func (hsm *LocalFileSystemHSM) Retrieve(key string) ([]byte, error) {
	encryptedData, err := os.ReadFile(path.Join(hsm.storePath, key))
	if err != nil {
		return nil, fmt.Errorf("key %s is not stored", key)
	}

	// why? only gods know
	nonceSize := hsm.gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("unable to retrieve value")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	data, err := hsm.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve value: %s", err.Error())
	}

	return data, nil
}

func genKey() string {
	//unixTimestamp := time.Now().Unix()
	//sum := sha256.Sum256([]byte(strconv.FormatInt(unixTimestamp, 10)))
	//return hex.EncodeToString(sum[:])
	return uuid.New().String()
}
