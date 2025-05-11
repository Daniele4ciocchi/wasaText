package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	//session
	rt.router.POST("/session", rt.doLogin) // doLogin

	//users
	rt.router.GET("/user", rt.getUsers)                   // getListOfUsers
	rt.router.GET("/user/:userID", rt.getUser)            // getUser
	rt.router.GET("/conversation", rt.getMyConversations) // getMyConversations

	//conversation
	rt.router.POST("/conversation", rt.addConversation)                     // addConversation ?
	rt.router.GET("/conversation/:conversationID", rt.getConversation)      // getConversation
	rt.router.GET("/conversation/:conversationID/message", rt.getMessages)  // getMessages
	rt.router.POST("/conversation/:conversationID/message", rt.sendMessage) // sendMessage

	//rt.router.POST("/conversation/:c_id/message/:m_id", rt.forwardMessage) // forwardMessage ?

	//messages
	//rt.router.GET("/message/:m_id", rt.getMessage) // getMessage
	//rt.router.DELETE("/message/:m_id", rt.deleteMessage) // deleteMessage

	//comments
	//rt.router.POST("/message/:m_id/comment", rt.commentMessage) // commentMessage
	//rt.router.GET("/message/:m_id/comment", rt.getComments) // getComments
	//rt.router.DELETE("/message/:m_id/comment/:c_id", rt.deleteComment) // deleteComment
	//rt.router.GET("/comment/:c_id", rt.getComment) // getComment  DA FARE URGENTE

	//groups
	//rt.router.GET("/group/:g_id", rt.getGroup) // getGroup
	//rt.router.POST("/group", rt.createGroup) // createGroup
	//rt.router.GET("/group/:g_id/name", rt.getGroupName) // getGroupName
	//rt.router.POST("/group/:g_id/name", rt.setGroupName) // setGroupName
	//rt.router.GET("/group/:g_id/members", rt.getGroupMembers) // getGroupMembers
	//rt.router.GET("/group/:g_id/photo", rt.getGroupPhoto) // getGroupPhoto DA FARE URGENTE
	//rt.router.POST("/group/:g_id/photo", rt.setGroupPhoto) // setGroupPhoto
	//rt.router.POST("/group/:g_id/members", rt

	//manca il metodo per eliminarsi da un gruppo

	//rt.router.PUT("/me/username", rt.SetMyUsername) // SetMyUsername
	//rt.router.GET("/me/username", rt.GetMyUsername) // GetMyUsername
	//rt.router.POST("/me/photo", rt.SetMyPhoto) // SetMyPhoto

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
