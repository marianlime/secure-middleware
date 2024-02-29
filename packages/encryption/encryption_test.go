package encryption

import (
	"encoding/base64"
	"encoding/hex"
	"os"
	"testing"
)


func TestKeyGenerator(t *testing.T){
	os.Remove(keyFilepath)
	KeyGenerator()

	content, err := os.ReadFile(keyFilepath)
	if err != nil {
		t.Fatalf("failed to read key fil: %v", err)
	}
	if _, err := hex.DecodeString(string(content)); err != nil {
		t.Errorf("Stored key is not a valid hex: %v", err)
	}
}

func TestDataEncryptionAndDecryption(t *testing.T){
	originalText := []byte("Message")
	KeyGenerator()
	key, _ := RetrieveKey()
	defer os.Remove(keyFilepath)

	cipherText, err := DataEncryption(originalText, key)
	if err != nil {
		t.Fatalf("Data Encryption has failed : %v", err)
	}
	decryptedText, err := DataDecryption(cipherText, key)
	if err != nil{
		t.Fatalf("Data Decryption has failed: %v", err)
	}

	if string(decryptedText) != string(originalText) {
		t.Errorf("Decryption has failed and does not match the sent request. got %s, instead of %s", decryptedText, originalText)
	}
}

func TestInvalidKeyLength(t *testing.T){
	originalText := []byte("Message")
	invalidKey := []byte("notakey")

	_, err := DataEncryption(originalText, invalidKey)
	if err == nil {
		t.Errorf("DataEncryption should fail with key provided")
	}

	_, err = DataDecryption(originalText, invalidKey)
	if err == nil{
		t.Errorf("DataDecryption has failed with key provided")
	}
}

func TestUniqueCiphertextEncryption (t *testing.T){
	plainText := []byte("Secret Message")
	KeyGenerator()
	key, err := RetrieveKey()
	defer os.Remove(keyFilepath)
	if err != nil{
		t.Fatalf("Failed to retrieve encryption key: %v", err)
	}

	ciphertexts := make(map[string]bool)
	iterations := 100

	for i := 0; i <iterations; i++{
		ciphertext, err := DataEncryption(plainText, key)
		if err != nil{
			t.Fatalf("Encryption failed: %v", err)
		}
		CipheredText := base64.StdEncoding.EncodeToString(ciphertext)
		if _, exists := ciphertexts[CipheredText]; exists {
			t.Fatalf("Duplicate cipheredtext has been found")
		}
		ciphertexts[CipheredText] = true
	}
}