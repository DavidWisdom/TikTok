package controller

import (
  /*
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
  */
	"github.com/gin-gonic/gin"
)

//var tempChat = map[string][]Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
  /*
	token := c.Query("token")
	toUserId, err := strconv.Atoi(c.Query("to_user_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "id no"})
	}
	content := c.Query("content")
	Mclaims, err := ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	Muserid := Mclaims.UserId
	MDB, err := Sql_server.DB()
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "发送失败"})
	}
	var um string
	e := MDB.QueryRow("select username from user where user_id=?", Muserid).Scan(um)
	if e != sql.ErrNoRows { //存在  chat_table
		atomic.AddInt64(&messageIdSequence, 1)
		Mnowtime := time.Now() //创建时间
		sqlstr := `insert into chat_table(chat_id,sender_id,receiver_id,message,timestamp) value(?,?,?,?,?)`
		Sql_server.Exec(sqlstr, messageIdSequence, Muserid, toUserId, content, Mnowtime)

	} else { //不存在
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	// token := c.Query("token")
	// toUserId := c.Query("to_user_id")
	// content := c.Query("content")
	// if user, exist := usersLoginInfo[token]; exist {
	// 	userIdB, _ := strconv.Atoi(toUserId)
	// 	chatKey := genChatKey(user.Id, int64(userIdB))

	// 	atomic.AddInt64(&messageIdSequence, 1)
	// 	curMessage := Message{
	// 		Id:         messageIdSequence,
	// 		Content:    content,
	// 		CreateTime: time.Now().Format(time.Kitchen),
	// 	}

	// 	if messages, exist := tempChat[chatKey]; exist {
	// 		tempChat[chatKey] = append(messages, curMessage)
	// 	} else {
	// 		tempChat[chatKey] = []Message{curMessage}
	// 	}
	// 	c.JSON(http.StatusOK, Response{StatusCode: 0})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }
 */
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
  /*
	token := c.Query("token")
	toUserid1 := c.Query("to_user_id")
	toUserId, err := strconv.Atoi(c.Query("to_user_id"))
	if err != nil {
		fmt.Println(toUserid1)
		fmt.Println("111111:", err, "    ", toUserId)
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	Mtoken, e := ParseToken(token)
	if e != nil {
		fmt.Println("22222")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	Muserid := Mtoken.UserId
	MDB, err := Sql_server.DB()
	if err != nil {
		fmt.Println("33333")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "连接失败"})
		return
	}

	var um string
	Me := MDB.QueryRow(`select username from user where user_id=?`, Muserid).Scan(um)
	if Me != sql.ErrNoRows { //存在    id content  time
		rows, err := MDB.Query("select chat_id,message,timestamp from chat_table where sender_id=? and receiver_id=?", Muserid, toUserId)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		} else {
			var messagelist []Message
			for rows.Next() {
				var id int
				var content string
				var time string
				rows.Scan(&id, &content, &time)
				messagelist = append(messagelist, Message{
					Id:         int64(id),
					Content:    content,
					CreateTime: time,
				})
			}
			c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: messagelist})
			fmt.Println(messagelist)
		}
	} else {
		fmt.Println("44444")
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	// token := c.Query("token")
	// toUserId := c.Query("to_user_id")

	// if user, exist := usersLoginInfo[token]; exist {
	// 	userIdB, _ := strconv.Atoi(toUserId)
	// 	chatKey := genChatKey(user.Id, int64(userIdB))

	// 	c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }
 */
}

// func genChatKey(userIdA int64, userIdB int64) string {
// 	if userIdA > userIdB {
// 		return fmt.Sprintf("%d_%d", userIdB, userIdA)
// 	}
// 	return fmt.Sprintf("%d_%d", userIdA, userIdB)
// }
