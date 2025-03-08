package main

import (
	"fmt"
	"net/http"
	"smsforwarder/message"
	"smsforwarder/notify"
	"smsforwarder/utils"
)

func main() {
	go message.ListenMessage()
	go notify.Notify()
	fmt.Println("开始监听短信......")
	http.HandleFunc("/api/sendMessage", message.SendMessage)
	http.HandleFunc("/api/getMessage", utils.GetMessageInfo)
	//http.HandleFunc("/api/getNumber", message.GetInfo)

	// exec cmd
	http.HandleFunc("/api/cmd", utils.TodoCMD)

	http.ListenAndServe(":801", nil)
	select {}
}
