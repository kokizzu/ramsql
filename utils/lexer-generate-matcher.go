package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

// Usage:
// go run --tags=generate .../lexer-generate-matcher.go --init
// go run --tags=generate .../lexer-generate-matcher.go {--lexeme <STRING>}+ [--name <STRING>]

var (
	initFlag   bool
	lexemeFlag []string
	nameFlag   string
)

const outputFileName = "lexer.matcher.go"

func init() {
	flag.BoolVarP(&initFlag, "init", "", false, "Indicates whether to ")
	flag.StringArrayVarP(&lexemeFlag, "lexeme", "", []string{}, "The lexeme/s to include in the matcher function")
	flag.StringVarP(&nameFlag, "name", "", "", "The name of the matcher function (default: the first lexeme name)")
}

func main() {
	flag.Parse()

	// Check lexeme flag
	if !initFlag && len(lexemeFlag) <= 0 {
		log.Fatal("at least one \"--lexeme <STRING>\" option is required")
	}

	// Check name flag
	if !initFlag && nameFlag == "" {
		nameFlag = strings.Title(lexemeFlag[0])
	}

	if initFlag {
		err := initOutputFile(outputFileName)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"outputFile": outputFileName}).Fatal("failed to init output file")
		}
	} else {
		err := appendOutputFile(outputFileName, nameFlag, lexemeFlag)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{"outputFile": outputFileName, "lexemes": lexemeFlag, "name": nameFlag}).Fatal("failed to append to output file")
		}
	}
}

func initOutputFile(fileName string) error {
	fileBytes := bytes.NewBuffer(nil)

	fileBytes.WriteString("package parser\n")

	// Write file
	return ioutil.WriteFile(fileName, fileBytes.Bytes(), 0644)
}

func appendOutputFile(fileName string, name string, lexemes []string) error {
	fileBytes := bytes.NewBuffer(nil)

	fileBytes.WriteString("\n")
	fileBytes.WriteString(fmt.Sprintf("func (l *lexer) Match%sToken() bool {\n", name))
	fileBytes.WriteString(fmt.Sprintf("  return %s\n", strings.Join(lexemesToMatchFuncInvocations(name, lexemes), " ||\n    ")))
	fileBytes.WriteString(fmt.Sprintf("}\n"))

	// Write file
	fh, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.Write(fileBytes.Bytes())
	return err
}

func lexemesToMatchFuncInvocations(name string, lexemes []string) []string {
	rLexemes := make([]string, len(lexemes))

	for i, v := range lexemes {
		if len(v) > 1 {
			rLexemes[i] = fmt.Sprintf("l.Match([]byte(\"%s\"), %sToken)", v, name)
		} else {
			rLexemes[i] = fmt.Sprintf("l.MatchSingle('%s', %sToken)", v, name)
		}
	}

	return rLexemes
}
