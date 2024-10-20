package codegen

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type Generator interface {
	Generate()
}

type GenerateBuffer interface {
	CreateBuffer(destPath string) (io.Writer, error)
	Close()
}

type GenerateBufferImpl struct {
	f *os.File
}

var BufferGenerator GenerateBuffer = &GenerateBufferImpl{}

func (g *GenerateBufferImpl) CreateBuffer(destPath string) (io.Writer, error) {
	f, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}

	g.f = f

	return f, nil
}

func (g *GenerateBufferImpl) Close() {
	if g.f == nil {
		return
	}

	_ = g.f.Close()
	g.f = nil
}

func BeginTest(g GenerateBuffer) {
	BufferGenerator = g
}

var templateDir string = ""

func SetTemplateDir() {
	_, file, _, _ := runtime.Caller(1)
	dir := filepath.Dir(file)

	templateDir = filepath.Join(dir, "template")
}

func GetTemplateDir() string {
	if len(templateDir) == 0 {
		SetTemplateDir()
	}

	return templateDir
}

type GeneratorFunc func()

func GeneratorFactory(templatePath, destPath string) GeneratorFunc {
	return func() {
		file := fmt.Sprintf("%s/%s", GetTemplateDir(), templatePath)

		t, err := template.New(templatePath).ParseFiles(file)
		if err != nil {
			log.Printf("could not generate %s: %v", destPath, err)
			return
		}

		f, err := BufferGenerator.CreateBuffer(destPath)
		if err != nil {
			log.Printf("could not generate %s: %v", destPath, err)
			return
		}
		defer BufferGenerator.Close()

		err = t.Execute(f, nil)
		if err != nil {
			log.Printf("could not generate %s: %v", destPath, err)
			return
		}
	}
}

func toTitle(text string) string {
	if len(text) == 0 {
		return ""
	}
	return strings.ToUpper(string(text[0])) + text[1:]
}

func ToCamelCase(text string) string {
	words := make([]string, 0)

	for i, word := range strings.Split(text, "_") {
		if i == 0 {
			words = append(words, word)
		} else {
			words = append(words, toTitle(word))
		}
	}

	return strings.Join(words, "")
}

func ToPascalCase(text string) string {
	words := make([]string, 0)

	for _, word := range strings.Split(text, "_") {
		words = append(words, toTitle(word))
	}

	return strings.Join(words, "")
}
