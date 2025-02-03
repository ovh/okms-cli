package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

func main() {
	outFile, err := os.Create("dependencies.txt")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	out := bufio.NewWriter(outFile)
	defer out.Flush()

	dir := os.DirFS("LICENSES/go")

	toUpdate := []string{}

	err = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			fmt.Println(path)
			contentBytes, err := fs.ReadFile(dir, path)
			if err != nil {
				return err
			}
			content := string(contentBytes)

			component := strings.TrimSuffix(path, "/LICENSE")
			re := regexp.MustCompile(`Copyright ((\(([cC])\)|©) )?\d+?.+`)
			copyright := "Copyright <YEAR> <COPYRIGHT HOLDER>\n"
			if !re.MatchString(content) {
				content = strings.ReplaceAll(content, "== "+component+" ==", "== "+component+" ==\n"+copyright)
				toUpdate = append(toUpdate, component)
			}

			if _, err := out.WriteString(content + "\n\n"); err != nil {
				return err
			}
		}
		return nil
	})

	fmt.Println("\nTo update manually:")
	for _, c := range toUpdate {
		fmt.Printf("- %s\n", c)
	}

	if err != nil {
		panic(err)
	}
}
