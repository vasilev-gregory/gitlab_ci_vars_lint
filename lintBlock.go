package main

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/fatih/color"
)

var (
	reUsedVars = regexp.MustCompile(`\$[{]?(\w+)[}]?`)
)

func lintBlock(block map[string]interface{}) {
	script := getScript(block)
	variables := getVariables(block)
	varsChecked := make(map[string]struct{})
	printAttention := color.New(color.FgRed, color.BgBlack, color.Bold).PrintfFunc()
	usedVarsMatches := reUsedVars.FindAllStringSubmatch(script, -1)
	for _, varDecl := range variables {
		switch t := varDecl.(type) {
		case string:
			usedVarsMatches = append(usedVarsMatches, reUsedVars.FindAllStringSubmatch(t, -1)...)
		}
	}
	for _, match := range usedVarsMatches {
		varName := match[1]
		declared := false
		if _, ok := varsChecked[varName]; ok {
			continue
		}
		varsChecked[varName] = struct{}{}

		reDecl := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*(%s=.*)$`, varName))
		decls := reDecl.FindAllStringSubmatch(script, -1)
		if len(decls) > 0 {
			for _, decl := range decls {
				fmt.Println(decl[1])
				declared = true
			}
		}
		if decl, ok := variables[varName]; ok {
			fmt.Printf("%s:%#v\n", varName, decl)
			//fmt.Println(varName+":", decl)
			if "$"+varName != fmt.Sprint(decl) {
				declared = true
			}
		}
		if doc, ok := gitlabVars[varName]; ok {
			fmt.Printf("%s - gitlab predefined var: \"%s\"\n", varName, doc.descr)
			declared = true
		}
		if !declared {
			printAttention("\"%s\" - NOT DECLARED IN THIS STAGE", varName)
			fmt.Println()
		}
	}
	fmt.Println("")
}

func getScript(block map[string]interface{}) string {
	scriptInterface, ok := block["script"]
	script := ""
	if !ok {
		fmt.Println("no script here")
		return ""
	}
	switch t := scriptInterface.(type) {
	case string:
		script = t
	case []interface{}:
		for _, n := range t {
			script += n.(string) + "\n"
		}
	default:
		var r = reflect.TypeOf(t)
		// todo probably should panic here
		fmt.Printf("Other:%v\n", r)
	}
	return script
}

func getVariables(block map[string]interface{}) map[string]interface{} {
	variablesInterface, ok := block["variables"]
	if !ok {
		return map[string]interface{}{}
	}
	return variablesInterface.(map[string]interface{})
}
