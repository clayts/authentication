package authentication

import (
	"encoding/hex"

	"github.com/clayts/database"
	"github.com/clayts/insist"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func initSession() {
	if baseURL == "" {
		return
	}
	insist.IsNil(database.Execute(func(t database.Transaction) error {
		var auth string
		err := t.Read("session/authentication", &auth)
		if err != nil {
			auth = hex.EncodeToString(securecookie.GenerateRandomKey(64))
			insist.IsNil(t.Write("session/authentication", auth))
		}
		var enc string
		err = t.Read("session/encryption", &enc)
		if err != nil {
			enc = hex.EncodeToString(securecookie.GenerateRandomKey(32))
			insist.IsNil(t.Write("session/encryption", enc))
		}
		gothic.Store = sessions.NewCookieStore(
			insist.OnByteSlice(hex.DecodeString(auth)),
			insist.OnByteSlice(hex.DecodeString(enc)),
		)
		return nil
	}))
}
