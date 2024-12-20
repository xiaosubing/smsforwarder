package conf

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Smsforwarder *smsforwarderSetting
var Message = make(chan string, 5)

type smsforwarderSetting struct {
	MessageTemplate string
	Notify          *notify
}

type notify struct {
	NotifyType           []string
	NotifyWebHookUrl     string
	NotifyWebHookType    string
	NotifyWebHookPayload string

	NotifyMailAccount  string
	NotifyMailPassword string
	NotifyMailSendTo   string
	NotifyMailSmtpHost string
	NotifyMailSmtpPort int
	NotifyMailSubject  string
	// more notify ...

}

func init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("未找到配置，已生成默认配置文件， 请编辑url后重新运行！！！")
		createConf()
		os.Exit(0)
	}

	Smsforwarder = NewSmsforwarder()
}

func NewSmsforwarder() *smsforwarderSetting {

	return &smsforwarderSetting{
		MessageTemplate: viper.GetString("template"),
		Notify: &notify{
			NotifyType:           viper.GetStringSlice("notify.type"),
			NotifyWebHookUrl:     viper.GetString("notify.webhook.url"),
			NotifyWebHookType:    viper.GetString("notify.webhook.type"),
			NotifyWebHookPayload: viper.GetString("notify.webhook.payload"),
			NotifyMailAccount:    viper.GetString("notify.mail.account"),
			NotifyMailPassword:   viper.GetString("notify.mail.password"),
			NotifyMailSendTo:     viper.GetString("notify.mail.sendTo"),
			NotifyMailSmtpHost:   viper.GetString("notify.mail.smtpHost"),
			NotifyMailSmtpPort:   viper.GetInt("notify.mail.smtpPort"),
			NotifyMailSubject:    viper.GetString("notify.mail.subject"),
		},
	}
}
func createConf() {
	file, err := os.OpenFile("./conf.yml", os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		fmt.Println("创建配置文件失败，请手动创建")
		os.Exit(0)
	}
	defer file.Close()

	s := `
# 消息模板
# 1 需要验证码 2需要收信人手机号最后两位 3 短信原文， 顺序不能乱！
template: "\n验证码: 1\n收信人: 2\n\n短信原文:\n3"

# 配置通知渠道
notify:
  # 通知渠道，必填！！！！！！！！！！！！！！！！！！！！！！！！！
  # 可以配置的值为： qq 、webhook 和 mail
  # 填写完成后请完善对应渠道的详细信息！！！
  type:
    - mail
    # 多渠道支持
    #- qq
    #- webhook

  webhook:
    # 钉钉url示例
    url: "https://oapi.dingtalk.com/robot/send?access_token=钉钉软件里面复制token"
    # 请根据自己渠道配置请求类型， 钉钉为post
    type: "post"
    # 钉钉的payload, "1" 表示短信内容，
    payload: '{"msgtype": "text","text": {"content": "1"}}'
    # 其他post请求的payload， 根据自己的请求配置
    #payload: '{"message": "1", "to": "梅干菜小酥饼"}'


  # 邮箱通知
  mail:
    account: "3286276407@qq.com"
    password: "非邮箱登陆密码，请自行获取邮箱的凭证"
    sendTo: "建议使用运营商邮箱为收件人，打开短信提醒，能够自动识别验证码"
    # 邮件主题，%s表示主动识别的验证码
    # 沃邮箱建议使用 "2"
    subject: "2"
    # 其他邮箱可以把程序内自动识别的验证码加上，但是识别有可能不准哦
    # subject: "验证码：1"

    # qq邮箱的服务器信息
    # 默认使用qq邮箱发送, 可自行替换其他邮箱
    smtpHost: "smtp.qq.com"
    smtpPort: 587
`
	write := bufio.NewWriter(file)
	_, err = write.WriteString(s)
	if err != nil {
		fmt.Println("写入配置失败，请手动创建并配置配置文件")
		os.Exit(0)
	}

	write.Flush()
}
