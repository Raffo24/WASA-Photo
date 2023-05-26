package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// POST REQUEST
	rt.router.POST("/session", rt.loginHandler)
	rt.router.POST("/photos", rt.uploadPhotoHandler)
	rt.router.POST("/photos/:id/comments", rt.addCommentHandler)
	// PUT REQUEST
	rt.router.PUT("/users/:id", rt.changeMyNameHandler)
	rt.router.PUT("/users/:id/follow/:followId", rt.followUserHandler)
	rt.router.PUT("/users/:id/ban/:banId", rt.banUserHandler)
	rt.router.PUT("/photos/:id/like/:userId", rt.likePhotoHandler)
	// GET REQUEST
	rt.router.GET("/feed", rt.getFeedHandler)
	rt.router.GET("/users/:id", rt.getUserHandler)
	rt.router.GET("/users/:id/photos", rt.getUserPhotosHandler)
	rt.router.GET("/users/:id/following", rt.getFollowingHandler)
	rt.router.GET("/users/:id/followers", rt.getFollowersHandler)
	rt.router.GET("/users", rt.searchUserHandler)
	rt.router.GET("/photos/:id", rt.getPhotoHandler)
	rt.router.GET("/photos/:id/comments", rt.getAllCommentsHandler)
	// DELETE REQUEST
	//rt.router.DELETE("/users/:id", rt.deleteUserHandler)
	rt.router.DELETE("/users/:id/follow/:followId", rt.unfollowUserHandler)
	rt.router.DELETE("/users/:id/ban/:banId", rt.unbanUserHandler)
	rt.router.DELETE("/photos/:id", rt.deletePhotoHandler)
	rt.router.DELETE("/photos/:id/like/:userId", rt.unlikePhotoHandler)
	rt.router.DELETE("/comments/:commentId", rt.deleteCommentHandler)
	// Special routes
	return rt.router
}
