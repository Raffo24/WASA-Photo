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
)

// User struct
type User struct {
	ID        int
	Username  string
	Following []User
	Followers []User
	Photos    []Photo
}

// Photo struct
type Photo struct {
	ID        int
	UserID    int
	User      User
	LikeCount int
	Comment   []Comment
	PhotoURL  string
	CreatedAt time.Time
}

// Comment struct
type Comment struct {
	ID        int
	PhotoID   int
	UserID    int
	User      User
	Comment   string
	CreatedAt time.Time
}
type Ban struct {
	ID       int
	BannedID int
	BannerID int
}
type Like struct {
	ID       int
	UserID   int
	PhotoID  int
	Username string
}
type Follow struct {
	ID          int
	FollowerID  int
	FollowingID int
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserID(username string) (int, error)
	GetUsername(userID int) (string, error)
	getPhotos(userID int) ([]Photo, error)
	getPhoto(photoID int) (Photo, error)
	getComments(photoID int) ([]Comment, error)
	getLikes(photoID int) ([]int, error)
	getFollowers(userID int) ([]int, error)
	getFollowing(userID int) ([]int, error)
	getFeed(userID int) ([]Photo, error)
	getBans(userID int) ([]int, error)
	getLogin(username string) (string, error)
	getUser(username string) (User, error)
	Ping() error
	GetUsers() ([]User, error)
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

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createTables(db)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	return &appdbimpl{
		c: db,
	}, nil
}
func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
		);
		CREATE TABLE IF NOT EXISTS photos (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userid INTEGER NOT NULL,
			likecount INTEGER NOT NULL,
			photourl TEXT NOT NULL,
			createdat INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			photoid INTEGER NOT NULL,
			userid INTEGER NOT NULL,
			comment TEXT NOT NULL,
			createdat INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS likes (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			photoid INTEGER NOT NULL,
			userid INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS follows (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			followerid INTEGER NOT NULL,
			followingid INTEGER NOT NULL
		);
		CREATE TABLE IF NOT EXISTS bans (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			bannedid INTEGER NOT NULL,
			bannerid INTEGER NOT NULL
		);
	`)
	return err
}
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

func (db *appdbimpl) GetUserID(username string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&id)
	return id, err
}
func (db *appdbimpl) GetUsername(userID int) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT username FROM users WHERE id=?", userID).Scan(&username)
	return username, err
}
func (db *appdbimpl) GetUsers() ([]User, error) {
	rows, err := db.c.Query("SELECT id, username, name FROM users")
	if err != nil {
		return nil, err
	}
	var users []User
	defer rows.Close()
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

func (db *appdbimpl) getPhotos(userID int) ([]Photo, error) {
	rows, err := db.c.Query("SELECT id, userid, likecount, photourl, createdat FROM photos WHERE userid=?", userID)
	if err != nil {
		return nil, err
	}
	var photos []Photo
	defer rows.Close()
	for rows.Next() {
		var id int
		var likecount int
		var photourl string
		var createdat int64
		err = rows.Scan(&id, &userID, &likecount, &photourl, &createdat)
		if err != nil {
			return nil, err
		}
		photos = append(photos, Photo{
			ID:        id,
			UserID:    userID,
			LikeCount: likecount,
			PhotoURL:  photourl,
			CreatedAt: time.Unix(createdat, 0),
		})
	}
	return photos, nil
}
func (db *appdbimpl) getPhoto(photoID int) (Photo, error) {
	var photo Photo
	err := db.c.QueryRow("SELECT id, userid, likecount, photourl, createdat FROM photos WHERE id=?", photoID).Scan(&photo.ID, &photo.UserID, &photo.LikeCount, &photo.PhotoURL, &photo.CreatedAt)
	return photo, err
}
func (db *appdbimpl) getComments(photoID int) ([]Comment, error) {
	rows, err := db.c.Query("SELECT id, photoid, userid, comment, createdat FROM comments WHERE photoid=?", photoID)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	defer rows.Close()
	for rows.Next() {
		var id int
		var comment string
		var createdat int64
		err = rows.Scan(&id, &photoID, &photoID, &comment, &createdat)
		if err != nil {
			return nil, err
		}
		comments = append(comments, Comment{
			ID:        id,
			PhotoID:   photoID,
			UserID:    photoID,
			Comment:   comment,
			CreatedAt: time.Unix(createdat, 0),
		})
	}
	return comments, nil
}
func (db *appdbimpl) getLikes(photoID int) ([]Like, error) {
	rows, err := db.c.Query("SELECT id, photoid, userid FROM likes WHERE photoid=?", photoID)
	if err != nil {
		return nil, err
	}
	var likes []Like
	defer rows.Close()
	for rows.Next() {
		var id int
		var userID int
		err = rows.Scan(&id, &photoID, &userID)
		if err != nil {
			return nil, err
		}
		likes = append(likes, Like{
			ID:      id,
			PhotoID: photoID,
			UserID:  userID,
		})
	}
	return likes, nil
}
func (db *appdbimpl) getFollows(userID int) ([]Follow, error) {
	rows, err := db.c.Query("SELECT id, followerid, followingid FROM follows WHERE followerid=?", userID)
	if err != nil {
		return nil, err
	}
	var follows []Follow
	defer rows.Close()
	for rows.Next() {
		var id int
		var followingID int
		err = rows.Scan(&id, &userID, &followingID)
		if err != nil {
			return nil, err
		}
		follows = append(follows, Follow{
			ID:          id,
			FollowerID:  userID,
			FollowingID: followingID,
		})
	}
	return follows, nil
}
func (db *appdbimpl) getBans(userID int) ([]Ban, error) {
	rows, err := db.c.Query("SELECT id, bannedid, bannerid FROM bans WHERE bannerid=?", userID)
	if err != nil {
		return nil, err
	}
	var bans []Ban
	defer rows.Close()
	for rows.Next() {
		var id int
		var bannedID int
		err = rows.Scan(&id, &bannedID, &userID)
		if err != nil {
			return nil, err
		}
		bans = append(bans, Ban{
			ID:       id,
			BannedID: bannedID,
			BannerID: userID,
		})
	}
	return bans, nil
}
