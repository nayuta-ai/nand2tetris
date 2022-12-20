package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"jackcompiler/analyzer"

	"github.com/sirupsen/logrus"
)

var format string

func init() {
	flag.StringVar(&format, "format", "", "log format")
	flag.StringVar(&format, "f", "", "log format")
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// The below line parses arguments.
	// The arguments contain a file path that indicates a file with jack file extensions.
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		logrus.Fatal(err)
	}
	if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	path := flag.Args()
	if len(path) < 1 || len(path) > 1 {
		logrus.Fatal("The length of arguments should be 1.")
	}
	files, err := filepath.Glob(path[0] + "/*.jack")
	if err != nil {
		logrus.Fatal(err)
	}
	if len(files) < 1 {
		logrus.Fatal(fmt.Sprintf("There is no file in %s", path[0]))
	}
	for i, file := range files {
		fp, err := os.Open(file)
		if err != nil {
			logrus.Fatal(err)
		}
		tk, err := analyzer.SyntaxAnalyzer(fp)
		if err != nil {
			logrus.Fatal(err)
		}
		err = convertToTXMLformat(i, tk)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func convertToTXMLformat(i int, token []analyzer.Token) error {
	var contents string
	contents += "<tokens>\n"
	for _, tk := range token {
		contents += fmt.Sprintf("<%s> %s </%s>\n", tk.Type, tk.Value, tk.Type)
	}
	contents += "</tokens>"
	f, err := os.Create(fmt.Sprintf("sample%dT.xml", i))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	if err != nil {
		return err
	}
	return nil
}
