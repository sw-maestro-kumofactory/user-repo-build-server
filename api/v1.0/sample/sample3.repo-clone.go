package sample

import (
	"fmt"

	"github.com/gin-gonic/gin"
	repo "github.com/sw-maestro-kumofactory/miz-ball/utils/repomanagement"
)

func SAMPLE_TEST_CLONE(c *gin.Context) {

	err := repo.RepoDownload("sample.tar.gz", "mkthebea", "CCKK_Run", "main")
	if err != nil {
		fmt.Println("파일 다운로드에 실패했습니다:", err)
	} else {
		fmt.Println("파일 다운로드가 완료되었습니다.")
	}
}
