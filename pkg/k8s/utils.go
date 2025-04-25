package k8s

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetNameSuffix(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	nameSuffix := hex.EncodeToString(hasher.Sum(nil))
	return nameSuffix[:8]
}

func GetPortName(s string) string {
	var name string
	if len(s) >= 15 {
		name = s[:7]
		if name[len(name)-1] == '-' {
			return name[:len(name)-1]
		}
		return name
	}
	return s
}
