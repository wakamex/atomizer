package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Global debug mode flag
var debugMode bool

// SetDebugMode sets the global debug mode
func SetDebugMode(enabled bool) {
	debugMode = enabled
}

// debugLog logs a message only if debug mode is enabled
func debugLog(format string, args ...interface{}) {
	if debugMode {
		log.Printf(format, args...)
	}
}

// DeriveInstrumentDetails represents instrument details from Derive API
type DeriveInstrumentDetails struct {
	InstrumentName    string `json:"instrument_name"`
	BaseAssetAddress  string `json:"base_asset_address"`
	BaseAssetSubID    string `json:"base_asset_sub_id"`
}

// DeriveAction represents a signed action for order placement
type DeriveAction struct {
	SubaccountID       uint64
	Owner              string
	Signer             string
	SignatureExpirySec int64
	Nonce              uint64
	ModuleAddress      string
	AssetAddress       string
	SubID              string
	LimitPrice         *big.Int
	Amount             *big.Int
	MaxFee             *big.Int
	RecipientID        uint64
	IsBid              bool
	Signature          string
}

// Sign signs the action using EIP-712
func (a *DeriveAction) Sign(privateKey *ecdsa.PrivateKey) error {
	// Protocol constants for mainnet
	domainSeparator := common.HexToHash("0xd96e5f90797da7ec8dc4e276260c7f3f87fedf68775fbe1ef116e996fc60441b")
	actionTypehash := common.HexToHash("0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17")
	
	debugLog("[EIP-712] Using domain separator: %s", domainSeparator.Hex())
	debugLog("[EIP-712] Using action typehash: %s", actionTypehash.Hex())
	
	// Encode module data
	subID, _ := new(big.Int).SetString(a.SubID, 10)
	moduleData, err := encodeTradeModuleData(
		common.HexToAddress(a.AssetAddress),
		subID,
		a.LimitPrice,
		a.Amount,
		a.MaxFee,
		a.RecipientID,
		a.IsBid,
	)
	if err != nil {
		return err
	}
	
	// Hash module data
	moduleDataHash := crypto.Keccak256Hash(moduleData)
	debugLog("[EIP-712] Module data hash: %s", moduleDataHash.Hex())
	
	// Encode action
	actionData, err := encodeAction(
		actionTypehash,
		a.SubaccountID,
		a.Nonce,
		common.HexToAddress(a.ModuleAddress),
		moduleDataHash,
		a.SignatureExpirySec,
		common.HexToAddress(a.Owner),
		common.HexToAddress(a.Signer),
	)
	debugLog("[EIP-712] Action data encoded, length: %d", len(actionData))
	if err != nil {
		return err
	}
	
	// Create typed data hash
	actionHash := crypto.Keccak256Hash(actionData)
	debugLog("[EIP-712] Action hash: %s", actionHash.Hex())
	
	message := append([]byte{0x19, 0x01}, domainSeparator.Bytes()...)
	message = append(message, actionHash.Bytes()...)
	typedDataHash := crypto.Keccak256Hash(message)
	debugLog("[EIP-712] Final typed data hash to sign: %s", typedDataHash.Hex())
	
	// Sign
	signature, err := crypto.Sign(typedDataHash.Bytes(), privateKey)
	if err != nil {
		return err
	}
	
	// Transform V from 0/1 to 27/28
	signature[64] += 27
	
	a.Signature = "0x" + hex.EncodeToString(signature)
	return nil
}

