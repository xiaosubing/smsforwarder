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

func Notify() {
	for v := range conf.Message {
		text = strings.Split(v, "-")
		code = utils.GetMessageCode(text[1])
		number := "**" + string(text[0][9]) + string(text[0][10])
		tmp := conf.Smsforwarder.MessageTemplate
		// 根据模板替换
		templates := strings.Replace(strings.Replace(strings.Replace(tmp, "[验证码]", "%s", -1), "[收信人]", "%s", -1), "[短信原文]", "%s", -1)

		if strings.Contains(tmp, "[验证码]") && strings.Contains(tmp, "[收信人]") && strings.Contains(tmp, "[短信原文]") {
			content = fmt.Sprintf(templates, code, number, text[1])
		} else if strings.Contains(tmp, "[收信人]") && strings.Contains(tmp, "[短信原文]") {
			content = fmt.Sprintf(templates, number, text[1])
		} else if strings.Contains(tmp, "[验证码]") && strings.Contains(tmp, "[短信原文]") {
			content = fmt.Sprintf(templates, code, text[1])
		} else {
			content = text[1]
		}

		fmt.Printf("短信原文: %s\n", text[1])

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
	subject = strings.Replace(strings.Replace(tmp, "1", "%s", -1), "2", "%s", -1)
	if strings.Contains(tmp, "1") && strings.Contains(tmp, "2") {
		subject = fmt.Sprintf(subject, "", content)
	} else if strings.Contains(tmp, "1") {
		subject = fmt.Sprintf(subject, code)
	} else if strings.Contains(tmp, "2") {
		subject = fmt.Sprintf(subject, content)
	} else {
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
