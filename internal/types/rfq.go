package types

// RFQNotification represents the structure of an incoming RFQ message from the server
type RFQNotification struct {
	JsonRPC string       `json:"jsonrpc"`
	ID      string       `json:"id"`
	Result  RFQResult    `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	Method  string       `json:"method,omitempty"`
	Params  RFQResult    `json:"params,omitempty"`
}

// RFQResult contains the actual details of the RFQ
type RFQResult struct {
	ID         string `json:"id,omitempty"`
	Asset      string `json:"asset"`
	AssetName  string `json:"assetName,omitempty"`
	OptionType string `json:"optionType,omitempty"`
	Bid        bool   `json:"bid"`
	Ask        bool   `json:"ask"`
	Quantity   string `json:"quantity"`
	Direction  string `json:"direction,omitempty"`
	LegID      string `json:"legId,omitempty"`
	MaxFee     string `json:"maxFee,omitempty"`
	Subaccount int    `json:"subaccount,omitempty"`
	ChainID    int    `json:"chainId,omitempty"`
	Expiry     int    `json:"expiry,omitempty"`
	IsPut      bool   `json:"isPut,omitempty"`
	Strike     string `json:"strike,omitempty"`
	IsTakerBuy bool   `json:"isTakerBuy,omitempty"`
}

// RFQConfirmation represents the structure of an incoming RFQ confirmation message
type RFQConfirmation struct {
	ID              string  `json:"id,omitempty"`	
	Maker           string  `json:"maker"`
	AssetAddress    string  `json:"assetAddress"`
	ChainID         int     `json:"chainId"`
	Expiry          int     `json:"expiry"`
	IsPut           bool    `json:"isPut"`
	Nonce           string  `json:"nonce"`
	Price           string  `json:"price"`
	Quantity        string  `json:"quantity"`
	QuoteNonce      string  `json:"quoteNonce"`
	QuoteValidUntil int     `json:"quoteValidUntil"`
	QuoteSignature  string  `json:"quoteSignature"`
	Strike          string  `json:"strike"`
	Taker           string  `json:"taker"`
	IsTakerBuy      bool    `json:"isTakerBuy"`
	Signature       string  `json:"signature"`
	ValidUntil      int     `json:"validUntil"`
	CreatedAt       int     `json:"createdAt,omitempty"`
	APR             float64 `json:"apr,omitempty"`
}