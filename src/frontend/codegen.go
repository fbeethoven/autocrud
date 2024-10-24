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

func generateComponent(conf config.Config, destDir, fileName string) error {
	templatePath := fmt.Sprintf(
		"%s/%s.tsx.tmpl",
		codegen.GetTemplateDir(),
		fileName,
	)
	templateName := fmt.Sprintf("%s.tsx.tmpl", fileName)

	t, err := template.New(templateName).
		Funcs(template.FuncMap{
			"toPascalCase": codegen.ToPascalCase,
			"toCamelCase":  codegen.ToCamelCase,
		}).ParseFiles(templatePath)
	if err != nil {
		return err
	}

	f, err := codegen.BufferGenerator.CreateBuffer(
		fmt.Sprintf(
			"%s/%s.tsx",
			destDir,
			fileName,
		))
	if err != nil {
		return err
	}
	defer codegen.BufferGenerator.Close()

	err = t.Execute(f, map[string]any{
		"ProjectName": conf.Name,
		"Tables":      conf.Schema.Tables,
	})
	if err != nil {
		return err
	}

	return nil
}

func GenerateResourceTables(destDir string, conf config.Config) error {
	err := generateComponent(conf, destDir, "columns")
	if err != nil {
		return err
	}

	err = generateComponent(conf, destDir, "Navbar")
	if err != nil {
		return err
	}

	err = generateComponent(conf, destDir, "ResourceBar")
	if err != nil {
		return err
	}

	err = generateComponent(conf, destDir, "page")
	if err != nil {
		return err
	}

	err = generateComponent(conf, destDir, "DialogResource")
	if err != nil {
		return err
	}

	err = generateComponent(conf, destDir, "DisplayResource")
	if err != nil {
		return err
	}

	return nil
}
