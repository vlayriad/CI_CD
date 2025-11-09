package helper

import (
	"fmt"
	"kaffein/config"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// getPath menelusuri path misal "postgres.host" dan mengembalikan nilainya
func getPath(path string) (any, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, fmt.Errorf("config not initialized — call SetupApp() first")
	}

	val := reflect.ValueOf(*cfg)
	parts := strings.Split(path, ".")

	for _, part := range parts {
		val = reflect.Indirect(val)
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("invalid path: %s", path)
		}

		found := false
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == part {
				val = val.Field(i)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("invalid key: %s", part)
		}
	}

	return val.Interface(), nil
}

func GetString(path string) string {
	v, err := getPath(path)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("GetString: failed to get '%s'", path)
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func GetInt(path string) int {
	v, err := getPath(path)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("GetInt: failed to get '%s'", path)
		return 0
	}

	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}

func GetBool(path string) bool {
	v, err := getPath(path)
	if err != nil {
		log.Warn().
			Err(err).
			Msgf("GetBool: failed to get '%s'", path)
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		b, _ := strconv.ParseBool(val)
		return b
	default:
		return false
	}
}
