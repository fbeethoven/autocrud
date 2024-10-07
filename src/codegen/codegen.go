package codegen

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"autocrud/src/config"
)

type MainData struct {
	Version     string
	ProjectName string
	Controllers []string
}

type GenerateBuffer interface {
	CreateBuffer(destPath string) (io.Writer, error)
	Close()
}

type GenerateBufferImpl struct {
	f *os.File
}

var internalGenerateBuffer GenerateBuffer = &GenerateBufferImpl{}

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

	return
}

func BeginTest(g GenerateBuffer) {
	internalGenerateBuffer = g
}

func GenerateMain(destPath, projName string, conf config.Config) error {
	file := getTemplateDir() + "/main.tmpl"

	t, err := template.New("main.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	if err := t.Execute(f, getMainData(conf, projName)); err != nil {
		return err
	}

	log.Printf("Generated file %s", destPath)

	return nil
}

type FieldData struct {
	Name string
	Type string
}

type ModelsData struct {
	Version      string
	ImportTime   bool
	ResourceName string
	Fields       []FieldData
}

var templateDir string = ""

func getTemplateDir() string {
	if len(templateDir) > 0 {
		return templateDir
	}

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	templateDir = filepath.Join(dir, "template")
	return templateDir
}

func GenerateModel(destPath string, table config.TableSchema) error {
	file := getTemplateDir() + "/resource.tmpl"

	t, err := template.New("resource.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	if err := t.Execute(f, generateModelData(table)); err != nil {
		return err
	}

	return nil
}

func toTitle(text string) string {
	if len(text) == 0 {
		return ""
	}
	return strings.ToUpper(string(text[0])) + text[1:]
}

func getMainData(conf config.Config, projName string) MainData {
	controllers := make([]string, 0, len(conf.Schema.Tables))
	for _, table := range conf.Schema.Tables {
		controllers = append(controllers, toTitle(table.Name))
	}

	return MainData{
		Version:     config.Version,
		ProjectName: projName,
		Controllers: controllers,
	}
}

func generateModelData(table config.TableSchema) ModelsData {
	fields := make([]FieldData, 0, len(table.Fields))

	importTime := false
	for _, field := range table.Fields {
		fieldType := field.Type
		if field.Type == config.FieldTimestamp {
			importTime = true
			fieldType = "time.Time"
		}

		fields = append(fields, FieldData{
			Name: toTitle(field.Name),
			Type: fieldType,
		})
	}

	return ModelsData{
		Version:      config.Version,
		ImportTime:   importTime,
		ResourceName: toTitle(table.Name),
		Fields:       fields,
	}
}

type DAOData struct {
	Version     string
	ProjectName string
	Resource    string
}

func GenerateDAO(destPath, projName string, table config.TableSchema) error {
	file := getTemplateDir() + "/resourceDAO.tmpl"

	t, err := template.New("resourceDAO.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	daoData := DAOData{
		Version:     config.Version,
		ProjectName: projName,
		Resource:    toTitle(table.Name),
	}

	if err := t.Execute(f, daoData); err != nil {
		return err
	}

	return nil
}

func GenerateControllerRouter(destPath, projName string) error {
	file := getTemplateDir() + "/controller.tmpl"

	t, err := template.New("controller.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f,
		struct {
			Version     string
			ProjectName string
		}{
			Version:     config.Version,
			ProjectName: projName,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

type ControllerData struct {
	Version     string
	ProjectName string
	Resource    string
	ResourceUrl string
}

func GenerateController(destPath, projName string, table config.TableSchema) error {
	file := getTemplateDir() + "/resourceController.tmpl"

	t, err := template.New("resourceController.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f, ControllerData{
		Version:     config.Version,
		ProjectName: projName,
		Resource:    toTitle(table.Name),
		ResourceUrl: table.Name,
	})
	if err != nil {
		return err
	}

	return nil
}
