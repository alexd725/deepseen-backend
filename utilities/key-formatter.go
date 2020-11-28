package utilities

// Create a key for Redis
func KeyFormatter(prefix string, value string) string {
	return prefix + "-" + value
}
