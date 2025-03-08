package notify

import (
	"fmt"
	"gopkg.in/gomail.v2"
	url2 "net/url"
	"smsforwarder/conf"
	"smsforwarder/utils"
	"strings"
)

// vars
var content string
var code string
var text []string
var sender string

func Notify() {
	for v := range conf.Message {
		text = strings.Split(v, "---")
		number := "**" + string(text[0][9]) + string(text[0][10])
		code = text[3]
		sender = text[2]

		tmp := conf.Smsforwarder.MessageTemplate
		if strings.Contains(tmp, "[验证码]") {
			content = strings.Replace(tmp, "[验证码]", code, -1)
		}
		if strings.Contains(content, "[收信人]") {
			content = strings.Replace(content, "[收信人]", number, -1)
		}
		if strings.Contains(tmp, "[发信人]") {
			content = strings.Replace(content, "[发信人]", sender, -1)
		}
		if strings.Contains(tmp, "[短信原文]") {
			content = strings.Replace(content, "[短信原文]", text[1], -1)
		}

		// send
		for _, v1 := range conf.Smsforwarder.Notify.NotifyType {
			moreType := strings.ToUpper(v1)
			if moreType == "QQ" {
				sendQQMessage()
			}

			if moreType == "WEBHOOK" {
				sendWebhookMessage()
			}

			if moreType == "MAIL" {
				sendMailMessage()
			}
		}
	}

}

func sendMailMessage() {
	// subject
	var subject string
	tmp := conf.Smsforwarder.Notify.NotifyMailSubject

	// 根据模板替换

	if strings.Contains(tmp, "[验证码]") {
		subject = strings.Replace(tmp, "[验证码]", code, -1)
	}

	if strings.Contains(tmp, "[发信人]") {
		subject = strings.Replace(subject, "[发信人]", sender, -1)
	}
	if strings.Contains(tmp, "[短信原文]") {
		subject = strings.Replace(subject, "[短信原文]", text[1], -1)
	}

	if subject == "" {
		subject = "短信转发"
	}

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

func sendWebhookMessage() {
	if strings.ToUpper(conf.Smsforwarder.Notify.NotifyWebHookType) == "GET" {
		url := fmt.Sprintf("%s%s", conf.Smsforwarder.Notify.NotifyWebHookUrl, content)
		utils.HttpGet(url)
	} else {
		content = strings.Replace(strings.Replace(content, "\n", "\\n", -1), "\r", "", -1)
		payload := strings.Replace(strings.Replace(conf.Smsforwarder.Notify.NotifyWebHookPayload, "[短信原文]", content, -1), "[验证码]", code, -1)
		utils.HttpPost(conf.Smsforwarder.Notify.NotifyWebHookUrl, payload)
	}

}

func sendQQMessage() {
	url := fmt.Sprintf("%s%s", conf.Smsforwarder.Notify.NotifyWebHookUrl, url2.QueryEscape(content))
	utils.HttpGet(url)
}
