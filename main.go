package main

import (
	"fmt"
	"log"
	"middleware/packages/authentication"
	"middleware/packages/encryption"
	"middleware/packages/logmon"
	"net/http"
)

func costumerEndpointHandler(w http.ResponseWriter, r  *http.Request){

}

func main(){
	encryption.LoadEnv()

	key, err := encryption.KeyGenerator()
	if err != nil {
		log.Fatalf("Error generating key : %v", err)
	}

	encryptedData, err := encryption.DataEncryption([]byte(""), key)
	if err != nil {
		log.Fatalf("Error encrypting data: %v", err)
	}
	
	http.HandleFunc("/encrypt", encryption.HandleEncryption)
	}