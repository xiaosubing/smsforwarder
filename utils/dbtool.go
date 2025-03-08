package utils

import (
	"errors"
	"fmt"
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
	ID        uint   `gorm:"primaryKey"`
	Phone     string `gorm:"size:20"`
	Sender    string `gorm:"size:128"`
	SMS       string `gorm:"type:text"`
	Code      string `gorm:"size:10"`
	CreatedAt string `gorm:"size:20"`
}

var db *gorm.DB

func init() {
	var err error
	// 1. 连接数据库（启用WAL模式提升并发性能）
	db, err = gorm.Open(sqlite.Open(conf.Smsforwarder.Db.DbName), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				//SlowThreshold: time.Second, // 慢查询阈值
				//LogLevel:      logger.Info, // 日志级别
				//Colorful:      true,        // 彩色日志
				LogLevel: logger.Silent,
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

func InsertData(phoneNumber, sms, sender, code string) {
	//if conf.Smsforwarder.Db.SaveType == "local" {
	//
	//} else {
	//	// mysql,等更新...
	//}

	message := []Message{
		{Phone: phoneNumber, Sender: sender, SMS: sms, Code: code},
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("开启事务失败: %v", tx.Error)
	}

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

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("提交事务失败: %v", err)
	}

}

func QueryData(phoneNumber string) string {
	var msg Message
	if err := db.Where("phone = ?", phoneNumber).Last(&msg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return `{"result": 404 }`
		} else {
			return `{"result": 500 }`
		}
	}

	tmp := fmt.Sprintf(`{"phone":"%s","sender":"%s","sms":"%s","code":"%s"}`, msg.Phone, msg.Sender, msg.SMS, msg.Code)
	return tmp

}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	m.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (m *Message) BeforeUpdate(tx *gorm.DB) error {
	m.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}
