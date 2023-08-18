package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// / FollowList all users have same follow list
// 我的关注
func FollowList(c *gin.Context) {
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

	// 查询当前用户对应的关注
	// 假设用户id为10(数据库模拟)
	// 通过当前用户id去查询关系表内的fans，只要当前用户为fans，就能拿到它对应的关注
	current_id := 10
	var follows []Follow

	if err := db.Where("fans = ?", current_id).Find(&follows).Error; err != nil {
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
}

// FollowerList all users have same follower list
// 我的粉丝
func FollowerList(c *gin.Context) {
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

	// 查询当前用户对应的关注
	// 查询当前用户对应的关注
	// 假设用户id为10(数据库模拟)
	// 通过当前用户id去查询关系表内的fans，只要当前用户为fans，就能拿到它对应的关注
	current_id := 10
	var follows []Follow

	if err := db.Where("attention = ?", current_id).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "数据库查询出错",
			},
		})
	}
	var FansNumber []int64

	for _, follow := range follows {
		FansNumber = append(FansNumber, follow.Fans)
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
}

// FriendList all users have same friend list
// 互相关注
// 思路：查询该用户外键对应的粉丝或关注是否存在关联字段
func FriendList(c *gin.Context) {
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
}
