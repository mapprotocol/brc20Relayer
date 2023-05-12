package utils

import (
	"encoding/json"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

func JSON(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) < 1
}

func IsDuplicateError(err string) bool {
	return strings.Contains(err, "Duplicate entry")
}

func Go(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("recover failed", "error", r)
			}
		}()

		fn()
	}()
}
