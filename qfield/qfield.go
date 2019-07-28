package qfield

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/ast"
)

/*QField is a convienient way of creating a graphQL library with the gqlgen library. It allows for logical blocks of code to be created, that can call
each other to allow for nested objects in nested objects for a truly easy GraphQL experience that meets expectations.

The QField implementation was origonally written by John Zarate working on an API for GALIoT Systems. The code was copied and provided to github by
Samuel Archibald (IoTPanic) Examples on how to use QFields are provdied in the examples folder.
*/
type QField struct {
	Name      string // useful for debugging
	Children  map[string]QField
	Arguments map[string]interface{}
}

// ConstructQField to get a starting object
func constructQField() QField {
	ret := QField{}
	ret.Children = make(map[string]QField)
	ret.Arguments = make(map[string]interface{})
	return ret
}

// ContainsDirectChild will return boolean if the child field is within a qfild object
func (QF QField) ContainsDirectChild(k string) bool {
	_, ret := QF.Children[k]
	return ret
}

// GetDirectChild gets the qfield for a child field. Returns a boolean if the field exists, as well as said qfield.
func (QF QField) GetDirectChild(k string) (bool, QField) {
	if QF.ContainsDirectChild(k) {
		ret, _ := QF.Children[k]
		return true, ret
	}
	return false, constructQField()
}

// HasArg tests if a qfield has an argument
func (QF QField) HasArg(k string) bool {
	_, b := QF.Arguments[k]
	return b
}

// GetArg returns a boolean if the argument is present in a qfield and the value as an interface
func (QF QField) GetArg(k string) (bool, interface{}) {
	ret, b := QF.Arguments[k]
	return b, ret
}

func (QF QField) GetArgAsListOfString(k string) []*string {
	_, val := QF.GetArg(k)
	switch v := val.(type) {
	case []*string:
		return v
	default:
		return nil
	}
}

// GetArgAsString Returns an argument as a string
func (QF QField) GetArgAsString(k string) string {
	_, val := QF.GetArg(k)
	switch v := val.(type) {
	case *ast.Value:
		return (*v).Raw
	case string:
		return v
	case json.Number:
		return string(v)
	default:
		cast, _ := v.(string)
		return cast
	}
	//Val, _ := val.(*ast.Value)
	//return (*Val).Raw
}

// GetArgAsInt returns a specified argument as an int
func (QF QField) GetArgAsInt(k string) int {
	val, _ := strconv.Atoi(QF.GetArgAsString(k))
	return val
}

// GetArgAsBool returns a specified argument as a boolean
func (QF QField) GetArgAsBool(k string) bool {
	_, val := QF.GetArg(k)
	switch v := val.(type) {
	case bool:
		return v
	default:
		str, _ := val.(string)
		b, _ := strconv.ParseBool(str)
		return b
	}
}

func toQField2(f ast.Field, vars map[string]interface{}) QField {
	ret := constructQField()
	ret.Name = f.Name
	for _, a := range f.Arguments {
		//print(a.Value)
		value := (*(a.Value))
		string_value := value.Raw
		if value.Kind == ast.Variable {
			ret.Arguments[a.Name] = vars[string_value]
		} else {
			ret.Arguments[a.Name] = string_value
		}
	}
	for _, a := range f.SelectionSet {
		switch p := a.(type) {
		case *ast.Field:
			curField := *p
			ret.Children[curField.Name] = toQField2(curField, vars)
		}
	}
	return ret
}

func toQField(f graphql.CollectedField, vars map[string]interface{}) QField {
	ret := constructQField()
	ret.Name = f.Name
	for _, a := range f.Arguments {
		value := (*(a.Value))
		string_value := value.Raw
		if value.Kind == ast.Variable {
			ret.Arguments[a.Name] = vars[string_value]
		} else {
			ret.Arguments[a.Name] = string_value
		}
	}
	for _, a := range f.SelectionSet {
		switch p := a.(type) {
		case *ast.Field:
			curField := *p
			ret.Children[curField.Name] = toQField2(curField, vars)
		}
	}
	return ret
}

func asQField(ctx context.Context, vars map[string]interface{}) QField {
	var ret = constructQField()
	for _, f := range graphql.CollectFieldsCtx(ctx, nil) {
		fmt.Println(f.Name)
		ret.Children[f.Name] = toQField(f, vars)
	}
	return ret
}

//GetQField allows for the gqlgen context to be passed and the representing qfield is returned
func GetQField(ctx context.Context, vars map[string]interface{}) QField {
	return asQField(ctx, vars)
}
