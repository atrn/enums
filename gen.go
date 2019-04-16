// Copyright (C) 2015 A.Newman
//
package main

import (
	"io"
	"text/template"
)

type CodeGenerator interface {
	Imports() []string
	Generate(Enums, io.Writer) error
}

type Generator struct {
	template *template.Template
	imports  []string
}

func NewGenerator(name, text string, imports []string) *Generator {
	return &Generator{
		template: template.Must(template.New(name).Parse(text)),
		imports:  imports,
	}
}

func (cg *Generator) Imports() []string {
	return cg.imports
}

func (cg *Generator) Generate(e Enums, w io.Writer) error {
	return cg.template.Execute(w, e)
}
