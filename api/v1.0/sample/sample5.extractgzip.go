package sample

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	repo "github.com/sw-maestro-kumofactory/miz-ball/utils/repomanagement"
)

func SAMPLE_EXTRACT(c *gin.Context) {
	r, err := os.Open("/app/repository/i-02b5064a1e36be086/repo.tar.gz")
	if err != nil {
		fmt.Println("error")
	}
	targetDir := "/app/repository/i-02b5064a1e36be086"
	repo.ExtractTarGz(r, targetDir)
}

func SAMPLE_ARCHIVE(c *gin.Context) {
	srcDir := "/app/repository/i-02b5064a1e36be086/coding-convention-sample-flask-f361fb2"
	dstDir := filepath.Dir(srcDir)
	err := repo.ArchiveToTarGz(srcDir, dstDir)
	if err != nil {
		fmt.Println(err.Error())
	}

}
