package consts

import "time"

const (
	ContextTimeout = time.Second * 10

	AccessTokenDuration  = time.Hour
	RefreshTokenDuration = time.Hour * 24 * 7
)

// Params Key
const (
	ParamProductID = "product_id"
	CartItemID     = "cart_item_id"
	ParamOrderID   = "order_id"
)

// Context Key
type contextKey string

const (
	UserClaimsKey contextKey = "user-claims-key"
	UserIDKey     contextKey = "user-id-key"
	SessionIDKey  contextKey = "session-id-key"
)
