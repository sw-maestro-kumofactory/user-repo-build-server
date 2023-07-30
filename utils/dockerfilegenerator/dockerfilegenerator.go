package dockerfilegenerator

import (
	"bytes"
	"fmt"
	"os"
)

type Builder struct {
	buffer bytes.Buffer
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) AddDirective(instruction, arguments string) {
	b.buffer.WriteString(fmt.Sprintf("%s %s\n", instruction, arguments))
}

func (b *Builder) AddCommand(command string) {
	b.buffer.WriteString(fmt.Sprintf("%s\n", command))
}

func (b *Builder) AddEnv(key, value string) {
	b.buffer.WriteString(fmt.Sprintf("ENV %s=%s\n", key, value))
}

func (b *Builder) Bytes() ([]byte, error) {
	return b.buffer.Bytes(), nil
}

func (b *Builder) AddDockerfile(dockerfile []byte) {
	b.buffer.Write(dockerfile)
}

func (b *Builder) CreateDockerfile(directory, filename string) error {
	content, err := b.Bytes()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", directory, filename)
	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
