package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// session
	rt.router.POST("/session", rt.doLogin) // doLogin

	// users
	rt.router.GET("/user", rt.getUsers)                   // getListOfUsers
	rt.router.GET("/user/:userID", rt.getUser)            // getUser ?? da sistemare al posto dell'id si richiede una stringa
	rt.router.GET("/conversation", rt.getMyConversations) // getMyConversations
	rt.router.GET("/user/:userID/photo", rt.getUserPhoto) // getUserPhoto

	rt.router.GET("/me", rt.getMe)                  // getMe
	rt.router.PUT("/me/username", rt.setMyUsername) // SetMyUsername
	rt.router.PUT("/me/photo", rt.setMyPhoto)       // SetMyPhoto

	// conversation
	rt.router.POST("/conversation", rt.addConversation)                     // addConversation ?
	rt.router.GET("/conversation/:conversationID", rt.getConversation)      // getConversation
	rt.router.GET("/conversation/:conversationID/message", rt.getMessages)  // getMessages
	rt.router.POST("/conversation/:conversationID/message", rt.sendMessage) // sendMessage
	rt.router.GET("/conversation/:conversationID/lastmessage", rt.getLastMessage)
	rt.router.POST("/message/:messageID/forwarded", rt.forwardMessage) // forwardMessage ?
	rt.router.GET("/me/newmessage", rt.getNewMessages)                 // getNewMessage
	// rt.router.DELETE("/messages/:messageID", rt.deleteMessage) // deleteMessage DA IMPLEMENTARE

	// messages // ho creato sti path ma non so se servono
	// rt.router.GET("/message/:m_id", rt.getMessage) // getMessage
	// rt.router.DELETE("/message/:m_id", rt.deleteMessage) // deleteMessage ?

	// reactions
	rt.router.POST("/message/:messageID/reaction", rt.addReaction) // addReaction
	rt.router.GET("/message/:messageID/reaction", rt.getReactions) // getReactions
	rt.router.DELETE("/reaction/:reactionID", rt.deleteReaction)   // deleteReaction

	// groups
	rt.router.GET("/group/:groupID", rt.getGroup)                 // getGroup
	rt.router.POST("/group", rt.createGroup)                      // createGroup
	rt.router.DELETE("/group/:groupID", rt.leaveGroup)            // leaveGroup
	rt.router.GET("/group/:groupID/members", rt.getGroupMembers)  // getGroupMembers
	rt.router.POST("/group/:groupID/members", rt.addGroupMembers) // addGroupMembers
	rt.router.POST("/group/:groupID/name", rt.setGroupName)       // setGroupName
	rt.router.GET("/group/:groupID/photo", rt.getGroupPhoto)      // getGroupPhoto DA FARE URGENTE
	rt.router.POST("/group/:groupID/photo", rt.setGroupPhoto)     // setGroupPhoto

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
