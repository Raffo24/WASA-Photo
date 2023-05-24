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
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const DELETED = "DELETED"

// Status Object
type Status struct {
	Status string
}

// User struct
type User struct {
	ID       int
	Username string
}

type UserBanFollow struct {
	ID       int
	Username string
	Followed bool
	Banned   bool
}

type UserExtended struct {
	ID        int
	Username  string
	Followers int
	Following int
	Photos    int
	Banned    int
}

// Photo struct
type Photo struct {
	ID          int
	UserID      int
	Username    string
	Photourl    string
	Title       string
	Description string
	CreatedAt   time.Time
	Comments    int
	Likes       int
	Liked       bool
}

// Comment struct
type Comment struct {
	ID        int
	PhotoID   int
	UserID    int
	Username  string
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
type JsonificaUsersBanFollow struct{ Items []UserBanFollow }
type JsonificaPhotos struct{ Items []Photo }
type JsonificaComments struct{ Items []Comment }

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserByUsername(username string) (User, error)
	GetUserByID(id int) (User, error)
	GetUserExtendedByID(id int) (UserExtended, error)
	GetUsers() ([]User, error)
	GetFollowersID(userID int) ([]User, error)
	GetFollowingID(userID int) ([]User, error)
	Ping() error
	GetPhotos(userID int) ([]Photo, error)
	GetPhoto(photoID int) (Photo, error)
	GetCommentsByPhotoID(photoID int) ([]Comment, error)
	GetLikes(photoID int) ([]User, error)
	GetFeed(userID int) ([]Photo, error)
	GetBansID(userID int) ([]User, error)
	AddUser(username string) (User, error)
	AddComment(photoID int, userID int, comment string) (Comment, error)
	AddLike(photoID int, userID int) (Like, error)
	AddFollow(followerID int, followingID int) (Follow, error)
	AddBan(bannedID int, bannerID int) (Ban, error)
	DeleteUser(id int) (Status, error)
	DeletePhoto(id int) (Status, error)
	DeleteComment(id int) (Status, error)
	DeleteLike(photoID int, userID int) (Status, error)
	DeleteFollow(followerID int, followingID int) (Status, error)
	DeleteBan(bannedID int, bannerID int) (Status, error)
	UpdateUser(id int, username string) (User, error)
	AddPhoto(id int, photourl string, title string, description string) (Photo, error)
	SearchUser(username string, UserID int) ([]UserBanFollow, error)
	GetCommentByID(id int) (Comment, error)
	UserIsPresent(id int) (bool, error)
	UserIsBanned(bannerID int, bannedID int) (bool, error)
	JsonificaUsersFun(users []UserBanFollow) JsonificaUsersBanFollow
	JsonificaPhotosFun(photos []Photo) JsonificaPhotos
	JsonificaCommentsFun(comments []Comment) JsonificaComments
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
	// crea la cartella per le foto se non esiste
	_, err = os.Stat("./service/api/images/")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("./service/api/images/", 0755)
		if errDir != nil {
			logrus.WithError(err).Error("error creating photos folder")
			return nil, fmt.Errorf("error creating photos folder: %w", err)
		}
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
			photourl varchar(1000) NOT NULL,
			title varchar(1000) NOT NULL,
			description varchar(1000) NOT NULL,
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

		CREATE TRIGGER IF NOT EXISTS delete_photos_on_user_delete
		AFTER DELETE ON users
		BEGIN
			DELETE FROM photos WHERE userid = OLD.id;
		END;

		CREATE TRIGGER IF NOT EXISTS delete_comments_on_user_delete
		AFTER DELETE ON users
		BEGIN
			DELETE FROM comments WHERE userid = OLD.id;
		END;

		CREATE TRIGGER IF NOT EXISTS delete_likes_on_user_delete
		AFTER DELETE ON users
		BEGIN
			DELETE FROM likes WHERE userid = OLD.id;
		END;

		CREATE TRIGGER IF NOT EXISTS delete_follows_on_user_delete
		AFTER DELETE ON users
		BEGIN
			DELETE FROM follows WHERE followerid = OLD.id;
			DELETE FROM follows WHERE followingid = OLD.id;
		END;

		CREATE TRIGGER IF NOT EXISTS delete_bans_on_user_delete
		AFTER DELETE ON users
		BEGIN
			DELETE FROM bans WHERE bannedid = OLD.id;
			DELETE FROM bans WHERE bannerid = OLD.id;
		END;

		CREATE TRIGGER IF NOT EXISTS delete_comments_on_photo_delete
		AFTER DELETE ON photos
		BEGIN
			DELETE FROM comments WHERE photoid = OLD.id;
		END;
			
		CREATE TRIGGER IF NOT EXISTS delete_likes_on_photo_delete
		AFTER DELETE ON photos
		BEGIN
			DELETE FROM likes WHERE photoid = OLD.id;
		END;

		CREATE TRIGGER IF NOT EXISTS unfollow_on_ban
		AFTER INSERT ON bans
		BEGIN
			DELETE FROM follows WHERE followerid = NEW.bannedid;
			DELETE FROM follows WHERE followingid = NEW.bannedid;
		END;
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
func (db *appdbimpl) JsonificaUsersFun(users []UserBanFollow) JsonificaUsersBanFollow {
	return JsonificaUsersBanFollow{
		Items: users,
	}
}
func (db *appdbimpl) JsonificaPhotosFun(photos []Photo) JsonificaPhotos {
	return JsonificaPhotos{
		Items: photos,
	}

}
func (db *appdbimpl) JsonificaCommentsFun(comments []Comment) JsonificaComments {
	return JsonificaComments{
		Items: comments,
	}
}
func (db *appdbimpl) GetPhotos(userID int) ([]Photo, error) {
	rows, err := db.c.Query("SELECT id, userid, photourl, title, description, createdat FROM photos where userid = ?", userID)
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
		var title string
		var description string
		var createdat time.Time
		err = rows.Scan(&id, &userid, &photourl, &title, &description, &createdat)
		if err != nil {
			return nil, err
		}
		var likes int
		var comments int
		var liked bool
		var username string
		err = db.c.QueryRow("SELECT count(*) FROM likes WHERE photoid = ?", id).Scan(&likes)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT count(*) FROM comments WHERE photoid = ?", id).Scan(&comments)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT count(*) FROM likes WHERE photoid = ? AND userid = ?", id, userID).Scan(&liked)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&username)
		if err != nil {
			return nil, err
		}
		photos = append(photos, Photo{
			ID:          id,
			UserID:      userid,
			Username:    username,
			Photourl:    photourl,
			Title:       title,
			Description: description,
			CreatedAt:   createdat,
			Likes:       likes,
			Comments:    comments,
			Liked:       liked,
		})
	}
	return photos, nil
}

