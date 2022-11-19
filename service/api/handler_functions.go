package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GET REQUEST
func (rt *_router) getApiStatusHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("L'api sta funzionando correttamente!"))
}

func (rt *_router) getUsersInfoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := rt.db.GetUsers()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) getMyInfoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := rt.db.getMyInfo(r, rt)
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
func (rt *_router) getUserInfoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := getUserInfo(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) getUserPhotosHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := getUserPhotos(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func (rt *_router) searchUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := searchUser(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) getFeedHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := getFeed(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) getPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := getPhoto(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) getCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := getComment(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) getAllCommentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := getAllComments(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
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
	if username == "" {
	}
	var userID int
	userID, err := rt.db.GetUserID(username)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		output := map[string]interface{}{"id": userID}
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func (rt *_router) changeMyNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := changeMyName(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) uploadPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := uploadPhoto(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func (rt *_router) addCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := addComment(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

// PUT REQUEST
func (rt *_router) likeCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := likeComment(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) likePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := likePhoto(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) followUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := followUser(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func (rt *_router) banUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := banUser(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

// DELETE REQUEST
func (rt *_router) deleteCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := deleteComment(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) deletePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := deletePhoto(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) unlikePhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := unlikePhoto(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) unfollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := unfollowUser(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

// / riguardare idea
func (rt *_router) unbanUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := unbanUser(r, rt)
	input := &unbanUserInput{}
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)

		json.NewDecoder(r.Body).Decode(&input)
		_ = r.Body.Close()
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
func (rt *_router) unlikeCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	output, err := unlikeComment(r, rt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	} else {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}
