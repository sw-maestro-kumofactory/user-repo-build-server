package dockerclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func BuildImage(dockerClient *client.Client, tarPath string, tags []string, dockerfilePath string) error {
	ctx := context.Background()

	tar, err := ioutil.ReadFile(tarPath)
	if err != nil {
		return err
	}

	buildOptions := types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: dockerfilePath,
	}
	buildCtx := &buildOptions

	resp, err := dockerClient.ImageBuild(ctx, strings.NewReader(string(tar)), *buildCtx)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}



//  curl -v -X POST -H "Content-Type:application/tar" --data-binary '@node-hello.tar.gz' --unix-socket /var/run/docker.sock http://v1.41/build?t=build_test\&dockerfile=node-hello/Dockerfile