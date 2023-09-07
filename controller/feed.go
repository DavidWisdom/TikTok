package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		}})
		return
	}	
	tx := db.Begin()
	token := c.PostForm("token")
	var id int64
	id = -1
	if len(token) != 0 {
		username, password, err := GetInfo(token)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, FeedResponse { Response: Response{
					StatusCode: 1, 
					StatusMsg: "用户鉴权出错",
			}})
			return
		}
		var user DBUser
		querySql := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user)
		if querySql.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, FeedResponse { Response: Response{
					StatusCode: 1,
					StatusMsg: "服务器异常错误",
			}})
			return
		}
		if querySql.RowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusOK, FeedResponse { Response: Response{
					StatusCode: 1,
					StatusMsg: "用户鉴权出错",
			}})
			return
		}
		user_id := c.Query("user_id")
		id, err = strconv.ParseInt(user_id, 10, 64)	
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, FeedResponse { Response: Response{
					StatusCode: 1, 
					StatusMsg: "非法用户标识符",
			}})
			return
		}
	}
	var tempVideos []DBVideo
	err = tx.Table("Video").Order("created_time DESC").Limit(30).Find(&tempVideos).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, FeedResponse { 
			Response: Response{
			StatusCode: 1,
			StatusMsg:  "查询视频信息错误",
		}})
		return
	}
	var videos []Video
	for _, tempVideo := range tempVideos {
    AuthorId := tempVideo.AuthorId
		var user DBUser
	  tx.Table("User").Where("user_id = ?", AuthorId).First(&user)
		queryIsFollow := tx.Table("Follow").Where("follower_id = ? AND followee_id = ?", user.Id, id)
		if queryIsFollow.Error != nil {
			c.JSON(http.StatusOK, FeedResponse { 
				Response: Response{
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
		result := tx.Table("Likes").Where("user_id = ? AND video_id = ?", AuthorId, tempVideo.Id).Select("1").Limit(1)
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
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
