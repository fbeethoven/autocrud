package frontend

import (
	"fmt"
	"log"
	"text/template"

	"autocrud/src/codegen"
	"autocrud/src/config"
)

func GenerateResources(destDir string, tables []config.TableSchema) error {
	t, err := template.New("resources.ts.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
		}).ParseFiles(codegen.GetTemplateDir() + "/resources.ts.tmpl")
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(destDir + "/resources.ts")
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, tables)
	if err != nil {
		log.Printf("[CODEGEN] error with resources")
		return err
	}

	for _, table := range tables {
		log.Printf("[CODEGEN] generating table %v", table)

		t, err = template.New("types.ts.tmpl").
			Funcs(template.FuncMap{
				"toPascalCase": codegen.ToPascalCase,
				"toCamelCase":  codegen.ToCamelCase,
				"getType": func(inType string) string {
					return config.TypeMap[inType]
				},
			}).ParseFiles(codegen.GetTemplateDir() + "/types.ts.tmpl")
		if err != nil {
			return err
		}

		f, err = codegen.BufferGenerator.CreateBuffer(
			fmt.Sprintf(
				"%s/%s.ts",
				destDir,
				codegen.ToPascalCase(table.Name),
			))
		if err != nil {
			return err
		}
		defer codegen.BufferGenerator.Close()

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
	file := codegen.GetTemplateDir() + "/columns.tsx.tmpl"

	t, err := template.New("columns.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
			"toCamelCase":  codegen.ToCamelCase,
		}).ParseFiles(file)
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(destDir + "/columns.tsx")
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, tables)
	if err != nil {
		return err
	}

	file = codegen.GetTemplateDir() + "/Navbar.tsx.tmpl"

	t, err = template.New("Navbar.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
		}).ParseFiles(file)
	if err != nil {
		return err
	}

	f, err = codegen.BufferGenerator.CreateBuffer(destDir + "/Navbar.tsx")
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, map[string]any{
		"ProjectName": conf.Name,
		"Tables":      tables,
	})
	if err != nil {
		return err
	}

	Tablefile := codegen.GetTemplateDir() + "/page.tsx.tmpl"
	t, err = template.New("page.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
			"toCamelCase":  codegen.ToCamelCase,
		}).ParseFiles(Tablefile)
	if err != nil {
		return err
	}

	f, err = codegen.BufferGenerator.CreateBuffer(destDir + "/page.tsx")
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, tables)
	if err != nil {
		return err
	}

	t, err = template.New("DialogResource.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
			"toCamelCase":  codegen.ToCamelCase,
		}).ParseFiles(codegen.GetTemplateDir() + "/DialogResource.tsx.tmpl")
	if err != nil {
		return err
	}

	f, err = codegen.BufferGenerator.CreateBuffer(destDir + "/DialogResource.tsx")
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, tables)
	if err != nil {
		return err
	}

	t, err = template.New("DisplayResource.tsx.tmpl").
		Funcs(template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
			"toCamelCase":  codegen.ToCamelCase,
		}).ParseFiles(codegen.GetTemplateDir() + "/DisplayResource.tsx.tmpl")
	if err != nil {
		return err
	}

	f, err = codegen.BufferGenerator.CreateBuffer(destDir + "/DisplayResource.tsx")
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, tables)
	if err != nil {
		return err
	}

	return nil
}
