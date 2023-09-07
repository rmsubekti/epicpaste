package helper

import "os"

func GetEnv(key, defaultVal string) string {
	v := os.Getenv(key)
	if len(v) > 0 {
		return v
	}
	return defaultVal
}
