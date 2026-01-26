package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Routes
	rt.router.GET("/mainpage/:id/users", rt.getUsers)
	//rt.router.GET("/mainpage/:id/users", rt.getUsers)
	//rt.router.GET("/mainpage/:id/users", rt.getUsers)
	//rt.router.GET("/mainpage/:id/users", rt.getUsers)
	//rt.router.GET("/mainpage/:id/users", rt.getUsers)
	//rt.router.GET("/mainpage/:id/users", rt.getUsers)

	return rt.router
}
