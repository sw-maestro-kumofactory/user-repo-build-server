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

func ArchiveToTarGz(srcDir, destDir string) error {
	// 소스 디렉토리 열기
	source, err := os.Open(srcDir)
	if err != nil {
		return err
	}
	defer source.Close()

	// 목적지 디렉토리로 tar.gz 파일 생성
	fileName := filepath.Base(srcDir) + ".tar.gz"
	destFile, err := os.Create(filepath.Join(destDir, fileName))
	if err != nil {
		return err
	}
	defer destFile.Close()

	// gzip writer 생성
	gw := gzip.NewWriter(destFile)
	defer gw.Close()

	// tar writer 생성
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 디렉토리 내부 파일들을 순회하며 tar 아카이브 작성
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 현재 아이템의 상대 경로 얻기
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 디렉토리 자체는 무시
		if relPath == "." {
			return nil
		}

		// 파일 정보 구하기
		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}

		// header 쓰기
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// 파일 내용 쓰기
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tw, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
