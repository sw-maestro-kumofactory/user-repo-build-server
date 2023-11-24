package deploy

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"

	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator"
	rep "github.com/sw-maestro-kumofactory/miz-ball/utils/repomanagement"

	samplebuilder "github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator/sample-builder"
)

type Deployer struct {
	ctx *gin.Context
}

func NewDeployer(c *gin.Context) *Deployer {
	return &Deployer{ctx: c}
}

func (d *Deployer) ApplicationDeployAsync() {
	defer d.ctx.Request.Body.Close()
	defer d.sendEvent("finish", "finish")

	var info DeployInfo
	var dockerClient *client.Client

	var repoDir, tarPath string

	var err error

	err = d.ctx.ShouldBindJSON(&info)
	if d.handleError(err, http.StatusBadRequest) {
		return
	}
	dockerClient, err = initDockerClient()
	if d.handleError(err, http.StatusBadRequest) {
		return
	}

	fmt.Println(info.InstanceId)
	repoDir, err = createRepositoryDirectory(info.InstanceId)
	if d.handleError(err, http.StatusBadRequest) {
		return
	}
	// defer os.RemoveAll(repoDir)

	tarPath = filepath.Join(repoDir, "repo.tar.gz")
	err = cloneGitHubRepository(tarPath, info.User, info.Repo, info.Branch)
	if d.handleError(err, http.StatusBadRequest) {
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
	if d.handleError(err, http.StatusBadRequest) {
		return
	}

	targetTarPath := filepath.Join(repoDir, folderName+".tar.gz")

	d.sendEvent("message", "Build start")

	imageName, err := buildContainer(dockerClient, targetTarPath, info.InstanceId, "Dockerfile")
	if err != nil {
		d.sendEvent("fail", err.Error())
		return
	}
	d.sendEvent("message", "Build success")

	d.sendEvent("message", "Push start")
	err = pushOnECR(dockerClient, imageName)
	if d.handleError(err, http.StatusBadRequest) {
		return
	}
	d.sendEvent("message", "Push success")

	// TODO: save info to redis

	// if info.PortBind.Count > 0 {
	// 	saveInfo(info.InstanceId, info.PortBind)
	// }
	d.sendEvent("success", "success")
	dockerClient.ImageRemove(context.Background(), imageName, types.ImageRemoveOptions{})
}

func (d *Deployer) handleError(err error, statusCode int) bool {
	if err != nil {
		d.ctx.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		d.sendEvent("error", err.Error())
		return true
	}
	return false
}

func (d *Deployer) sendEvent(eventType string, message string) {

	fmt.Println("Event Type: ", eventType, ", Message: ", message)
	d.ctx.SSEvent(eventType, message)
	d.ctx.Writer.Flush()
}
