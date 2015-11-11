package youtubeapi

import (
	"errors"
)

// Enumeration of Errors
var (
	ErrorInvalidServiceAccount = errors.New("Invalid service account")
	ErrorInvalidClientSecrets  = errors.New("Invalid client secrets configuration")
	ErrorMissingContentOwner   = errors.New("Missing content owner")
	ErrorCacheTokenRead        = errors.New("Invalid Cache Token")
	ErrorCacheTokenWrite       = errors.New("Unable to create cache token")
	ErrorTokenExchange         = errors.New("Token Exchange Error")
	ErrorResponse              = errors.New("Bad Response")
	ErrorBadParameter          = errors.New("Invalid Parameter")
	ErrorMissingChannelFlag    = errors.New("Missing --channel flag")
	ErrorMissingBindFlags      = errors.New("Both --stream and --video flags required for bind operation")
	ErrorMissingVideoFlag      = errors.New("Missing --video flag")
)
