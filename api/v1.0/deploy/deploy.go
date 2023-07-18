package deploy

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"

	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerclient"
	"github.com/sw-maestro-kumofactory/miz-ball/utils/ecr"
	rep "github.com/sw-maestro-kumofactory/miz-ball/utils/repomanagement"
)

type DeployInfo struct {
	InstanceId  string       `json:"target-instance" binding:"required"`
	GitHubToken string       `json:"github-token"`
	User        string       `json:"user" binding:"required"`
	Repo        string       `json:"repo" binding:"required"`
	Branch      string       `json:"branch" binding:"required"`
	Dockerfile  bool         `json:"Dockerfile" binding:"required"`
	PortBind    PortBindInfo `json:"portbind"`
	Runtime     string       `json:"Runtime"`
	Compiler    string       `json:"Compiler"`
}

type PortBindInfo struct {
	Count  int      `json:"count"`
	Expose []int    `json:"expose"`
	Bind   []string `json:"bind"`
}

var rootDir = "/app/repository/"

func ApplicationDeploy(c *gin.Context) {
	var info DeployInfo
	var dockerClient *client.Client

	var repoDir, tarPath, dockerfilePath string

	var err error

	err = c.ShouldBindJSON(&info)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}
	dockerClient, err = initDockerClient()
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	repoDir, err = createRepositoryDirectory(info.InstanceId)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	tarPath = filepath.Join(repoDir, "repo.tar.gz")
	err = cloneGitHubRepository(tarPath, info.User, info.Repo, info.Branch)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	if info.Dockerfile {
		dockerfilePath, err = rep.FindDockerfileInTar(tarPath)
		if handleError(c, err, http.StatusBadRequest) {
			return
		}

		imageName, err := buildContainer(dockerClient, tarPath, info.InstanceId, dockerfilePath)
		if handleError(c, err, http.StatusBadRequest) {
			return
		}

		err = pushOnECR(dockerClient, imageName)
		if handleError(c, err, http.StatusBadRequest) {
			return
		}

	} else {
		// TODO: Case when Dockerfile is not exist
		dockerfilePath = ""
	}
}

func handleError(c *gin.Context, err error, status int) bool {
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return true
	}
	return false
}

func createRepositoryDirectory(instanceID string) (string, error) {
	repoDir := filepath.Join(rootDir, instanceID)

	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return "", err
	}

	return repoDir, nil
}

func cloneGitHubRepository(filePath string, gh_user string, repo string, branch string) error {
	err := rep.RepoDownload(filePath, gh_user, repo, branch)
	return err
}

func initDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return cli, err
}

func buildContainer(cli *client.Client, tarPath string, instanceId string, dockerfilePath string) (string, error) {
	template := "434126037102.dkr.ecr.ap-northeast-2.amazonaws.com/kumo-customer:%s"
	imageName := fmt.Sprintf(template, instanceId)
	tags := []string{imageName}

	err := dockerclient.BuildImage(cli, tarPath, tags, dockerfilePath)
	if err != nil {
		return "", err
	}
	return imageName, nil
}

func pushOnECR(cli *client.Client, imageName string) error {
	return ecr.Push(cli, imageName)
}
