package cmd

import (
	"github.com/rohankmr414/argodocs/logger"
	"github.com/rohankmr414/argodocs/markdown"
	"github.com/rohankmr414/argodocs/mdgen"
	"github.com/rohankmr414/argodocs/workflow"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputPrefix string
)

func NewGenerateCommand() *cobra.Command {
	// generateCmd represents the generate command
	var generateCmd = &cobra.Command{
		Use:   "generate PATH --output-prefix=PREFIX",
		Short: "Generate docs from workflow manifest.",
		Long:  `Generate reference docs from argo workflows.`,
		Run:   generate,
	}

	generateCmd.Flags().StringVar(
		&outputPrefix,
		"output-prefix",
		"",
		"Markdown output path prefix absolute or relative to the input YAML file",
	)

	return generateCmd
}

func generate(cmd *cobra.Command, args []string) {
	var LOGGER = logger.GetLogger("[Command] ")

	for _, arg := range args {
		LOGGER.Printf("Parsing file: %v", arg)
		parsedTemplateFiles, err := workflow.ParseFiles(arg)
		if err != nil {
			panic(err)
		}
		for _, parsedTemplateFile := range parsedTemplateFiles {
			var path string
			if outputPrefix == "" {
				yamlFileNameSplit := strings.Split(parsedTemplateFile.FilePath, "/")
				mdFileName := strings.Replace(yamlFileNameSplit[len(yamlFileNameSplit)-1], ".yaml", ".md", 1)
				mdFileName = strings.Replace(mdFileName, ".yml", ".md", 1)
				path = "./" + "docs" + "/" + mdFileName
			} else {
				if strings.HasPrefix(outputPrefix, ".") {
					yamlFileNameSplit := strings.Split(parsedTemplateFile.FilePath, "/")
					mdFileName := strings.Replace(yamlFileNameSplit[len(yamlFileNameSplit)-1], ".yaml", ".md", 1)
					mdFileName = strings.Replace(mdFileName, ".yml", ".md", 1)

					yamlPathSplit := strings.Split(parsedTemplateFile.FilePath, "/")
					mdFullPath := strings.Join(yamlPathSplit[:len(yamlPathSplit)-1], "/") + "/" + outputPrefix + "/" + mdFileName
					path = mdFullPath
				} else {
					yamlFileNameSplit := strings.Split(parsedTemplateFile.FilePath, "/")
					mdFileName := strings.Replace(yamlFileNameSplit[len(yamlFileNameSplit)-1], ".yaml", ".md", 1)
					mdFileName = strings.Replace(mdFileName, ".yml", ".md", 1)
					path = outputPrefix + "/" + mdFileName
				}
			}
			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				LOGGER.Panicln(err)
			}

			var doc *markdown.Doc
			doc, err = mdgen.GetMdDoc(parsedTemplateFile)
			LOGGER.Printf("Writing File: %v", path)
			err = doc.Export(path)
			if err != nil {
				LOGGER.Panicln(err)
			}
		}
	}
}
