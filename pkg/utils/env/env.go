package env

import (
	"os"
	"regexp"
)

var envPattern = regexp.MustCompile(`\$\{([^:}]+)(:([^}]+))?\}`)

func ExpandEnvWithDefault(s string) string {
	return envPattern.ReplaceAllStringFunc(s, func(sub string) string {
		matches := envPattern.FindStringSubmatch(sub)

		key := matches[1]
		def := matches[3]

		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		return def
	})
}
