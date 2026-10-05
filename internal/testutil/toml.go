package testutil

import (
	"os"
)

func SaveTomlEmpty(file string) error {
	return os.WriteFile(file, []byte{}, 0644)
}
