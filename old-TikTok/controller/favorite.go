package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm/logger"
	
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	usr := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", usr, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: exec,
		// DSN: "root:root12345@tcp(localhost:3306)/db_test?charset=utf8&parseTime=True&loc=Local",
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Info level log
	})
    token := c.Query("token")
    toke, err := ParseJWTToken(token)
    if err != nil {
      // 处理解析错误
      return
    }
    claims, ok := toke.Claims.(jwt.MapClaims)
    if !ok {
      // 无法获取声明信息
      return
    }
    username := claims["username"].(string)
    password := claims["password"].(string)
	  var user User
	  result := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user) //改成表字段一致
	if result.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Error querying user data"},
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
    video_id := c.Query("video_id")
    var tempVideo []TempVideo
	result = db.Table("Video").Where("video_id = ?", video_id).First(&tempVideo)
	if result.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Error querying video data"},
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Video doesn't exist"},
		})
		return
	}
  actionType := c.Query("action_type")
	if actionType == "1" {
		insertQuery := `
			INSERT INTO Likes (user_id, video_id)
			VALUES (?, ?)
		`	
		result = db.Exec(insertQuery, user.Id, video_id)
	} else {
		 deleteQuery := `
	 		DELETE FROM Likes
		  WHERE user_id = ? AND video_id = ?
	 	`
		result = db.Exec(deleteQuery, user.Id, video_id)
	}
	c.JSON(http.StatusOK, Response{
			StatusCode: 0,
	})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	/*
	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	usr := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", usr, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: exec,
		// DSN: "root:root12345@tcp(localhost:3306)/db_test?charset=utf8&parseTime=True&loc=Local",
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Info level log
	})
  
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
 */
}
