package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fbeethoven/autocrud/src/codegen"
	"github.com/fbeethoven/autocrud/src/config"
	"github.com/fbeethoven/autocrud/src/database"
)

type FrontendGeneratorImpl struct {
	Config      config.Config
	Directories database.ProjectDirectories
}

func New(
	config config.Config,
	directories database.ProjectDirectories,
) FrontendGeneratorImpl {
	return FrontendGeneratorImpl{
		Config:      config,
		Directories: directories,
	}
}

func (f FrontendGeneratorImpl) Generate() {
	codegen.SetTemplateDir()

	err := config.MultiRunCmdInDir(
		f.Directories.Frontend,
		config.Command{
			Cmd: "npm",
			Args: []string{
				"create", "vite", ".", "--", "--template", "react-ts",
			},
		},
		config.Command{
			Cmd: "npm",
			Args: []string{
				"install", "-D", "tailwindcss@3", "postcss", "autoprefixer",
			},
		},
		config.Command{
			Cmd: "npx",
			Args: []string{
				"tailwindcss", "init", "-p",
			},
		},
		config.Command{
			Cmd:  "internal",
			Func: codegen.GeneratorFactory("tailwind.config.tmpl", "tailwind.config.js"),
		},
		config.Command{
			Cmd:  "internal",
			Func: codegen.GeneratorFactory("postcss.config.js.tmpl", "postcss.config.js"),
		},
		config.Command{
			Cmd:  "internal",
			Func: codegen.GeneratorFactory("index.css.tmpl", "src/index.css"),
		},
		config.Command{
			Cmd:  "internal",
			Func: addCompilerOptions,
		},
		config.Command{
			Cmd: "npm",
			Args: []string{
				"i", "-D", "@types/node",
			},
		},
		config.Command{
			Cmd:  "internal",
			Func: codegen.GeneratorFactory("vite.config.tmpl", "vite.config.ts"),
		},
		config.Command{
			Cmd: "npx",
			Args: []string{
				"shadcn@2.3.0", "init", "--defaults",
			},
		},
		config.Command{
			Cmd: "npm",
			Args: []string{
				"install",
				"react-hook-form", "@hookform/resolvers", "zod",
				"@tanstack/react-table", "lucide-react", "sonner",
			},
		},
		config.Command{
			Cmd: "npx",
			Args: []string{
				"shadcn@2.3.0", "add",
				"button", "input", "label", "select",
				"dialog", "table",
			},
		},
	)
	if err != nil {
		log.Printf("Error initiating frontend module: %v", err)
		return
	}

	config.MakeRelativeDir(
		fmt.Sprintf("%s/%s", f.Directories.Frontend, "src"),
		"components",
	)

	componentsDir := fmt.Sprintf(
		"%s/%s/%s",
		f.Directories.Frontend,
		"src",
		"components",
	)

	config.MakeRelativeDir(
		fmt.Sprintf("%s/%s", f.Directories.Frontend, "src"),
		"types",
	)

	typesDir := fmt.Sprintf(
		"%s/%s/%s",
		f.Directories.Frontend,
		"src",
		"types",
	)

	type FileInput struct {
		Tmpl string
		Dest string
	}

	files := []FileInput{
		{"data-table.tsx.tmpl", componentsDir + "/data-table.tsx"},
		{"theme-provider.tsx.tmpl", componentsDir + "/theme-provider.tsx"},
		{"App.tsx.tmpl", f.Directories.Frontend + "/src/App.tsx"},
	}

	for _, file := range files {
		codegen.GeneratorFactory(file.Tmpl, file.Dest)()
		log.Printf("[FRONTEND] generated file for %v", file)
	}

	err = GenerateResources(typesDir, f.Config.Schema.Tables)
	if err != nil {
		log.Printf("[FRONTEND] error %v", err)
	}

	err = GenerateResourceTables(componentsDir, f.Config)
	if err != nil {
		log.Printf("[FRONTEND] error %v", err)
	}
}

func addCompilerOptions() {
	tsConfigsFiles := []string{
		"tsconfig.json",
		"tsconfig.node.json",
	}

	for _, file := range tsConfigsFiles {
		tsConfigData, err := readJsonWithComments(file)
		if err != nil {
			log.Printf("Error could not read frontend configs: %v", err)
			return
		}

		log.Printf("unmarshalling the configs")

		var tsConfig map[string]interface{}
		err = json.Unmarshal(tsConfigData, &tsConfig)
		if err != nil {
			log.Printf("Error could not read frontend configs: %v", err)
			return
		}

		opt, ok := tsConfig["compilerOptions"].(map[string]interface{})
		if !ok {
			opt = make(map[string]interface{})
			tsConfig["compilerOptions"] = opt
		}

		opt["baseUrl"] = "."
		opt["paths"] = map[string]interface{}{
			"@/*": []string{"./src/*"},
		}

		tsConfigJson, err := json.MarshalIndent(tsConfig, "", "  ")
		if err != nil {
			log.Printf("Error could not read frontend configs: %v", err)
			return
		}

		err = os.WriteFile(file, tsConfigJson, 0755)
		if err != nil {
			log.Printf("Error could not read frontend configs: %v", err)
			return
		}
	}
}

func readJsonWithComments(filePath string) ([]byte, error) {
	tsConfigData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	result := make([][]byte, 0)

	for _, line := range bytes.Split(tsConfigData, []byte("\n")) {
		if bytes.Contains(line, []byte("/*")) {
			continue
		}

		result = append(result, line)
	}

	return bytes.Join(result, []byte("\n")), nil

}
