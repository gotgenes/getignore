package getignore

import "fmt"

// Version is the version of getignore.
// It should be populated through ldflags, e.g.,
// -ldflags "-X constants/constants.Version=${VERSION}"
var Version = "development"

// UserAgentString identifies getignore as the user agent for outbound HTTP requests
var UserAgentString = fmt.Sprintf("getignore/%s", Version)