func (db *appdbimpl) GetPhoto(photoID int) (Photo, error) {
	var photo Photo
	err := db.c.QueryRow("SELECT id, userid, photourl, title, description, createdat FROM photos WHERE id=?", photoID).Scan(&photo.ID, &photo.UserID, &photo.Photourl, &photo.Title, &photo.Description, &photo.CreatedAt)
	if err != nil {
		return Photo{}, err
	}
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
		var userID int
		var createdat time.Time
		err = rows.Scan(&id, &photoID, &userID, &comment, &createdat)
		if err != nil {
			return nil, err
		}
		var Username string
		err = db.c.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, Comment{
			ID:        id,
			PhotoID:   photoID,
			UserID:    userID,
			Username:  Username,
			Content:   comment,
			CreatedAt: createdat,
		})
	}
	return comments, nil
}
func (db *appdbimpl) GetLikes(photoID int) ([]User, error) {
	rows, err := db.c.Query("SELECT userid FROM likes WHERE photoid=?", photoID)
	if err != nil {
		return nil, err
	}
	var usersThatLiked []User
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
		user, err := db.GetUserByID(userID)
		if err != nil {
			return nil, err
		}
		usersThatLiked = append(usersThatLiked, user)
	}
	return usersThatLiked, nil
}
func (db *appdbimpl) GetFollowersID(userID int) ([]User, error) {
	rows, err := db.c.Query("SELECT followerid FROM follows WHERE followingid=?", userID)
	if err != nil {
		return nil, err
	}
	var seguono []User
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var followerID int
		err = rows.Scan(&followerID)
		if err != nil {
			return nil, err
		}
		user, err := db.GetUserByID(followerID)
		if err != nil {
			return nil, err
		}
		seguono = append(seguono, user)
	}
	return seguono, nil
}

