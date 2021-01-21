package util

import (
	"os"
	"strconv"
)

// GetenvOrDefault ...
func GetenvOrDefault(name, d string) string {
	value := os.Getenv(name)
	if value == "" {
		return d
	}
	return value
}

// GetChannelBuffer ...
func GetChannelBuffer() int {
	env := os.Getenv("CHANNEL_BUFFER_SIZE")
	size, err := strconv.Atoi(env)
	if err != nil {
		return ChannelBufferSize
	}
	return size
}
