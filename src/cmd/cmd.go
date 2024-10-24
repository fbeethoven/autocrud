package cmd

import (
	"flag"
)

type Args struct {
	ConfigName     string
	GenerateConfig bool
}

const configDefault = "config.yaml"

func GetArgs() Args {
	fileName := flag.String("f", configDefault, "Name of config file")
	fileNameLong := flag.String("file", configDefault, "Name of config file")

	generateConfig := flag.Bool("g", false, "Generate a default config file")
	generateConfigLong := flag.Bool("generate", false, "Generate a default config file")

	flag.Parse()

	args := Args{
		ConfigName:     *fileName,
		GenerateConfig: *generateConfig,
	}

	if *fileNameLong != configDefault {
		args.ConfigName = *fileNameLong
	}

	if *generateConfigLong {
		args.GenerateConfig = *generateConfigLong
	}

	return args
}
