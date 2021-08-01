package transforms

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	nameRegexp      = regexp.MustCompile(`^[\w-]{1,32}$`)
	snowflakeRegexp = regexp.MustCompile(`^[0-9]{1,}$`)
)

// ValidateSnowflake ensures the input is a snowflake ID.
func ValidateSnowflake(val interface{}, key string) (warns []string, errs []error) {
	value := val.(string)

	if !snowflakeRegexp.MatchString(value) {
		errs = append(errs, fmt.Errorf("%s is not a snowflake, got: %s", key, value))
	}

	return
}

// ValidateName ensures the input is lowercase, a-z, may include dashes, and a length of 1-32
func ValidateName(val interface{}, key string) (warns []string, errs []error) {
	value := val.(string)

	if !nameRegexp.MatchString(value) {
		errs = append(errs, fmt.Errorf("command name unacceptable: `%s`, refer to documentation: https://discord.com/developers/docs/interactions/slash-commands#create-global-application-command-json-params", value))
	}

	if value != strings.ToLower(value) {
		errs = append(errs, fmt.Errorf("command name not lower case: `%s`, refer to documentation: https://discord.com/developers/docs/interactions/slash-commands#create-global-application-command-json-params", value))
	}

	return
}

// ValidateDescription ensures the input is 1-100 characters.
// This may be used for more than descriptions, but is the primary use case. Option choice values also use this.
func ValidateDescription(val interface{}, key string) (warns []string, errs []error) {
	value := val.(string)
	length := len(value)

	if length < 1 || length > 100 {
		errs = append(errs, fmt.Errorf("command descriptions must be 1-100 characters, got: `%s` (length %d)", value, length))
	}

	return
}
