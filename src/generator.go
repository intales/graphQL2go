package main

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (graphQLType GraphQLType) GenerateGoStruct(jsonTag, dynamodbTag bool) (strings.Builder, strings.Builder) {
	var goCode strings.Builder
	var goCodeHeader strings.Builder

	importTime := false

	formattedType := fmt.Sprintf("type %s struct \n", graphQLType.Name)

	goCode.WriteString(formattedType + "{")

	for _, field := range graphQLType.Fields {
		if !importTime && strings.HasPrefix(field.Type, "time.") {
			importTime = true
		}
		capitalizedName := title(field.Name)
		nameTags := tags(field.Name, jsonTag, dynamodbTag)
		formattedString := fmt.Sprintf("\t%s %s %s\n", capitalizedName, field.Type, nameTags)
		goCode.WriteString(formattedString)
	}
	goCode.WriteString("}")

	if importTime {
		goCodeHeader.WriteString("import \"time\"\n\n")
	}

	return goCodeHeader, goCode
}

func (graphQLEnum GraphQLEnum) GenerateEnumCode() strings.Builder {
	var goCode strings.Builder
	goCode.WriteString(fmt.Sprintf("type %s string \n\n", graphQLEnum.Name))
	goCode.WriteString("const ( \n")
	for _, value := range graphQLEnum.Values {
		formattedString := fmt.Sprintf("%s %s = \"%s\" \n", title(value), graphQLEnum.Name, value)
		goCode.WriteString(formattedString)
	}
	goCode.WriteString(") \n\n")
	return goCode
}

func title(value string) string {
	return cases.Title(language.Und, cases.NoLower).String(value)
}

func tags(value string, json, dynamodb bool) string {
	tag := make([]string, 0)

	if json {
		tag = append(tag, fmt.Sprintf("json:\"%s\"", value))
	}
	if dynamodb {
		tag = append(tag, fmt.Sprintf("dynamodbav:\"%s\"", value))
	}

	if len(tag) > 0 {
		return "`" + strings.Join(tag, " ") + "`"
	}
	return ""
}
