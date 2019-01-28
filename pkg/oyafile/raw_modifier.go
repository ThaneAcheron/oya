package oyafile

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var importKey = "Import:"
var projectKey = "Project:"
var uriVal = "  %s: %s"
var importRegxp = regexp.MustCompile("(?m)^" + importKey + "$")
var projectRegxp = regexp.MustCompile("^" + projectKey)

type OyafileRawModifier struct {
	filePath string
	file     []byte
}

func NewOyafileRawModifier(oyafilePath string) (OyafileRawModifier, error) {
	file, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return OyafileRawModifier{}, err
	}

	return OyafileRawModifier{
		filePath: oyafilePath,
		file:     file,
	}, nil
}

func (o *OyafileRawModifier) addImport(name string, uri string) error {
	uriStr := fmt.Sprintf(uriVal, name, uri)
	fileContent := string(o.file)
	var output []string
	if importRegxp.MatchString(fileContent) {
		output = o.appendAfter(importRegxp, []string{uriStr})
	} else if projectRegxp.MatchString(fileContent) {
		output = o.appendAfter(projectRegxp, []string{importKey, uriStr, ""})
	} else {
		output = []string{importKey, uriStr}
		output = append(output, strings.Split(fileContent, "\n")...)
	}
	if err := writeToFile(o.filePath, output); err != nil {
		return err
	}

	return nil
}

func (o *OyafileRawModifier) appendAfter(find *regexp.Regexp, data []string) []string {
	var output []string
	fileArr := strings.Split(string(o.file), "\n")
	for _, line := range fileArr {
		output = append(output, line)
		if find.MatchString(line) {
			output = append(output, data...)
		}
	}
	return output
}

func writeToFile(filePath string, content []string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, []byte(strings.Join(content, "\n")), info.Mode()); err != nil {
		return err
	}
	return nil
}
