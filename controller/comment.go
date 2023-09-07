package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		}})
		return
	}	
	token := c.Query("token")
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse { Response: Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
		}})
		return
	}
	var user DBUser
	tx := db.Begin()
	querySql := tx.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusOK, CommentActionResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		}})
		return
	}
	if querySql.RowsAffected == 0 {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusOK, CommentActionResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
		}})
		return
	}
	video_id := c.Query("video_id")
  var tempVideo []DBVideo
	result := tx.Table("Video").Where("video_id = ?", video_id).First(&tempVideo)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "服务器异常错误"},
		})
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "视频不存在"},
		})
		return
	}
	actionType := c.Query("action_type")
  if actionType == "1" {
    comment_text := c.Query("comment_text")
    timeObj := time.Now()
    month := timeObj.Month()
    day := timeObj.Day()
    date := fmt.Sprintf("%02d-%02d", month, day)
		// 执行插入评论的 SQL 语句
		if err := tx.Exec("INSERT INTO Comment (user_id, video_id, content) VALUES (?, ?, ?)",
		    user.Id, video_id, comment_text).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
			  return
		}
		// 获取插入后的 comment_id
		var comment_id int64
		if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&comment_id).Error; err != nil {
		    tx.Rollback() // 回滚事务
			  c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
			  return
		}	
		sql := `
			UPDATE Video SET comment_count = comment_count + 1 WHERE video_id = ?
		`
		if err := tx.Exec(sql, video_id).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
			  return
		}
		// 提交事务
		if err := tx.Commit().Error; err != nil {
				tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
				return
		}
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: Comment{
				Id: comment_id,
				User: User{
					Id: user.Id,
					Name: user.Name,
					FollowCount: user.FollowCount,
					FollowerCount: user.FollowerCount,
					IsFollow: false,
					Avatar: user.Avatar,
					BackGroundImage: user.BackGroundImage,
					Signature: user.Signature,
					TotalFavorited: user.TotalFavorited,
					WorkCount: user.WorkCount,
					FavoriteCount: user.FavoriteCount,
				},
				Content: comment_text,
				CreateDate: date,
		}})
    return
  }
	comment_id := c.Query("comment_id")
	sql := `
	 		DELETE FROM Comment
		  WHERE comment_id = ?
	 `
	var comment []DBComment
	result = tx.Table("Comment").Where("comment_id = ?", comment_id).First(&comment)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "服务器异常错误"},
		})
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 1, StatusMsg: "评论不存在"},
		})
		return
	}
	// 执行插入评论的 SQL 语句
	if err := tx.Exec(sql, comment_id).Error; err != nil {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
			return
	}
	sql = `
		UPDATE Video SET comment_count = comment_count - 1 WHERE video_id = ?
	`
	if err := tx.Exec(sql, video_id).Error; err != nil {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
			return
	}
	if err := tx.Commit().Error; err != nil {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
			return
	}
	c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},})
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		}})
		return
	}	
	video_id := c.Query("video_id")
	tx := db.Begin()
	var tempComments []DBComment
	err = tx.Table("Comment").Where("video_id = ?", video_id).Order("created_time DESC").Find(&tempComments).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusInternalServerError, CommentListResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "查询评论信息错误",
		}})
		return
	}
	var id int64
	id = -1
	token := c.Query("token")
	if len(token) != 0 {
		username, password, err := GetInfo(token)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse { Response: Response{
					StatusCode: 1, 
					StatusMsg: "用户鉴权出错",
			}})
			return
		}
		var user DBUser
		querySql := tx.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
		if querySql.Error != nil {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, FeedResponse { Response: Response{
					StatusCode: 1,
					StatusMsg: "服务器异常错误",
			}})
			return
		}
		if querySql.RowsAffected == 0 {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, FeedResponse { Response: Response{
					StatusCode: 1,
					StatusMsg: "用户鉴权出错",
			}})
			return
		}
		id = user.Id
	}
  var comments []Comment
  for _, tempComment := range tempComments {
    // Author
    AuthorId := tempComment.UserId
    var user DBUser
	  tx.Table("User").Where("user_id = ?", AuthorId).First(&user)
		queryIsFollow := tx.Table("Follow").Where("follower_id = ? AND followee_id = ?", user.Id, id)
		if queryIsFollow.Error != nil {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, CommentListResponse{
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
		createDate := tempComment.Date.Format("01-02")
		newComment := Comment{
			Id: tempComment.UserId,
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
			Content: tempComment.Content,
			CreateDate: createDate,
		}
		comments = append(comments, newComment)
  }
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
			},
		})
		return
	}
	response := CommentListResponse{
			Response: Response{StatusCode: 0},
			CommentList: comments,
	}
	// Return the response
	c.JSON(http.StatusOK, response)
}
