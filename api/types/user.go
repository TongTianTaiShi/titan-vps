package types

// UserInfo represents information about an user
type UserInfo struct {
	UserID        string `db:"user_id"`
	Balance       string `db:"balance"`
	LockedBalance string
}
