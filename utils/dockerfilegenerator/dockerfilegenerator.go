package dockerfilegenerator

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func (b *Builder) Bytes() ([]byte, error) {
	return b.buffer.Bytes(), nil
}

func (b *Builder) CreateDockerfile(directory, filename string) error {
	content, err := b.Bytes()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", directory, filename)
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
