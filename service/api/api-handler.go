package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Routes
	rt.router.POST("/session", rt.doLogin)

	// Handle sending a message

	// Start the broadcaster goroutine

	rt.router.GET("/mainpage/:Id/users/:specificId", rt.wrap(wrapHandler(rt.getSpecificUser)))
	rt.router.GET("/mainpage/:Id/users", rt.wrap(wrapHandler(rt.getUsers)))
	rt.router.PUT("/mainpage/:Id/settings/name", rt.wrap(wrapHandler(rt.setMyUserName)))
	rt.router.PUT("/mainpage/:Id/settings/picture", rt.wrap(wrapHandler(rt.setMyPhoto)))

	rt.router.GET("/mainpage/:Id/conversations", rt.wrap(wrapHandler(rt.getMyConversations)))
	rt.router.POST("/mainpage/:Id/conversations", rt.wrap(wrapHandler(rt.startNewConversation)))

	rt.router.GET("/mainpage/:Id/conversations/:conversationId//:messageId/comments/:commentId", rt.wrap(wrapHandler(rt.getSpecificComment)))
	rt.router.DELETE("/mainpage/:Id/conversations/:conversationId//:messageId/comments/:commentId", rt.wrap(wrapHandler(rt.uncommentMessage)))

	rt.router.POST("/mainpage/:Id/conversations/:conversationId//:messageId/comments", rt.wrap(wrapHandler(rt.commentMessage)))
	rt.router.GET("/mainpage/:Id/conversations/:conversationId//:messageId/comments", rt.wrap(wrapHandler(rt.getComments)))

	rt.router.GET("/mainpage/:Id/conversations/:conversationId//:messageId", rt.wrap(wrapHandler(rt.getMessage)))
	rt.router.POST("/mainpage/:Id/conversations/:conversationId//:messageId", rt.wrap(wrapHandler(rt.forwardMessage)))
	rt.router.DELETE("/mainpage/:Id/conversations/:conversationId//:messageId", rt.wrap(wrapHandler(rt.deleteMessage)))

	rt.router.GET("/mainpage/:Id/conversations/:conversationId/groupMembers", rt.wrap(wrapHandler(rt.listGroupMembers)))
	rt.router.POST("/mainpage/:Id/conversations/:conversationId/groupMembers", rt.wrap(wrapHandler(rt.addToGroup)))
	rt.router.DELETE("/mainpage/:Id/conversations/:conversationId/groupMembers", rt.wrap(wrapHandler(rt.leaveGroup)))

	rt.router.PUT("/mainpage/:Id/conversations/:conversationId/conversationsettings/groupname", rt.wrap(wrapHandler(rt.setGroupName)))
	rt.router.PUT("/mainpage/:Id/conversations/:conversationId/conversationsettings/grouppicture", rt.wrap(wrapHandler(rt.setGroupPhoto)))

	rt.router.GET("/mainpage/:Id/conversations/:conversationId", rt.wrap(wrapHandler(rt.getConversation)))
	rt.router.POST("/mainpage/:Id/conversations/:conversationId", rt.wrap(wrapHandler(rt.sendMessage)))

	rt.router.GET("/mainpage/:Id", rt.wrap(wrapHandler(rt.getSelf)))

	return rt.router

}
