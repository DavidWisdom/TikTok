package controller

import (
	"github.com/gin-gonic/gin"
  "os"
  "fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
  //"fmt"
  "time"
  //"strconv"
	"log"
  "github.com/dgrijalva/jwt-go"
  "gorm.io/gorm/logger"
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
		Logger: logger.Default.LogMode(logger.Info), // Info level log
	})
    token := c.Query("token")
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
  video_id := c.Query("video_id")
  action_type := c.Query("action_type")
  if action_type == "1" {
    comment_text := c.Query("comment_text")
    timeObj := time.Now()
    month := timeObj.Month()
    day := timeObj.Day()
    date := fmt.Sprintf("%02d-%02d", month, day)
	  tx := db.Begin()
	// 执行插入评论的 SQL 语句
	if err := tx.Exec("INSERT INTO Comment (user_id, video_id, content) VALUES (?, ?, ?)",
	    user.Id, video_id, comment_text).Error; err != nil {
	    tx.Rollback() // 回滚事务
	    log.Fatal(err)
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
		  return
	}
	// 获取插入后的 comment_id
	var commentID int64
	if err := tx.Raw("SELECT LAST_INSERT_ID()").Scan(&commentID).Error; err != nil {
	    tx.Rollback() // 回滚事务
	    log.Fatal(err)
		  c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
		  return
	}	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
	    log.Fatal(err)
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1},})
			return
	}
    c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
      Comment: Comment{
        Id: commentID, // TODO: 评论ID
        User: user,
        Content: comment_text,
        CreateDate: date,
      }})
    return
  }
  comment_id := c.Query("comment_id")
	deleteQuery := `
	 		DELETE FROM Comment
		  WHERE comment_id = ?
	 `
	result = db.Exec(deleteQuery, comment_id)
  c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},})
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
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
		Logger: logger.Default.LogMode(logger.Info), // Info level log
	})
	video_id := c.Query("video_id")
	var tempComments []TempComment
	err = db.Table("Comment").Where("video_id = ?", video_id).Order("create_date DESC").Find(&tempComments).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "查询评论信息错误",
		})
		return
	}
  var comments []Comment
  for _, tempComment := range tempComments {
    // Author
    AuthorId := tempComment.UserId
    var author User
	  db.Table("User").Where("id = ?", AuthorId).First(&author)
		newComment := Comment{
			Id: tempComment.UserId,
			User: author,
			Content: tempComment.Content,
			CreateDate: tempComment.Date,
		}
		comments = append(comments, newComment)
  }
	response := CommentListResponse{
			Response: Response{StatusCode: 0, StatusMsg: "Success"},
			CommentList: comments,
	}
	// Return the response
	c.JSON(http.StatusOK, response)
}
