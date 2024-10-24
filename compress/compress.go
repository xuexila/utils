package compress

import (
	"archive/zip"
	"fmt"
	"github.com/xuexila/utils/close/osClose"
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
	defer osClose.CloseFile(f)

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
			defer osClose.CloseFile(file)

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

// UnCompressZip 解压zip包
func UnCompressZip(zipFilePath string, targetDir string) error {
	reader, err := zip.OpenReader(zipFilePath)
	defer CloseZipReader(reader)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	// 遍历zip文件中的所有条目
	for _, file := range reader.File {
		err := func(file *zip.File) error {
			// 获取条目的相对路径
			filePath := filepath.Join(targetDir, file.Name)
			// 如果是目录，则创建它
			if file.FileInfo().IsDir() {
				if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
					return err
				}
				return nil
			}
			// 创建目标文件
			outputFile, err := os.Create(filePath)
			defer osClose.CloseFile(outputFile)
			if err != nil {
				return err
			}
			// 从zip文件中打开条目的读取流
			zipFile, err := file.Open()
			defer CloseIoReader(zipFile)
			if err != nil {
				return err
			}
			// 将条目内容复制到目标文件
			_, err = io.Copy(outputFile, zipFile)
			return err
		}(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseZipReader(f *zip.ReadCloser) {
	if f != nil {
		_ = f.Close()
	}
}

// CloseIoReader 关闭IoReader
func CloseIoReader(f io.ReadCloser) {
	if f != nil {
		_ = f.Close()
	}
}
