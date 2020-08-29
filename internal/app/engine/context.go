package engine

type contextKey string

// CTX Keys for engine
const (
	CTXKeySessionState contextKey = "session"
	CTXKeyCookieState  contextKey = "cookie"

	// CTXKeyData is a context key for the accumulating
	// map[string]interface{} (engine.HTMLData) to pass to the
	// renderer
	CTXKeyData contextKey = "data"

	// CTXKeyValues is to pass the data submitted from API request or form
	// along in the context in case modules need it. The only module that needs
	// user information currently is remember so only auth/oauth2 are currently
	// going to use this.
	CTXKeyValues contextKey = "values"
)

func (c contextKey) String() string {
	return "engine ctx key " + string(c)
}
