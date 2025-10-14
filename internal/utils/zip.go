package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()
    for _, f := range r.File {
        fpath := filepath.Join(dest, f.Name)
        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }
        if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return err
        }
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()
        outFile, err := os.Create(fpath)
        if err != nil {
            return err
        }
        defer outFile.Close()
        _, err = io.Copy(outFile, rc)
        if err != nil {
            return err
        }
    }
    return nil
}