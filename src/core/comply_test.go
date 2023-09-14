package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_client_getSorScoreByContest(t *testing.T) {
	db, _ := gorm.Open(mysql.New(mysql.Config{DSN: "root:my123456@tcp(127.0.0.1:3306)/mycube2?charset=utf8&parseTime=True&loc=Local"}), &gorm.Config{
		Logger: logger.Discard,
	})

	c := &client{
		debug: false,
		db:    db,
		cache: cache.New(time.Minute*15, time.Minute*15),
	}
	gotSingle, gotAvg := c.getSorScoreByContest(21)
	fmt.Println(gotSingle, gotAvg)
}
