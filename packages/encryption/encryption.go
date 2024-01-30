package crypto

import(
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func keyGenerator() ([]byte, error){
	key := make([]byte, 32)
	if _,err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

