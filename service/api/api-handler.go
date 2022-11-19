package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getApiStatusHandler)
	rt.router.POST("/login", rt.loginHandler)
	rt.router.GET("/feed", rt.getFeedHandler)
	rt.router.GET("/users/{id}", rt.getUserInfoHandler)
	rt.router.GET("/users/me", rt.getMyInfoHandler)
	rt.router.POST("/users/me/change_name", rt.changeMyNameHandler)
	rt.router.PUT("/users/{id}/follow", rt.followUserHandler)
	rt.router.DELETE("/users/{id}/unfollow", rt.unfollowUserHandler)
	rt.router.PUT("/users/{id}/ban", rt.banUserHandler)
	rt.router.DELETE("/users/{id}/unban", rt.unbanUserHandler)
	rt.router.GET("/users/{id}/photos", rt.getUserPhotosHandler)
	rt.router.GET("/users/search/{search_username}", rt.searchUserHandler)
	rt.router.POST("/photos/upload", rt.uploadPhotoHandler)
	rt.router.GET("/photos/{id}", rt.getPhotoHandler)
	rt.router.DELETE("/photos/{id}", rt.deletePhotoHandler)
	rt.router.PUT("/photos/{id}/like", rt.likePhotoHandler)
	rt.router.DELETE("/photos/{id}/unlike", rt.unlikePhotoHandler)
	rt.router.GET("/photos/{id}/comments", rt.getAllCommentsHandler)
	rt.router.POST("/photos/{id}/comments/upload", rt.addCommentHandler)
	rt.router.GET("/photos/{photoId}/comments/{commentId}", rt.getCommentHandler)
	rt.router.DELETE("/photos/{photoId}/comments/{commentId}", rt.deleteCommentHandler)
	rt.router.PUT("/photos/{photoId}/comments/{commentId}/like", rt.likeCommentHandler)
	rt.router.DELETE("/photos/{photoId}/comments/{commentId}/unlike", rt.unlikeCommentHandler)
	// Special routes
	rt.router.GET("/liveness", rt.liveness)
	//context?
	//rt.router.GET("/context", rt.wrap(rt.getContextReply))
	return rt.router
}
