package model
import "github.com/RaymondCode/simple-demo/database"
type Comment struct {
    CommentID   int    `gorm:"primaryKey"`
    UserID      int    `gorm:"not null;index;onDelete:CASCADE"`
    VideoID     int    `gorm:"not null;index;onDelete:CASCADE"`
    Content     string `gorm:"type:varchar(256);not null"`
    CreateDate  int64  `gorm:"autoCreateTime:milli"`
}

func CreateComment(comment *Comment) (err error) {
  err = database.DB.Model(&Comment{}).Create(comment).Error
  return 
}

func GetCommentListByVideoID(videoID string) ([]Comment, error) {
    var comments []Comment
    err := database.DB.Where("video_id = ?", videoID).Order("create_date desc").Find(&comments).Error
    return comments, err
}