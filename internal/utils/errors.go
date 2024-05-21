package utils

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ApiResponse interface {
	StatusCode() int
}

func CheckCreateError(name string, response ApiResponse, err error) *diag.ErrorDiagnostic {
	if err != nil {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error creating %s", name),
			fmt.Sprintf("Could not create %s, unexpected error: %s", name, err.Error()))
		return &d
	}

	if response.StatusCode() != http.StatusOK {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error creating %s", name),
			fmt.Sprintf("Could not create %s, status code: %d (%s)",
				name, response.StatusCode(), readResponseBody(response)))
		return &d
	}

	return nil
}

func CheckGetError(name string, id string, response ApiResponse, err error) *diag.ErrorDiagnostic {
	if err != nil {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error retrieving %s with id %s", name, id),
			fmt.Sprintf("Could not retrieve %s with id %s, unexpected error: %s", name, id, err.Error()))
		return &d
	}

	if response.StatusCode() != http.StatusOK {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error retrieving %s with id %s", name, id),
			fmt.Sprintf("Could not retrieve %s with id %s, status code: %d (%s)",
				name, id, response.StatusCode(), readResponseBody(response)))
		return &d
	}

	return nil
}

func CheckUpdateError(name string, response ApiResponse, err error) *diag.ErrorDiagnostic {
	if err != nil {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error updating %s", name),
			fmt.Sprintf("Could not update %s, unexpected error: %s", name, err.Error()))
		return &d
	}

	if response.StatusCode() != http.StatusOK {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error updating %s", name),
			fmt.Sprintf("Could not update %s, status code: %d (%s)",
				name, response.StatusCode(), readResponseBody(response)))
		return &d
	}

	return nil
}

func CheckDeleteError(name string, response ApiResponse, err error) *diag.ErrorDiagnostic {
	if err != nil {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error deleting %s", name),
			fmt.Sprintf("Could not delete %s, unexpected error: %s", name, err.Error()))
		return &d
	}

	if response.StatusCode() != http.StatusOK {
		d := diag.NewErrorDiagnostic(
			fmt.Sprintf("Error deleting %s", name),
			fmt.Sprintf("Could not delete %s, status code: %d (%s)",
				name, response.StatusCode(), readResponseBody(response)))
		return &d
	}

	return nil
}

func readResponseBody(input ApiResponse) string {
	// Use reflection to get the field value
	ref := reflect.ValueOf(input)
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	// Check if the field exists and is readable
	value := ref.FieldByName("Body")
	if value.IsValid() && value.CanInterface() {
		if fieldValue, ok := value.Interface().([]byte); ok {
			return string(fieldValue)
		}
	}
	return "(no response body)"
}
