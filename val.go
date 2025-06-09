package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/geico-private/gts/clients/go/tlm/lg"
)

// ValidationErrorHandler handles validation errors and makes them user-friendly
type ValidationErrorHandler struct {
	validator  *validator.Validate
	translator ut.Translator
}

// NewValidationErrorHandler creates a new validation error handler with translations
func NewValidationErrorHandler() (*ValidationErrorHandler, error) {
	// Create validator instance
	validate := validator.New()

	// Create English locale and universal translator
	en := en.New()
	uni := ut.New(en, en)

	// Get English translator
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, fmt.Errorf("failed to get English translator")
	}

	// Register default English translations
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, fmt.Errorf("failed to register default translations: %w", err)
	}

	// Register custom field name translations for better readability
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validation messages
	registerCustomTranslations(validate, trans)

	return &ValidationErrorHandler{
		validator:  validate,
		translator: trans,
	}, nil
}

// registerCustomTranslations adds custom, more user-friendly validation messages
func registerCustomTranslations(v *validator.Validate, trans ut.Translator) {
	// Custom translation for required fields
	v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", humanizeFieldName(fe.Field()))
		return t
	})

	// Custom translation for email validation
	v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email address", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", humanizeFieldName(fe.Field()))
		return t
	})

	// Custom translation for minimum length
	v.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} must be at least {1} characters long", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", humanizeFieldName(fe.Field()), fe.Param())
		return t
	})

	// Custom translation for maximum length
	v.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "{0} must be at most {1} characters long", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", humanizeFieldName(fe.Field()), fe.Param())
		return t
	})

	// Add more custom translations as needed for your specific validation rules
}

// humanizeFieldName converts field names to more readable format
// Example: "componentName" becomes "Component Name"
func humanizeFieldName(fieldName string) string {
	// Handle common field name patterns
	switch fieldName {
	case "componentName":
		return "Component Name"
	case "systemEntryId":
		return "System Entry ID"
	case "subEnvironment":
		return "Sub Environment"
	case "deploymentName":
		return "Deployment Name"
	case "correlationId":
		return "Correlation ID"
	default:
		// Generic conversion: add spaces before capital letters and capitalize first letter
		result := ""
		for i, r := range fieldName {
			if i > 0 && r >= 'A' && r <= 'Z' {
				result += " "
			}
			if i == 0 {
				result += strings.ToUpper(string(r))
			} else {
				result += string(r)
			}
		}
		return result
	}
}

// ValidateStruct validates a struct and returns user-friendly error messages
func (v *ValidationErrorHandler) ValidateStruct(s interface{}) []string {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	var errorMessages []string

	// Handle validation errors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			// Translate each field error to a user-friendly message
			translatedError := fieldError.Translate(v.translator)
			errorMessages = append(errorMessages, translatedError)
		}
	} else {
		// Fallback for non-validation errors
		errorMessages = append(errorMessages, "Invalid request format")
	}

	return errorMessages
}

// FormatValidationErrors creates a formatted error response for API consumers
func (v *ValidationErrorHandler) FormatValidationErrors(errors []string) map[string]interface{} {
	return map[string]interface{}{
		"error":   "Validation failed",
		"message": "Please correct the following errors and try again:",
		"details": errors,
	}
}
