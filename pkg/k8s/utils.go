package k8s

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetNameSuffix(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	nameSuffix := hex.EncodeToString(hasher.Sum(nil))
	return nameSuffix[:7]
}
