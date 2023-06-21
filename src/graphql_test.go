package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/graphql-go/graphql/gqlerrors"
)

func TestParseGraphQLSchemaInvalidSchema(t *testing.T) {
	/* Arrange */
	expectedErrorStart := "Syntax Error GraphQL (1:6) Expected Name, found EOF"
	/* Act */
	_, _, err := ParseGraphQLSchema("type ")
	/* Assert */
	if err == nil {
		t.Errorf("Error should not be null")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(&gqlerrors.Error{}) {
		t.Errorf("Error should be `gqlerrors.Error`, instead is %v", reflect.TypeOf(err))
	}

	if !strings.HasPrefix(err.Error(), expectedErrorStart) {
		t.Errorf("Error should start with `%v`, instead is %v", expectedErrorStart, err.Error())
	}

}
func TestParseGraphQLSchemaEmptySchema(t *testing.T) {
	/* Arrange */

	/* Act */
	types, enums, err := ParseGraphQLSchema("")
	/* Assert */
	if err != nil {
		t.Errorf("Error should be null, instead %v", err)
	}

	if len(types) > 0 {
		t.Error("Types length should be more than 0")
	}

	if len(enums) > 0 {
		t.Error("Enums length should be more than 0")
	}
}
func TestParseGraphQLSchemaTypeSchema(t *testing.T) {
	/* Arrange */
	expectedType := GraphQLType{
		Name:   "User",
		Fields: []GraphQLField{{Name: "id", Type: "string"}, {Name: "name", Type: "string"}, {Name: "age", Type: "int"}},
	}
	scheme := `
	type User {
		id: ID
		name: String
		age: Int
	}
	`
	var actualType GraphQLType
	/* Act */
	types, enums, err := ParseGraphQLSchema(scheme)
	/* Assert */
	if err != nil {
		t.Errorf("Error should be null, instead %v", err)
	}

	if len(types) != 1 {
		t.Error("Types length should be 1")
	}

	actualType = types[0]

	if actualType.Name != expectedType.Name {
		t.Errorf("Expected GraphQLType.Name to be %s instead got %s", types[0].Name, expectedType.Name)
	}

	if len(actualType.Fields) != len(expectedType.Fields) {
		t.Errorf("Expected GraphQLType.Fields length to be %d instead got %d", len(expectedType.Fields), len(actualType.Fields))
	}

	for i := 0; i < len(expectedType.Fields); i++ {
		currentActualType := actualType.Fields[i]
		currentExpectedType := expectedType.Fields[i]
		if currentActualType.Name != currentExpectedType.Name {
			t.Errorf("Expected GraphQLType.Field.Name to be %s instead got %s", currentActualType.Name, currentExpectedType.Name)
		}
		if currentActualType.Type != currentExpectedType.Type {
			t.Errorf("Expected GraphQLType.Field.Type to be %s instead got %s", currentActualType.Type, currentExpectedType.Type)
		}
	}

	if len(enums) > 0 {
		t.Error("Enums length should be more than 0")
	}
}
