package dockerclient

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func BuildImage(dockerClient *client.Client, tarPath string, tags []string, dockerfilePath string) (string, error) {
	ctx := context.Background()

	buildCtx, err := os.Open(tarPath)
	if err != nil {
		return "", err
	}
	defer buildCtx.Close()

	buildOptions := types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: dockerfilePath,
	}

	resp, err := dockerClient.ImageBuild(ctx, buildCtx, buildOptions)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(body))

	body, _ := io.ReadAll(resp.Body)

	return string(body), nil
}

func BuildImage2(dockerClient *client.Client, contextPath string, tags []string) error {
	ctx := context.Background()

	buildCtx, err := os.Open(contextPath)
	if err != nil {
		return err
	}
	defer buildCtx.Close()

	buildOptions := types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: filepath.Join(contextPath, "Dockerfile"),
	}

	resp, err := dockerClient.ImageBuild(ctx, buildCtx, buildOptions)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
