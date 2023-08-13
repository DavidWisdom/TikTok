package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"sync/atomic"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// // connect database

// func init() {
//   dsn := "root:gFFIZVKe@tcp(172.16.32.27:13306)/tiktok?charset=utf8&parseTime=True&loc=Local"
//   db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
//   if err != nil {
//     log.Fatal(err)
//   }
//   db.AutoMigrate(&User{})

// }

func Register(c *gin.Context) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:zjh97867860@tcp(127.0.0.1:13306)/tiktok?charset=utf8&parseTime=True&loc=Local",
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Info level log
	})
	if err != nil {
		// 引发异常
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库连接异常",
			},
		})
	}

	//TODO：保存数据
	// check strong password
	// jwt token
	// check whether exist
	// register successful! Store in database
	username := c.Query("username")
	password := c.Query("password")
	// Perform strong password check
	if !isStrongPassword(password) {
		c.JSON(http.StatusBadRequest, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Please set a strong password"},
		})
		return
	}

	token, err := generateJWTToken(username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Error generating token"},
		})
		return
	}

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    token,
		})
	}
	db.AutoMigrate(&User{})
	//TODO: 如何设置变量
	var newUser User
	result := db.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Error saving user data"},
		})
		return
	}

}

// check strong password
func isStrongPassword(password string) bool {
	upperExist := 0
	lowerExist := 0
	digitExist := 0
	specialExist := 0
	if len(password) < 8 {
		return false
	}

	for i := 0; i < len(password); i++ {
		singleChar := password[i]
		if singleChar >= 'A' && singleChar <= 'Z' {
			upperExist = 1
		} else if singleChar >= 'a' && singleChar <= 'z' {
			lowerExist = 1
		} else if singleChar >= '0' && singleChar <= '9' {
			digitExist = 1
		} else if singleChar == '!' || singleChar == '@' || singleChar == '#' || singleChar == '?' {
			specialExist = 1
		}
	}
	return upperExist == 1 && lowerExist == 1 && digitExist == 1 && specialExist == 1
}

// 生成JWT token
func generateJWTToken(username, password string) (string, error) {
	// Define the claims
	claims := jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	// create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with a secret key
	signedToken, err := token.SignedString([]byte("my_secret_key"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func Login(c *gin.Context) {
	// 1. init
	// 2.
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:zjh97867860@tcp(127.0.0.1:13306)/tiktok?charset=utf8&parseTime=True&loc=Local",
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Info level log
	})
	if err != nil {
		// 引发异常
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库连接异常",
			},
		})
	}

	// 获取用户名，密码， token
	username := c.Query("username")
	password := c.Query("password")
	token := c.Query("token")
	// // TODO: DONE
	// token, err := generateJWTToken(username, password)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, UserLoginResponse{
	// 		Response: Response{StatusCode: 1, StatusMsg: "Error generating token"},
	// 	})
	// 	return
	// }
	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
	//TODO: 查询数据
	// type User struct {
	// 	// omitempty：如果没有赋值，就忽略
	// 	Id            int64  `json:"id,omitempty" gorm:"column:user_id"`
	// 	Name          string `json:"name,omitempty" gorm:"column:username"`
	// 	PassWord      string `json:"password,omitempty" gorm:"column:password"`
	// 	FollowCount   int64  `json:"follow_count,omitempty" gorm:"column:follow_count"`
	// 	FollowerCount int64  `json:"follower_count,omitempty" gorm:"column:fans_count"`
	// 	IsFollow      bool   `json:"is_follow,omitempty" `
	// }

	var user User
	result := db.Where("username = ? AND password = ?", username, password).First(&user) //改成表字段一致
	if result.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Error querying user data"},
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Invalid username or password"},
		})
		return
	}
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})
}

func Logout(c *gin.Context) {
	// get token
	token := c.Query("token")
	if _, exist := usersLoginInfo[token]; exist {
		delete(usersLoginInfo, token)
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "Logout successful",
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Please check whether login or correct token!"},
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
