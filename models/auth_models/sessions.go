package auth_models

type SessionsPrimaryKey struct {
	UserID    int64
	AuthToken string
}
