package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

func findLoggerFields(filePath string) ([]string, error) {
    file, err := os.Open(filePath)

    if err != nil {
        return nil, err
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var output []string

    keyValueRegexp := regexp.MustCompile("\"([^:]+)\": [^,}]+")
    withFieldRegexp := regexp.MustCompile("logger.WithField\\(\"([^,]+)\",")
    var process = false
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "logger.Fields{") || process {
            all := keyValueRegexp.FindAllStringSubmatch(line, -1)
            for _, match := range all {
                output = append(output, match[1])
            }
            process = !strings.Contains(line, "}")
        } else {
            match := withFieldRegexp.FindStringSubmatch(line)
            if match != nil {
                output = append(output, match[1])
            }
        }
    }

    _ = file.Close()
    return output, nil
}

func main() {
    if len(os.Args) != 2 {
        panic(fmt.Errorf("must have one arugment"))
    }
    args := os.Args[1:]
    root := args[0]
    var filePaths []string

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if strings.HasSuffix(path, ".go") {
            filePaths = append(filePaths, path)
        }
        return nil
    })
    if err != nil {
        panic(err)
    }

    fieldSet := make(map[string]bool)  // to dedupe

    for _, filePath := range filePaths {
        matches, _ := findLoggerFields(filePath)
        for _, match := range matches {
            fieldSet[match] = true
        }
    }
    for field := range fieldSet {
        fmt.Println(field)
    }
}
