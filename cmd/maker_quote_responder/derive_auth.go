package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// DeriveAuth handles authentication for Derive API
type DeriveAuth struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

// NewDeriveAuth creates a new Derive authenticator from a private key
func NewDeriveAuth(privateKeyHex string) (*DeriveAuth, error) {
	// Remove 0x prefix if present
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key hex: %w", err)
	}
	
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	
	return &DeriveAuth{
		privateKey: privateKey,
		address:    address,
	}, nil
}

// GetAuthHeaders generates authentication headers for Derive API requests
func (d *DeriveAuth) GetAuthHeaders(deriveWalletAddress string) (map[string]string, error) {
	// Current UTC timestamp in milliseconds
	timestamp := time.Now().UTC().UnixMilli()
	timestampStr := strconv.FormatInt(timestamp, 10)
	
	// Sign the timestamp
	signature, err := d.SignMessage(timestampStr)
	if err != nil {
		return nil, fmt.Errorf("failed to sign timestamp: %w", err)
	}
	
	headers := map[string]string{
		"X-LyraWallet":    deriveWalletAddress,
		"X-LyraTimestamp": timestampStr,
		"X-LyraSignature": signature,
	}
	
	return headers, nil
}

// SignMessage signs a message using Ethereum's personal_sign method
func (d *DeriveAuth) SignMessage(message string) (string, error) {
	// Ethereum personal_sign prefixes the message
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	
	debugLog("[Derive Auth] Signing message: '%s'", message)
	debugLog("[Derive Auth] Prefixed message: '%s'", prefixedMessage)
	debugLog("[Derive Auth] Message hash: %s", hash.Hex())
	
	signature, err := crypto.Sign(hash.Bytes(), d.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %w", err)
	}
	
	// Transform V from 0/1 to 27/28 according to the yellow paper
	signature[64] += 27
	
	sigHex := hex.EncodeToString(signature)
	debugLog("[Derive Auth] Raw signature: %s", sigHex)
	
	return "0x" + sigHex, nil
}

// GetAddress returns the Ethereum address for this private key
func (d *DeriveAuth) GetAddress() string {
	return d.address.Hex()
}

// GetPrivateKey returns the private key for signing
func (d *DeriveAuth) GetPrivateKey() *ecdsa.PrivateKey {
	return d.privateKey
}

// SignOrderPayload signs an order payload for self-custodial requests
// This is used for the second auth step in order placement
func (d *DeriveAuth) SignOrderPayload(payload interface{}) (string, error) {
	// TODO: Implement EIP-712 structured data signing for order payloads
	// This requires the exact structure from Derive's matching contract
	// For now, this is a placeholder
	return "", fmt.Errorf("order payload signing not yet implemented - requires EIP-712 structure")
}