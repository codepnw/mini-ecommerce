package errs

import "errors"

// User
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserCredentials    = errors.New("invalid email or password")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenRevoked  = errors.New("token revoked")
	ErrTokenExpires  = errors.New("token expires")

	ErrUnauthorized  = errors.New("unauthorized")
	ErrNoPermissions = errors.New("no permissions")
)

// Product
var (
	ErrProductNotFound     = errors.New("product not found")
	ErrNoFieldsToUpdate    = errors.New("no fields to update")
	ErrProductStockInvalid = errors.New("product stock greater than zero")
	ErrProductPriceInvalid = errors.New("product price greater than zero")
	ErrProductSKUExists    = errors.New("sku already exists")
	ErrProductNotEnough    = errors.New("product not enough stock")
)

// Cart
var (
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrItemNotInCart   = errors.New("item not in cart")
	ErrCartIsEmpty     = errors.New("cart is empty")
)

// Order
var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrCannotCancelOrder   = errors.New("cannot cancel order")
	ErrInvalidStatusChange = errors.New("invalid status change")
)
