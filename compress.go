package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CompressDirectoryToZip compresses the given directory recursively into a ZIP file.
func CompressDirectoryToZip(dirPath string, zipFilePath string) error {
	// Create a new ZIP file.
	f, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file: %w", err)
	}
	defer CloseFile(f)

	zw := zip.NewWriter(f)
	defer CloseZipWriter(zw)

	// Walk the directory tree, adding files and directories to the ZIP file.
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}
		// Calculate the relative path within the ZIP file.
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		if relPath == filepath.Base(zipFilePath) {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create file info header: %w", err)
		}
		header.Name = relPath

		// If it's a directory, set the appropriate flags.
		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store
		} else {
			// For regular files, use Deflate compression method.
			header.Method = zip.Deflate
		}

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create ZIP entry: %w", err)
		}

		if !info.IsDir() {
			// Open the file and copy its contents to the ZIP writer.
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer CloseFile(file)

			_, err = io.Copy(writer, file)
			if err != nil {
				return fmt.Errorf("failed to copy file contents: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}
	return nil
}

func CloseZipWriter(f *zip.Writer) {
	if f != nil {
		_ = f.Close()
	}
}
