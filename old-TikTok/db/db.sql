DROP DATABASE IF EXISTS tiktok;
CREATE DATABASE tiktok;

USE tiktok;
DROP TABLE IF EXISTS User;
CREATE TABLE User (
    id INT NOT NULL AUTO_INCREMENT, -- user_id自动加一
    username VARCHAR(255) NOT NULL,
    nickname VARCHAR(255) DEFAULT "abc",
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255) UNIQUE,
    is_follow INT DEFAULT 0,
    follow_count INT DEFAULT 0,
    fans_count INT DEFAULT 0,
    video_count INT DEFAULT 0,
    avatar VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (id)
);

USE tiktok;
DROP TABLE IF EXISTS Video;
CREATE TABLE Video (
   video_id INT NOT NULL AUTO_INCREMENT, -- video_id自动加一
   id INT NOT NULL,
   title VARCHAR(255) NOT NULL, -- 视频名称
   play_url VARCHAR (255) NOT NULL,
   cover_image VARCHAR(255) NOT NULL,
   views_count INT DEFAULT 0, -- 观看数量
   likes_count INT DEFAULT 0, -- 喜爱数
   comments_count INT DEFAULT 0,-- 评论数
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 视频发布时间
   PRIMARY KEY (video_id),
   FOREIGN KEY (id) REFERENCES User(id) ON DELETE CASCADE
);

USE tiktok;
DROP TABLE IF EXISTS Likes;
CREATE TABLE Likes (
    user_id INT,
    video_id INT,
    PRIMARY KEY(user_id, video_id),
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE,
    FOREIGN KEY (video_id) REFERENCES Video(video_id) ON DELETE CASCADE
);
USE tiktok;
DROP TABLE IF EXISTS Comment;
CREATE TABLE Comment (
    comment_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    video_id INT,
    content VARCHAR(256) NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE,
    FOREIGN KEY (video_id) REFERENCES Video(video_id) ON DELETE CASCADE
);
-- 聊天表1.0
USE tiktok;
DROP TABLE IF EXISTS chat_table;
CREATE TABLE chat_table (
    -- 每一次聊天消息的唯一标识
    chat_id INT PRIMARY KEY,
    -- 发送方的用户id
    sender_id INT,
    -- 接收方的用户id
    receiver_id INT,
    -- 聊天消息的内容
    message VARCHAR(255),
    -- 发送消息的时间戳，根据时间戳来排序
    timestamp TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES User(id),
    FOREIGN KEY (receiver_id) REFERENCES User(id)
);
-- 关注表1.0
USE tiktok;
DROP TABLE IF EXISTS following_table;
CREATE TABLE following_table (
    -- 关注
    attention INT,
    -- 粉丝
    fans INT,
    -- -- 互相关注为true
    mutual BOOLEAN,
    -- 关注和粉丝都作为主键，唯一不为空且不重复
    PRIMARY KEY (attention, fans),
    FOREIGN KEY (attention) REFERENCES User(id),
    FOREIGN KEY (fans) REFERENCES User(id)
);