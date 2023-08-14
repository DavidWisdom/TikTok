package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
  "fmt"
  "time"
  "strconv"
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
  token := c.Query("token")
  _, err := strconv.Atoi(c.Query("video_id"))
  if err != nil {
    c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Invalid number"})
    return
  }
  action_type := c.Query("action_type")
	if user, exist := usersLoginInfo[token]; exist {
		if action_type == "1" {
      comment_text := c.Query("comment_text")
      timeObj := time.Now()
      month := timeObj.Month()
      day := timeObj.Day()
      date := fmt.Sprintf("%02d-%02d", month, day)
      c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
        Comment: Comment{
          Id:         1, // TODO: 评论ID
          User:       user,
          Content:    comment_text,
          CreateDate: date,
        }})
      return
		}
  comment_id := c.Query("comment_id")
  fmt.Println("comment_id: %v", comment_id)
	c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
