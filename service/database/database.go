/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const DELETED = "DELETED"

// User struct
type User struct {
	ID       int
	Username string
}

// Photo struct
type Photo struct {
	ID        int
	UserID    int
	Photourl  string
	CreatedAt time.Time
}

// Comment struct
type Comment struct {
	ID        int
	PhotoID   int
	UserID    int
	Content   string
	CreatedAt time.Time
}
type Ban struct {
	BannedID int
	BannerID int
}
type Like struct {
	UserID  int
	PhotoID int
}
type Follow struct {
	FollowerID  int
	FollowingID int
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserByUsername(username string) (User, error)
	GetUserByID(id int) (User, error)
	GetUsers() ([]User, error)
	GetFollowersID(userID int) ([]int, error)
	GetFollowingID(userID int) ([]int, error)
	Ping() error
	GetPhotos(userID int) ([]Photo, error)
	GetPhoto(photoID int) (Photo, error)
	GetCommentsByPhotoID(photoID int) ([]Comment, error)
	GetLikes(photoID int) ([]int, error)
	GetFeed(userID int) ([]Photo, error)
	GetBansID(userID int) ([]int, error)
	AddUser(username string) (User, error)
	AddComment(photoID int, userID int, comment string) (Comment, error)
	AddLike(photoID int, userID int) (Like, error)
	AddFollow(followerID int, followingID int) (Follow, error)
	AddBan(bannedID int, bannerID int) (Ban, error)
	DeleteUser(id int) (string, error)
	DeletePhoto(id int) (string, error)
	DeleteComment(id int) (string, error)
	DeleteLike(photoID int, userID int) (string, error)
	DeleteFollow(followerID int, followingID int) (string, error)
	DeleteBan(bannedID int, bannerID int) (string, error)
	UpdateUser(id int, username string) (User, error)
	AddPhoto(id int, photourl string) (Photo, error)
	SearchUser(username string) ([]User, error)
	GetCommentByID(id int) (Comment, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	// check if the database is ready
	err := createTables(db)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}
	return &appdbimpl{
		c: db,
	}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS photos (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			photourl varchar(100) NOT NULL,
			createdat DATETIME DEFAULT CURRENT_TIMESTAMP,
			foreign key (userid) references users(id)
		);
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			photoid INTEGER NOT NULL,
			userid INTEGER NOT NULL,
			comment TEXT NOT NULL,
			createdat DATETIME DEFAULT CURRENT_TIMESTAMP,
			foreign key (photoid) references photos(id),
			foreign key (userid) references users(id)
		);
		CREATE TABLE IF NOT EXISTS likes (
			photoid INTEGER NOT NULL,
			userid INTEGER NOT NULL,
			PRIMARY KEY (photoID, userID)
			foreign key (photoid) references photos(id),
			foreign key (userid) references users(id)
		);
		CREATE TABLE IF NOT EXISTS follows (
			followerid INTEGER NOT NULL,
			followingid INTEGER NOT NULL,
			PRIMARY KEY (followerid, followingid)
			foreign key (followerId) references users(id),
			foreign key (followingId) references users(id)
		);
		CREATE TABLE IF NOT EXISTS bans (
			bannedid INTEGER NOT NULL,
			bannerid INTEGER NOT NULL,
			PRIMARY KEY (bannedid, bannerid)
			foreign key (bannedId) references users(id),
			foreign key (bannerId) references users(id)
		);
	`)
	return err
}
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
func (db *appdbimpl) GetUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT id, username, name FROM users")
	if err != nil {
		return nil, err
	}
	var users []User
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var id int
		var username string
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, err
		}
		users = append(users, User{
			ID:       id,
			Username: username,
		})
	}
	return users, nil
}

func (db *appdbimpl) GetPhotos(userID int) ([]Photo, error) {
	rows, err := db.c.Query("SELECT id, userid, photourl, createdat FROM photos where userid = ?", userID)
	if err != nil {
		return nil, err
	}
	var photos []Photo
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var id int
		var userid int
		var photourl string
		var createdat time.Time
		err = rows.Scan(&id, &userid, &photourl, &createdat)
		if err != nil {
			return nil, err
		}
		photos = append(photos, Photo{
			ID:        id,
			UserID:    userid,
			Photourl:  photourl,
			CreatedAt: createdat,
		})
	}
	return photos, nil
}

func (db *appdbimpl) GetPhoto(photoID int) (Photo, error) {
	var photo Photo
	err := db.c.QueryRow("SELECT id, userid, photourl, createdat FROM photos WHERE id=?", photoID).Scan(&photo.ID, &photo.UserID, &photo.Photourl, &photo.CreatedAt)
	return photo, err
}
func (db *appdbimpl) GetCommentsByPhotoID(photoID int) ([]Comment, error) {
	rows, err := db.c.Query("SELECT id, photoid, userid, comment, createdat FROM comments WHERE photoid=?", photoID)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var id int
		var comment string
		var createdat time.Time
		err = rows.Scan(&id, &photoID, &photoID, &comment, &createdat)
		if err != nil {
			return nil, err
		}
		comments = append(comments, Comment{
			ID:        id,
			PhotoID:   photoID,
			UserID:    photoID,
			Content:   comment,
			CreatedAt: createdat,
		})
	}
	return comments, nil
}
func (db *appdbimpl) GetLikes(photoID int) ([]int, error) {
	rows, err := db.c.Query("SELECT userid FROM likes WHERE photoid=?", photoID)
	if err != nil {
		return nil, err
	}
	var usersIDThatLiked []int
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var userID int
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		usersIDThatLiked = append(usersIDThatLiked, userID)
	}
	return usersIDThatLiked, nil
}
func (db *appdbimpl) GetFollowersID(userID int) ([]int, error) {
	rows, err := db.c.Query("SELECT followerid FROM follows WHERE followingid=?", userID)
	if err != nil {
		return nil, err
	}
	var seguono []int
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var followingID int
		err = rows.Scan(&followingID)
		if err != nil {
			return nil, err
		}
		seguono = append(seguono, followingID)
	}
	return seguono, nil
}

func (db *appdbimpl) GetFollowingID(userID int) ([]int, error) {
	rows, err := db.c.Query("SELECT followingid FROM follows WHERE followerid=?", userID)
	if err != nil {
		return nil, err
	}
	var seguiti []int
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var followingID int
		err = rows.Scan(&followingID)
		if err != nil {
			return nil, err
		}
		seguiti = append(seguiti, followingID)
	}
	return seguiti, nil
}

func (db *appdbimpl) GetBansID(userID int) ([]int, error) {
	rows, err := db.c.Query("SELECT bannedid FROM bans WHERE bannerid=?", userID)
	if err != nil {
		return nil, err
	}
	var bannati []int
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var bannedID int
		err = rows.Scan(&bannedID)
		if err != nil {
			return nil, err
		}
		bannati = append(bannati, bannedID)
	}
	return bannati, nil
}

func (db *appdbimpl) GetFeed(userID int) ([]Photo, error) {
	rows, err := db.c.Query("SELECT id, userid, photo, createdat FROM photos WHERE userid IN (SELECT followingid FROM follows WHERE followerid=?);", userID)
	if err != nil {
		return nil, err
	}
	var photos []Photo
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var id int
		var userID int
		var photourl string
		var createdat time.Time
		err = rows.Scan(&id, &userID, &photourl, &createdat)
		if err != nil {
			return nil, err
		}
		photos = append(photos, Photo{
			ID:        id,
			UserID:    userID,
			Photourl:  photourl,
			CreatedAt: createdat,
		})
	}
	return photos, nil
}

func (db *appdbimpl) GetUserByID(id int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, username FROM users WHERE id=?", id).Scan(&user.ID, &user.Username)
	return user, err
}
func (db *appdbimpl) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, username FROM users WHERE username=?", username).Scan(&user.ID, &user.Username)
	return user, err
}

func (db *appdbimpl) AddUser(username string) (User, error) {
	res, err := db.c.Exec("INSERT INTO users (username) VALUES (?)", username)
	if err != nil {
		return User{}, err
	}
	id, err := res.LastInsertId()
	return User{
		ID:       int(id),
		Username: username,
	}, err
}

func (db *appdbimpl) AddPhoto(userID int, photourl string) (Photo, error) {
	res, err := db.c.Exec("INSERT INTO photos (userid, photourl) VALUES (?, ?)", userID, photourl)
	if err != nil {
		return Photo{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Photo{}, err
	}
	photok, err := db.GetPhoto(int(id))
	return photok, err
}

func (db *appdbimpl) AddComment(photoID int, userID int, comment string) (Comment, error) {
	res, err := db.c.Exec("INSERT INTO comments (photoid, userid, comment) VALUES (?, ?, ?)", photoID, userID, comment)
	if err != nil {
		return Comment{}, err
	}
	id, err := res.LastInsertId()
	return Comment{
		ID:      int(id),
		PhotoID: photoID,
		UserID:  userID,
		Content: comment,
	}, err
}

func (db *appdbimpl) AddLike(photoID int, userID int) (Like, error) {
	_, err := db.c.Exec("INSERT INTO likes (photoid, userid) VALUES (?, ?)", photoID, userID)
	if err != nil {
		return Like{}, err
	}
	return Like{
		PhotoID: photoID,
		UserID:  userID,
	}, err
}

func (db *appdbimpl) AddFollow(followerID int, followingID int) (Follow, error) {
	_, err := db.c.Exec("INSERT INTO follows (followerid, followingid) VALUES (?, ?)", followerID, followingID)
	if err != nil {
		return Follow{}, err
	}
	return Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}, err
}

func (db *appdbimpl) AddBan(bannedID int, bannerID int) (Ban, error) {
	_, err := db.c.Exec("INSERT INTO bans (bannedid, bannerid) VALUES (?, ?)", bannedID, bannerID)
	if err != nil {
		return Ban{}, err
	}
	return Ban{
		BannedID: bannedID,
		BannerID: bannerID,
	}, err
}

func (db *appdbimpl) DeleteUser(id int) (string, error) {
	_, err := db.c.Exec("DELETE FROM users WHERE id=?", id)
	return DELETED, err
}

func (db *appdbimpl) DeletePhoto(id int) (string, error) {
	_, err := db.c.Exec("DELETE FROM photos WHERE id=?", id)
	return DELETED, err
}

func (db *appdbimpl) DeleteComment(id int) (string, error) {
	_, err := db.c.Exec("DELETE FROM comments WHERE id=?", id)
	return DELETED, err
}

func (db *appdbimpl) DeleteLike(photoID int, userID int) (string, error) {
	_, err := db.c.Exec("DELETE FROM likes WHERE photoid=? AND userid=?", photoID, userID)
	return DELETED, err
}

func (db *appdbimpl) DeleteFollow(followerID int, followingID int) (string, error) {
	if followerID == followingID {
		return "", errors.New("CAN'T UNFOLLOW YOURSELF")
	}
	_, err := db.c.Exec("DELETE FROM follows WHERE followerid=? AND followingid=?", followerID, followingID)
	return DELETED, err
}

func (db *appdbimpl) DeleteBan(bannedID int, bannerID int) (string, error) {
	if bannedID == bannerID {
		return "", errors.New("CAN'T UNBAN YOURSELF")
	}
	_, err := db.c.Exec("DELETE FROM bans WHERE bannedid=? AND bannerid=?", bannedID, bannerID)
	return DELETED, err
}

func (db *appdbimpl) UpdateUser(id int, username string) (User, error) {
	_, err := db.c.Exec("UPDATE users SET username=? WHERE id=?", username, id)
	if err != nil {
		logrus.Error(err)
		return User{}, err
	}
	user, err := db.GetUserByID(id)
	return user, err
}

func (db *appdbimpl) SearchUser(search_username string) ([]User, error) {
	rows, err := db.c.Query("SELECT id, username FROM users WHERE username LIKE ?", "%"+search_username+"%")
	if err != nil {
		return nil, err
	}
	var users []User
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var id int
		var username string
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, err
		}
		users = append(users, User{
			ID:       id,
			Username: username,
		})
	}
	return users, nil
}

func (db *appdbimpl) GetCommentByID(id int) (Comment, error) {
	var comment Comment
	err := db.c.QueryRow("SELECT id, photoid, userid, comment, createdat FROM comments WHERE id=?", id).Scan(&comment.ID, &comment.PhotoID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	return comment, err
}
