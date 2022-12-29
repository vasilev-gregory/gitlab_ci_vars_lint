package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

var (
	printCurrYml = color.New(color.FgGreen, color.BgBlack, color.Bold).PrintfFunc()
	printStage   = color.New(color.FgBlue, color.BgBlack, color.Bold).PrintfFunc()
)

func lintYml(path string) {
	printCurrYml("# " + path)
	fmt.Println("")
	fmt.Println("")

	ci := make(map[string]interface{}, 0)
	data, _ := os.ReadFile(path)
	if err := yaml.Unmarshal([]byte(data), &ci); err != nil {
		log.Fatalf("error: %v", err)
	}

	file, _ := os.Open(path)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == ' ' {
			continue
		}
		lineSl := strings.Split(line, ":")
		if len(lineSl) < 2 {
			log.Fatalf("bad line: ", line)
		}
		blockKye := lineSl[0]
		if _, ok := reservedCINames[blockKye]; ok {
			continue
		}
		block := ci[blockKye].(map[string]interface{})
		printStage(`### analizing stage "%+v"`, blockKye)
		fmt.Println()
		lintBlock(block)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
}
