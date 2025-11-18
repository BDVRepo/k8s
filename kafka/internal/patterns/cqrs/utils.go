package cqrs

import (
	"crypto/rand"
	"fmt"
	"time"
)

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getCurrentTime() time.Time {
	return time.Now()
}

