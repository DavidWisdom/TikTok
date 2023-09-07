package controller

import (
	// "fmt"
	"net/http"
	"os"
	"fmt"
	// "sync/atomic"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
  // "log"
  // "database/sql"
)
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func Register(c *gin.Context) {
	// 连接数据库
	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", user, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: exec,
		// DSN: "root:root12345@tcp(localhost:3306)/db_test?charset=utf8&parseTime=True&loc=Local",
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
	username := c.Query("username")
	password := c.Query("password")
	var existingUser User
	// result := db.Where("username = ?", username).First(&existingUser)
	tableName := "User" // 替换为您想要查询的表名
	result := db.Table(tableName).Where("username = ?", username).First(&existingUser)
	// fmt.Printf("existingUser: %+v\n", existingUser)
	// fmt.Printf("result: %+v\n", result)
	if result.Error != nil {
    // gorm.ErrRecordNotFound指的是当结果为空时，会出现ErrRecordNotFound红条
    // TODO: 不影响该程序！！！记录！！！
		if result.Error == gorm.ErrRecordNotFound {
			// 用户不存在，继续处理注册逻辑
			token, err := generateJWTToken(username, password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error generating token"},
				})
				return
			}
			// atomic.AddInt64(&userIdSequence, 1)
      createUser := User{Name: username, PassWord: password}
			// usersLoginInfo[token] = newUser
			rst := db.Table(tableName).Create(&createUser)
			// fmt.Printf("Create result: %+v\n", rst)
			if rst.Error != nil {
				c.JSON(http.StatusInternalServerError, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error saving user data"},
				})
				return
			}
			// 返回成功状态
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "用户注册成功",
				},
				UserId: createUser.Id,
				Token:  token,
			})
		} else {
			c.JSON(http.StatusInternalServerError, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Error checking existing user"},
			})
			return
		}
	} else {
		// 用户已存在，返回错误
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exists"},
		})
		return
	}

}
// 生成JWT token
func generateJWTToken(username, password string) (string, error) {
	// Define the claims
	claims := jwt.MapClaims{
		"username": username,
		"password": password,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
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

func ParseJWTToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the custom claims type
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used for signing
		return []byte("my_secret_key"), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func Login(c *gin.Context) {
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
	if err != nil {
		// 引发异常
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库连接异常",
			},
		})
	}
	// 获取用户名，密码， 生成token
	username := c.Query("username")
	password := c.Query("password")
	token, err := generateJWTToken(username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Error generating token"},
		})
		return
	}
	var user User
  tableName := "User" // 替换为您想要查询的表名
	// result := db.Table(tableName).Where("username = ?", username).First(&existingUser)
	result := db.Table(tableName).Where("username = ? AND password = ?", username, password).First(&user) //改成表字段一致
	// result := db.Where("username = ? AND password = ?", username, password).First(&user) //改成表字段一致
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
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "用户登录成功",
		},
		UserId: user.Id,
		Token:  token,
	})
}

// func Logout(c *gin.Context) {
// 	// get token
// 	token := c.Query("token")
// 	if _, exist := usersLoginInfo[token]; exist {
// 		delete(usersLoginInfo, token)
// 		c.JSON(http.StatusOK, Response{
// 			StatusCode: 0,
// 			StatusMsg:  "Logout successful",
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, UserLoginResponse{
// 			Response: Response{StatusCode: 1, StatusMsg: "Please check whether login or correct token!"},
// 		})
// 	}

// }

func UserInfo(c *gin.Context) {
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
		Logger: logger.Default.LogMode(logger.Info),
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
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Invalid user ID",
		})
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid user ID",
		})
		return
	}
	var user User
	result := db.Table("User").Select("id, nickname, username, follow_count, fans_count, is_follow").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message: "Error retrieving user information",
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User: user,
	})
}
