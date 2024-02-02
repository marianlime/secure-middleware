package crypto

import(
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func KeyGenerator() ([]byte, error){
	key := make([]byte, 32)
	if _,err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to create new key: %w", err)
	}
	return key, nil
}

func DataEncryption(plainText []byte, key []byte) ([]byte, error){
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

func DataDecryption(cipherText []byte, key []byte) ([]byte, error){
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < gcm.NonceSize(){
		return nil, errors.New("cipherText too short")
	}

	nonce, cipherText := cipherText[:gcm.NonceSize()], cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil{
		return nil, err
	}

	return plainText, nil
}