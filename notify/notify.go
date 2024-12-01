package notify

import (
	"fmt"
	"gopkg.in/gomail.v2"
	url2 "net/url"
	"smsforwarder/conf"
	"smsforwarder/utils"
	"strings"
)

func Notify() {
	var content string
	tempMessage := conf.Smsforwarder.MessageTemplate
	for v := range conf.Message {
		text := strings.Split(v, "\n")
		fmt.Println("开始转发短信: ", text[1])
		code := utils.GetMessageCode(text[1])
		number := "**" + string(text[0][9]) + string(text[0][10])
		content = fmt.Sprintf(tempMessage, code, number, text[1])

		if code == "None" {
			content = strings.Replace(content, "验证码: None\n", "", -1)
			fmt.Println("处理过后的内容: ", content)
		}

		// forwarder
		if conf.Smsforwarder.Forwarder.Enable == true {

			return
		}

		// send
		if strings.ToUpper(conf.Smsforwarder.Notify.NotifyType) == "QQ" {
			sendQQMessage(content)
		}

		//if strings.ToUpper(conf.Smsforwarder.Notify.NotifyType) == "WX" {
		//	sendWxMessage(content)
		//}

		if strings.ToUpper(conf.Smsforwarder.Notify.NotifyType) == "WEBHOOK" {
			sendWebhookMessage(content)
		}

		if strings.ToUpper(conf.Smsforwarder.Notify.NotifyType) == "MAIL" {
			var subject string
			if len(code) == 0 {
				subject = "短信转发"
			} else {
				subject = code
			}

			sendMailMessage(subject, content)
		}
	}

}

//func sendWxMessage(content string) {
//	content = strings.Replace(content, "\n", "\\n", -1)
//	payload := strings.Replace(conf.Smsforwarder.Notify.NotifyWebHookPayload, "1", url2.QueryEscape(content), -1)
//	fmt.Println(payload)
//	utils.HttpPost(conf.Smsforwarder.Notify.NotifyWebHookUrl, payload)
//}

func sendMailMessage(subject, content string) {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Smsforwarder.Notify.NotifyMailAccount)
	m.SetHeader("To", conf.Smsforwarder.Notify.NotifyMailSendTo)

	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	d := gomail.NewDialer(
		conf.Smsforwarder.Notify.NotifyMailSmtpHost,
		conf.Smsforwarder.Notify.NotifyMailSmtpPort,
		conf.Smsforwarder.Notify.NotifyMailAccount,
		conf.Smsforwarder.Notify.NotifyMailPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}

func sendWebhookMessage(content string) {
	if strings.ToUpper(conf.Smsforwarder.Notify.NotifyWebHookType) == "GET" {
		url := fmt.Sprintf("%s%s", conf.Smsforwarder.Notify.NotifyWebHookUrl, content)
		utils.HttpGet(url)
	} else {
		content = strings.Replace(content, "\n", "\\n", -1)
		payload := strings.Replace(conf.Smsforwarder.Notify.NotifyWebHookPayload, "1", content, -1)
		fmt.Println(payload)
		utils.HttpPost(conf.Smsforwarder.Notify.NotifyWebHookUrl, payload)
	}

}

func sendQQMessage(content string) {
	url := fmt.Sprintf("%s%s", conf.Smsforwarder.Notify.NotifyWebHookUrl, url2.QueryEscape(content))
	utils.HttpGet(url)
}
