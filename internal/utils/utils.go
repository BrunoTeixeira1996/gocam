package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func SaveFFMPEGOutput(logLocation string, recordID string, output []byte) error {
	f := logLocation + recordID + ".log"

	file, err := os.Create(f)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to create file %s:%s\n", f, err)
	}

	defer file.Close()

	if _, err = file.Write(output); err != nil {
		return fmt.Errorf("[ERROR] Failed to write to file %s:%s", f, err)
	}

	return nil
}
