package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
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

	// 解析请求中的上传文件
	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  "上传文件错误",
			},
		})
		return
	}
	// 将上传的文件保存到指定位置
	uploadedFilePath := fmt.Sprintf("path/to/save/%s", file.Filename)
	err = c.SaveUploadedFile(file, uploadedFilePath)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 3,
				StatusMsg:  "保存文件错误",
			},
		})
		return
	}

	// type Video struct {
	// 	Id            int64  `json:"id,omitempty"`
	// 	Author        User   `json:"author"`
	//  Title         string `json:"title"`
	// 	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	// 	CoverUrl      string `json:"cover_url,omitempty"`
	// 	FavoriteCount int64  `json:"favorite_count,omitempty"`
	// 	CommentCount  int64  `json:"comment_count,omitempty"`
	// 	IsFavorite    bool   `json:"is_favorite,omitempty"`
	// }

	// 在数据库中保存上传的视频信息
	video := Video{
		PlayUrl: uploadedFilePath,
		// 其他字段赋值
	}

	err = db.Create(&video).Error
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 4,
				StatusMsg:  "保存视频信息错误",
			},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}
