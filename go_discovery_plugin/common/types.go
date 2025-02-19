package common

// AuthPair represents SSH credentials with an ID.
type AuthPair struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
