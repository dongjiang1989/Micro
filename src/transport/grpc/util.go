package grpc

import (
	"encoding/base64"
	"strings"
)

// EncodeKeyValue sanitizes a key-value pair for use in gRPC metadata headers.
func EncodeKeyValue(key, val string) (string, string) {
	key = strings.ToLower(key)
	if strings.HasSuffix(key, binHdrSuffix) {
		v := base64.StdEncoding.EncodeToString([]byte(val))
		val = string(v)
	}
	return key, val
}
