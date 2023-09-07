package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type TempVideo struct {
	Id            int64  `json:"id,omitempty" gorm:"column:video_id"`
	AuthorId      int64  `json:"author" gorm:"column:id"`
  Title         string `json:"title,omitempty" gorm:"column:title"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty" gorm:"column:play_url"`
	CoverUrl      string `json:"cover_url,omitempty" gorm:"column:cover_image"`
	FavoriteCount int64  `json:"favorite_count,omitempty" gorm:"column:likes_count"`
	CommentCount  int64  `json:"comment_count,omitempty" gorm:"column:comments_count"`
}
type Likes struct {
  UserId int64 `gorm:"column:user_id"`
  VideoId int64 `gorm:"column:video_id"`
}
type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
  Title         string `json:"title,omitempty"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}
type TempComment struct {
	CommentId int64 `gorm:"column:comment_id"`
	UserId int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
	Content string `gorm:"column:content"`
	Date string `gorm:"column:create_date"`
}
type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	// omitempty：如果没有赋值，就忽略
	Id int64 `json:"id,omitempty"`
	// Id   int64  `json:"id,omitempty" gorm:"column:user_id"`
	Name string `json:"name,omitempty" gorm:"column:username"` // 用户名唯一
	// Name          string `json:"name,omitempty"`
	Nickname      string `json:"nickname,omitempty"` // 新加
	PassWord      string `json:"password,omitempty" gorm:"column:password"`
	FollowCount   int64  `json:"follow_count,omitempty" gorm:"column:follow_count"`
	FollowerCount int64  `json:"follower_count,omitempty" gorm:"column:fans_count"`
	IsFollow      bool   `json:"is_follow,omitempty" `
}


type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

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
	Attention int64 `json:"attention,omitempty"`
	Fans      int64 `json:"fans,omitempty"`
	Mutual    bool  `json:"is_follow,omitempty"`
}
