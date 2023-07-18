package sample

import (
	"fmt"

	"github.com/gin-gonic/gin"

	repo "github.com/sw-maestro-kumofactory/miz-ball/utils/repomanagement"
)

func FindDockerfile(c *gin.Context) {
	filePath := "/app/repository/my_archive.tar.gz"

	dockerfilePath, err := repo.FindDockerfileInTar(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Dockerfile path:", dockerfilePath)
}
