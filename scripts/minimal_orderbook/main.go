package main

import ("fmt";"github.com/gorilla/websocket")

func main() {
    conn, _, _ := websocket.DefaultDialer.Dial("wss://api.lyra.finance/ws", nil)
    defer conn.Close()
    conn.WriteJSON(map[string]interface{}{"method": "subscribe","params": map[string]interface{}{"channels": []string{"orderbook.ETH-PERP.1.10"}}})
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Error:", err)
            break
        }
        fmt.Println(string(msg))
    }
}