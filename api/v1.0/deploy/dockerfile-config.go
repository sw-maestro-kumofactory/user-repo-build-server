package deploy

import (
	"os"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	generator "github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator"
)

func injectEnvToDockerfile(dockerfilePath string, env []EnvInfo) *generator.Builder {
	dockerfile, _ := os.ReadFile(dockerfilePath)

	builder := generator.NewBuilder()

	result, _ := parser.Parse(strings.NewReader(string(dockerfile)))
	for _, node := range result.AST.Children {
		builder.AddCommand(node.Original)
		if node.Value == "FROM" {
			for _, env := range env {
				builder.AddEnv(env.Key, env.Value)
			}
		}
	}

	return builder
}
