package samplebuilder

import (
	"github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator"
	dfenum "github.com/sw-maestro-kumofactory/miz-ball/utils/dockerfilegenerator/enums"
)

func NodeApplication(targetDir string) {
	builder := dockerfilegenerator.NewBuilder()

	builder.AddDirective(dfenum.FROM, "node:16-alpine AS builder")
	builder.AddDirective(dfenum.WORKDIR, "/app")
	builder.AddDirective(dfenum.COPY, ". .")
	builder.AddDirective(dfenum.RUN, "npm install")
	builder.AddDirective(dfenum.RUN, "npm run build")

	builder.AddDirective(dfenum.FROM, "node:16-alpine")
	builder.AddDirective(dfenum.WORKDIR, "/app")
	builder.AddDirective(dfenum.COPY, "--from=builder /app /app")

	builder.AddDirective(dfenum.EXPOSE, "3000")
	builder.AddDirective(dfenum.ENTRYPOINT, `["npm","run","start"]`)

	err := builder.CreateDockerfile(targetDir, "Dockerfile")
	if err != nil {
		panic(err)
	}
}

func AddNodeBuilder(builder *dockerfilegenerator.Builder) {

	builder.AddDirective(dfenum.FROM, "node:16-alpine AS builder")
	builder.AddDirective(dfenum.WORKDIR, "/app")
	builder.AddDirective(dfenum.COPY, ". .")
	builder.AddDirective(dfenum.RUN, "npm install")
	// builder.AddDirective(dfenum.RUN, "npm run build")

	builder.AddDirective(dfenum.FROM, "node:16-alpine")
	builder.AddDirective(dfenum.WORKDIR, "/app")
	builder.AddDirective(dfenum.COPY, "--from=builder /app /app")

	builder.AddDirective(dfenum.EXPOSE, "3000")
	builder.AddDirective(dfenum.ENTRYPOINT, `["npm","run","start"]`)
}
