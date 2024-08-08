package util

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

type FileTreeEntry struct {
	Name     string
	IsDir    bool
	Content  string
	Children []*FileTreeEntry
}

func ParseTemplate(text string, params map[string]interface{}) (string, error) { //nolint
	tf := template.FuncMap{
		"isInt": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			//nolint
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				return true
			default:
				return false
			}
		},
		"isString": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			//nolint
			switch v.Kind() {
			case reflect.String:
				return true
			default:
				return false
			}
		},
		"isSlice": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			//nolint
			switch v.Kind() {
			case reflect.Slice:
				return true
			default:
				return false
			}
		},
		"isArray": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			//nolint
			switch v.Kind() {
			case reflect.Array:
				return true
			default:
				return false
			}
		},
		"isMap": func(i interface{}) bool {
			v := reflect.ValueOf(i)
			//nolint
			switch v.Kind() {
			case reflect.Map:
				return true
			default:
				return false
			}
		},
		"add": func(i, add int) int {
			return i + add
		},
	}

	tmpl, err := template.New("tpl").Delims("{{%", "%}}").Funcs(tf).Parse(text)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, params); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func ReplaceGoPackages(content, owner, newName string) string {
	return strings.ReplaceAll(
		content,
		fmt.Sprintf("%s/verathread-app-template", owner),
		fmt.Sprintf("%s/%s", owner, newName),
	)
}
