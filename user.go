package authentication

import (
	"errors"

	"github.com/clayts/database"
	"github.com/markbates/goth"
)

//User is a string consisting of some identifier unique amongst users (in practice, email address is used), or a blank string in the case of an anonymous user.
type User string

//IsAnonymous returns true if the user is anonymous
func (u User) IsAnonymous() bool { return u == "" }

//ID returns a string consisting of the user's unique identifier, prefixed by "user/".
func (u User) ID() string {
	if u.IsAnonymous() {
		return "user/anonymous"
	}
	return "user/" + string(u)
}

//Profiles returns the profiles a user has logged in with
func (u User) Profiles(t database.Transaction) []goth.User {
	if u.IsAnonymous() {
		return nil
	}
	var profileSlice []goth.User
	profiles := make(map[string]map[string]goth.User)
	err := t.Read(u.ID()+"/profiles", &profiles)
	if err != nil {
		return nil
	}
	for provider := range profiles {
		for account := range profiles[provider] {
			profileSlice = append(profileSlice, profiles[provider][account])
		}
	}
	return profileSlice
}

func (u User) updateProfile(t database.Transaction, gu goth.User) error {
	if u.IsAnonymous() {
		return errors.New("cannot assign profile to anonymous user")
	}
	profiles := make(map[string]map[string]goth.User)
	if err := t.Read(u.ID()+"/profiles", &profiles); err != nil && err != database.ErrNotFound {
		return err
	}
	if _, ok := profiles[gu.Provider]; !ok {
		profiles[gu.Provider] = make(map[string]goth.User)
	}
	profiles[gu.Provider][gu.UserID] = gu
	return t.Write(u.ID()+"/profiles", profiles)
}
