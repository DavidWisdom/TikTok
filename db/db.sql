DROP DATABASE IF EXISTS tiktok;
CREATE DATABASE tiktok;
USE tiktok;
DROP TABLE IF EXISTS User;
CREATE TABLE User (
    user_id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(32) UNIQUE NOT NULL,
    password VARCHAR(32) NOT NULL,
    follow_count INT DEFAULT 0,
    follower_count INT DEFAULT 0,
    avatar VARCHAR(256) DEFAULT NULL,
    background_image VARCHAR(256) DEFAULT NULL,
    signature VARCHAR(32) DEFAULT NULL,
    total_favorited INT DEFAULT 0,
    work_count INT DEFAULT 0,
    favorite_count INT DEFAULT 0,
    PRIMARY KEY (user_id)
);
DROP TABLE IF EXISTS Video;
CREATE TABLE Video (
    video_id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    play_url VARCHAR(512) NOT NULL,
    cover_url VARCHAR(512) NOT NULL,
    favorite_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    title VARCHAR(32) NOT NULL,
    PRIMARY KEY (video_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE
);
USE tiktok;
DROP TABLE IF EXISTS Likes;
CREATE TABLE Likes (
    user_id INT,
    video_id INT,
    PRIMARY KEY (user_id, video_id),
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (video_id) REFERENCES Video(video_id) ON DELETE CASCADE
);
USE tiktok;
DROP TABLE IF EXISTS Comment;
CREATE TABLE Comment (
    comment_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    video_id INT,
    content VARCHAR(256) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (video_id) REFERENCES Video(video_id) ON DELETE CASCADE
);
USE tiktok;
DROP TABLE IF EXISTS Message;
CREATE TABLE Message (
    message_id INT PRIMARY KEY AUTO_INCREMENT,
    from_user_id INT,
    to_user_id INT,
    content VARCHAR(256) NOT NULL,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (to_user_id) REFERENCES User(user_id) ON DELETE CASCADE
);
USE tiktok;
DROP TABLE IF EXISTS Follow;
CREATE TABLE Follow (
    from_user_id INT,
    to_user_id INT,
    is_mutual BOOLEAN,
    PRIMARY KEY (from_user_id, to_user_id),
    FOREIGN KEY (from_user_id) REFERENCES User(user_id) ON DELETE CASCADE,
    FOREIGN KEY (to_user_id) REFERENCES User(user_id) ON DELETE CASCADE
);