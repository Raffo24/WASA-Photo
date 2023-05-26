package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"wasaPhoto/service/database"

	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

// log the error wrapper
func logerr(n int, err error) {
	if err != nil {
		logrus.Printf("Write failed: %v", err)
	}
}

// PRIMITIVE FUNCTIONS TO HELP THE HANDLERS
// add responce code of success in the parameter
func finalize(output interface{}, err error, w http.ResponseWriter, code int) {
	// set response code of success
	w.WriteHeader(code)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		if json.NewEncoder(w).Encode(output) != nil {
			w.WriteHeader(500)
			logerr(w.Write([]byte("Internal Server Error")))
		}
	} else {
		w.WriteHeader(400)
		logerr(w.Write([]byte(err.Error())))
	}
}
func (rt *_router) youAreLogged(r *http.Request, w http.ResponseWriter) (bool, int) {
	myID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	if err != nil {
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non sei loggato")))
		return true, 0
	}
	bool, err := rt.db.UserIsPresent(myID)
	if err != nil {
		w.WriteHeader(500)
		logerr(w.Write([]byte("server error")))
		return true, 0
	}
	if !bool {
		w.WriteHeader(401)
		logerr(w.Write([]byte("non sei loggato")))
		return true, 0
	}
	return false, myID
}

func (rt *_router) youAreBanned(myID int, banner int, r *http.Request, w http.ResponseWriter) bool {
	bool, err := rt.db.UserIsBanned(banner, myID)
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("server error")))
		return false
	}
	return bool
}
func (rt *_router) securityChecker(bannerID int, r *http.Request, w http.ResponseWriter) bool {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return true
	}
	if rt.youAreBanned(myID, bannerID, r, w) {
		w.WriteHeader(403)
		logerr(w.Write([]byte("sei bannato")))
		return true
	}
	return false
}

func (rt *_router) getUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreLogged(r, w)
	if flag {
		return
	}
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if rt.securityChecker(userID, r, w) {
		return
	}
	user, err := rt.db.GetUserExtendedByID(userID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("user id not exist")))
		return
	}
	finalize(user, err, w, 200)
}
func (rt *_router) getUserPhotosHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreLogged(r, w)
	if flag {
		return
	}
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if rt.securityChecker(userID, r, w) {
		return
	}
	output, err := rt.db.GetPhotos(userID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("user id not exist")))
		return
	}
	finalize(rt.db.JsonificaPhotosFun(output), err, w, 200)
}
func (rt *_router) getFollowersHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreLogged(r, w)
	if flag {
		return
	}
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if rt.securityChecker(userID, r, w) {
		return
	}
	output, err := rt.db.GetFollowersID(userID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("user id not exist")))
		return
	}
	finalize(output, err, w, 200)
}

func (rt *_router) getFollowingHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreLogged(r, w)
	if flag {
		return
	}
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if rt.securityChecker(userID, r, w) {
		return
	}
	output, err := rt.db.GetFollowingID(userID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("user id not exist")))
		return
	}
	finalize(output, err, w, 200)
}

func (rt *_router) searchUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// read params query
	flag, userID := rt.youAreLogged(r, w)
	if flag {
		return
	}

	username_searched := r.URL.Query().Get("query")
	if username_searched == "" {
		w.WriteHeader(400)
		logerr(w.Write([]byte("query is empty")))
		return
	}
	output, err := rt.db.SearchUser(username_searched, userID)
	finalize(rt.db.JsonificaUsersFun(output), err, w, 200)
}
func (rt *_router) getFeedHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, userID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	output, err := rt.db.GetFeed(userID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("user id not exist")))
		return
	}
	finalize(output, err, w, 200)
}
func (rt *_router) getPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreLogged(r, w)
	if flag {
		return
	}
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	photo, err := rt.db.GetPhoto(id)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("photo not found in db")))
		return
	}
	// check if user is authorized to see the photo
	if rt.securityChecker(photo.UserID, r, w) {
		return
	}
	// read photo from disk
	file, err := os.Open(photo.Photourl)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("photo not found")))
		return
	}
	fi, err := file.Stat()
	if err != nil {
		w.WriteHeader(500)
		logerr(w.Write([]byte("internal error reading data of the photo")))
		return
	}
	defer file.Close()

	// write form-data to response with these field : BytePhoto, description, title
	w.Header().Set("Content-Disposition", "attachment; filename="+photo.Title)
	w.Header().Set("Filename", photo.Title)
	w.Header().Set("Description", photo.Description)
	w.Header().Set("CreatedAt", photo.CreatedAt.String())
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(int(fi.Size())))
	if _, err := io.Copy(w, file); err != nil {
		w.WriteHeader(500)
		logerr(w.Write([]byte("internal error writing photo to response")))
		return
	}
}
func (rt *_router) getAllCommentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreLogged(r, w)
	if flag {
		return
	}
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("photo id is empty")))
		return
	}
	output, err := rt.db.GetCommentsByPhotoID(id)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("photo not found in db")))
		return
	}
	finalize(rt.db.JsonificaCommentsFun(output), err, w, 200)
}

// POST REQUEST
func (rt *_router) loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user1 := database.User{}
	err := json.NewDecoder(r.Body).Decode(&user1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	code := 200
	user, err := rt.db.GetUserByUsername(user1.Username)
	if err != nil {
		user, err = rt.db.AddUser(user1.Username)
		code = 201
	}
	finalize(user, err, w, code)
}

