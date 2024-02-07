package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyRecursive(src, dst string) error {
	return filepath.WalkDir(src, func(pathname string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if pathname, err = filepath.Rel(src, pathname); err != nil {
			return err
		}

		if de.IsDir() {
			return os.Mkdir(filepath.Join(dst, pathname), 0777)
		}

		r, err := os.Open(filepath.Join(src, pathname))
		if err != nil {
			return err
		}
		w, err := os.Create(filepath.Join(dst, pathname))
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
