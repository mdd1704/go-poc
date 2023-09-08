package respond

import "github.com/gin-gonic/gin"

func Success(w *gin.Context, trx string, statusCode int, data interface{}) {
	w.JSON(
		statusCode,
		APIResponse{
			TransactionID: trx,
			Success:       true,
			Data:          data, Error: nil,
		},
	)
}

func Error(w *gin.Context, trx string, statusCode int, code, desc string) {
	w.JSON(
		statusCode,
		APIResponse{
			TransactionID: trx,
			Success:       false,
			Data:          nil,
			Error:         &ErrorAPIModel{code, desc},
		},
	)
}

func Invalid(w *gin.Context, trx string, statusCode int, outputs interface{}) {
	w.JSON(
		statusCode,
		APIResponse{
			TransactionID: trx,
			Success:       false,
			Data:          outputs,
			Error:         nil,
		},
	)
}
