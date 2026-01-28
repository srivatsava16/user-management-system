package shared

import (
	"errors"
	"strings"
)

var (
	ErrTokenExpired          = errors.New("TOKEN_EXPIRED")
	ErrTokenInvalid          = errors.New("TOKEN_INVALID")
	ErrTokenMalformed        = errors.New("TOKEN_MALFORMED")
	ErrTokenInvalidSignature = errors.New("TOKEN_INVALID_SIGNATURE")
)

func HandleDatabaseError(err error) string {
	if err == nil {
		return ""
	}

	errStr := err.Error()

	if strings.Contains(errStr, "Duplicate entry") {
		return parseDuplicateError(errStr)
	}

	if strings.Contains(errStr, "foreign key constraint") {
		return "Cannot delete this record as it is being used by other records."
	}

	return "An error occurred while processing your request."
}

func parseDuplicateError(errStr string) string {

	if strings.Contains(errStr, "idx_division_name_bu") {
		return "A division with this name already exists in this Business Unit."
	}

	if strings.Contains(errStr, "idx_role_name_div") {
		return "A Role with this name already exists in this Division."
	}

	if strings.Contains(errStr, "uni_APT_CUSTOM_ZXDS_BUSINESS_UNITS_name") {
		return "A Business Unit with this name already exists."
	}

	if strings.Contains(errStr, "uni_APT_CUSTOM_ZXDS_BUSINESS_UNITS_code") {
		return "A Business Unit with this code already exists."
	}

	if strings.Contains(errStr, "uni_APT_CUSTOM_ZXDS_USERS_email") {
		return "A User with this email already exists."
	}
	if strings.Contains(errStr, "uni_APT_CUSTOM_ZXDS_USERS_user_name") {
		return "A User with this username already exists."
	}

	return "This record already exists.Please use different values."

}
