package cmd

import (
	"fmt"
	"os"
)

func makeDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o700); err != nil { // If path is already a directory, MkdirAll does nothing
			return fmt.Errorf("can't make directory %s: %w", dir, err)
		}
	}
	return nil
}
