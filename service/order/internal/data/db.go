package data

import (
	"mshop/service/order/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

// gormWriter 实现 logger.Writer 接口
type gormWriter struct {
	helper *log.Helper
}

func (w *gormWriter) Printf(format string, args ...interface{}) {
	w.helper.Infof(format, args...)
}

func NewDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch c.Database.Driver {
	case "mysql":
		dialector = mysql.Open(c.Database.Source)
	case "postgres":
		dialector = postgres.Open(c.Database.Source)
	default:
		dialector = mysql.Open(c.Database.Source)
	}

	// 创建 GORM logger 适配器
	helper := log.NewHelper(logger)
	writer := &gormWriter{helper: helper}
	gormLogger := gormlog.New(
		writer,
		gormlog.Config{
			SlowThreshold: 0,
			LogLevel:      gormlog.Info,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	return db, err
}