func encodeTradeModuleData(asset common.Address, subID, limitPrice, amount, maxFee *big.Int, recipientID uint64, isBid bool) ([]byte, error) {
	addressType, _ := abi.NewType("address", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	intType, _ := abi.NewType("int256", "", nil)
	boolType, _ := abi.NewType("bool", "", nil)
	
	args := abi.Arguments{
		{Type: addressType},
		{Type: uintType},
		{Type: intType},
		{Type: intType},
		{Type: uintType},
		{Type: uintType},
		{Type: boolType},
	}
	
	return args.Pack(asset, subID, limitPrice, amount, maxFee, big.NewInt(int64(recipientID)), isBid)
}

func encodeAction(typehash common.Hash, subaccountID, nonce uint64, moduleAddr common.Address, moduleDataHash common.Hash, expiry int64, owner, signer common.Address) ([]byte, error) {
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	
	args := abi.Arguments{
		{Type: bytes32Type},
		{Type: uintType},
		{Type: uintType},
		{Type: addressType},
		{Type: bytes32Type},
		{Type: uintType},
		{Type: addressType},
		{Type: addressType},
	}
	
	return args.Pack(typehash, big.NewInt(int64(subaccountID)), big.NewInt(int64(nonce)), moduleAddr, moduleDataHash, big.NewInt(expiry), owner, signer)
}

// DeriveOrderResponse represents the order response from Derive API
type DeriveOrderResponse struct {
	Result struct {
		OrderID        string  `json:"order_id"`
		InstrumentName string  `json:"instrument_name"`
		Side           string  `json:"side"`
		OrderType      string  `json:"order_type"`
		Price          float64 `json:"price"`
		Amount         float64 `json:"amount"`
		FilledAmount   float64 `json:"filled_amount"`
		Status         string  `json:"status"`
		CreatedAt      int64   `json:"created_at"`
	} `json:"result"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// PlaceDeriveOrder places an order directly using Derive WebSocket API
func PlaceDeriveOrder(instrumentName string, side string, orderType string, price float64, amount float64, privateKey string, deriveWalletAddress string, subaccountID uint64) (*DeriveOrderResponse, error) {
	// Create WebSocket client
	wsClient, err := NewDeriveWSClient(privateKey, deriveWalletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket client: %w", err)
	}
	defer wsClient.Close()
	
	auth, err := NewDeriveAuth(privateKey)
	if err != nil {
		return nil, err
	}
	
	// Get instrument details
	instrument, err := FetchDeriveInstrumentDetails(instrumentName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch instrument: %w", err)
	}
	
	log.Printf("[Derive Order] Instrument details - Name: %s, BaseAssetAddress: %s, SubID: %s", 
		instrument.InstrumentName, instrument.BaseAssetAddress, instrument.BaseAssetSubID)
	
	// Create signed action
	signerEOA := auth.GetAddress()
	log.Printf("[Derive Order] Creating action - Owner: %s (Derive wallet), Signer: %s (EOA)", deriveWalletAddress, signerEOA)
	
	action := &DeriveAction{
		SubaccountID:       subaccountID,
		Owner:              deriveWalletAddress, // Derive wallet is the owner!
		Signer:             signerEOA,            // EOA is the signer
		SignatureExpirySec: time.Now().Unix() + 3600, // 1 hour
		Nonce:              uint64(time.Now().UnixMilli())*1000 + 1,
		ModuleAddress:      "0xB8D20c2B7a1Ad2EE33Bc50eF10876eD3035b5e7b", // Trade module
		AssetAddress:       instrument.BaseAssetAddress,
		SubID:              instrument.BaseAssetSubID,
		// Convert price and amount to big.Int properly to avoid overflow
		LimitPrice:         func() *big.Int { v, _ := new(big.Float).Mul(big.NewFloat(price), big.NewFloat(1e18)).Int(nil); return v }(),
		Amount:             func() *big.Int { v, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(nil); return v }(),
		MaxFee:             new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18)), // 1000 USDC max fee (matching Python SDK)
		RecipientID:        subaccountID,
		IsBid:              side == "buy",
	}
	
	log.Printf("[Derive Order] Action details before signing:")
	log.Printf("  SubaccountID: %d", action.SubaccountID)
	log.Printf("  Owner: %s", action.Owner)
	log.Printf("  Signer: %s", action.Signer)
	log.Printf("  Nonce: %d", action.Nonce)
	log.Printf("  LimitPrice: %s", action.LimitPrice.String())
	log.Printf("  Amount: %s", action.Amount.String())
	log.Printf("  MaxFee: %s", action.MaxFee.String())
	log.Printf("  IsBid: %v", action.IsBid)
	
	if err := action.Sign(auth.privateKey); err != nil {
		return nil, err
	}
	
	// Prepare order request
	orderReq := map[string]interface{}{
		"instrument_name":      instrumentName,
		"direction":           side,
		"order_type":         orderType,
		"time_in_force":      "gtc",
		"mmp":                false,
		"subaccount_id":      subaccountID,
		"nonce":              action.Nonce,
		"signer":             action.Signer,
		"signature_expiry_sec": action.SignatureExpirySec,
		"signature":          action.Signature,
		"limit_price":        fmt.Sprintf("%.6f", price),
		"amount":             fmt.Sprintf("%.6f", amount),
		"max_fee":            "1000", // Must match what we signed (1000 USDC)
	}
	
	// Submit order via WebSocket
	orderResp, err := wsClient.SubmitOrder(orderReq)
	if err != nil {
		return nil, err
	}
	
	// Query open orders to verify
	openOrders, err := wsClient.GetOpenOrders(subaccountID)
	if err != nil {
		log.Printf("[Derive Order] Failed to query open orders: %v", err)
	} else {
		log.Printf("[Derive Order] Open orders after submission: %d orders", len(openOrders))
		for i, order := range openOrders {
			if instrumentName, ok := order["instrument_name"].(string); ok {
				if side, ok := order["direction"].(string); ok {
					if amount, ok := order["amount"].(float64); ok {
						if price, ok := order["price"].(float64); ok {
							log.Printf("[Derive Order] Order %d: %s %s %.4f @ %.2f", i+1, instrumentName, side, amount, price)
						}
					}
				}
			}
		}
	}
	
	return orderResp, nil
}

// FetchDeriveInstrumentDetails fetches instrument details from Derive API
func FetchDeriveInstrumentDetails(instrumentName string) (*DeriveInstrumentDetails, error) {
	url := "https://api.lyra.finance/public/get_instrument"
	payload := map[string]interface{}{
		"instrument_name": instrumentName,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Result DeriveInstrumentDetails `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return &result.Result, nil
}