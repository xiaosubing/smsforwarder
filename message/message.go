package message

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/godbus/dbus/v5"
	"net/http"
	"os"
	"smsforwarder/conf"
	"smsforwarder/utils"
	"strings"
	"time"
)

// vars
var (
	c       = make(chan *dbus.Signal, 2)
	conn, _ = dbus.SystemBus()
	rule    = "type='signal',interface='org.freedesktop.ModemManager1.Modem.Messaging',member='Added'"
	phone   []string
)

type Response struct {
	Message string `json:"phone"`
}

func init() {
	numberService := conn.Object("org.freedesktop.ModemManager1", dbus.ObjectPath("/org/freedesktop/ModemManager1/Modem/0"))
	numberService.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Modem", "OwnNumbers").Store(&phone)
}

func ListenMessage() {
	call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)
	if call.Err != nil {
		fmt.Println("Failed to add match: %v", call.Err)
		os.Exit(1)
	}
	conn.Signal(c)
	defer conn.Close()

	var tmp string
	for v := range c {
		if v.Body[1] == true {
			var text string
			var sender string
			service := conn.Object("org.freedesktop.ModemManager1", dbus.ObjectPath(fmt.Sprintf("%s", v.Body[0])))
			service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "Text").Store(&text)
			if len(text) == 0 {
				for i := 0; i <= 20; i++ {
					time.Sleep(1 * time.Second)
					service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "Text").Store(&text)
					service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "number").Store(&sender)
					if len(text) != 0 {
						break
					}
				}
			}
			if text != tmp {
				conf.Message <- fmt.Sprintf("%s---%s", phone[0][2:], text)
				saveMessage(phone[0][2:], text, sender)
				tmp = text
			}
		}

	}

}

// SendMessage :801
func SendMessage(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	number := query.Get("number")
	content := query.Get("message")

	messagingObj := conn.Object("org.freedesktop.ModemManager1", dbus.ObjectPath("/org/freedesktop/ModemManager1/Modem/0"))
	smsProps := map[string]dbus.Variant{
		"number": dbus.MakeVariant(number),
		"text":   dbus.MakeVariant(content),
	}

	var smsPath dbus.ObjectPath
	err := messagingObj.Call("org.freedesktop.ModemManager1.Modem.Messaging.Create", 0, smsProps).Store(&smsPath)
	if err != nil {
		fmt.Println("Failed to create SMS message: %v", err)
	}

	smsObj := conn.Object("org.freedesktop.ModemManager1", smsPath)
	err = smsObj.Call("org.freedesktop.ModemManager1.Sms.Send", 0).Err
	if err != nil {
		fmt.Println("Failed to send SMS: %v", err)
	}

}

// GetMessage  801
func GetMessage(writer http.ResponseWriter, request *http.Request) {
	// 打开文件
	file, err := os.Open("message.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	code := utils.GetMessageCode(lastLine)

	if code == "None" {
		writer.Write([]byte(fmt.Sprintf(`{"message": "%s" }`, lastLine)))
	} else {
		writer.Write([]byte(fmt.Sprintf(`{"code":"%s", "message": "%s" }`, code, lastLine)))
	}

}

func saveMessage(phone, text string, sender string) {
	// 写入数据库
	utils.InsertData(phone, text, sender)
	fmt.Println("开始写入本地文件")
	file, err := os.OpenFile("message.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := []byte(text + "\n")
	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}

}

func GetInfo(writer http.ResponseWriter, request *http.Request) {
	response := Response{
		Message: strings.Join(phone, ""),
	}
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

}
