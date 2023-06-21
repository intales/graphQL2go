package main

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
)

// GraphQLType represents a GraphQL type
type GraphQLType struct {
	Name   string         `json:"name"`
	Fields []GraphQLField `json:"fields"`
}

// GraphQLField represents a GraphQL field
type GraphQLField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// GraphQLEnum represents a GraphQL enum type
type GraphQLEnum struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func ParseGraphQLSchema(schemaString string) ([]GraphQLType, []GraphQLEnum, error) {
	var schemaTypes []GraphQLType
	var schemaEnums []GraphQLEnum

	astDoc, err := parser.Parse(parser.ParseParams{
		Source: schemaString,
	})

	if err != nil {
		return schemaTypes, schemaEnums, err
	}

	for _, definition := range astDoc.Definitions {
		switch typeDefinition := definition.(type) {
		case *ast.ObjectDefinition:
			schemaTypes = append(schemaTypes, parseGraphQLFields(typeDefinition))
		case *ast.EnumDefinition:
			schemaEnums = append(schemaEnums, parseGraphQLEnums(typeDefinition))
		}
	}

	return schemaTypes, schemaEnums, nil
}

func parseGraphQLFields(def *ast.ObjectDefinition) GraphQLType {
	var fields []GraphQLField
	for _, field := range def.Fields {
		fields = append(fields, GraphQLField{
			Name: field.Name.Value,
			Type: getTypeString(field.Type),
		})
	}
	return GraphQLType{
		Name:   def.Name.Value,
		Fields: fields,
	}
}

func parseGraphQLEnums(def *ast.EnumDefinition) GraphQLEnum {
	var values []string
	for _, enumValue := range def.Values {
		values = append(values, enumValue.Name.Value)
	}
	return GraphQLEnum{
		Name:   def.Name.Value,
		Values: values,
	}
}

func getTypeString(t ast.Type) string {
	switch typ := t.(type) {
	case *ast.Named:
		return getNamedType(typ.Name.Value)
	case *ast.NonNull:
		return getTypeString(typ.Type)
	case *ast.List:
		return "[]" + getTypeString(typ.Type)
	default:
		return ""
	}
}

func getNamedType(value string) string {
	switch value {
	case "String":
		return "string"
	case "ID":
		return "string"
	case "Int":
		return "int"
	case "Boolean":
		return "bool"
	case "AWSDateTime":
		return "time.Time"
	case "AWSTimestamp":
		return "string"
	default:
		return value
	}
}
