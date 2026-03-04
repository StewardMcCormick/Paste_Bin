package validation

import (
	"strings"
)

func ValidateAPIKey(key string) bool {
	splitKey := strings.Split(key, "_")

	return len(splitKey) == 4 && splitKey[0] == "pb" &&
		(splitKey[1] == "test" || splitKey[1] == "live") &&
		len(splitKey[2]) == 4 && len(splitKey[3]) == 12
}
