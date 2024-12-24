package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyRecursive(src, dst string) error {
	return filepath.WalkDir(src, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path, err = filepath.Rel(src, path); err != nil {
			return err
		}

		if de.IsDir() {
			return os.Mkdir(filepath.Join(dst, path), 0777)
		}

		r, err := os.Open(filepath.Join(src, path))
		if err != nil {
			return err
		}
		w, err := os.Create(filepath.Join(dst, path))
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, r); err != nil {
			return err
		}
		if err := r.Close(); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}

		return nil
	})
}
