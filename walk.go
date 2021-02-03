package chain

import (
	"log"
	"os"
	"path/filepath"
	// m "github.com/therealfakemoot/gomarkov"
)

// W builds a closure that fits the WalkFunc signature so you can recursively load corpus files.
func W(c *Chain) filepath.WalkFunc {
	wf := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			log.Printf("Unable to walk %s.\n", path)
			return err
		}

		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			log.Printf("Unable to open file %s: %s\n", path, err)
			return err
		}
		c.Build(f)

		return nil
	}
	return wf
}