func (db *appdbimpl) GetFollowingID(userID int) ([]User, error) {
	rows, err := db.c.Query("SELECT followingid FROM follows WHERE followerid=?", userID)
	if err != nil {
		return nil, err
	}
	var seguiti []User
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
		user, err := db.GetUserByID(followingID)
		if err != nil {
			return nil, err
		}
		seguiti = append(seguiti, user)
	}
	return seguiti, nil
}

func (db *appdbimpl) GetBansID(userID int) ([]User, error) {
	rows, err := db.c.Query("SELECT bannedid FROM bans WHERE bannerid=?", userID)
	if err != nil {
		return nil, err
	}
	var bannati []User
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
		user, err := db.GetUserByID(bannedID)
		if err != nil {
			return nil, err
		}
		bannati = append(bannati, user)
	}
	return bannati, nil
}

func (db *appdbimpl) GetFeed(userID int) ([]Photo, error) {
	rows, err := db.c.Query("SELECT id, userid, photourl, title, description, createdat FROM photos WHERE userid IN (SELECT followingid FROM follows WHERE followerid=?);", userID)
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
		var title string
		var description string
		var createdat time.Time
		err = rows.Scan(&id, &userID, &photourl, &title, &description, &createdat)
		if err != nil {
			return nil, err
		}
		var likes int
		var comments int
		var liked bool
		var username string
		err = db.c.QueryRow("SELECT count(*) FROM likes WHERE photoid = ?", id).Scan(&likes)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT count(*) FROM comments WHERE photoid = ?", id).Scan(&comments)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT count(*) FROM likes WHERE photoid = ? AND userid = ?", id, userID).Scan(&liked)
		if err != nil {
			return nil, err
		}
		err = db.c.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
		if err != nil {
			return nil, err
		}
		photos = append(photos, Photo{
			ID:          id,
			UserID:      userID,
			Username:    username,
			Photourl:    photourl,
			Title:       title,
			Description: description,
			CreatedAt:   createdat,
			Likes:       likes,
			Comments:    comments,
			Liked:       liked,
		})
	}
	return photos, nil
}

func (db *appdbimpl) GetUserByID(id int) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, username FROM users WHERE id=?", id).Scan(&user.ID, &user.Username)
	if err != nil {
		return User{}, err
	}
	return user, err
}
func (db *appdbimpl) GetUserExtendedByID(id int) (UserExtended, error) {
	var user UserExtended
	err := db.c.QueryRow("SELECT id, username FROM users WHERE id=?", id).Scan(&user.ID, &user.Username)
	if err != nil {
		return UserExtended{}, err
	}
	followers, err := db.GetFollowersID(id)
	if err != nil {
		return UserExtended{}, err
	}
	following, err := db.GetFollowingID(id)
	if err != nil {
		return UserExtended{}, err
	}
	photos, err := db.GetPhotos(id)
	if err != nil {
		return UserExtended{}, err
	}
	banned, err := db.GetBansID(id)
	if err != nil {
		return UserExtended{}, err
	}
	user.Photos = len(photos)
	user.Followers = len(followers)
	user.Following = len(following)
	user.Banned = len(banned)
	return user, err
}
func (db *appdbimpl) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.c.QueryRow("SELECT id, username FROM users WHERE username=?", username).Scan(&user.ID, &user.Username)
	if err != nil {
		return User{}, err
	}
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

