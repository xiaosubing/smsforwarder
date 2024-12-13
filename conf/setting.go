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
	NotifyType           string
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
			NotifyType:           viper.GetString("notify.type"),
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
# 消息模板, 暂不支持自定义，等后续更新吧
template: "验证码: %s\n收信人: %s\n\n短信原文:\n%s"
# 消息示例：

# 配置通知渠道
notify:
  # 通知渠道，必填！！！！！！！！！！！！！！！！！！！！！！！！！
  # 可以配置的值为： qq 、wx、webhook 和 mail
  # 填写完成后请完善对应渠道的详细信息！！！
  type: qq

  webhook:
    #url: "http://127.0.0.1:3000/send_private_msg?user_id=3286276407&message="
    # wx的url,请在wxbot里面对消息内容进行urlcode解码操作
    #url: "http://192.168.86.78:2802/api/sendMessage"
    # 钉钉url示例
    #url: "https://oapi.dingtalk.com/robot/send?access_token=钉钉软件里面复制token"
    # 请根据自己渠道配置请求类型， 钉钉为post
    type: "post"
    # 钉钉的payload, "1" 表示短信内容，
    #payload: '{"msgtype": "text","text": {"content": "1"}}'
    # 其他post请求的payload， 根据自己的请求配置
    #payload: '{"message": "1", "to": "梅干菜小酥饼"}'


  # 邮箱通知
  mail:
    account: "3286276407@qq.com"
    password: "非邮箱登陆密码，请自行获取邮箱的凭证"
    sendTo: "邮件接收人"
    # 邮件主题，%s表示主动识别的验证码
    subject: "验证码: %s"
    
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
