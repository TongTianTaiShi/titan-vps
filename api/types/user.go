package types

// UserInfo represents information about an user
type UserInfo struct {
	User  string `db:"user_addr"`
	Token string `db:"token"`
}
