package cmd

import (
	"github.com/rohankmr414/argodocs/logger"
	"github.com/rohankmr414/argodocs/markdown"
	"github.com/rohankmr414/argodocs/mdgen"
	"github.com/rohankmr414/argodocs/workflow"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var outputPrefix string

func NewGenerateCommand() *cobra.Command {
	// generateCmd represents the generate command
	generateCmd := &cobra.Command{
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
	LOGGER := logger.GetLogger("[Command] ")

	for _, arg := range args {
		LOGGER.Printf("Parsing file: %v", arg)
		parsedTemplateFiles, err := workflow.ParseFiles(arg)
		if err != nil {
			panic(err)
		}
		for _, parsedTemplateFile := range parsedTemplateFiles {
			var path string

			path = parsedTemplateFile.FilePath + "/docs/" + parsedTemplateFile.Name + "/docs/" + parsedTemplateFile.Name + ".md"
			if len(outputPrefix) > 0 {
				path = outputPrefix + "/docs/" + parsedTemplateFile.Name + "/docs/" + parsedTemplateFile.Name + ".md"
			}
			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				LOGGER.Panicln(err)
			}

			var doc *markdown.Doc
			doc, err = mdgen.GetMdDoc(parsedTemplateFile)
			if err != nil {
				LOGGER.Panicln(err)
			}
			LOGGER.Printf("Writing File: %v", path)
			err = doc.Export(path)
			if err != nil {
				LOGGER.Panicln(err)
			}
		}
	}
}
