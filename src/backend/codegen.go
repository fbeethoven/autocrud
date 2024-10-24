package backend

import (
	"log"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/fbeethoven/autocrud/src/codegen"
	"github.com/fbeethoven/autocrud/src/config"
	"github.com/fbeethoven/autocrud/src/database"
)

type MainData struct {
	Version     string
	ProjectName string
	Controllers []string
}

func GenerateMain(destPath, projName string, conf config.Config) error {
	file := getTemplateDir() + "/main.tmpl"

	t, err := template.New("main.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	if err := t.Execute(f, getMainData(conf, projName)); err != nil {
		return err
	}

	log.Printf("Generated file %s", destPath)

	return nil
}

type FieldData struct {
	Name    string
	Type    string
	Exclude bool
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

	t, err := template.New("resource.tmpl").Funcs(
		template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
		}).ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	if err := t.Execute(f, generateModelData(table)); err != nil {
		return err
	}

	return nil
}

func getMainData(conf config.Config, projName string) MainData {
	controllers := make([]string, 0, len(conf.Schema.Tables))
	for _, table := range conf.Schema.Tables {
		controllers = append(controllers, codegen.ToPascalCase(table.Name))
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
			Name:    field.Name,
			Type:    fieldType,
			Exclude: field.HasDefault || field.IsPrimaryKey,
		})
	}

	return ModelsData{
		Version:      config.Version,
		ImportTime:   importTime,
		ResourceName: codegen.ToPascalCase(table.Name),
		Fields:       fields,
	}
}

type DAOData struct {
	ProjectName  string
	Table        config.TableSchema
	DatabasePath string
}

type DAOTmplData struct {
	Version           string
	ProjectName       string
	Resource          string
	Fields            []config.FieldSchema
	TableName         string
	TableColumns      []string
	DatabasePath      string
	QueryResource     string
	QueryResourceById string
	TableIdField      string
	ImportTime        bool
	ImportStrconv     bool
}

func GenerateDAO(destPath string, daoData DAOData) error {
	file := getTemplateDir() + "/resourceDAO.tmpl"

	t, err := template.New("resourceDAO.tmpl").Funcs(
		template.FuncMap{
			"sub":          func(a, b int) int { return a - b },
			"toPascalCase": codegen.ToPascalCase,
		}).ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, generateDAOTmplData(daoData))
	if err != nil {
		return err
	}

	return nil
}

func generateDAOTmplData(daoData DAOData) DAOTmplData {
	columns := make([]string, 0, len(daoData.Table.Fields))
	importTime := false
	importStrconv := false
	for _, field := range daoData.Table.Fields {
		if field.IsPrimaryKey || field.HasDefault {
			continue
		}

		if field.Type == config.FieldTimestamp {
			importTime = true
		}

		if field.Type == config.FieldInt {
			importStrconv = true
		}

		columns = append(columns, field.Name)
	}

	return DAOTmplData{
		Version:           config.Version,
		ProjectName:       daoData.ProjectName,
		Resource:          codegen.ToPascalCase(daoData.Table.Name),
		Fields:            daoData.Table.Fields,
		TableName:         daoData.Table.Name,
		TableColumns:      columns,
		DatabasePath:      daoData.DatabasePath,
		QueryResource:     database.GetResourceQuery(daoData.Table),
		QueryResourceById: database.GetResourceByIdQuery(daoData.Table),
		TableIdField:      getTableIdField(daoData.Table),
		ImportTime:        importTime,
		ImportStrconv:     importStrconv,
	}
}

func getTableFields(table config.TableSchema) []string {
	fields := make([]string, 0, len(table.Fields))

	for _, field := range table.Fields {
		fields = append(fields, codegen.ToPascalCase(field.Name))
	}

	return fields
}

func getTableIdField(table config.TableSchema) string {
	idField := ""

	for _, field := range table.Fields {
		if field.IsPrimaryKey {
			idField = field.Name
		}
	}

	return idField
}

func GenerateControllerRouter(destPath, projName string) error {
	file := getTemplateDir() + "/controller.tmpl"

	t, err := template.New("controller.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

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
	Version      string
	ProjectName  string
	Resource     string
	ResourceUrl  string
	TableIdField string
}

func GenerateController(destPath, projName string, table config.TableSchema) error {
	file := getTemplateDir() + "/resourceController.tmpl"

	t, err := template.New("resourceController.tmpl").ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, ControllerData{
		Version:      config.Version,
		ProjectName:  projName,
		Resource:     codegen.ToPascalCase(table.Name),
		ResourceUrl:  table.Name,
		TableIdField: codegen.ToPascalCase(getTableIdField(table)),
	})
	if err != nil {
		return err
	}

	return nil
}
