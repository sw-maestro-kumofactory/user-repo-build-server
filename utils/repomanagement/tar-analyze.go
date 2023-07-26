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

// TODO: Dockerfile 추출 & EXPOSE 구문
/*
func ExtractFileFromTarGz(tarGzFilePath, targetFilePath, destinationDirPath string) error {
	file, err := os.Open(tarGzFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	absTargetFilePath, err := filepath.Abs(targetFilePath)
	if err != nil {
		return err
	}

	absDestDirPath, err := filepath.Abs(destinationDirPath)
	if err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		if header.Name == absTargetFilePath {
			extractedFilePath := filepath.Join(absDestDirPath, filepath.Base(targetFilePath))

			extractedFile, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer extractedFile.Close()

			_, err = io.Copy(extractedFile, tarReader)
			if err != nil {
				return err
			}

			fmt.Printf("File %s extracted to %s\n", targetFilePath, extractedFilePath)
			return nil
		}
	}

	return fmt.Errorf("file not found")
}
*/
