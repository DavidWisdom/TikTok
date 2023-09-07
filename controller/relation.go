package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		})
		return
	}	
	token := c.Query("token")
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	var user DBUser
	querySql := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		})
		return
	}
	if querySql.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	to_user := c.Query("to_user_id")
	to_user_id, err := strconv.ParseInt(to_user, 10, 64)
	var newUser DBUser
	querySql = db.Table("User").Where("user_id", to_user_id).First(&newUser)
	if querySql.Error != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		})
		return
	}
	if querySql.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "用户不存在",
		})
		return
	}
	if user.Id == to_user_id {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "不能关注自己"},
		})
		return
	}
	actionType := c.Query("action_type")
	if actionType == "1" {
		  tx := db.Begin()
			var follow Follow
			result := tx.Table("Follow").Where("from_user_id = ? AND to_user_id = ?", to_user_id, user.Id).Limit(1).Scan(&follow)
			if result.Error != nil {
					tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, UserLoginResponse{
						Response: Response{StatusCode: 1, StatusMsg: result.Error.Error()},
					})
					return
			}		
			state := false 
			if result.RowsAffected > 0 {
					state = true
					sql := `UPDATE Follow SET is_mutual = TRUE WHERE from_user_id = ? AND to_user_id = ?`
					result = tx.Exec(sql, to_user_id, user.Id)
					if result.Error != nil {
						tx.Rollback() // 回滚事务
						c.JSON(http.StatusOK, UserLoginResponse{
							Response: Response{StatusCode: 1, StatusMsg: "服务器异常错误"},
						})
						return
					}		
			}
			sql := `
				INSERT INTO Follow (from_user_id, to_user_id, is_mutual)
				VALUES (?, ?, ?)
			`
			// 执行插入操作
			result = tx.Exec(sql, user.Id, to_user_id, state)
			if result.Error != nil {
					tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, UserLoginResponse{
						Response: Response{StatusCode: 1, StatusMsg: "服务器异常错误"},
					})
				return
			}
			sql = `
				UPDATE User SET follow_count = follow_count + 1 WHERE user_id = ?
			`
			if err := tx.Exec(sql, user.Id).Error; err != nil {
			    tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
				  return
			}
			sql = `
				UPDATE User Set follower_count = follower_count + 1 WHERE user_id = ?
			`
			if err := tx.Exec(sql, to_user_id).Error; err != nil {
			    tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
				  return
			}
			if err := tx.Commit().Error; err != nil {
					c.JSON(http.StatusOK, UserLoginResponse{
						Response: Response{StatusCode: 1, StatusMsg: "Error"},
					})
					return
			}
	} else {
		  tx := db.Begin()
			sql := `DELETE FROM Follow WHERE from_user_id = ? AND to_user_id = ?`
			result := tx.Exec(sql, user.Id, to_user_id)
			if result.Error != nil {
					tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, UserLoginResponse{
						Response: Response{StatusCode: 1, StatusMsg: "服务器异常错误"},
					})
					return
			}		
			sql = `
				UPDATE Follow SET is_mutual = FALSE WHERE from_user_id = ? AND to_user_id = ?
			`
			// 执行插入操作
			result = tx.Exec(sql, to_user_id, user.Id)
			if result.Error != nil {
					tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, UserLoginResponse{
						Response: Response{StatusCode: 1, StatusMsg: "服务器异常错误"},
					})
				return
			}
			sql = `
				UPDATE User SET follow_count = follow_count - 1 WHERE user_id = ?
			`
			if err := tx.Exec(sql, user.Id).Error; err != nil {
			    tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
				  return
			}
			sql = `
				UPDATE User Set follower_count = follower_count - 1 WHERE user_id = ?
			`
			if err := tx.Exec(sql, to_user_id).Error; err != nil {
			    tx.Rollback() // 回滚事务
					c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
				  return
			}
			if err := tx.Commit().Error; err != nil {
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

func FollowerList(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		})
		return
	}	
	token := c.Query("token")
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	var user DBUser
	querySql := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		})
		return
	}
	if querySql.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	user_id := c.Query("user_id")
	var id int64
	id, err = strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1, 
				StatusMsg: "非法用户标识符",
		}})
		return
	}
	var follows []Follow
	if err := db.Table("Follow").Where("to_user_id = ?", id).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	var FansNumber []int64
	for _, follow := range follows {
		FansNumber = append(FansNumber, follow.FollowerId)
	}
	var users []User
	for _, usr_id := range FansNumber {
			var usr DBUser
	 		db.Table("User").Where("user_id = ?", usr_id).First(&usr)
			queryIsFollow := db.Table("Follow").Where("from_user_id = ? AND to_user_id = ?", usr.Id, id).Select("1").Limit(1)
			if queryIsFollow.Error != nil {
				c.JSON(http.StatusOK, UserListResponse { Response: Response{
						StatusCode: 1,
						StatusMsg: "服务器异常错误",
				}})
				return
			}
			subscribe := false
			if queryIsFollow.RowsAffected > 0 {
				subscribe = true
			}
			newUser := User{
				Id: usr.Id,
				Name: usr.Name,
				FollowCount: usr.FollowCount,
				FollowerCount: usr.FollowerCount,
				IsFollow: subscribe,
				Avatar: usr.Avatar,
				BackGroundImage: usr.BackGroundImage,
				Signature: usr.Signature,
				TotalFavorited: usr.TotalFavorited,
				WorkCount: usr.WorkCount,
				FavoriteCount: usr.FavoriteCount,
			}
			users = append(users, newUser)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

func FollowList(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		})
		return
	}	
	token := c.Query("token")
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	var user DBUser
	querySql := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		})
		return
	}
	if querySql.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	user_id := c.Query("user_id")
	var id int64
	id, err = strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1, 
				StatusMsg: "非法用户标识符",
		}})
		return
	}
	var follows []Follow
	if err := db.Table("Follow").Where("from_user_id = ?", id).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	var FansNumber []int64
	for _, follow := range follows {
		FansNumber = append(FansNumber, follow.FollowId)
	}
	var users []User
	for _, usr_id := range FansNumber {
			var usr DBUser
	 		db.Table("User").Where("user_id = ?", usr_id).First(&usr)
			queryIsFollow := db.Table("Follow").Where("from_user_id = ? AND to_user_id = ?", usr.Id, id).Select("1").Limit(1)
			if queryIsFollow.Error != nil {
				c.JSON(http.StatusOK, UserListResponse { Response: Response{
						StatusCode: 1,
						StatusMsg: "服务器异常错误",
				}})
				return
			}
			subscribe := false
			if queryIsFollow.RowsAffected > 0 {
				subscribe = true
			}
			newUser := User{
				Id: usr.Id,
				Name: usr.Name,
				FollowCount: usr.FollowCount,
				FollowerCount: usr.FollowerCount,
				IsFollow: subscribe,
				Avatar: usr.Avatar,
				BackGroundImage: usr.BackGroundImage,
				Signature: usr.Signature,
				TotalFavorited: usr.TotalFavorited,
				WorkCount: usr.WorkCount,
				FavoriteCount: usr.FavoriteCount,
			}
			users = append(users, newUser)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		})
		return
	}	
	token := c.Query("token")
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	var user DBUser
	querySql := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		})
		return
	}
	if querySql.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	user_id := c.Query("user_id")
	var id int64
	id, err = strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1, 
				StatusMsg: "非法用户标识符",
		}})
		return
	}


	var follows []Follow
	if err := db.Table("Follow").Where("from_user_id = ? AND is_mutual = TRUE", id).Find(&follows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	var FansNumber []int64
	for _, follow := range follows {
		FansNumber = append(FansNumber, follow.FollowId)
	}
	var users []User
	for _, usr_id := range FansNumber {
			var usr DBUser
	 		db.Table("User").Where("user_id = ?", usr_id).First(&usr)
			queryIsFollow := db.Table("Follow").Where("from_user_id = ? AND to_user_id = ?", usr.Id, id).Select("1").Limit(1)
			if queryIsFollow.Error != nil {
				c.JSON(http.StatusOK, UserListResponse { Response: Response{
						StatusCode: 1,
						StatusMsg: "服务器异常错误",
				}})
				return
			}
			subscribe := false
			if queryIsFollow.RowsAffected > 0 {
				subscribe = true
			}
			newUser := User{
				Id: usr.Id,
				Name: usr.Name,
				FollowCount: usr.FollowCount,
				FollowerCount: usr.FollowerCount,
				IsFollow: subscribe,
				Avatar: usr.Avatar,
				BackGroundImage: usr.BackGroundImage,
				Signature: usr.Signature,
				TotalFavorited: usr.TotalFavorited,
				WorkCount: usr.WorkCount,
				FavoriteCount: usr.FavoriteCount,
			}
			users = append(users, newUser)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: users,
	})
}
