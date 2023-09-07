package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	//"gorm.io/gorm/schema"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm/logger"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}


// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) { //已完成
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
	to_user_id := c.Query("to_user_id")
	ID, err := strconv.ParseInt(to_user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	var new_user User
	result = db.Table("User").Where("id = ?", ID).First(&new_user) 
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
	if user.Id == ID {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "不能关注自己"},
		})
		return
	}
	actionType := c.Query("action_type")
	if actionType == "1" {
		// 开启事务
		tx := db.Begin()
		updateQuery := `SELECT * FROM following_table WHERE attention = ? AND fans = ?`
		result = tx.Exec(updateQuery, to_user_id, user.Id)
		if result.Error != nil {
				tx.Rollback() // 回滚事务
				log.Fatal(result.Error)
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error"},
				})
				return
		}		
		state := false 
		if result.RowsAffected > 0 {
				state = true
				updateQuery := `UPDATE following_table SET mutual = TRUE WHERE attention = ? AND fans = ?`
				result = tx.Exec(updateQuery, to_user_id, user.Id)
				if result.Error != nil {
					tx.Rollback() // 回滚事务
					log.Fatal(result.Error)
					c.JSON(http.StatusOK, UserLoginResponse{
						Response: Response{StatusCode: 1, StatusMsg: "Error"},
					})
					return
				}		
		}
		insertQuery := `
			INSERT INTO following_table (attention, fans, mutual)
			VALUES (?, ?, ?)
		`
		// 执行插入操作
		result := tx.Exec(insertQuery, to_user_id, user.Id, state)
		if result.Error != nil {
				tx.Rollback() // 回滚事务
				log.Fatal(result.Error)
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error"},
				})
			return
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
				log.Fatal(err)
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error"},
				})
				return
		}
	} else {
		tx := db.Begin()
		 deleteQuery := `
			DELETE FROM following_table
			WHERE attention = ? AND fans = ?
		`
		result = db.Exec(deleteQuery, to_user_id, user.Id)

		if result.Error != nil {
				tx.Rollback() // 回滚事务
				log.Fatal(result.Error)
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error"},
				})
			return
		}
		updateQuery := `UPDATE following_table SET mutual = FALSE WHERE attention = ? AND fans = ?`
		result = tx.Exec(updateQuery, to_user_id, user.Id)
		if result.Error != nil {
				tx.Rollback() // 回滚事务
				log.Fatal(result.Error)
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error"},
				})
				return
		}		
		// 提交事务
		if err := tx.Commit().Error; err != nil {
				log.Fatal(err)
				c.JSON(http.StatusOK, UserLoginResponse{
					Response: Response{StatusCode: 1, StatusMsg: "Error"},
				})
				return
		}	
	}
	c.JSON(http.StatusOK, Response{
			StatusCode: 0,
	})
}

// / FollowList all users have same follow list
// 我的关注
func FollowList(c *gin.Context) {

	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	usr := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", usr, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       exec,
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
	// 查询当前用户对应的关注
	// 假设用户id为10(数据库模拟)
	// 通过当前用户id去查询关系表内的fans，只要当前用户为fans，就能拿到它对应的关注
	current_id := c.Query("user_id")
	ID, err := strconv.ParseInt(current_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	if ID != user.Id {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	var follows []Follow

	
	if err := db.Table("following_table").Where("fans = ?", current_id).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库查询出错",
			},
		})
	}
	var FansNumber []int64

	for _, follow := range follows {
		FansNumber = append(FansNumber, follow.Attention)
	}
	var users []User
	for _, user_id := range FansNumber {
			var user User
	 		db.Table("User").Where("id = ?", user_id).First(&user)
			query := db.Table("following_table").Where("attention = ? AND fans = ?", user_id, user.Id)
			subscribe := false
			if query.RowsAffected == 0 {
					
			} else {
				subscribe = true
			}
			user.IsFollow = subscribe
			users = append(users, user)
	}
	// 查询得到数据
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "查询成功",
		},
		UserList: users,
	})
}

// FollowerList all users have same follower list
// 我的粉丝
func FollowerList(c *gin.Context) {
	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	usr := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", usr, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       exec,
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
	// 查询当前用户对应的关注
	// 假设用户id为10(数据库模拟)
	// 通过当前用户id去查询关系表内的fans，只要当前用户为fans，就能拿到它对应的关注
	current_id := c.Query("user_id")
	ID, err := strconv.ParseInt(current_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	if ID != user.Id {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	var follows []Follow

	
	if err := db.Table("following_table").Where("attention = ?", current_id).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库查询出错",
			},
		})
	}
	var FansNumber []int64

	for _, follow := range follows {
		FansNumber = append(FansNumber, follow.Attention)
	}
	var users []User
	for _, user_id := range FansNumber {
			var user User
	 		db.Table("User").Where("id = ?", user_id).First(&user)
			query := db.Table("following_table").Where("attention = ? AND fans = ?", user_id, user.Id)
			subscribe := false
			if query.RowsAffected == 0 {
					
			} else {
				subscribe = true
			}
			user.IsFollow = subscribe
			users = append(users, user)
	}
	// 查询得到数据
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "查询成功",
		},
		UserList: users,
	})
}

// FriendList all users have same friend list
// 互相关注
// 思路：查询该用户外键对应的粉丝或关注是否存在关联字段
func FriendList(c *gin.Context) {
  /*
	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", user, password, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       exec,
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

	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		// 登录成功
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
		})
	} else {
		// 登录失败
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	// 查询当前用户对应互相关注
	var current_id int64
	current_id = 10
	var follows []Follow

	if err := db.Where("(fans = ? OR attention = ?) AND mutual = ?", current_id, current_id, 1).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库查询出错",
			},
		})
	}
	var FansNumber []int64

	for _, follow := range follows {
		if follow.Attention == current_id {
			FansNumber = append(FansNumber, follow.Fans)
		} else {
			FansNumber = append(FansNumber, follow.Attention)
		}
	}

	var users []User
	// 模拟，查找所有
	if err := db.Find(&users, FansNumber).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库查询出错",
			},
		})
	}

	// 查询得到数据
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "查询成功",
		},
		UserList: users,
	})
 */
}
