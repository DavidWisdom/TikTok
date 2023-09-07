package database
import (
  "os"
  "fmt"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)
var DB *gorm.DB
func connect() {
    server := os.Getenv("MYSQL_HOST")
    port := os.Getenv("MYSQL_PORT")
    user := os.Getenv("MYSQL_USER")
    password := os.Getenv("MYSQL_PASSWORD")
    dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok", user, password, server, port)
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Connect database failed: " + err.Error())
    }
    fmt.Println("Connect database succeed")
}
