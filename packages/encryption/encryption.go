// The encryption package provides encryption and decryption functionaliities
// Utiliizng the AES-256-GCM Algorithm. It allows for key generation as well as secure key storage,
// as well as offering a simple API for encryption and decryption of data.
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const keyFilepath = "./encryption_keys.txt"

// The Key generator method generates a AES-256 encryption key and stores it in a text file with predefined roles.
// The method returns an error for when the generation of the key has failed or writing it into the file has failed.
func KeyGenerator()  error{
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		 return fmt.Errorf("failed to create new key: %v", err)
	}

	keyHex := hex.EncodeToString(key)
	if err := ioutil.WriteFile(keyFilepath, []byte(keyHex), 0600); err != nil {
		return fmt.Errorf("failed to write key to file: %v", err)
	}
	return nil
}

func RetrieveKey() ([]byte, error){
	keyHex, err := ioutil.ReadFile(keyFilepath)
	if err != nil {
		if os.IsNotExist(err){
			fmt.Println("Key file not found, generating a new key")
			if err := KeyGenerator(); err != nil {
				return nil, fmt.Errorf("Failed to generate key: %v", err)
			}
			keyHex, err = ioutil.ReadFile(keyFilepath)
			if err != nil {
				return nil, fmt.Errorf("failed to read the new key : %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read key from file: %v", err)
		}
	}
	key, err := hex.DecodeString(string(keyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %v", err)
	}
	return key, nil
}

func DataEncryption(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return cipherText, nil
}

func DataDecryption(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < gcm.NonceSize() {
		return nil, errors.New("cipherText too short")
	}

	nonce, cipherText := cipherText[:gcm.NonceSize()], cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func HandleEncryption(w http.ResponseWriter, r *http.Request) {
	action := r.Header.Get("X-Encryption-Action")
	key, err := RetrieveKey()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve encryption key: %v", err), http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var processedData []byte
	switch action {
	case "encrypt":
		processedData, err = DataEncryption(body, key)
	case "decrypt":
		processedData, err = DataDecryption(body, key)
	default:
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to %s data : %v", action, err), http.StatusInternalServerError)
		return
	}
	w.Write(processedData)
}
