package errors

import "errors"

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrDeliveryNotFound = errors.New("delivery not found")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrItemScanFailed   = errors.New("failed to scan order items")
	ErrTxBegin          = errors.New("failed to begin transaction")
	ErrTxCommit         = errors.New("failed to commit transaction")
	ErrInsertOrder      = errors.New("failed to insert into orders")
	ErrInsertDelivery   = errors.New("failed to insert into delivery")
	ErrInsertPayment    = errors.New("failed to insert into payment")
	ErrInsertItem       = errors.New("failed to insert into items")
	ErrInvalidJSON      = errors.New("invalid JSON format")
	ErrNilOrder         = errors.New("order is nil")
	ErrValidation       = errors.New("validation error")
	ErrCreateOrder      = errors.New("error creating order")
)
