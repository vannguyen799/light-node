package utils

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetWalletAddress() (*string, error) {
	privKey := GetEnv("PRIVATE_KEY", "")

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.PublicKey

	walletAddress := crypto.PubkeyToAddress(publicKey).Hex()
	return &walletAddress, nil
}

func SignMessage(message string) (*string, error) {
	privKey := GetEnv("PRIVATE_KEY", "")

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	data := []byte(message)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	hexSign := hexutil.Encode(signature)
	return &hexSign, nil
}

func VerifyMessage(sign string, message string, expectedAddress string) error {
	messageHash := crypto.Keccak256Hash([]byte(message))
	signature, err := hexutil.Decode(sign)
	if err != nil {
		log.Printf("invalid signature")
		return err
	}
	if len(signature) != 65 {
		log.Println("invalid signature length")
		return fmt.Errorf("invalid signature length")
	}

	// signature[64] -= 27

	pubKey, err := crypto.SigToPub(messageHash.Bytes(), signature)
	if err != nil {
		log.Println("Failed to recover public key:", err)
		return fmt.Errorf("Failed to recover public key:", err)
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	log.Printf("recovered address: %s, expected address: %s", recoveredAddress, expectedAddress)
	return nil
}
