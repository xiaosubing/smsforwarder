package message

import (
	"bufio"
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
)

func ListenMessage() {
	rule := "type='signal',interface='org.freedesktop.ModemManager1.Modem.Messaging',member='Added'"
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
			service := conn.Object("org.freedesktop.ModemManager1", dbus.ObjectPath(fmt.Sprintf("%s", v.Body[0])))
			service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "Text").Store(&text)
			if len(text) == 0 {
				for i := 0; i <= 20; i++ {
					fmt.Println("未获取到短信，重新获取......")
					time.Sleep(1 * time.Second)
					service.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Sms", "Text").Store(&text)
					if len(text) != 0 {
						break
					}
				}
			}
			if text != tmp {
				var number []string
				numberService := conn.Object("org.freedesktop.ModemManager1", dbus.ObjectPath("/org/freedesktop/ModemManager1/Modem/0"))
				numberService.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.ModemManager1.Modem", "OwnNumbers").Store(&number)
				conf.Message <- fmt.Sprintf("%s\n%s", number[0][2:], text)
				saveMessage(text)
				tmp = text
			}

		}

	}

}

// sendMessage  2801
func SendMessage(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	number := query.Get("number")
	content := query.Get("message")
	fmt.Println("获取到收信人号码: ", number)
	fmt.Println("内容： ", content)
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

// sendMessage  2801
func GetMessage(writer http.ResponseWriter, request *http.Request) {
	// 打开文件
	file, err := os.Open("message.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 创建一个扫描器用于读取文件
	scanner := bufio.NewScanner(file)

	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text() // 更新最后一行的内容
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

func saveMessage(content string) {
	// 打开文件，如果不存在则创建，设置为追加模式和写入模式
	file, err := os.OpenFile("message.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 要写入的数据
	data := []byte(content + "\n")

	// 写入数据
	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}

}
