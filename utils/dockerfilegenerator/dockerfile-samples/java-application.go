package dockerfilesamples

import (
	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator"
	dfenum "github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator/enums"
)

func JavaApplication(targetDir string) {
	builder := dockerfilegenerator.NewBuilder()

	// Build the first stage
	builder.AddDirective(dfenum.FROM, "gradle:jdk17-alpine AS builder")
	builder.AddDirective(dfenum.WORKDIR, "/src")
	builder.AddDirective(dfenum.COPY, ". .")
	builder.AddDirective(dfenum.RUN, "gradle build && \\\n    cd build/libs && \\\n    ls && \\\n    rm $(ls *plain.jar) && \\\n    mv $(ls *.jar) app.jar")

	// Build the second stage
	builder.AddDirective(dfenum.FROM, "openjdk:17-alpine")
	builder.AddDirective(dfenum.WORKDIR, "/app")
	builder.AddDirective(dfenum.COPY, "--from=builder /src/build/libs /app")
	builder.AddDirective(dfenum.EXPOSE, "8080")
	builder.AddDirective(dfenum.ENTRYPOINT, `["java","-jar","./app.jar"]`)

	err := builder.CreateDockerfile(targetDir, "Dockerfile")
	if err != nil {
		panic(err)
	}
}