func (db *appdbimpl) AddPhoto(userID int, photourl string, title string, description string) (Photo, error) {
	res, err := db.c.Exec("INSERT INTO photos (userid, photourl, title, description) VALUES (?, ?, ?, ?)", userID, photourl, title, description)
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

func (db *appdbimpl) DeleteUser(id int) (Status, error) {
	flag, err := db.UserIsPresent(id)
	if err != nil {
		logrus.Error(err)
	}
	if !flag {
		return Status{}, errors.New("USER NOT FOUND")
	}
	// delete all photos
	photos, err := db.GetPhotos(id)
	if err == nil {
		for idx := range photos {
			photo, err := db.GetPhoto(photos[idx].ID)
			if err != nil {
				logrus.Error(err)
			}
			if os.Remove(photo.Photourl) != nil {
				logrus.Error("IMAGE NOT found")
			}
		}
	}
	_, err = db.c.Exec("DELETE FROM users WHERE id=?", id)
	return Status{Status: DELETED}, err
}

func (db *appdbimpl) DeletePhoto(id int) (Status, error) {
	photo, err := db.GetPhoto(id)
	if err != nil {
		return Status{}, err
	}
	if os.Remove(photo.Photourl) != nil {
		return Status{}, errors.New("IMAGE NOT found")
	}
	_, err = db.c.Exec("DELETE FROM photos WHERE id=?", id)
	return Status{Status: DELETED}, err
}

func (db *appdbimpl) DeleteComment(id int) (Status, error) {
	_, err := db.c.Exec("DELETE FROM comments WHERE id=?", id)
	return Status{Status: DELETED}, err
}

func (db *appdbimpl) DeleteLike(photoID int, userID int) (Status, error) {
	_, err := db.c.Exec("DELETE FROM likes WHERE photoid=? AND userid=?", photoID, userID)
	return Status{Status: DELETED}, err
}

func (db *appdbimpl) DeleteFollow(followerID int, followingID int) (Status, error) {
	if followerID == followingID {
		return Status{}, errors.New("CAN'T UNFOLLOW YOURSELF")
	}
	_, err := db.c.Exec("DELETE FROM follows WHERE followerid=? AND followingid=?", followerID, followingID)
	return Status{Status: DELETED}, err
}

func (db *appdbimpl) DeleteBan(bannedID int, bannerID int) (Status, error) {
	if bannedID == bannerID {
		return Status{}, errors.New("CAN'T UNBAN YOURSELF")
	}
	_, err := db.c.Exec("DELETE FROM bans WHERE bannedid=? AND bannerid=?", bannedID, bannerID)
	return Status{Status: DELETED}, err
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

func (db *appdbimpl) SearchUser(search_username string, userID int) ([]UserBanFollow, error) {
	rows, err := db.c.Query("SELECT id, username FROM users WHERE username LIKE ?", "%"+search_username+"%")
	if err != nil {
		return nil, err
	}
	var users []UserBanFollow
	defer func() {
		_ = rows.Close()
		_ = rows.Err() // or modify return value
	}()
	for rows.Next() {
		var id int
		var username string
		var banned bool
		var followed bool
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, err
		}

		banned, err := db.UserIsBanned(userID, id)
		if err != nil {
			return nil, err
		}
		db.c.QueryRow("SELECT COUNT(*) FROM follows WHERE followerid=? AND followingid=?", userID, id).Scan(&followed)
		if err != nil {
			return nil, err
		}

		users = append(users, UserBanFollow{
			ID:       id,
			Username: username,
			Banned:   banned,
			Followed: followed,
		})
	}
	return users, nil
}

func (db *appdbimpl) GetCommentByID(id int) (Comment, error) {
	var comment Comment
	err := db.c.QueryRow("SELECT id, photoid, userid, comment, createdat FROM comments WHERE id=?", id).Scan(&comment.ID, &comment.PhotoID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		return Comment{}, err
	}
	return comment, err
}

func (db *appdbimpl) UserIsPresent(id int) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE id=?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (db *appdbimpl) UserIsBanned(idBanner int, idBanned int) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM bans WHERE bannerid=? AND bannedid=?", idBanner, idBanned).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}
