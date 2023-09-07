package controller
import "time"
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	Title 				string `json:"title"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}
type DBVideo struct {
    Id            int64  `gorm:"column:video_id"`
    AuthorId      int64  `gorm:"column:user_id"`
    Title         string `gorm:"column:title"`
    PlayUrl       string `gorm:"column:play_url"`
    CoverUrl      string `gorm:"column:cover_url"`
    FavoriteCount int64  `gorm:"column:favorite_count"`
    CommentCount  int64  `gorm:"column:comment_count"`
}
type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
type DBComment struct {
	CommentId int64 `gorm:"column:comment_id"`
	UserId int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
	Content string `gorm:"column:content"`
	Date time.Time `gorm:"column:created_time"`
}
type User struct {
    Id            int64  `json:"id,omitempty"`
    Name          string `json:"name,omitempty"`
    FollowCount   int64  `json:"follow_count,omitempty"`
    FollowerCount int64  `json:"follower_count,omitempty"`
    IsFollow      bool   `json:"is_follow,omitempty"`
    Avatar        string `json:"avatar,omitempty"`
    BackGroundImage string `json:"background_image,omitempty"`
    Signature     string `json:"signature,omitempty"`
    TotalFavorited int64  `json:"total_favorited,omitempty"`
    WorkCount     int64  `json:"work_count,omitempty"`
    FavoriteCount int64  `json:"favorite_count,omitempty"`
}
type DBUser struct {
    Id              int64  `gorm:"column:user_id"`
    Name            string `gorm:"column:username"`
    Pwd             string `gorm:"column:password"`
    FollowCount     int64  `gorm:"column:follow_count"`
    FollowerCount   int64  `gorm:"column:follower_count"`
    Avatar          string `gorm:"column:avatar"`
    BackGroundImage string `gorm:"column:background_image"`
    Signature       string `gorm:"column:signature"`
    TotalFavorited  int64  `gorm:"column:total_favorited"`
    WorkCount       int64  `gorm:"column:work_count"`
    FavoriteCount   int64  `gorm:"column:favorite_count"`
}
type Likes struct {
  UserId int64 `gorm:"column:user_id"`
  VideoId int64 `gorm:"column:video_id"`
}
// type Message struct {
// 	Id int64  `json:"id,omitempty" gorm:"column:message_id"`
// 	UserId int64 `json:"from_user_id,omitempty" gorm:"column:from_user_id"`
// 	ToUserId int64 `json:"to_user_id,omitempty" gorm:"column:to_user_id"`
// 	Content string `json:"content,omitempty" gorm:"column:content"`
// 	CreateTime time.Time `json:"create_time,omitempty" gorm:"column:created_time"`
// }
type Message struct {
	Id         int64     `json:"id,omitempty"`
	UserId     int64     `json:"from_user_id,omitempty"`
	ToUserId   int64     `json:"to_user_id,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreateTime int64 	 `json:"create_time,omitempty"`
}

type DBMessage struct {
	Id   int64     `gorm:"column:message_id"`
	UserId  int64     `gorm:"column:from_user_id"`
	ToUserId    int64     `gorm:"column:to_user_id"`
	Content     string    `gorm:"column:content"`
	CreatedTime time.Time `gorm:"column:created_time"`
}
// type Message struct {
// 	Id         int64  `json:"id,omitempty"`
// 	Content    string `json:"content,omitempty"`
// 	CreateTime string `json:"create_time,omitempty"`
// }

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
type Follow struct {
	FollowId int64  `gorm:"column:to_user_id"`
	FollowerId int64  `gorm:"column:from_user_id"`
	Mutual bool  `gorm:"column:is_mutual"`
}

