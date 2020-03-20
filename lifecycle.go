package authentication

import (
	"github.com/clayts/database"
)

//TODO try and stop adblocker preventing social login links (js onclick?)

func init() {
	initSession()
	initTemplates()
	initAuth()
	initRouter()
}

//Terminate must be called before the process exits
func Terminate() {
	database.Terminate()
}
