package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"

	conf "github.com/sw-maestro-kumofactory/miz-ball/config"

	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerclient"
	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator"
	"github.com/sw-maestro-kumofactory/miz-ball/utils/ecr"
	rep "github.com/sw-maestro-kumofactory/miz-ball/utils/repomanagement"

	samplebuilder "github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator/sample-builder"
)

// TODO: add validation
//	- check github repo
//	- check branch
//	- check dockerfile
//	- check portbind
//	- check language
//	- check runtime
//	- check compiler

// TODO: add error handling
//	- check github repo
//	- check branch
//	- check dockerfile
//	- check portbind
//	- check language
//	- check runtime
//	- check compiler

type DeployInfo struct {
	InstanceId  string       `json:"targetInstance" binding:"required"`
	GitHubToken string       `json:"githubToken"`
	User        string       `json:"user" binding:"required"`
	Repo        string       `json:"repo" binding:"required"`
	Branch      string       `json:"branch" binding:"required"`
	Dockerfile  bool         `json:"Dockerfile"`
	PortBind    PortBindInfo `json:"portbind"`
	Language    string       `json:"language"`
	Runtime     string       `json:"runtime"`
	Compiler    string       `json:"compiler"`
	Env         []EnvInfo    `json:"env"`
}

type PortBindInfo struct {
	Count  int      `json:"count"`
	Expose []int    `json:"expose"`
	Bind   []string `json:"bind"`
}

type EnvInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var rootDir = "/app/repository/"

func ApplicationDeploy2(c *gin.Context) {
	var info DeployInfo
	var dockerClient *client.Client

	var repoDir, tarPath string

	var err error

	err = c.ShouldBindJSON(&info)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}
	dockerClient, err = initDockerClient()
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	fmt.Println(info.InstanceId)
	repoDir, err = createRepositoryDirectory(info.InstanceId)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}
	// defer os.RemoveAll(repoDir)

	tarPath = filepath.Join(repoDir, "repo.tar.gz")
	err = cloneGitHubRepository(tarPath, info.User, info.Repo, info.Branch)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	r, err := os.Open(tarPath)
	if err != nil {
		fmt.Println("error")
	}
	rep.ExtractTarGz(r, repoDir)

	folderName, _ := rep.GetFolderNameFromTar(tarPath)
	fmt.Println(folderName)

	srcDir := filepath.Join(repoDir, folderName)
	dstDir := repoDir

	// TODO: wrap this code
	builder := dockerfilegenerator.NewBuilder()
	// if info.Env != nil {
	// 	for _, env := range info.Env {
	// 		builder.AddEnv(env.Key, env.Value)
	// 	}
	// }
	if !info.Dockerfile {
		if info.Language == "node" {
			samplebuilder.AddNodeBuilder(builder)
			// samplebuilder.NodeApplication(srcDir)

		} else if info.Language == "java" {
			samplebuilder.AddJavaBuilder(builder)
			// samplebuilder.JavaApplication(srcDir)
		}
	} else if info.Dockerfile {
		dockerfilePath := filepath.Join(srcDir, "Dockerfile")
		dockerfileStream, err := os.ReadFile(dockerfilePath)
		if err != nil {
			fmt.Println("error")
		}
		builder.AddDockerfile(dockerfileStream)
		os.Remove(dockerfilePath)
	}
	builder.CreateDockerfile(srcDir, "Dockerfile")

	if info.Env != nil {
		builder = injectEnvToDockerfile(filepath.Join(srcDir, "Dockerfile"), info.Env)
		os.Remove(filepath.Join(srcDir, "Dockerfile"))
		builder.CreateDockerfile(srcDir, "Dockerfile")
	}

	// until here

	err = rep.CompressToTarGz(srcDir, dstDir)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	targetTarPath := filepath.Join(repoDir, folderName+".tar.gz")

	imageName, err := buildContainer(dockerClient, targetTarPath, info.InstanceId, "Dockerfile")
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	err = pushOnECR(dockerClient, imageName)
	if handleError(c, err, http.StatusBadRequest) {
		return
	}

	// TODO: save info to redis

	// if info.PortBind.Count > 0 {
	// 	saveInfo(info.InstanceId, info.PortBind)
	// }
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

	buildLog, err := dockerclient.BuildImage(cli, tarPath, tags, dockerfilePath)
	if err != nil {
		return "", err
	}

	images, _ := cli.ImageList(context.Background(), types.ImageListOptions{})

	imageExists := false

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			if repoTag == imageName {
				imageExists = true
				break
			}
		}
		if imageExists {
			break
		}
	}

	if !imageExists {
		return "", fmt.Errorf(buildLog)
	}

	return imageName, nil
}

func pushOnECR(cli *client.Client, imageName string) error {
	return ecr.Push(cli, imageName)
}

func saveInfo(key string, value PortBindInfo) {
	redisCli, _ := conf.InitRedisClient()
	// redisCli.Set("hello", "world", 0).Err()
	jsonData, _ := json.Marshal(value)
	redisCli.Set(key, jsonData, 0).Err()
}
