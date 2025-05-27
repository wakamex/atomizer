package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

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
	LimitPrice         int64
	Amount             int64
	MaxFee             int64
	RecipientID        uint64
	IsBid              bool
	Signature          string
}

// Sign signs the action using EIP-712
func (a *DeriveAction) Sign(privateKey *ecdsa.PrivateKey) error {
	// Protocol constants for mainnet
	domainSeparator := common.HexToHash("0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105")
	actionTypehash := common.HexToHash("0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17")
	
	// Encode module data
	subID, _ := new(big.Int).SetString(a.SubID, 10)
	moduleData, err := encodeTradeModuleData(
		common.HexToAddress(a.AssetAddress),
		subID,
		big.NewInt(a.LimitPrice),
		big.NewInt(a.Amount),
		big.NewInt(a.MaxFee),
		a.RecipientID,
		a.IsBid,
	)
	if err != nil {
		return err
	}
	
	// Hash module data
	moduleDataHash := crypto.Keccak256Hash(moduleData)
	
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
	if err != nil {
		return err
	}
	
	// Create typed data hash
	actionHash := crypto.Keccak256Hash(actionData)
	message := append([]byte{0x19, 0x01}, domainSeparator.Bytes()...)
	message = append(message, actionHash.Bytes()...)
	typedDataHash := crypto.Keccak256Hash(message)
	
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

// PlaceDeriveOrder places an order directly using Derive API
func PlaceDeriveOrder(instrumentName string, side string, orderType string, price float64, amount float64, privateKey string, deriveWalletAddress string, subaccountID uint64) (*DeriveOrderResponse, error) {
	auth, err := NewDeriveAuth(privateKey)
	if err != nil {
		return nil, err
	}
	
	// Get instrument details
	instrument, err := FetchDeriveInstrumentDetails(instrumentName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch instrument: %w", err)
	}
	
	// Create signed action
	action := &DeriveAction{
		SubaccountID:       subaccountID,
		Owner:              deriveWalletAddress,
		Signer:             auth.GetAddress(),
		SignatureExpirySec: time.Now().Unix() + 3600, // 1 hour
		Nonce:              uint64(time.Now().UnixMilli())*1000 + 1,
		ModuleAddress:      "0x87F2863866D85E3192a35A73b388BD625D83f2be", // Trade module
		AssetAddress:       instrument.BaseAssetAddress,
		SubID:              instrument.BaseAssetSubID,
		LimitPrice:         int64(price * 1e18),
		Amount:             int64(amount * 1e18),
		MaxFee:             int64(1e18), // 1 USDC max fee
		RecipientID:        subaccountID,
		IsBid:              side == "buy",
	}
	
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
		"max_fee":            "1000",
	}
	
	jsonData, err := json.Marshal(orderReq)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", "https://api.lyra.finance/private/order", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	authHeaders, _ := auth.GetAuthHeaders(deriveWalletAddress)
	for k, v := range authHeaders {
		req.Header.Set(k, v)
	}
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[Derive Order] Response: %s", string(body))
	
	var orderResp DeriveOrderResponse
	json.Unmarshal(body, &orderResp)
	if orderResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", orderResp.Error.Message)
	}
	
	return &orderResp, nil
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