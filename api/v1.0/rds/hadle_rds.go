package rds

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sw-maestro-kumofactory/miz-ball/utils/rdsutil"
)

func HandleRdsRequest(c *gin.Context) {
	rdsName := c.Param("rds-name")
	dbUsername := c.PostForm("dbUsername")
	dbPassword := c.PostForm("dbPassword")
	sqlFile, _, err := c.Request.FormFile("sqlFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SQL file not provided"})
		return
	}
	defer sqlFile.Close()

	tempSQLFile, err := saveTempSQLFile(sqlFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save SQL file"})
		return
	}
	defer os.Remove(tempSQLFile)

	rdsClient, err := rdsutil.NewRDSClient()
	if err != nil {
		log.Fatal("Failed to create RDS client:", err)
		return
	}

	endpoint, err := rdsClient.GetEndpointByDBName(rdsName)
	if err != nil {
		log.Fatal("Failed to get RDS endpoint:", err)
		return
	}

	rdsEngine, err := rdsClient.GetRdsTypeByDBName(rdsName)
	if err != nil {
		log.Fatal("Failed to get RDS engine:", err)
		return
	}

	dbClient, err := rdsutil.NewDBClient(endpoint, "3306", dbUsername, dbPassword)
	if err != nil {
		log.Fatal("Failed to create DB client:", err)
		return
	}

	err = dbClient.ExecuteSQL(tempSQLFile)
	if err != nil {
		log.Fatal("Failed to execute SQL:", err)
		return
	}

	fmt.Println("RDS Endpoint:", endpoint)
	fmt.Println("RDS Engine:", rdsEngine)
}

func saveTempSQLFile(file io.Reader) (string, error) {
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "temp-sql-*.sql")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
