package authentication

import (
	"log"

	"github.com/clayts/database"
)

//TODO try and stop adblocker preventing social login links (js onclick?)

func init() {
	initSession()
	initTemplates()
	initAuth()
	initRouter()
	log.Println("authentication initialisation complete")
}

//Terminate must be called before the process exits
func Terminate() {
	database.Terminate()
}
