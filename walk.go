package chain

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	// m "github.com/therealfakemoot/gomarkov"
	"bufio"
)

// W builds a closure that fits the WalkFunc signature so you can recursively load corpus files.
func W(c *Chain) filepath.WalkFunc {
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
		// c.Build(f)
		s := bufio.NewScanner(f)
		for s.Scan() {
			c.Build(strings.NewReader(s.Text()))
		}
		if err := s.Err(); err != nil {
			return fmt.Errorf("error scanning corpus content: %w", err)
		}

		return nil
	}
	return wf
}
