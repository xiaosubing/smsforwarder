# jd_uz801

## uz801配置详情
| 项目 | 参数 |
|-----|------|
|CPU|MSM8916(高通410)|
|RAM|385MiB|
|ROM|3.82GiB|
|购入价|4.4 (我是大冤种，我哭死)|



## 使用简介

如果是刷的Archlinux，可以一键执行

```shell
pacman -S smsforwarder-beta-3-aarch64.pkg.tar.xz
```

写的比较垃圾，仅供学习。不要喷我行不行呀 好哥哥

## 刷入后优化
### 1、关闭所有led灯 

如果是Archlinux 安装后会自动关闭

```shell
echo 0 > /sys/class/leds/green:internet/brightness
echo 0 > /sys/class/leds/blue:wifi/brightness
echo 0 > /sys/class/leds/mmc0::/brightness
```



## 版本说明

已实现功能：

* 接收短信

* 转发短信

  * qq
  * wx
  * mail
  * 自定义的GET/POST 请求

* 发送短信

  * 端口 801
  * 路径 /api/sendMessage
  * GET参数
    * number
    * message

* 获取验证码

  * 端口 801
  * 路径： /api/getMessage

* 保存短信,暂时存在message.txt中

* 清理短信(规划中...)

  ...

![img](README.assets/1.png)

![img](README.assets/2.png)
