package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OperationResult struct {
	OperationName       string
	SuccessCode         int
	ErrorCode           int
	SuccessMessage      string
	ErrorMessage        string
	VerboseErrorMessage string
	Error               error
	Data                interface{}
}

func HandleError(operation func() (interface{}, error), operationName string) *OperationResult {
	result := &OperationResult{
		OperationName: operationName,
		SuccessCode:   http.StatusOK,
		ErrorCode:     http.StatusInternalServerError,
	}

	data, err := operation()

	if err != nil {
		if statusCode, hasCode := GetStatusCode(err); hasCode {
			result.ErrorCode = statusCode
		}
		if message, hasMessage := GetErrorMessage(err); hasMessage {
			result.ErrorMessage = message
		}

		result.Error = err
		result.VerboseErrorMessage = fmt.Sprintf("Detailed error information: %v", err)
	} else {
		result.Data = data
		result.SuccessMessage = fmt.Sprintf("%s completed successfully", operationName)
	}

	return result
}

func HandleErrorWithStatusCode(operation func() (interface{}, error), operationName string, successCode int) *OperationResult {
	result := HandleError(operation, operationName)
	result.SuccessCode = successCode
	return result
}

func (r *OperationResult) RespondWithJSON(c *gin.Context) {
	if r.Error != nil {
		c.JSON(r.ErrorCode, gin.H{
			"success":        false,
			"error":          r.ErrorMessage,
			"verbose_error":  r.VerboseErrorMessage,
			"status_code":    r.ErrorCode,
		})
	} else {
		response := gin.H{
			"success":     true,
			"message":     r.SuccessMessage,
			"status_code": r.SuccessCode,
		}

		if r.Data != nil {
			response["data"] = r.Data
		}

		c.JSON(r.SuccessCode, response)
	}
}

func (r *OperationResult) IsSuccess() bool {
	return r.Error == nil
}

func (r *OperationResult) GetStatusCode() int {
	if r.Error != nil {
		return r.ErrorCode
	}
	return r.SuccessCode
}

func (r *OperationResult) GetMessage() string {
	if r.Error != nil {
		return r.ErrorMessage
	}
	return r.SuccessMessage
}