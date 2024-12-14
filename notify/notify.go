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
		code := utils.GetMessageCode(text[1])
		number := "**" + string(text[0][9]) + string(text[0][10])
		content = fmt.Sprintf(tempMessage, code, number, text[1])

		if code == "None" {
			content = strings.Replace(content, "验证码: None\n", "", -1)
			// 有验证码的发到邮箱，非验证码发送到其他渠道(也可以是邮箱)
			if len(conf.Smsforwarder.Notify.NotifySecType) != 0 {
				forwarderMoreType("None", strings.ToUpper(conf.Smsforwarder.Notify.NotifySecType), content)
			}
			return
		}

		forwarderMoreType(code, strings.ToUpper(conf.Smsforwarder.Notify.NotifyType), content)
	}

}

func forwarderMoreType(code, forwarderType, content string) {
	// send
	if forwarderType == "QQ" {
		sendQQMessage(content)
	}

	if forwarderType == "WEBHOOK" {
		sendWebhookMessage(content)
	}

	if forwarderType == "MAIL" {
		var subject string
		if code == "None" {
			subject = "短信转发"
		} else {
			subject = fmt.Sprintf(conf.Smsforwarder.Notify.NotifyMailSubject, code)
		}
		sendMailMessage(subject, content)
	}
}

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
		content = strings.Replace(strings.Replace(content, "\n", "\\n", -1), "\r", "", -1)
		payload := strings.Replace(conf.Smsforwarder.Notify.NotifyWebHookPayload, "1", content, -1)
		fmt.Println(payload)
		utils.HttpPost(conf.Smsforwarder.Notify.NotifyWebHookUrl, payload)
	}

}

func sendQQMessage(content string) {
	url := fmt.Sprintf("%s%s", conf.Smsforwarder.Notify.NotifyWebHookUrl, url2.QueryEscape(content))
	utils.HttpGet(url)
}
