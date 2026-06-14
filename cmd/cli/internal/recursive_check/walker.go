package recursivecheck

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func WalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".yaml" || ext == ".json" {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
