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
	rt.router.POST("/session", rt.doLogin) // doLogin OK

	// users
	rt.router.GET("/user", rt.getUsers)                   // getListOfUsers OK
	rt.router.GET("/user/:userID", rt.getUser)            // getUser OK
	rt.router.GET("/conversation", rt.getMyConversations) // getMyConversations OK
	rt.router.GET("/user/:userID/photo", rt.getUserPhoto) // getUserPhoto OK

	rt.router.GET("/me", rt.getMe)                  // getMe OK
	rt.router.PUT("/me/username", rt.setMyUserName) // SetMyUsername OK
	rt.router.PUT("/me/photo", rt.setMyPhoto)       // SetMyPhoto OK

	// conversation
	rt.router.POST("/conversation", rt.addConversation)                                      // addConversation OK
	rt.router.GET("/conversation/:conversationID", rt.getConversation)                       // getConversation OK
	rt.router.GET("/conversation/:conversationID/member", rt.getConversationMembers)         // getGroupMembers OK
	rt.router.POST("/conversation/:conversationID/member", rt.addToGroup)                    // addToGroup OK VA RIVISTO IN QUANTO SOTTO UN'ALTRO PATH
	rt.router.GET("/conversation/:conversationID/message", rt.getMessages)                   // getMessages OK
	rt.router.POST("/conversation/:conversationID/message", rt.sendMessage)                  // sendMessage OK
	rt.router.POST("/conversation/:conversationID/photo", rt.sendPhoto)                      // sendPhoto OK
	rt.router.GET("/conversation/:conversationID/photo/:messageID", rt.getConversationPhoto) // getConversationPhoto OK
	rt.router.GET("/conversation/:conversationID/lastmessage", rt.getLastMessage)            // getLastMessage OK
	rt.router.POST("/message/:messageID/forwarded", rt.forwardMessage)                       // forwardMessage OK
	rt.router.GET("/me/newmessage", rt.getNewMessages)                                       // getNewMessage OK
	rt.router.DELETE("/message/:messageID", rt.deleteMessage)                                // deleteMessage OK

	// reactions
	rt.router.POST("/message/:messageID/reaction", rt.commentMessage)                 // addReaction OK
	rt.router.GET("/message/:messageID/reaction", rt.getComments)                     // getReactions OK
	rt.router.DELETE("/message/:messageID/reaction/:reactionID", rt.uncommentMessage) // deleteReaction OK

	// groups
	rt.router.GET("/group/:groupID", rt.getGroup)             // getGroup OK
	rt.router.POST("/group", rt.createGroup)                  // createGroup OK
	rt.router.DELETE("/group/:groupID", rt.leaveGroup)        // leaveGroup OK
	rt.router.PUT("/group/:groupID/name", rt.setGroupName)    // setGroupName OK
	rt.router.GET("/group/:groupID/photo", rt.getGroupPhoto)  // getGroupPhoto OK
	rt.router.POST("/group/:groupID/photo", rt.setGroupPhoto) // setGroupPhoto OK

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
