package ecr

import (
	"context"
	"fmt"
	"io/ioutil"

	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/sw-maestro-kumofactory/miz-ball/config"
)

func ecrAuthConfig() (string, error) {
	password, err := config.ReadECRPassword()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	var authConfig = types.AuthConfig{
		Username:      "AWS",
		Password:      password,
		ServerAddress: "434126037102.dkr.ecr.ap-northeast-2.amazonaws.com/kumo-customer",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	return authConfigEncoded, nil
}

func Push(dockerClient *client.Client, imageName string) error {
	authConfigEncoded, err := ecrAuthConfig()
	if err != nil {
		return err
	}

	pushOptions := types.ImagePushOptions{
		RegistryAuth: authConfigEncoded,
	}

	pushResponse, err := dockerClient.ImagePush(context.Background(), imageName, pushOptions)
	if err != nil {
		return err
	}
	defer pushResponse.Close()

	pushResult, err := ioutil.ReadAll(pushResponse)
	if err != nil {
		return err
	}

	fmt.Println(string(pushResult))

	return nil
}
