package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"os"
	"time"
	"strings"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	username, password, err := GetInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1, 
				StatusMsg: "用户鉴权出错",
		})
		return
	}
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
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
	data, err := c.FormFile("data")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg: "上传文件读取错误",
		})
		return
	}
	currentTime := time.Now()
	timestamp := currentTime.Unix()
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%d_%s", user.Id, timestamp, filename)
	dir, err := os.Getwd()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg: "上传文件失败",
		})
		return
	}
	saveFile := filepath.Join(dir, "/public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg: "上传文件失败",
		})
		return
	}
	 ext := filepath.Ext(saveFile)
	 saveImage := strings.TrimSuffix(saveFile, ext)
	 _, err = GetSnapshot(saveFile, saveImage, 1)
	 if err != nil {
		 	tx.Rollback()
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: "无法获取视频图像",
			})
			return
	 }
	url := os.Getenv("paas_url")
	imageUrl := fmt.Sprintf("%d_%d_%s", user.Id, timestamp, strings.TrimSuffix(filename, ext))
  title := c.PostForm("title")
  coverImage := "https://" + url + "/static/" + imageUrl + ".png"
	saveFile = "https://" + url + "/static/" + imageUrl + ".mp4"
	sql := `
 		INSERT INTO Video (user_id, play_url, cover_url, title)
		VALUES (?, ?, ?, ?)
	`
	if err := tx.Exec(sql, user.Id, saveFile, coverImage, title).Error; err != nil {
	    tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg: "服务器异常错误",
			})
			return
	}
	sql = `
		UPDATE User SET work_count = work_count + 1 WHERE user_id = ?
	`
	if err := tx.Exec(sql, user.Id).Error; err != nil {
	    tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg: "服务器异常错误",
			})
			return
	}
	if err := tx.Commit().Error; err != nil {
			tx.Rollback() // 回滚事务
			c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg: "服务器异常错误",
			})
			return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg: "视频已成功上传",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	db, err := Connect()
	if err != nil {
		c.JSON(http.StatusOK,VideoListResponse { Response: Response{
				StatusCode: 1,
				StatusMsg: "数据库连接异常",
		}})
		return
	}	
	tx := db.Begin()
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
	var tempVideos []DBVideo
	err = tx.Table("Video").Where("user_id = ?", id).Order("created_time DESC").Find(&tempVideos).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, VideoListResponse { Response: Response{
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
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
