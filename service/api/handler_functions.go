package api

import (
	"bytes"
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
func finalize(output interface{}, err error, w http.ResponseWriter) {
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		if json.NewEncoder(w).Encode(output) != nil {
			w.WriteHeader(500)
			logerr(w.Write([]byte("Internal Server Error")))
		}
	} else {
		w.WriteHeader(400)
		// rt.baseLogger.WithError(err).Warning("getMyInfoHandler failed")
		// w.writeHeader(http.statusBadRequest)
		logerr(w.Write([]byte(err.Error())))
	}
}
func (rt *_router) youAreNotAuthorized(r *http.Request, w http.ResponseWriter) (bool, int) {
	myID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("Non sei loggato")))
		return true, 0
	}
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("Non sei loggato")))
		return true, 0
	}
	bool, err := rt.db.UserIsPresent(myID)
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("server error")))
		return true, 0
	}
	if !bool {
		w.WriteHeader(400)
		logerr(w.Write([]byte("utente non presente")))
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
	flag, myID := rt.youAreNotAuthorized(r, w)
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

// GET REQUEST
func (rt *_router) getApiStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode("l'api sta funzionando correttamente") != nil {
		w.WriteHeader(500)
		logerr(w.Write([]byte("Internal Server Error")))
	}
}

func (rt *_router) getUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	if rt.securityChecker(userID, r, w) {
		return
	}
	user, err := rt.db.GetUserByID(userID)
	finalize(user, err, w)
}
func (rt *_router) getUserPhotosHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	finalize(output, err, w)
}

func (rt *_router) searchUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username_searched := ps.ByName("id")
	output, err := rt.db.SearchUser(username_searched)
	finalize(output, err, w)
}
func (rt *_router) getFeedHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, userID := rt.youAreNotAuthorized(r, w)
	if flag {
		return
	}
	output, err := rt.db.GetFeed(userID)
	finalize(output, err, w)
}
func (rt *_router) getPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("id is empty")))
		return
	}
	photo, err := rt.db.GetPhoto(id)
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("photo not found in db")))
		return
	}
	// check if user is authorized to see the photo
	if rt.securityChecker(photo.UserID, r, w) {
		return
	}
	// read photo from disk

	bytePhoto, err := os.ReadFile(photo.Photourl)
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("photo not found on disk")))
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytePhoto)))
	w.Write(bytePhoto)
}

func (rt *_router) getCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreNotAuthorized(r, w)
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
	finalize(comment, err, w)
}
func (rt *_router) getAllCommentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, _ := rt.youAreNotAuthorized(r, w)
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
	finalize(output, err, w)
}

// SPECIAL GET REQUEST
// liveness is an HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
func (rt *_router) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := rt.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// POST REQUEST
func (rt *_router) loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user1 := database.User{}
	err := json.NewDecoder(r.Body).Decode(&user1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := rt.db.GetUserByUsername(user1.Username)
	if err != nil {
		user, err = rt.db.AddUser(user1.Username)
	}
	finalize(user, err, w)
}

func (rt *_router) changeMyNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
	if flag {
		return
	}
	id, err := strconv.Atoi(ps.ByName("id"))
	if myID != id {
		w.WriteHeader(401)
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
	finalize(output, err, w)
}
func (rt *_router) uploadPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, userID := rt.youAreNotAuthorized(r, w)
	if flag {
		return
	}
	buf := make([]byte, r.ContentLength)
	if len(buf) == 0 {
		http.Error(w, "photo is empty", http.StatusBadRequest)
		return
	}
	// Legge l'intero contenuto del corpo della richiesta in una variabile di tipo []byte
	data, err := readAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("error reading photo")))
		return
	}
	imageUrl := "service/api/images/" + strconv.Itoa(userID) + "_" + strconv.Itoa(int(time.Now().Unix())) + ".jpg"
	f, err := os.Create(imageUrl)
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("error create file")))
		return
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("error writing to file")))
		return
	}
	output, err := rt.db.AddPhoto(userID, imageUrl)
	finalize(output, err, w)
}
func readAll(r io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (rt *_router) addCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, userID := rt.youAreNotAuthorized(r, w)
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
	finalize(output, err, w)
}

// PUT REQUEST
func (rt *_router) likePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(401)
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
	finalize(output, err, w)
}
func (rt *_router) followUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non puoi seguire utilizzando un altro utente")))
		return
	}
	if rt.youAreBanned(myID, followingID, r, w) {
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non puoi seguire un utente che ti ha bannato")))
		return
	}
	output, err := rt.db.AddFollow(followerID, followingID)
	finalize(output, err, w)
}

func (rt *_router) banUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(401)
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
	finalize(output, err, w)
}

// DELETE REQUEST
func (rt *_router) deleteUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non puoi cancellare un utente che non ti appartiene")))
		return
	}
	output, err := rt.db.DeleteUser(userID)
	finalize(output, err, w)
}

func (rt *_router) deleteCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(400)
		logerr(w.Write([]byte("comment not found")))
		return
	}
	if comment.UserID != myID {
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non puoi cancellare un commento che non ti appartiene")))
		return
	}
	output, err := rt.db.DeleteComment(commentID)
	finalize(output, err, w)
}
func (rt *_router) deletePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(400)
		logerr(w.Write([]byte("photo not found")))
		return
	}
	if photo.UserID != myID {
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non puoi cancellare una foto che non ti appartiene")))
		return
	}
	output, err := rt.db.DeletePhoto(photoID)
	finalize(output, err, w)
}
func (rt *_router) unlikePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non puoi levare un like che non hai messo tu")))
		return
	}
	output, err := rt.db.DeleteLike(photoID, userID)
	finalize(output, err, w)
}
func (rt *_router) unfollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
	if flag {
		return
	}
	followerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("follow id is empty")))
		return
	}
	if myID != followerID {
		w.WriteHeader(401)
		logerr(w.Write([]byte("Non puoi smettere di seguire da un account che non ti appartiene")))
		return
	}
	followingID, err := strconv.Atoi(ps.ByName("followId"))
	if err != nil {
		w.WriteHeader(400)
		logerr(w.Write([]byte("follow id is empty")))
		return
	}
	output, err := rt.db.DeleteFollow(followerID, followingID)
	finalize(output, err, w)
}

func (rt *_router) unbanUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flag, myID := rt.youAreNotAuthorized(r, w)
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
		w.WriteHeader(401)
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
	finalize(output, err, w)
}
