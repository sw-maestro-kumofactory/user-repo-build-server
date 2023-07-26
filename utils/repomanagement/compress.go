package repomanagement

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ExtractTarGz(gzipStream io.Reader, targetDir string) {
	uncompressedStream, _ := gzip.NewReader(gzipStream)
	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		fullPath := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				log.Fatalf("ExtractTarGz: MkdirAll() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(fullPath)
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()
		default:
			log.Printf("ExtractTarGz: unknown type: %c in %s", header.Typeflag, header.Name)
		}
	}
}

func CompressToTarGz(srcDir, dstDir string) error {
	// 폴더 이름 설정
	fileName := filepath.Base(srcDir) + ".tar.gz"
	dstFilePath := filepath.Join(dstDir, fileName)

	// 파일 생성
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// gzip 생성
	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	// tar 생성
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// srcDir 내부 파일들을 가져와서 tar에 추가
	err = filepath.Walk(srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 디렉토리는 제외하고 파일만 처리
		if !info.Mode().IsRegular() {
			return nil
		}

		// 상대 경로 구하기
		relPath, err := filepath.Rel(srcDir, filePath)
		if err != nil {
			return err
		}

		// tar 헤더 작성
		header := &tar.Header{
			Name: relPath,
			Mode: int64(info.Mode()),
			Size: info.Size(),
		}

		// tar에 파일 추가
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// 파일 내용 추가
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
