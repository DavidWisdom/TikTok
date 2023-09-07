package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
  "os"
  "fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

var db *gorm.DB

func Feed(c *gin.Context) {
	var err error
  server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	usr := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", usr, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: exec,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
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
    return
	}
	var tempVideos []TempVideo
	//返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个 -- 待测试

  // TODO: 加入查询user对象
	err = db.Table("Video").Order("created_at DESC").Limit(30).Find(&tempVideos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "查询视频信息错误",
		})
		return
	}
  var videos []Video


    token := c.Query("token")
		var me_id int64
		me_id = -1
		if len(token) == 0 {
		} else {
			toke, err := ParseJWTToken(token)
			if err != nil {
				c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid token"})
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
				c.JSON(http.StatusOK, Response{
					StatusCode: 1, 
					StatusMsg: "Error querying user data",
				})
				return
			}
			if result.RowsAffected == 0 {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1, 
					StatusMsg: "User doesn't exist",
				})
				return
			}
			me_id = user.Id
		}
  for _, tempVideo := range tempVideos {
    // Author
    AuthorId := tempVideo.AuthorId
    var author User
	  db.Table("User").Where("id = ?", AuthorId).First(&author)
		query := db.Table("following_table").Where("attention = ? AND fans = ?", AuthorId, me_id)
		subscribe := false
		if query.RowsAffected == 0 {
				
		} else {
			subscribe = true
		}
		author.IsFollow = subscribe
    // IsFavorite
    var like Likes
    isFavorite := true
    result := db.Table("Likes").Where("user_id = ? AND video_id = ?", AuthorId, tempVideo.Id).First(&like)
    if result.RowsAffected == 0 {
      isFavorite = false
    }
    newVideo := Video{
      Id: tempVideo.Id,
      Author: author,
      Title: tempVideo.Title,
      PlayUrl: tempVideo.PlayUrl,
      CoverUrl: tempVideo.CoverUrl,
      FavoriteCount: tempVideo.FavoriteCount,
      CommentCount: tempVideo.CommentCount,
      IsFavorite: isFavorite,
    }

    videos = append(videos, newVideo)
  }
  
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
