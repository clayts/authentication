package authentication

import (
	"os"
)

var baseURL = os.Getenv("BASE_URL")
var port = os.Getenv("PORT")
var authPrefix = os.Getenv("AUTHENTICATION_PREFIX")
var sessionAuthenticationKey = os.Getenv("SESSION_AUTHENTICATION_KEY")
var sessionEncryptionKey = os.Getenv("SESSION_ENCRYPTION_KEY")
