package main

import (
	"fmt"
	"net/http"
	"smsforwarder/message"
	"smsforwarder/notify"
)

func main() {
	go message.ListenMessage()
	go notify.Notify()
	fmt.Println("开始监听短信......")
	http.HandleFunc("/api/sendMessage", message.SendMessage)
	http.HandleFunc("/api/getMessage", message.GetMessage)
	http.HandleFunc("/api/getNumber", message.GetInfo)
	http.ListenAndServe(":801", nil)
	select {}
}