func (rt *_router) changeMyNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	id, err := strconv.Atoi(ps.ByName("id"))
	if myID != id {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi cambiare il nome di un altro utente")))
		return
	}
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	user := database.User{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := rt.db.UpdateUser(id, user.Username)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("user not found in db")))
		return
	}
	finalize(output, err, w, 201)
}
func (rt *_router) uploadPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, userID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	// send a form multipart with 3 fields: title, description, photo(jpeg)
	// read the form
	err := r.ParseMultipartForm(100)
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("error parsing form")))
		return
	}
	// read the fields
	title := r.FormValue("title")
	description := r.FormValue("description")
	// write the file
	file, _, err := r.FormFile("photo")
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("error reading photo")))
		return
	}
	defer file.Close()
	// store uploaded file into local path
	imageUrl := "service/api/images/" + strconv.Itoa(userID) + "_" + strconv.Itoa(int(time.Now().Unix())) + ".jpg"
	f, err := os.Create(imageUrl)
	if err != nil {
		w.WriteHeader(500)
		logerr(w.Write([]byte("internal error create file")))
		return
	}
	defer f.Close()

	// copy the uploaded file into the local filesystem
	_, err = io.Copy(f, file)
	if err != nil {
		w.WriteHeader(500)
		logerr(w.Write([]byte("internal error copy file")))
		return
	}
	output, err := rt.db.AddPhoto(userID, imageUrl, title, description)
	finalize(output, err, w, 201)
}

func (rt *_router) addCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, userID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	comment := database.Comment{}
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := rt.db.AddComment(id, userID, comment.Content)
	finalize(output, err, w, 201)
}

// PUT REQUEST
func (rt *_router) likePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	userID, err := strconv.Atoi(ps.ByName("userId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if myID != userID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi mettere like da un altro utente")))
		return
	}
	photoID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	output, err := rt.db.AddLike(photoID, userID)
	finalize(output, err, w, 201)
}
func (rt *_router) followUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	followerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	followingID, err := strconv.Atoi(ps.ByName("followId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("follow id is empty")))
		return
	}
	if myID != followerID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi seguire utilizzando un altro utente")))
		return
	}
	if rt.youAreBanned(myID, followingID, r, w) {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi seguire un utente che ti ha bannato")))
		return
	}
	output, err := rt.db.AddFollow(followerID, followingID)
	finalize(output, err, w, 201)
}

func (rt *_router) banUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	bannerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if myID != bannerID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi bannare utilizzando un altro utente")))
		return
	}
	bannedID, err := strconv.Atoi(ps.ByName("banId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("ban id is empty")))
		return
	}
	output, err := rt.db.AddBan(bannedID, bannerID)
	finalize(output, err, w, 201)
}

// DELETE REQUEST
/*func (rt *_router) deleteUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if myID != userID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi cancellare un utente che non ti appartiene")))
		return
	}
	output, err := rt.db.DeleteUser(userID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("user id not exist")))
		return
	}
	finalize(output, err, w, 200)
}
*/

func (rt *_router) deleteCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	commentID, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("comment id is empty")))
		return
	}
	comment, err := rt.db.GetCommentByID(commentID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("comment not found")))
		return
	}
	if comment.UserID != myID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi cancellare un commento che non ti appartiene")))
		return
	}
	output, err := rt.db.DeleteComment(commentID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("comment not found")))
		return
	}
	finalize(output, err, w, 200)
}
func (rt *_router) deletePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	photoID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("photo id is empty")))
		return
	}
	photo, err := rt.db.GetPhoto(photoID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("photo not found")))
		return
	}
	if photo.UserID != myID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi cancellare una foto che non ti appartiene")))
		return
	}
	output, err := rt.db.DeletePhoto(photoID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("photo not found")))
		return
	}
	finalize(output, err, w, 200)
}
func (rt *_router) unlikePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	photoID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("photo id is empty")))
		return
	}
	userID, err := strconv.Atoi(ps.ByName("userId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("user id is empty")))
		return
	}
	if myID != userID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi levare un like che non hai messo tu")))
		return
	}
	output, err := rt.db.DeleteLike(photoID, userID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("like not found")))
		return
	}
	finalize(output, err, w, 200)
}
func (rt *_router) unfollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	followerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("follower id is empty")))
		return
	}
	if myID != followerID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi smettere di seguire da un account che non ti appartiene")))
		return
	}
	followingID, err := strconv.Atoi(ps.ByName("followId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("following id is empty")))
		return
	}
	output, err := rt.db.DeleteFollow(followerID, followingID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("follow not exist")))
		return
	}
	finalize(output, err, w, 200)
}

func (rt *_router) unbanUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreLogged(r, w)
	if flag {
		return
	}
	bannerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("banner id is empty")))
		return
	}
	if myID != bannerID {
		w.WriteHeader(403)
		logerr(w.Write([]byte("Non puoi sbannare da un account che non ti appartiene")))
		return
	}
	bannedID, err := strconv.Atoi(ps.ByName("banId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("banned id is empty")))
		return
	}
	output, err := rt.db.DeleteBan(bannedID, bannerID)
	if err != nil {
		w.WriteHeader(404)
		logerr(w.Write([]byte("ban not exist")))
		return
	}
	finalize(output, err, w, 200)
}
