package authModels

type SessionsPrimaryKey struct {
	UserID    int64
	AuthToken string
}
