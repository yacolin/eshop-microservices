package errcode

import "fmt"

type BizError struct {
	Code    int
	Message string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// 常用业务错误
var (
	ErrProductNotFound = &BizError{Code: 1001, Message: "product not found"}
	ErrStockNotEnough  = &BizError{Code: 1002, Message: "stock not enough"}
	ErrInvalidParams   = &BizError{Code: 1003, Message: "invalid parameters"}
	ErrPaginationQuery = &BizError{Code: 1004, Message: "invalid pagination query"}
	ErrUnauthorized    = &BizError{Code: 1005, Message: "unauthorized"}
	ErrUserNotFound    = &BizError{Code: 1006, Message: "user not found"}
	ErrOrderNotFound   = &BizError{Code: 1007, Message: "order not found"}
	ErrDuplicateOrder  = &BizError{Code: 1008, Message: "duplicate order"}
	ErrPaymentFailed   = &BizError{Code: 1009, Message: "payment failed"}

	ErrInvalidCredentials     = &BizError{Code: 1010, Message: "invalid credentials"}
	ErrEmailAlreadyRegistered = &BizError{Code: 1011, Message: "email already registered"}
	ErrUserAlreadyRegistered  = &BizError{Code: 1012, Message: "user already registered"}

	ErrNotFound = &BizError{Code: 1013, Message: "resource not found"}
)
