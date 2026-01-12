package export

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func AddToZip(srcFile, zipFile string) error {
	// Buka file sumber
	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Buat file zip
	zipF, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipF.Close()

	zipWriter := zip.NewWriter(zipF)
	defer zipWriter.Close()

	// Buat file di dalam zip
	zipEntry, err := zipWriter.Create(filepath.Base(srcFile))
	if err != nil {
		return err
	}

	// Copy isi file ke zip
	_, err = io.Copy(zipEntry, file)
	if err != nil {
		return err
	}

	// Tutup zip writer dulu sebelum delete file
	if err := zipWriter.Close(); err != nil {
		return err
	}

	// Hapus file original
	if err := os.Remove(srcFile); err != nil {
		return err
	}

	return nil
}
