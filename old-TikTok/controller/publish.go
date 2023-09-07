package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"github.com/dgrijalva/jwt-go"
  "os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
   ffmpeg "github.com/u2takey/ffmpeg-go"
	 "bytes"
	 "strings"
	 "github.com/disintegration/imaging"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	 buf := bytes.NewBuffer(nil)
	 err = ffmpeg.Input(videoPath).
			 Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
			 Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
			 WithOutput(buf, os.Stdout).
			 Run()
	 if err != nil {
			 log.Fatal("生成缩略图失败：", err)
			 return "", err
	 }
	 img, err := imaging.Decode(buf)
	 if err != nil {
			 log.Fatal("生成缩略图失败：", err)
			 return "", err
	 }
	 err = imaging.Save(img, snapshotPath+".png")
	 if err != nil {
			 log.Fatal("生成缩略图失败：", err)
			 return "", err
	 }
	 fmt.Println("--snapshotPath--", snapshotPath)
	 // --snapshotPath-- ./assets/testImage
	 names := strings.Split(snapshotPath, "/")
	 fmt.Println("----names----", names)
	 // ----names---- [./assets/testImage]
	 // 这里把 snapshotPath 的 string 类型转换成 []string
	 snapshotName = names[len(names)-1] + ".png"
	 fmt.Println("----snapshotName----", snapshotName)
	 // ----snapshotName---- ./assets/testImage.png
	 return snapshotName, nil
}

func Publish(c *gin.Context) {
	token := c.PostForm("token")
  // 通过token获得username和password
  toke, err := ParseJWTToken(token)
  if err != nil {
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

  // 连接数据库
	server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	usr := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", usr, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: exec,
		// DSN: "root:root12345@tcp(localhost:3306)/db_test?charset=utf8&parseTime=True&loc=Local",
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

// 获取其他声明信息
	var user User
	result := db.Table("User").Where("username = ? AND password = ?", username, password).First(&user) //改成表字段一致
	if result.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Error querying user data"},
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	// c.JSON(http.StatusOK, UserLoginResponse{
	// 	Response: Response{
	// 		StatusCode: 0,
	// 		StatusMsg:  "用户登录成功",
	// 	},
	// 	UserId: user.Id,
	// 	Token:  token,
	// })

	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	currentTime := time.Now()
	timestamp := currentTime.Unix()
  // fmt.Println(data)
	filename := filepath.Base(data.Filename)
	// fmt.Println("Hello,world")
	// user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%d_%s", user.Id, timestamp, filename)
	dir, err := os.Getwd()
	if err != nil {
		// fmt.Println("获取当前路径失败:", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg: err.Error(),
		})
		return
	}
	saveFile := filepath.Join(dir, "/public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		// fmt.Println(err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg: err.Error(),
		})
		return
	} 


	// videoPath := "视频文件路径"
	// outputImagePath := "输出图片路径"

	// cmd := exec.Command("ffmpeg", "-i", videoPath, "-vframes", "1", "-f", "image2", outputImagePath)
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("成功提取视频的第一帧图像。")
	 ext := filepath.Ext(saveFile)
	 saveImage := strings.TrimSuffix(saveFile, ext)
	 fmt.Println(saveFile)
	 fmt.Println(saveImage)
	 _, err = GetSnapshot(saveFile, saveImage, 1)
	 if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg: err.Error(),
			})
			 return
	 }
	// ext = filepath.Ext(filename)
	// imageUrl := strings.TrimSuffix(filename, ext)
	url := os.Getenv("paas_url")
	imageUrl := fmt.Sprintf("%d_%d_%s", user.Id, timestamp, strings.TrimSuffix(filename, ext))
  title := c.PostForm("title")
  coverImage := "https://" + url + "/static/" + imageUrl + ".png"
	saveFile = "https://" + url + "/static/" + imageUrl + ".mp4"
	insertQuery := `
		INSERT INTO Video (id, play_url, cover_image, title)
		VALUES (?, ?, ?, ?)
	`	
  // fmt.Println("user_id: ", user.Id)
	resultt := db.Exec(insertQuery, user.Id, saveFile, coverImage, title)
  if resultt.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "Error inserting video data",
		})
		return
	}
	// fmt.Println("Hello,world")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
  server := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	usr := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PASSWORD")
	exec := fmt.Sprintf("%v:%v@tcp(%v:%v)/tiktok?charset=utf8&parseTime=True&loc=Local", usr, pwd, server, port)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: exec,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	// fmt.Println("Hello,world")
	if err != nil {
		c.JSON(http.StatusInternalServerError, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Error connecting to database",
			},
		})
		return
	}
  userid := c.Query("user_id")
	var tempVideos []TempVideo
	err = db.Table("Video").Where("id = ?", userid).Order("created_at DESC").Find(&tempVideos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "查询视频信息错误",
		})
		return
	}
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
	
  var videos []Video

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
	})
}


