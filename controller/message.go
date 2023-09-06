package controller

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	// "sync/atomic"
	// "time"
	//"time"
)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
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
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, 
			StatusMsg: "不能发送消息给自己",
		})
		return
	}
	action_type := c.Query("action_type")
	content := c.Query("content")
	if action_type == "1" {
			sql := `
		 		INSERT INTO Message (from_user_id, to_user_id, content)
				VALUES (?, ?, ?)
			`
			result := db.Exec(sql, user.Id, to_user_id, content)
			if result.Error != nil {
				c.JSON(http.StatusOK, Response{
						StatusCode: 1,
						StatusMsg: "服务器异常错误",
				})
				return
			}
			c.JSON(http.StatusOK, Response{
					StatusCode: 0,
			})
	} else {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "操作类型非法",
		})
	}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
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
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, 
			StatusMsg: "不能发送消息给自己",
		})
		return
	}
	var tempMessages []DBMessage
	tx := db.Begin()
	err = tx.Table("Message").Where("from_user_id = ? AND to_user_id = ?", user.Id, to_user_id).Order("created_time").Find(&tempMessages).Error
	if err != nil {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "查询聊天信息失败",
			})
			return
	}
	// sql := "UPDATE Message SET pre_msg_time = TRUE WHERE from_user_id = ? AND to_user_id = ?"
	// if err := tx.Exec(sql, user.Id, to_user_id).Error; err != nil {
	// 		tx.Rollback() // 回滚事务
	// 		c.JSON(http.StatusOK, Response{
	// 			StatusCode: 1,
	// 			StatusMsg: "查询聊天信息失败",
	// 		})
	// 		return
	// }
	// if err := tx.Commit().Error; err != nil {
	// 		c.JSON(http.StatusOK, Response{
	// 			StatusCode: 1,
	// 			StatusMsg: "查询聊天信息失败",
	// 		})
	// 		return
	// }
	var messages []Message
	for _, message := range tempMessages {
		newMessage := Message {
			Id: message.Id,
			UserId: message.UserId,
			ToUserId: message.ToUserId,
			Content: message.Content,
			// CreateTime: message.CreatedTime.Format(time.Kitchen),
		}
		messages = append(messages, newMessage)
	}
	c.JSON(http.StatusOK, ChatResponse{ Response: Response{StatusCode: 0}, MessageList: messages})
}
