// Copyright (C) 2015 A.Newman
//
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	genScan := flag.Bool("fmt-scanner", false, "Generate a fmt.Scanner Scan() method.")
	genDatabaseSQL := flag.Bool("sql-scanner", false, "Generate a sql.Scanner Scan() method.")
	genEncodingJSON := flag.Bool("json", false, "Output JSON marshal/unmarshal methods.")
	genEncodingXML := flag.Bool("xml", false, "Output XML marshal/unmarshal methods.")

	flag.Usage = usage
	flag.Parse()

	var generator CodeGenerator

	// Check if more than one option was supplied and
	// set our generator while we're at it. If more
	// than generator was requested we'll create more
	// than one but then exit with an error.
	//
	ngen := 0
	if *genScan {
		ngen++
		generator = NewGenerator("scan", scanTemplate, []string{"fmt"})
	}
	if *genDatabaseSQL {
		ngen++
		generator = NewGenerator("sql", sqlTemplate, []string{"fmt"})
	}
	if *genEncodingJSON {
		ngen++
		generator = NewGenerator("json", jsonTemplate, []string{"encoding/json"})
	}
	if *genEncodingXML {
		ngen++
		generator = NewGenerator("xml", xmlTemplate, []string{"encoding/xml"})
	}

	switch {
	case ngen == 0:
		generator = NewGenerator("std", stdTemplate, []string{})
	case ngen > 1:
		usage()
		os.Exit(1)
	}

	processAll := func(filenames []string) error {
		for _, filename := range filenames {
			if err := process(filename, generator); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return err
			}
		}
		return nil
	}

	var err error
	if flag.NArg() > 0 {
		err = processAll(flag.Args())
	} else {
		var filenames []string
		if filenames, err = filepath.Glob("*.enum"); err == nil {
			err = processAll(filenames)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "enums:", err)
		os.Exit(1)
	}
}

func process(filename string, generator CodeGenerator) error {
	e, err := parse(filename, generator)
	if err == nil {
		err = gen(e, filename, generator)
	}
	return err
}

func parse(filename string, generator CodeGenerator) (Enums, error) {
	var e Enums
	e.Filename = filename
	e.Time = time.Now().Format(time.ANSIC)
	e.User = "unknown"
	u, err := user.Current()
	if err == nil {
		e.User = u.Username
	}
	if file, err := os.Open(filename); err != nil {
		return Enums{}, err
	} else if err = ParseToEnd(file, &e); err != nil {
		return Enums{}, err
	}
	e.Imports = append(e.Imports, generator.Imports()...)
	return e, nil
}

func gen(e Enums, filename string, generator CodeGenerator) error {
	path := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".go"
	if file, err := os.Create(path); err == nil {
		defer func() {
			if err != nil {
				os.Remove(path)
			}
		}()
		if err = generator.Generate(e, file); err != nil {
			_ = file.Close()
			return err
		}
		err = file.Close()
		return err
	} else {
		return err
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: enums [option] [filename...]

A single option may be supplied to control code generation.
Option may be one of:
`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, `
In all cases the basic implementation of an enumerated type
is generated along with a String() method for textual output.
`)
}
