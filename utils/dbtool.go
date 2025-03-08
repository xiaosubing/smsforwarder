package utils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"smsforwarder/conf"
	"time"
)

// 定义消息模型
type Message struct {
	ID     uint   `gorm:"primaryKey"`
	Phone  string `gorm:"size:128"`
	Sender string `gorm:"size:128"`
	SMS    string `gorm:"type:text"`
	//Code   string `gorm:"size:128"`
}

var db *gorm.DB

func sqliteTool() {
	var err error
	// 1. 连接数据库（启用WAL模式提升并发性能）
	db, err = gorm.Open(sqlite.Open(conf.Smsforwarder.Db.DbName), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second, // 慢查询阈值
				LogLevel:      logger.Info, // 日志级别
				Colorful:      true,        // 彩色日志
			},
		),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
		return
	}

	// 2. 自动迁移表结构
	err = db.AutoMigrate(&Message{})
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}

}

func InsertData(phoneNumber, sms, sender string) {
	// 判断数据库类型
	if conf.Smsforwarder.Db.SaveType == "local" {
		sqliteTool()

	} else {
		// mysql
	}
	// 3. 批量插入数据（使用事务提升性能）
	message := []Message{
		{Phone: phoneNumber, Sender: sender, SMS: sms},
	}

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("开启事务失败: %v", tx.Error)
	}

	// 分批次插入（避免单次事务过大）
	batchSize := 1000
	for i := 0; i < len(message); i += batchSize {
		end := i + batchSize
		if end > len(message) {
			end = len(message)
		}
		batch := message[i:end]
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			log.Fatalf("插入数据失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("提交事务失败: %v", err)
	}

	// 4. 验证数据写入
	var result []Message
	db.Find(&result)
	for _, m := range result {
		log.Printf("用户信息: ID=%d, Phone=%s, Sender=%s, SMS=%s", m.ID, m.Phone, m.Sender, m.SMS)
	}
}
