package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	//rt.router.GET("/", rt.getHelloWorld)
	//rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	//rt.router.GET("/liveness", rt.liveness)

	// Routes
	rt.router.POST("/session", rt.doLogin)

	//rt.router.GET("/mainpage/:Id/users/specificUser", rt.getSpecificUser)
	rt.router.GET("/mainpage/:Id/users", rt.getUsers)
	rt.router.PUT("/mainpage/:Id/settings/name", rt.setMyUserName)
	rt.router.PUT("/mainpage/:Id/settings/picture", rt.setMyPhoto)

	rt.router.GET("/mainpage/:Id/conversations", rt.getMyConversations)
	rt.router.POST("/mainpage/:Id/conversations", rt.startNewConversation)

	rt.router.GET("/mainpage/:Id/conversations/:conversationId//:messageId/comments/:commentId", rt.getSpecificComment)
	rt.router.DELETE("/mainpage/:Id/conversations/:conversationId//:messageId/comments/:commentId", rt.uncommentMessage)

	rt.router.POST("/mainpage/:Id/conversations/:conversationId//:messageId/comments", rt.commentMessage)
	rt.router.GET("/mainpage/:Id/conversations/:conversationId//:messageId/comments", rt.getComments)

	rt.router.GET("/mainpage/:Id/conversations/:conversationId//:messageId", rt.getMessage)
	rt.router.POST("/mainpage/:Id/conversations/:conversationId//:messageId", rt.forwardMessage)
	rt.router.DELETE("/mainpage/:Id/conversations/:conversationId//:messageId", rt.deleteMessage)

	rt.router.GET("/mainpage/:Id/conversations/:conversationId/GroupMembers", rt.listGroupMembers)
	rt.router.POST("/mainpage/:Id/conversations/:conversationId/GroupMembers", rt.addToGroup)
	rt.router.DELETE("/mainpage/:Id/conversations/:conversationId/GroupMembers", rt.leaveGroup)

	rt.router.PUT("/mainpage/:Id/conversations/:conversationId/conversationsettings/groupname", rt.setGroupName)
	rt.router.PUT("/mainpage/:Id/conversations/:conversationId/conversationsettings/grouppicture", rt.setGroupPhoto)

	rt.router.GET("/mainpage/:Id/conversations/:conversationId", rt.getConversation)
	rt.router.POST("/mainpage/:Id/conversations/:conversationId", rt.sendMessage)

	rt.router.GET("/mainpage/:Id", rt.getSelf)

	return rt.router
}
