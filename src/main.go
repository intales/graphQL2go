package main

import (
	"errors"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	schemaFlag := cli.StringFlag{
		Name:     "schema",
		Aliases:  []string{"s"},
		Usage:    "Load graphQL schema from `FILE`",
		Required: true,
	}
	jsonFlag := cli.BoolFlag{
		Name:    "json-tag",
		Aliases: []string{"jt"},
		Usage:   "Append json: \"field\" to every field in each type",
	}
	dynamodbFlag := cli.BoolFlag{
		Name:    "dynamodb-tag",
		Aliases: []string{"dt"},
		Usage:   "Append dynamodbav: \"field\" to every field in each type",
	}
	outputFileFlag := cli.StringFlag{
		Name:    "output-file",
		Aliases: []string{"of"},
		Usage:   "Write schema into `FILE`",
		// Value:   "schema.go",
	}
	outputDirFlag := cli.StringFlag{
		Name:    "output-dir",
		Aliases: []string{"od"},
		Usage:   "Write schema into multiple files into `DIR`",
		// Value:   "entities",
	}
	packageNameFlag := cli.StringFlag{
		Name:    "package-name",
		Aliases: []string{"pn"},
		Usage:   "Package name to use instead of 'main'",
		Value:   "main",
	}

	app := &cli.App{
		Name:                 "graphQL2go",
		EnableBashCompletion: true,
		Version:              "v0.1",
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			{Name: "riccardo pacioni"},
		},
		Flags: []cli.Flag{
			&schemaFlag,
			&jsonFlag,
			&dynamodbFlag,
			&outputFileFlag,
			&outputDirFlag,
			&packageNameFlag,
		},
		Action: func(cCtx *cli.Context) error {
			fileName := cCtx.String(schemaFlag.Name)

			outputFile := cCtx.String(outputFileFlag.Name)
			outputDir := cCtx.String(outputDirFlag.Name)

			packageName := cCtx.String(packageNameFlag.Name)

			if outputFile != "" && outputDir != "" {
				return cli.Exit("Cannot specify both output file and output dir", 1)
			}

			jsonTag := cCtx.String(jsonFlag.Name) == "true"
			dynamodbTag := cCtx.String(dynamodbFlag.Name) == "true"

			packageName = "package " + packageName + "\n\n"

			actualParser(fileName, outputFile, outputDir, packageName, jsonTag, dynamodbTag)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

type Entity struct {
	code string
	path string
}

func actualParser(inputFileName, outputFileName, outputDir, packageName string, jsonTag, dynamodbTag bool) {

	// Read the GraphQL file
	fileContent, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	// Interpret the GraphQL schema
	schemaTypes, schemaEnums, err := ParseGraphQLSchema(string(fileContent))
	if err != nil {
		fmt.Printf("Failed to interpret GraphQL schema: %v", err)
		os.Exit(1)
	}

	// Generate the Go struct definitions
	if outputDir != "" {
		createDir(outputDir)
		entities := make([]Entity, 0, len(schemaTypes)+len(schemaEnums))

		for _, graphQLType := range schemaTypes {
			header, body := graphQLType.GenerateGoStruct(jsonTag, dynamodbTag)
			entities = append(entities, Entity{
				code: packageName + header.String() + body.String(),
				path: path.Join(outputDir, graphQLType.Name+".go"),
			})
		}

		for _, graphQLEnum := range schemaEnums {
			body := graphQLEnum.GenerateEnumCode()
			entities = append(entities, Entity{
				code: packageName + body.String(),
				path: path.Join(outputDir, graphQLEnum.Name+".go"),
			})
		}

		for _, entity := range entities {
			content := formatCode(entity.code)
			writeCode(entity.path, content)
		}

		fmt.Println("Go files generation successful!")
	} else {
		// writing all output into on file
		goStructs := packageName + ""
		for _, graphQLType := range schemaTypes {
			header, body := graphQLType.GenerateGoStruct(jsonTag, dynamodbTag)
			goStructs += header.String() + body.String()
		}
		for _, graphQLEnum := range schemaEnums {
			body := graphQLEnum.GenerateEnumCode()
			goStructs += body.String()
		}
		content := formatCode(goStructs)
		writeCode(outputFileName, content)
		fmt.Println("Go file generation successful!")
	}
}

func formatCode(code string) []byte {
	content, err := format.Source([]byte(code))
	if err != nil {
		fmt.Printf("Code: %s\n", code)
		fmt.Printf("Failed to format Go file: %v", err)
		os.Exit(1)
	}
	return content
}

func writeCode(outputFileName string, content []byte) {
	err := os.WriteFile(outputFileName, content, 0644)
	if err != nil {
		fmt.Printf("Failed to write Go file into %s: %v", outputFileName, err)
		os.Exit(1)
	}
}

func createDir(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create folder: %v", err)
			os.Exit(1)
		}
	}
}
