package controller
import (
  "os"
  "fmt"
  _ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
)
var db *gorm.DB
func connect() *gorm.DB {
    server := os.Getenv("MYSQL_HOST")
    port := os.Getenv("MYSQL_PORT")
    user := os.Getenv("MYSQL_USER")
    password := os.Getenv("MYSQL_PASSWORD")
    exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok", user, password, server, port)
    var err error
    db, err = gorm.Open("mysql", exec) 
    if err != nil {
        panic("Connect database failed: " + err.Error())
    }
    defer db.Close()
    db.DB().SetMaxOpenConns(100)
    db.DB().SetMaxIdleConns(20)
    return db
}