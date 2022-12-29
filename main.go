package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	reservedCINames = map[string]struct{}{
		"image":    struct{}{},
		"include":  struct{}{},
		"stages":   struct{}{},
		"workflow": struct{}{},
	}
)

func init() {
	parseGilabPredefinedVars()
}

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {

			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && (filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml") {
			lintYml(path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
