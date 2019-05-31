// Copyright (C) 2015 A.Newman
//
package main

// basicTemplate defines the basic type and a "stringer" function.
//
// All other templates start with this.
//
// This template is incomplete and requires an {{end}} to terminate a
// {{range}} performed over the enumerated types being processed.
//
const basicTemplate = `// DO NOT EDIT.
//
// Generated: {{.Time}}
// From:      {{.Filename}}
// By:        {{.User}}
//

package {{.Package}}

{{- if .Imports}}

import ({{range .Imports}}
	"{{.}}"
{{- end}}
)
{{- end}}
{{range .Enums}}
type {{.TypeName}} {{.BaseType}}

const (
	{{.TypeName}}_Zero_ {{.TypeName}} = iota
{{- $TypeName:=.TypeName -}}
{{- range .Enumerators }}
	{{$TypeName}}_{{.Enum}}
{{- end}}
)

func (v {{.TypeName}}) String() string {
	switch v {
	case {{.TypeName}}_Zero_:
		return "*!!!* UNINITIALIZED {{.TypeName}} VALUE *!!!*"
        {{- $TypeName:=.TypeName -}}
	{{- range .Enumerators}}
	case {{$TypeName}}_{{.Enum}}:
		return "{{.Tag}}"
	{{- end}}
	default:
		return "*!* INVALID {{.TypeName}} VALUE *!*"
	}
}
`

// stdTemplate defines the default generated output which is
// a basicTemplate that.
//
const stdTemplate = basicTemplate + `
{{end}}
`

// scanStringTemplate provides the scanString() method to convert a
// string to a T.
//
const scanStringTemplate = `
func (v *{{.TypeName}}) scanString(s string) error {
	switch s {
	case "0":
		*v = {{.TypeName}}_Zero_
        {{- $TypeName:=.TypeName -}}
	{{- range .Enumerators}}
	case "{{.Tag}}":
		*v = {{$TypeName}}_{{.Enum}}
	{{- end}}
	default:
		return fmt.Errorf("*!* INVALID {{.TypeName}} LITERAL: %q *!*", s)
	}
	return nil
}
`

// scanTemplate provides an implementation of the fmt Scanner interface's
// Scan() method allowing T values to be read from io.Readers as text.
//
const scanTemplate = basicTemplate + scanStringTemplate + `
func (v *{{.TypeName}}) Scan(state fmt.ScanState, verb rune) error {
    if token, err := state.Token(true, nil); err != nil {
        return err
    } else {
        return v.scanString(string(token))
    }
}
{{end}}
`

// sqlTemplate provides an implementation of the database/sql Scan
// method allowing T values to be stored textually in database columns.
//
const sqlTemplate = basicTemplate + scanStringTemplate + `
func (v *{{.TypeName}}) Scan(src interface{}) error {
        return v.scanString(src.(string))
}
{{end}}
`

//
const jsonTemplate = basicTemplate + `
// json.MarshalJSON/UnmarshalJSON implementation here
`

//
const xmlTemplate = basicTemplate + `
// xml.MarshalXML/UnmarshalXML implementation here
`
