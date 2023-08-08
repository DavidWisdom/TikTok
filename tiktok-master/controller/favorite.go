package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
  // "log"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
  // db, err := connect_mysql()
  // if err != nil {
  //   log.Fatal(err)
  //   c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "Internal Server Error"})
  //   return
  // }
  // defer db.Close()
  // token := c.Query("token")
  // // TODO: 查询token是否合法
  // query := "SELECT user_id FROM User WHERE token = ?"
  // rows, err := db.Query(query, token)
  // log.Println(query)
  // if err != nil {
  //   log.Fatal(err)
  //   c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "Internal Server Error"})
  //   return
  // }
  // log.Println(query, token)
  // defer rows.Close()
  // if rows.Next() {
	 c.JSON(http.StatusOK, Response{StatusCode: 0}) 
  // } else {
	 // c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
  // }
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
