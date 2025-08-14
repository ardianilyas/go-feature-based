package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
)

var validate = validator.New()

func ValidateStruct(c *gin.Context, req interface{}) bool {
	// Parsing JSON ke struct
	if err := c.ShouldBindJSON(req); err != nil {
		// Cek kalau error validasi
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.JSON(422, gin.H{"errors": formatValidationError(verr)})
			return false
		}

		// Kalau error karena tipe data salah
		var ute *json.UnmarshalTypeError
		if errors.As(err, &ute) {
			c.JSON(422, gin.H{"errors": gin.H{
				strings.ToLower(ute.Field): fmt.Sprintf("Invalid type for field %s", ute.Field),
			}})
			return false
		}

		// Kalau body kosong
		if err.Error() == "EOF" {
			c.JSON(422, gin.H{"errors": gin.H{"body": "Request body is empty"}})
			return false
		}

		// Selain itu, berarti memang JSON-nya rusak
		c.JSON(422, gin.H{"errors": gin.H{"json": "Invalid JSON format"}})
		return false
	}

	// Validasi tambahan (jika ada)
	if err := validate.Struct(req); err != nil {
		c.JSON(422, gin.H{"errors": formatValidationError(err.(validator.ValidationErrors))})
		return false
	}

	return true
}

func formatValidationError(errs validator.ValidationErrors) map[string]string {
	errorsMap := make(map[string]string)
	for _, e := range errs {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errorsMap[field] = fmt.Sprintf("%s is required", field)
		case "email":
			errorsMap[field] = fmt.Sprintf("%s must be a valid email", field)
		case "min":
			errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
		default:
			errorsMap[field] = fmt.Sprintf("%s is not valid", field)
		}
	}
	return errorsMap
}
