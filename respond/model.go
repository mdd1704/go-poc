package respond

var (
	ErrUnknown    = "ErrUnknown"
	ErrInternal   = "ErrInternal"
	ErrBadRequest = "ErrBadRequest"
	ErrNotFound   = "ErrNotFound"
)

type ErrorAPIModel struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

type APIResponse struct {
	TransactionID string         `json:"transaction_id"`
	Success       bool           `json:"success"`
	Data          interface{}    `json:"data"`
	Error         *ErrorAPIModel `json:"error"`
}
