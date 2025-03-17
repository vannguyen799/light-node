package utils

import (
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
