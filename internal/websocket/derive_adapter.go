package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wakamex/atomizer/internal/exchange/derive"
)

// DeriveAuthAdapter adapts Derive authentication to the generic AuthProvider interface
type DeriveAuthAdapter struct {
	auth   *derive.DeriveAuth
	wallet string
}

// NewDeriveAuthAdapter creates a new Derive auth adapter
func NewDeriveAuthAdapter(privateKey, wallet string) (*DeriveAuthAdapter, error) {
	auth, err := derive.NewDeriveAuth(privateKey)
	if err != nil {
		return nil, err
	}
	
	return &DeriveAuthAdapter{
		auth:   auth,
		wallet: wallet,
	}, nil
}

// RequiresAuth returns true as Derive requires authentication
func (d *DeriveAuthAdapter) RequiresAuth() bool {
	return true
}

// Authenticate performs Derive-specific authentication
func (d *DeriveAuthAdapter) Authenticate(conn *websocket.Conn) error {
	// Send login request
	loginReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "login",
		"method":  "public/login",
		"params": map[string]interface{}{
			"wallet": d.wallet,
		},
	}
	
	if err := conn.WriteJSON(loginReq); err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}
	
	// Read login response
	var loginResp struct {
		ID     string `json:"id"`
		Result struct {
			Nonce int64 `json:"nonce"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	
	if err := conn.ReadJSON(&loginResp); err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}
	
	if loginResp.Error != nil {
		return fmt.Errorf("login error: %s", loginResp.Error.Message)
	}
	
	// Sign the nonce (convert to string)
	nonceStr := fmt.Sprintf("%d", loginResp.Result.Nonce)
	signature, err := d.auth.SignMessage(nonceStr)
	if err != nil {
		return fmt.Errorf("failed to sign nonce: %w", err)
	}
	
	// Send authenticate request
	authReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "authenticate",
		"method":  "public/authenticate",
		"params": map[string]interface{}{
			"wallet":    d.wallet,
			"signature": signature,
		},
	}
	
	if err := conn.WriteJSON(authReq); err != nil {
		return fmt.Errorf("failed to send authenticate request: %w", err)
	}
	
	// Read auth response
	var authResp struct {
		ID     string `json:"id"`
		Result struct {
			AccessToken string `json:"access_token"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	
	if err := conn.ReadJSON(&authResp); err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}
	
	if authResp.Error != nil {
		return fmt.Errorf("auth error: %s", authResp.Error.Message)
	}
	
	log.Printf("[Derive Auth] Successfully authenticated")
	return nil
}

// CreateDeriveWebSocketClient creates a WebSocket client configured for Derive
func CreateDeriveWebSocketClient(privateKey, wallet string, handler MessageHandler) (*Client, error) {
	authAdapter, err := NewDeriveAuthAdapter(privateKey, wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth adapter: %w", err)
	}
	
	config := ClientConfig{
		URL:            "wss://api.lyra.finance/ws",
		Name:           "Derive",
		AuthProvider:   authAdapter,
		MessageHandler: handler,
	}
	
	return NewClient(config), nil
}