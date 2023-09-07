package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
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
	tx := db.Begin()
	var user DBUser
	querySql := tx.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
	if querySql.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		})
		return
	}
	if querySql.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	video_id := c.Query("video_id")
	var tempVideo DBVideo
	result := tx.Table("Video").Where("video_id = ?", video_id).First(&tempVideo)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
		return
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "视频不存在"})
		return
	}
	actionType := c.Query("action_type")
  if actionType == "1" {
		if err := tx.Exec("INSERT INTO Likes (user_id, video_id) VALUES (?, ?)", user.Id, video_id).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1})
			  return
		}
		sql := `
			UPDATE Video SET favorite_count = favorite_count + 1 WHERE video_id = ?
		`
		if err := tx.Exec(sql, video_id).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
			  return
		}
		sql = `
			UPDATE User SET total_favorited = total_favorited + 1 WHERE user_id = ?
		`
		if err := tx.Exec(sql, tempVideo.AuthorId).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
			  return
		}
		sql = `
			UPDATE User Set favorite_count = favorite_count + 1 WHERE user_id = ?
		`
		if err := tx.Exec(sql, user.Id).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
			  return
		}
		// 提交事务
		if err := tx.Commit().Error; err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
				return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return
  } else {
		if err := tx.Exec("DELETE FROM Likes WHERE user_id = ? AND video_id = ?", user.Id, video_id).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1})
			  return
		}
		sql := `
			UPDATE Video SET favorite_count = favorite_count - 1 WHERE video_id = ?
		`
		if err := tx.Exec(sql, video_id).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
			  return
		}
		sql = `
			UPDATE User SET total_favorited = total_favorited - 1 WHERE user_id = ?
		`
		if err := tx.Exec(sql, tempVideo.AuthorId).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
			  return
		}
		sql = `
			UPDATE User Set favorite_count = favorite_count - 1 WHERE user_id = ?
		`
		if err := tx.Exec(sql, user.Id).Error; err != nil {
		    tx.Rollback() // 回滚事务
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
			  return
		}
		// 提交事务
		if err := tx.Commit().Error; err != nil {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "服务器异常错误"})
				return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK,VideoListResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		}})
		return
	}	
	token := c.Query("token")
	var id int64
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
		}})
		return
	}
	tx := db.Begin()
	var me DBUser
	querySql := tx.Table("User").Where("username = ? AND password = ?", username, password).First(&me)
	if querySql.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		}})
		return
	}
	if querySql.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "用户鉴权出错",
		}})
		return
	}
	user_id := c.Query("user_id")
	id, err = strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1, 
				StatusMsg: "非法用户标识符",
		}})
		return
	}
	var likes []Likes
	err = tx.Table("Likes").Where("user_id = ?", id).Find(&likes).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		}})
		return
	}
	var videos []Video
	for _, like := range likes {
		video_id := like.VideoId
		var tempVideo DBVideo
		result := tx.Table("Video").Where("video_id = ?", video_id).First(&tempVideo)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, VideoListResponse { Response: Response{StatusCode: 1, StatusMsg: "服务器异常错误"}})
			return
		}
		if result.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusOK, VideoListResponse { Response: Response{StatusCode: 1, StatusMsg: "视频不存在"}})
			return
		}
		AuthorId := tempVideo.AuthorId
		var user DBUser
	  tx.Table("User").Where("user_id = ?", AuthorId).First(&user)
		queryIsFollow := tx.Table("Follow").Where("from_user_id = ? AND to_user_id = ?", user.Id, id)
		if queryIsFollow.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, VideoListResponse { Response: Response{
					StatusCode: 1,
					StatusMsg: "服务器异常错误",
			}})
			return
		}
		subscribe := false
		if queryIsFollow.RowsAffected > 0 {
			subscribe = true
		}
		isFavorite := true
		result = tx.Table("Likes").Where("user_id = ? AND video_id = ?", AuthorId, tempVideo.Id).Select("1").Limit(1)
    if result.RowsAffected == 0 {
      isFavorite = false
    }
		newVideo := Video{
      Id: tempVideo.Id,
      Author: User{
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
      Title: tempVideo.Title,
      PlayUrl: tempVideo.PlayUrl,
      CoverUrl: tempVideo.CoverUrl,
      FavoriteCount: tempVideo.FavoriteCount,
      CommentCount: tempVideo.CommentCount,
      IsFavorite: isFavorite,
    }
    videos = append(videos, newVideo)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // 回滚事务
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "服务器异常错误",
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
