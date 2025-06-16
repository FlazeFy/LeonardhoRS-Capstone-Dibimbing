package utils

import (
	"fmt"
	"net/http"

	"pelita/config"

	"github.com/gin-gonic/gin"
)

func BuildResponseMessage(c *gin.Context, typeResponse, contextKey, method string, statusCode int, data, metadata interface{}) {
	wording := config.ResponseMessages[method]

	var message string
	if typeResponse == "success" {
		message = fmt.Sprintf("%s %s", contextKey, wording)
	} else {
		message = fmt.Sprintf("Failed to %s %s", contextKey, wording)
	}

	response := gin.H{
		"message": Capitalize(message),
		"status":  typeResponse,
	}

	if typeResponse == "success" && data != nil {
		response["data"] = data
	}

	if typeResponse == "success" && metadata != nil {
		response["metadata"] = metadata
	}

	c.JSON(statusCode, response)
}

func BuildErrorMessage(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": err,
		"status":  "failed",
	})
}
