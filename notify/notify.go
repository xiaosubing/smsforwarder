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
		text = strings.Split(v, "\n")
		code = utils.GetMessageCode(text[1])
		number := "**" + string(text[0][9]) + string(text[0][10])
		tmp := conf.Smsforwarder.MessageTemplate
		templates := strings.Replace(strings.Replace(strings.Replace(tmp, "1", "%s", -1), "2", "%s", -1), "3", "%s", -1)
		if strings.Contains(tmp, "1") && strings.Contains(tmp, "2") && strings.Contains(tmp, "3") {
			content = fmt.Sprintf(templates, code, number, text[1])
		} else if strings.Contains(tmp, "2") && strings.Contains(tmp, "3") {
			content = fmt.Sprintf(templates, number, text[1])
		} else if strings.Contains(tmp, "1") && strings.Contains(tmp, "3") {
			content = fmt.Sprintf(templates, code, text[1])
		} else {
			content = text[1]
		}

		fmt.Printf("短信原文\n: %s\n", text[1])

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
	fmt.Println("读取到的subject:", tmp)
	subject = strings.Replace(strings.Replace(tmp, "1", "%s", -1), "2", "%s", -1)
	fmt.Println("替换的字符串主题subject:", subject)
	if strings.Contains(tmp, "1") && strings.Contains(tmp, "2") {
		subject = fmt.Sprintf(subject, "", content)
		fmt.Println("12:", subject)
	} else if strings.Contains(tmp, "1") {
		subject = fmt.Sprintf(subject, code)
		fmt.Println("1:", subject)
	} else if strings.Contains(tmp, "2") {
		subject = fmt.Sprintf(subject, content)
		fmt.Println("2:", subject)
	} else {
		subject = "短信转发"
	}
	fmt.Println("主题：", subject)
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
		payload := strings.Replace(conf.Smsforwarder.Notify.NotifyWebHookPayload, "1", content, -1)
		utils.HttpPost(conf.Smsforwarder.Notify.NotifyWebHookUrl, payload)
	}

}

func sendQQMessage() {
	url := fmt.Sprintf("%s%s", conf.Smsforwarder.Notify.NotifyWebHookUrl, url2.QueryEscape(content))
	utils.HttpGet(url)
}
