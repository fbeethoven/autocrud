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

	"autocrud/src/config"
	"autocrud/src/database"
)

type Generator interface {
	Generate()
}

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

	t, err := template.New("resource.tmpl").Funcs(
		template.FuncMap{
			"toPascalCase": toPascalCase,
		}).ParseFiles(file)
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

func toCamelCase(text string) string {
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

func toPascalCase(text string) string {
	words := make([]string, 0)

	for _, word := range strings.Split(text, "_") {
		words = append(words, toTitle(word))
	}

	return strings.Join(words, "")
}

func getMainData(conf config.Config, projName string) MainData {
	controllers := make([]string, 0, len(conf.Schema.Tables))
	for _, table := range conf.Schema.Tables {
		controllers = append(controllers, toPascalCase(table.Name))
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
			Name: field.Name,
			Type: fieldType,
		})
	}

	return ModelsData{
		Version:      config.Version,
		ImportTime:   importTime,
		ResourceName: toPascalCase(table.Name),
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
	Fields            []string
	TableName         string
	TableColumns      []string
	DatabasePath      string
	QueryResource     string
	QueryResourceById string
	TableIdField      string
}

func GenerateDAO(destPath string, daoData DAOData) error {
	file := getTemplateDir() + "/resourceDAO.tmpl"

	t, err := template.New("resourceDAO.tmpl").Funcs(
		template.FuncMap{
			"sub": func(a, b int) int { return a - b },
		}).ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f, generateDAOTmplData(daoData))
	if err != nil {
		return err
	}

	return nil
}

func generateDAOTmplData(daoData DAOData) DAOTmplData {
	columns := make([]string, 0, len(daoData.Table.Fields))
	for _, field := range daoData.Table.Fields {
		if field.IsPrimaryKey {
			continue
		}

		columns = append(columns, field.Name)
	}

	return DAOTmplData{
		Version:           config.Version,
		ProjectName:       daoData.ProjectName,
		Resource:          toPascalCase(daoData.Table.Name),
		Fields:            getTableFields(daoData.Table),
		TableName:         daoData.Table.Name,
		TableColumns:      columns,
		DatabasePath:      daoData.DatabasePath,
		QueryResource:     database.GetResourceQuery(daoData.Table),
		QueryResourceById: database.GetResourceByIdQuery(daoData.Table),
		TableIdField:      getTableIdField(daoData.Table),
	}
}

func getTableFields(table config.TableSchema) []string {
	fields := make([]string, 0, len(table.Fields))

	for _, field := range table.Fields {
		fields = append(fields, toPascalCase(field.Name))
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

	f, err := internalGenerateBuffer.CreateBuffer(destPath)
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f, ControllerData{
		Version:      config.Version,
		ProjectName:  projName,
		Resource:     toPascalCase(table.Name),
		ResourceUrl:  table.Name,
		TableIdField: toPascalCase(getTableIdField(table)),
	})
	if err != nil {
		return err
	}

	return nil
}

type GeneratorFunc func()

func GeneratorFactory(templatePath, destPath string) GeneratorFunc {
	return func() {
		file := fmt.Sprintf("%s/%s", getTemplateDir(), templatePath)

		t, err := template.New(templatePath).ParseFiles(file)
		if err != nil {
			log.Printf("could not generate %s: %v", destPath, err)
			return
		}

		f, err := internalGenerateBuffer.CreateBuffer(destPath)
		if err != nil {
			log.Printf("could not generate %s: %v", destPath, err)
			return
		}
		defer internalGenerateBuffer.Close()

		err = t.Execute(f, nil)
		if err != nil {
			log.Printf("could not generate %s: %v", destPath, err)
			return
		}
	}
}

func GenerateResources(destDir string, tables []config.TableSchema) error {
	t, err := template.New("resources.ts.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": toPascalCase,
		}).ParseFiles(getTemplateDir() + "/resources.ts.tmpl")
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destDir + "/resources.ts")
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f, tables)
	if err != nil {
		log.Printf("[CODEGEN] error with resources")
		return err
	}

	for _, table := range tables {
		log.Printf("[CODEGEN] generating table %v", table)

		t, err = template.New("types.ts.tmpl").
			Funcs(template.FuncMap{
				"toPascalCase": toPascalCase,
				"toCamelCase":  toCamelCase,
				"getType": func(inType string) string {
					return config.TypeMap[inType]
				},
			}).ParseFiles(getTemplateDir() + "/types.ts.tmpl")
		if err != nil {
			return err
		}

		f, err = internalGenerateBuffer.CreateBuffer(
			fmt.Sprintf(
				"%s/%s.ts",
				destDir,
				toPascalCase(table.Name),
			))
		if err != nil {
			return err
		}
		defer internalGenerateBuffer.Close()

		err = t.Execute(f, table)
		if err != nil {
			log.Printf("[CODEGEN] error with table %v", table)
			return err
		}
	}

	return nil
}

func GenerateResourceTables(destDir string, conf config.Config) error {
	tables := conf.Schema.Tables
	file := getTemplateDir() + "/columns.tsx.tmpl"

	t, err := template.New("columns.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": toPascalCase,
			"toCamelCase":  toCamelCase,
		}).ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := internalGenerateBuffer.CreateBuffer(destDir + "/columns.tsx")
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f, tables)
	if err != nil {
		return err
	}

	file = getTemplateDir() + "/Navbar.tsx.tmpl"

	t, err = template.New("Navbar.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": toPascalCase,
		}).ParseFiles(file)
	if err != nil {
		return err
	}

	f, err = internalGenerateBuffer.CreateBuffer(destDir + "/Navbar.tsx")
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f, map[string]any{
		"ProjectName": conf.Name,
		"Tables":      tables,
	})
	if err != nil {
		return err
	}

	Tablefile := getTemplateDir() + "/page.tsx.tmpl"
	t, err = template.New("page.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": toPascalCase,
			"toCamelCase":  toCamelCase,
		}).ParseFiles(Tablefile)
	if err != nil {
		return err
	}

	f, err = internalGenerateBuffer.CreateBuffer(destDir + "/page.tsx")
	if err != nil {
		return err
	}
	defer internalGenerateBuffer.Close()

	err = t.Execute(f, tables)
	if err != nil {
		return err
	}

	return nil
}
