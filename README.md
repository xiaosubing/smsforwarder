# jd_uz801

## uz801配置详情
| 项目 | 参数 |
|-----|------|
|CPU|MSM8916(高通410)|
|RAM|385MiB|
|ROM|3.82GiB|


## 1.0版本说明
已实现功能：
  * 接收短信

  * 转发短信(不过目前仅支持get请求的url，举个栗子： cq-http)

  * 获取验证码(规划中...)

  * 清理短信(规划中...)

  * 保存短信(规划中...)

    ...

## 使用简介

* 创建sms文件夹，当然不创建也行

  ```shell
  mkdir sms 
  ```

* 复制文件到uz801刚刚创建的sms文件夹里面，这个自行发挥神通，能复制上去就行

  ```shell
  scp  smsforwarder  notify  root@棒子Ip地址:/root/sms/
  ```


* 运行notify

  ```shell
  cd  /root/sms/
  ./notify
  
  # 如果提示无权限,执行下方操作
  chmod 777 *
  ```

  运行后会成成配置文件

  ```shell
  [root@SIM-9898 smserver]# ./notify_beta 
  未找到配置，已生成默认配置文件， 请编辑url后重新运行！！！
  ```

* 根据自己的转发渠道自行修改`conf.yml`文件

* 修改完成后在运行两个文件

* 最后为了方便实现开机自启：

  输入`nano /usr/lib/systemd/system/smsforwarder.service`复制粘贴下方示例：

  ```shell
  [Unit]
  Description=SMS Forwarder Service
  After=network.target
  
  [Service]
  ExecStart=/root/sms/smsforwarder
  WorkingDirectory=/root/sms
  User=root
  Group=root
  
  [Install]
  WantedBy=multi-user.target
  ```

  然后创建notfiy的开机自启

  ```shell
  nano /usr/lib/systemd/system/smsnotify.service
  
  # 粘贴
  [Unit]
  Description=SMS Forwarder Service
  After=network.target
  
  [Service]
  ExecStart=/root/sms/notify
  WorkingDirectory=/root/sms
  User=root
  Group=root
  
  [Install]
  WantedBy=multi-user.target
  ```

* 刷新单元文件并设置开机自启

  ```shell
  systemctl    daemon-reload 
  systemctl   enable   --now  smsforwarder 
  systemctl   enable   --now  smsnotify
  ```

  



