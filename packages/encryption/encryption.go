package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"log"
	"net/http"
)

const keyFilepath = "./encryption_keys.txt"

func KeyGenerator() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		 log.Fatalf("failed to create new key: %v", err)
	}

	keyHex := hex.EncodeToString(key)
	if err := os.WriteFile(keyFilepath, []byte(keyHex), 0600); err != nil {
		log.Fatalf("Failed ot write key in file : %v", err)
	}

}

func RetrieveKey() ([]byte, error){
	keyHex, err := os.ReadFile(keyFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key from file: %v", err)
	}
	key, err := hex.DecodeString(string(keyHex))
	if err != nil {
		return nil, fmt.Errorf("failed to decode key provided : %v", err)
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
