package sample

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"

	repo "github.com/sw-maestro-kumofactory/miz-ball/utils/repomanagement"
)

func FindDockerfile(c *gin.Context) {
	filePath := "/app/repository/i-02b5064a1e36be086/repo.tar.gz"

	dockerfilePath, err := repo.FindDockerfileInTar(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Dockerfile path:", dockerfilePath)
	fmt.Println("Parent path:", filepath.Dir(dockerfilePath))
}
