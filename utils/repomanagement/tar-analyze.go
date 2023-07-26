package repomanagement

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FindDockerfileInTar(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	var dockerfilePath string

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		headerPath := strings.TrimPrefix(header.Name, "./")
		segments := strings.Split(headerPath, "/")

		if len(segments) <= 2 {
			if header.Typeflag == tar.TypeReg && filepath.Base(headerPath) == "Dockerfile" {
				dockerfilePath = headerPath
				break
			}
		}

	}

	if dockerfilePath == "" {
		return "", fmt.Errorf("dockerfile not found")
	}

	return dockerfilePath, nil
}

func GetFolderNameFromTar(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 파일의 끝에 도달하면 종료
		}
		if err != nil {
			return "", err
		}

		if header.Typeflag == tar.TypeDir {
			return filepath.Base(header.Name), nil
		}
	}

	return "", fmt.Errorf("Folder not found")
}
