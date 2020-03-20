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
	insist.IsNil(database.Transact(func(t database.Transaction) error {
		var auth string
		err := t.Read("session/authentication")(&auth)
		if err != nil {
			auth = hex.EncodeToString(securecookie.GenerateRandomKey(64))
			t.Write("session/authentication")(auth)
		}
		var enc string
		err = t.Read("session/encryption")(&enc)
		if err != nil {
			enc = hex.EncodeToString(securecookie.GenerateRandomKey(32))
			t.Write("session/encryption")(enc)
		}
		gothic.Store = sessions.NewCookieStore(
			insist.OnByteSlice(hex.DecodeString(auth)),
			insist.OnByteSlice(hex.DecodeString(enc)),
		)
		return nil
	}))
}