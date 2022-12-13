package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func finalize(output interface{}, err error, w http.ResponseWriter) {
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		//rt.baseLogger.WithError(err).Warning("getMyInfoHandler failed")
		//w.writeHeader(http.statusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

// GET REQUEST
func (rt *_router) getApiStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("L'api sta funzionando correttamente!"))
}

func (rt *_router) getLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := r.FormValue("username")
	if name == "" {
		w.WriteHeader(400)
		w.Write([]byte("username is empty"))
		return
	}
	user, err := rt.db.AddUser(name)
	finalize(user, err, w)
}

func (rt *_router) getUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//read from body request userID
	userID, err := strconv.Atoi(ps.ByName("id"))
	user, err := rt.db.GetUserByID(userID)
	finalize(user, err, w)
}
func (rt *_router) getUserPhotosHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("id"))
	output, err := rt.db.GetPhotos(userID)
	finalize(output, err, w)
}

func (rt *_router) searchUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username_searched := ps.ByName("id")
	output, err := rt.db.SearchUser(username_searched)
	finalize(output, err, w)
}
func (rt *_router) getFeedHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	output, err := rt.db.GetFeed(userID)
	finalize(output, err, w)
}
func (rt *_router) getPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	photo, err := rt.db.GetPhoto(id)
	finalize(photo, err, w)
}
func (rt *_router) getCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	commentID, err := strconv.Atoi(ps.ByName("commentId"))
	output, err := rt.db.GetCommentByID(commentID)
	finalize(output, err, w)
}
func (rt *_router) getAllCommentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	id, err := strconv.Atoi(ps.ByName("id"))
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
	username := r.FormValue("username")
	user, err := rt.db.GetUserByUsername(username)
	if err != nil {
		user, err = rt.db.AddUser(username)
	}
	finalize(user, err, w)
}

func (rt *_router) changeMyNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	fmt.Print(id)
	username := r.FormValue("username")
	output, err := rt.db.UpdateUser(id, username)
	finalize(output, err, w)
}
func (rt *_router) uploadPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	bytefile, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if bytefile == nil || len(bytefile) == 0 { // check if file is empty
		http.Error(w, "photo is empty", http.StatusBadRequest)
		return
	}
	output, err := rt.db.AddPhoto(userID, bytefile)
	finalize(output, err, w)
}

func (rt *_router) addCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	id, err := strconv.Atoi(ps.ByName("id"))
	content := r.FormValue("content")
	fmt.Print(id, userID, content+"\n")
	output, err := rt.db.AddComment(id, userID, content)
	finalize(output, err, w)
}

// PUT REQUEST
func (rt *_router) likePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("userId"))
	photoID, err := strconv.Atoi(ps.ByName("id"))
	output, err := rt.db.AddLike(photoID, userID)
	finalize(output, err, w)
}
func (rt *_router) followUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	followerID, err := strconv.Atoi(ps.ByName("id"))
	followingID, err := strconv.Atoi(ps.ByName("followId"))
	output, err := rt.db.AddFollow(followerID, followingID)
	finalize(output, err, w)
}

func (rt *_router) banUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bannerID, err := strconv.Atoi(ps.ByName("id"))
	bannedID, err := strconv.Atoi(ps.ByName("banId"))
	output, err := rt.db.AddBan(bannedID, bannerID)
	finalize(output, err, w)
}

// DELETE REQUEST
func (rt *_router) deleteCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	commentID, err := strconv.Atoi(ps.ByName("commentId"))
	output, err := rt.db.DeleteComment(commentID)
	finalize(output, err, w)
}
func (rt *_router) deletePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	photoID, err := strconv.Atoi(ps.ByName("id"))
	output, err := rt.db.DeletePhoto(photoID)
	finalize(output, err, w)
}
func (rt *_router) unlikePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	photoID, err := strconv.Atoi(ps.ByName("id"))
	userID, err := strconv.Atoi(ps.ByName("userId"))
	output, err := rt.db.DeleteLike(photoID, userID)
	finalize(output, err, w)
}
func (rt *_router) unfollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	followerID, err := strconv.Atoi(ps.ByName("id"))
	followingID, err := strconv.Atoi(ps.ByName("followId"))
	output, err := rt.db.DeleteFollow(followerID, followingID)
	finalize(output, err, w)
}

// / riguardare idea
func (rt *_router) unbanUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bannerID, err := strconv.Atoi(ps.ByName("id"))
	bannedID, err := strconv.Atoi(ps.ByName("banId"))
	output, err := rt.db.DeleteBan(bannedID, bannerID)
	finalize(output, err, w)
}
