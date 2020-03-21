package authentication

import (
	"log"
	"net/http"

	"github.com/gorilla/pat"
)

//Router is the router on which the various login endpoints are set up. You should add your own routes and then call http.ListenAndServe
var Router *pat.Router

func initRouter() {
	log.Println("initialising router")
	Router = pat.New()
	Router.Get("/"+authPrefix+"/login", providerSelectionHandlerFunc)
	Router.Get("/"+authPrefix+"/logout", logoutHandlerFunc)
	Router.Get("/"+authPrefix+"/{provider}", loginHandlerFunc)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {
		ErrorHandlerFunc(w, r, err)
		return true
	}
	return false
}
