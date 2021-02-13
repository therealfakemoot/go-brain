package brain

import (
	"log"
	"os"
	"path/filepath"
)

// W builds a closure that fits the WalkFunc signature so you can recursively load corpus files.
func W(b Brain) filepath.WalkFunc {
	wf := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Unable to walk %s.\n", path)
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			log.Printf("Unable to open file %s: %s\n", path, err)
			return err
		}
		b.Learn(f)

		return nil
	}
	return wf
}
