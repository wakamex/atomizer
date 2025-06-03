package derive

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// TradeModuleData represents the data for a trade order
type TradeModuleData struct {
	AssetAddress string
	SubID        *big.Int
	LimitPrice   *big.Int
	Amount       *big.Int
	MaxFee       *big.Int
	RecipientID  uint64
	IsBid        bool
}

// ToABIEncoded encodes the trade module data according to the contract ABI
func (t *TradeModuleData) ToABIEncoded() ([]byte, error) {
	// Define the ABI types: ["address", "uint", "int", "int", "uint", "uint", "bool"]
	addressType, _ := abi.NewType("address", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	intType, _ := abi.NewType("int256", "", nil)
	boolType, _ := abi.NewType("bool", "", nil)

	arguments := abi.Arguments{
		{Type: addressType},
		{Type: uintType},
		{Type: intType},
		{Type: intType},
		{Type: uintType},
		{Type: uintType},
		{Type: boolType},
	}

	// Pack the data
	return arguments.Pack(
		common.HexToAddress(t.AssetAddress),
		t.SubID,
		t.LimitPrice,
		t.Amount,
		t.MaxFee,
		new(big.Int).SetUint64(t.RecipientID),
		t.IsBid,
	)
}

// Action represents a signed action for Derive
type Action struct {
	SubaccountID       uint64
	Owner              common.Address
	Signer             common.Address
	SignatureExpirySec uint64
	Nonce              uint64
	ModuleAddress      common.Address
	ModuleData         *TradeModuleData
	DomainSeparator    [32]byte
	ActionTypehash     [32]byte
	Signature          []byte
}

// GetActionHash returns the keccak256 hash of the action
func (a *Action) GetActionHash() ([32]byte, error) {
	// Encode module data
	moduleDataEncoded, err := a.ModuleData.ToABIEncoded()
	if err != nil {
		return [32]byte{}, err
	}

	// Hash the module data
	moduleDataHash := crypto.Keccak256Hash(moduleDataEncoded)

	// Define the ABI types for action encoding
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{
		{Type: bytes32Type}, // actionTypehash
		{Type: uintType},    // subaccountId
		{Type: uintType},    // nonce
		{Type: addressType}, // moduleAddress
		{Type: bytes32Type}, // moduleDataHash
		{Type: uintType},    // signatureExpirySec
		{Type: addressType}, // owner
		{Type: addressType}, // signer
	}

	// Pack the action data
	encoded, err := arguments.Pack(
		a.ActionTypehash,
		new(big.Int).SetUint64(a.SubaccountID),
		new(big.Int).SetUint64(a.Nonce),
		a.ModuleAddress,
		moduleDataHash,
		new(big.Int).SetUint64(a.SignatureExpirySec),
		a.Owner,
		a.Signer,
	)
	if err != nil {
		return [32]byte{}, err
	}

	// Return the hash
	return crypto.Keccak256Hash(encoded), nil
}

// GetTypedDataHash returns the EIP-712 typed data hash
func (a *Action) GetTypedDataHash() ([32]byte, error) {
	actionHash, err := a.GetActionHash()
	if err != nil {
		return [32]byte{}, err
	}

	// Construct the EIP-712 message
	// 0x1901 + domainSeparator + actionHash
	message := append([]byte{0x19, 0x01}, a.DomainSeparator[:]...)
	message = append(message, actionHash[:]...)

	return crypto.Keccak256Hash(message), nil
}

// Sign signs the action with the given private key
func (a *Action) Sign(privateKey *ecdsa.PrivateKey) error {
	typedDataHash, err := a.GetTypedDataHash()
	if err != nil {
		return fmt.Errorf("failed to get typed data hash: %w", err)
	}

	// Sign the hash
	signature, err := crypto.Sign(typedDataHash[:], privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign: %w", err)
	}

	// Transform V from 0/1 to 27/28
	signature[64] += 27

	a.Signature = signature
	return nil
}

// ToJSON returns the JSON representation needed for the API
func (a *Action) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"subaccount_id":        a.SubaccountID,
		"nonce":                a.Nonce,
		"signer":               a.Signer.Hex(),
		"signature_expiry_sec": a.SignatureExpirySec,
		"signature":            "0x" + common.Bytes2Hex(a.Signature),
		"limit_price":          a.ModuleData.LimitPrice.String(),
		"amount":               a.ModuleData.Amount.String(),
		"max_fee":              a.ModuleData.MaxFee.String(),
	}
}

// DecimalToBigInt converts a decimal amount to big.Int with 18 decimals
func DecimalToBigInt(amount float64) *big.Int {
	// Convert to wei (18 decimals)
	amountWei := amount * 1e18
	return new(big.Int).SetInt64(int64(amountWei))
}

// GetActionNonce generates a nonce based on timestamp + random suffix
func GetActionNonce() uint64 {
	// Use current timestamp in milliseconds + 3 digit suffix
	// In production, add random 3 digits, for now just use 001
	return uint64(time.Now().UnixMilli())*1000 + 1
}
