package authentication

import (
	"errors"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/clayts/database"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	uuid "github.com/satori/go.uuid"
)

//ENV (optional): AUTENTICATION_{PROVIDER}_{ARGUMENT}

//PostRegistrationURL is the URL the user is redirected to after successfully logging in for the first time. The default value is "/".
var PostRegistrationURL = "/"

//PostLoginURL is the URL the user is redirected to after successful login. The default value is "/".
var PostLoginURL = "/"

//PostLogoutURL is the URL the user is redirected to after successful logout. The default value is "/".
var PostLogoutURL = "/"

//ErrorHandlerFunc is executed when an error occurs. By default, the function reports to the user that an error has occurred and provides them with a unique error code. Internally, details of the error are logged along with the unique code.
var ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
	errUUID := uuid.NewV4().String()
	log.Println("ERROR "+errUUID+":", err)
	http.Error(w, "internal server error "+errUUID, http.StatusInternalServerError)
}

func initAuth() {
	log.Println("initialising authentication at", baseURL+"/"+authPrefix+"/")
	getParams := func(provider string, keys ...string) []string {
		var parameters []string
		for _, key := range keys {
			if key == "{{URL}}" {
				parameters = append(parameters, baseURL+"/"+authPrefix+"/"+provider)
			} else {
				envKey := "AUTHENTICATION_" + strings.ToUpper(provider) + "_" + strings.ToUpper(key)
				value := os.Getenv(envKey)
				if value == "" {
					return nil
				}
				parameters = append(parameters, value)
			}
		}
		return parameters
	}
	var gothProviders []goth.Provider
	initProvider := func(p goth.Provider) {
		gothProviders = append(gothProviders, p)
	}

	//initialise providers
	if p := getParams("facebook", "key", "secret", "{{URL}}"); p != nil {
		initProvider(facebook.New(p[0], p[1], p[2]))
	}
	if p := getParams("google", "key", "secret", "{{URL}}"); p != nil {
		initProvider(google.New(p[0], p[1], p[2]))
	}
	if p := getParams("twitter", "key", "secret", "{{URL}}"); p != nil {
		initProvider(twitter.NewAuthenticate(p[0], p[1], p[2]))
	}
	if p := getParams("discord", "key", "secret", "{{URL}}"); p != nil {
		initProvider(discord.New(p[0], p[1], p[2]))
	}
	if p := getParams("twitch", "key", "secret", "{{URL}}"); p != nil {
		initProvider(twitch.New(p[0], p[1], p[2]))
	}

	goth.UseProviders(gothProviders...)
	log.Println("initialising authentication providers:", Providers())
}

//Providers returns a list of active providers
func Providers() []string {
	var providers []string
	for name := range goth.GetProviders() {
		providers = append(providers, name)
	}
	sort.Strings(providers)
	return providers
}

//IdentifyUser returns the user associated with the request (possibly anonymous).
func IdentifyUser(r *http.Request) User {
	id, err := gothic.GetFromSession("user", r)
	if err != nil {
		return User("")
	}
	return User(id)
}
func providerSelectionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if IdentifyUser(r).IsAnonymous() {
		providerSelectionTemplate(w, r, Providers())
		return
	}
	http.Redirect(w, r, PostLoginURL, http.StatusTemporaryRedirect)
	return
}
func loginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if !IdentifyUser(r).IsAnonymous() {
		http.Redirect(w, r, PostLoginURL, http.StatusTemporaryRedirect)
		return
	}
	profile, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		url, err := gothic.GetAuthURL(w, r)
		if handleError(w, r, err) {
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	if profile.Email == "" {
		handleError(w, r, errors.New("profile invalid"))
		return
	}
	if handleError(w, r, gothic.StoreInSession("user", profile.Email, r, w)) {
		return
	}
	var unregistered bool
	if handleError(w, r, database.Execute(func(t database.Transaction) error {
		u := User(profile.Email)
		if len(u.Profiles(t)) == 0 {
			unregistered = true
		}
		return u.updateProfile(t, profile)
	})) {
		return
	}
	if unregistered {
		http.Redirect(w, r, PostRegistrationURL, http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, PostLoginURL, http.StatusTemporaryRedirect)
	return
}
func logoutHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if u := IdentifyUser(r); !u.IsAnonymous() {
		if handleError(w, r, gothic.Logout(w, r)) {
			return
		}
	}
	http.Redirect(w, r, PostLogoutURL, http.StatusTemporaryRedirect)
}
