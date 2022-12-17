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

func finalize(output interface{}, err error, w http.ResponseWriter) {
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		if json.NewEncoder(w).Encode(output) != nil {
			w.WriteHeader(500)
			_, err = w.Write([]byte("Internal Server Error"))
			if err != nil {
				logrus.Error(err)
			}
		}
	} else {
		w.WriteHeader(400)
		// rt.baseLogger.WithError(err).Warning("getMyInfoHandler failed")
		// w.writeHeader(http.statusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			logrus.Error(err)
		}
	}
}

// GET REQUEST
func (rt *_router) getApiStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("L'api sta funzionando correttamente!"))
}

func (rt *_router) getUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
		return
	}
	user, err := rt.db.GetUserByID(userID)
	finalize(user, err, w)
}
func (rt *_router) getUserPhotosHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
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
	userID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Non sei loggato"))
		return
	}
	output, err := rt.db.GetFeed(userID)
	finalize(output, err, w)
}
func (rt *_router) getPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
		return
	}
	photo, err := rt.db.GetPhoto(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("photo not found in db"))
		return
	}
	// read photo from disk

	bytePhoto, err := os.ReadFile("service/api/images/" + photo.Photourl)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("photo not found on disk"))
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytePhoto)))
	w.Write(bytePhoto)
}

func (rt *_router) getCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	commentID, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("comment id is empty"))
		return
	}
	output, err := rt.db.GetCommentByID(commentID)
	finalize(output, err, w)
}
func (rt *_router) getAllCommentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("photo id is empty"))
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
	user := database.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err = rt.db.GetUserByUsername(user.Username)
	if err != nil {
		user, err = rt.db.AddUser(user.Username)
	}
	finalize(user, err, w)
}

func (rt *_router) changeMyNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
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
	userID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Non sei loggato"))
		return
	}
	bytefile := make([]byte, r.ContentLength)
	_, err = io.Copy(bytes.NewBuffer(bytefile), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(bytefile) == 0 {
		http.Error(w, "photo is empty", http.StatusBadRequest)
		return
	}
	// write photo to disk in the folder "photos"
	filename := strconv.Itoa(userID) + "_" + strconv.Itoa(int(time.Now().Unix()))
	err = os.WriteFile("service/api/images/"+filename, bytefile, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	output, err := rt.db.AddPhoto(userID, filename)
	finalize(output, err, w)
}

func (rt *_router) addCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID, err := strconv.Atoi(strings.Split(r.Header.Get("Authorization"), " ")[1])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
		return
	}
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
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
	userID, err := strconv.Atoi(ps.ByName("userId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
		return
	}
	photoID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
		return
	}

	output, err := rt.db.AddLike(photoID, userID)
	finalize(output, err, w)
}
func (rt *_router) followUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	followerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
		return
	}
	followingID, err := strconv.Atoi(ps.ByName("followId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("follow id is empty"))
		return
	}
	output, err := rt.db.AddFollow(followerID, followingID)
	finalize(output, err, w)
}

func (rt *_router) banUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bannerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is empty"))
		return
	}
	bannedID, err := strconv.Atoi(ps.ByName("banId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ban id is empty"))
		return
	}
	output, err := rt.db.AddBan(bannedID, bannerID)
	finalize(output, err, w)
}

// DELETE REQUEST
func (rt *_router) deleteCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	commentID, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("comment id is empty"))
		return
	}

	output, err := rt.db.DeleteComment(commentID)
	finalize(output, err, w)
}
func (rt *_router) deletePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	photoID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("photo id is empty"))
		return
	}
	output, err := rt.db.DeletePhoto(photoID)
	finalize(output, err, w)
}
func (rt *_router) unlikePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	photoID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("photo id is empty"))
		return
	}
	userID, err := strconv.Atoi(ps.ByName("userId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("user id is empty"))
		return
	}
	output, err := rt.db.DeleteLike(photoID, userID)
	finalize(output, err, w)
}
func (rt *_router) unfollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	followerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("follow id is empty"))
		return
	}
	followingID, err := strconv.Atoi(ps.ByName("followId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("follow id is empty"))
		return
	}
	output, err := rt.db.DeleteFollow(followerID, followingID)
	finalize(output, err, w)
}

// / riguardare idea
func (rt *_router) unbanUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bannerID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("banner id is empty"))
		return
	}

	bannedID, err := strconv.Atoi(ps.ByName("banId"))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("banned id is empty"))
		return
	}
	output, err := rt.db.DeleteBan(bannedID, bannerID)
	finalize(output, err, w)
}
