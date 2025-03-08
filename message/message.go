package message

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"net/http"
	"os"
	"smsforwarder/conf"
	"smsforwarder/utils"
	"time"
)

// vars
var (
	c       = make(chan *dbus.Signal, 2)
	conn, _ = dbus.SystemBus()
	rule    = "type='signal',interface='org.freedesktop.ModemManager1.Modem.Messaging',member='Added'"
)

type Response struct {
	Message string `json:"phone"`
}

func init() {
	numberService := conn.Object("org.freedesktop.ModemManager1", dbus.ObjectPath("/org/freedesktop/ModemManager1/Modem/0"))
	numberService.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Modem", "OwnNumbers").Store(&conf.Smsforwarder.PhoneNumber)

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
			service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "Number").Store(&sender)
			if len(text) == 0 {
				for i := 0; i <= 20; i++ {
					time.Sleep(1 * time.Second)
					service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "Text").Store(&text)
					service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "Number").Store(&sender)
					if len(text) != 0 {
						break
					}
				}
			}

			if text != tmp {
				code := utils.GetMessageCode(text)
				conf.Message <- fmt.Sprintf("%s---%s---%s---%s", conf.Smsforwarder.PhoneNumber[0][2:], text, sender, code)
				saveMessage(text, sender, code)
				tmp = text
				fmt.Println("获取到的消息：", text)
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

func saveMessage(text, sender, code string) {
	// 是否加密，等更新...
	utils.InsertData(conf.Smsforwarder.PhoneNumber[0][2:], text, sender, code)

}
