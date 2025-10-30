package consts

import "time"

const (
	ContextTimeout = time.Second * 10

	AccessTokenDuration  = time.Hour
	RefreshTokenDuration = time.Hour * 24 * 7
)
