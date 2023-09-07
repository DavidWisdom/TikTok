package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	// "sync/atomic"
	"gorm.io/gorm"
	"strconv"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}
type UserRegisterResponse UserLoginResponse
type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
			},
		})
		return
	}
	username := c.Query("username")
	password := c.Query("password")
	var user DBUser
	querySql := db.Table("User").Where("username = ?", username).First(&user)
	if querySql.Error != nil {
		if querySql.Error == gorm.ErrRecordNotFound {
			token, err := GetToken(username, password)
			if err != nil {
				c.JSON(http.StatusOK, UserRegisterResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg: "服务器异常错误",
					},
				})
				return
			}
			newUser := DBUser{Name: username, Pwd: password}
			insertSql := db.Table("User").Create(&newUser)
			if insertSql.Error != nil {
				c.JSON(http.StatusOK, UserRegisterResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg: "服务器异常错误",
					},
				})
				return
			}
			c.JSON(http.StatusOK, UserRegisterResponse{
				Response: Response{
					StatusCode: 0,
				},
				UserId: newUser.Id,
				Token: token,
			})
	  } 
	} else {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "该用户名已被注册",
			},
		})
  }
}

func Login(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
			},
		})
		return
	}
	username := c.Query("username")
	password := c.Query("password")
	token, err := GetToken(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	var user DBUser
	querySql := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	if querySql.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "用户名或密码错误",
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserId: user.Id,
		Token: token,
	})
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
			},
		})
		return
	}
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
			},
		})
		return
	}
	var user DBUser
	querySql := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	if querySql.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
			},
		})
		return
	}
	user_id := c.Query("user_id")
	id, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1, 
				StatusMsg: "非法用户标识符",
			},
		})
		return
	}
	queryIsFollow := db.Table("Follow").Where("from_user_id = ? AND to_user_id = ?", user.Id, id)
	if queryIsFollow.Error != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	subscribe := false
	if queryIsFollow.RowsAffected > 0 {
		subscribe = true
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User: User{
			Id: user.Id,
			Name: user.Name,
			FollowCount: user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow: subscribe,
			Avatar: user.Avatar,
			BackGroundImage: user.BackGroundImage,
			Signature: user.Signature,
			TotalFavorited: user.TotalFavorited,
			WorkCount: user.WorkCount,
			FavoriteCount: user.FavoriteCount,
		},
	})
}
