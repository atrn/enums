# enums - enumerated types via `go generate`

Enums is a 'go generate' tool that implements simple enumerated types
for Go. The enums tool reads _.enum_ files containing definitions of
enumerated types written in a Go-like manner. The tool reads these
files and generates corresponding Go source code to implement the
types defined in each file.

## Usage

    enums [option] [filename...]

## Options

- -fmt-scanner  
Generate a `fmt.Scanner` Scan() method.
- -sql-scanner  
Generate an `sql.Scanner` Scan() method.
- -json  
Generate JSON Marshal/Unmarshal methods.
- -xml  
Generate XML Marshal/Unmarshal methods.

## Description

The enum tool reads the specified files or, if no files are named, all
the `.enum` files in the current directory, and parses their content,
enumerated type definitions written in a Go-like language.

Each `_.enum` input file generates a `.go` file of the same name which
contains Go code to implement the enumerated types defined in the
file.

### Language

The enumerated types are defined using a Go like syntax which is
really a variant of Go's `type` statement.

We use the _keyword_ `enum` as the type of thing being defined,
e.g.,

    type MyType enum {
        First
        Second
        Third
    }

With the following idiomatic Go form,

    type MyType int
    const (
        First MyType = iota
        Second
        Third
    )

By default enums generates a basic type implementation that defines
the type, its enumerators and a String() method for the type. Switches
may be used to generate implementations of some standard
interfaces. Note, because some of the interfaces use the same method
names enums only supports generating a single type of output.

## Syntax

Enumerated types are defined using a Go-like syntax.  As in Go each
source file starts with a package declaration and is followed by one
or more enumerated type declarations.

Enumerated types use a modified Go type declaration syntax that uses
the psuedo-keyword 'enum' as the base type. This is followed by a
brace-delimited list of identifiers, the type's enumerators.

E.g.

    package zoo
    
	type Animal enum {
		Bear
            Gorilla
            Lion
            Seal
         }
    }

By default enums are represented as an int, the above
definition generates the code,

    type Animal int
    
    const (
           Animal_0 Animal = iota
           Bear
           Lion
           Gorilla
           Seal
    )
    
    func (a Animal) String() string {
        ...  stringer implementation
    }

## Enum Base Type

Enumerated types may define the type used as their base type by
following the 'enum' keyword in the type declaration with the name of
a Go inegral type.  E.g.,

    type MessageCode enum int16 {
        mycode1,
        mycode2,
        ...
    }

## Code Generation Templates

fmt     Defines the method `Scan(fmt.ScanState, rune)` to implement the
        `fmt` package's scanning interface. This allows values of the
        type to be read by name using the fmt Scan functions.

sql     Defines the method `Scan(inteface{}) error` to implement the
        `database/sql` pacakge's `Scanner` inteface. This allows values of
        the type to be stored in textual columns in databases.
