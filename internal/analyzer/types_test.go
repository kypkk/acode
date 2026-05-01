package analyzer

import (
	"reflect"
	"testing"
)

func TestExtractTypes_NoTypes(t *testing.T) {
	tree, src := parseGo(t, "package x\nfunc F() {}\n")
	types := ExtractTypes(tree, src)
	if len(types) != 0 {
		t.Fatalf("got %d types, want 0: %+v", len(types), types)
	}
}

func TestExtractTypes_StructWithNamedFields(t *testing.T) {
	src := `package x

type Foo struct {
	Name string
	Age  int
}
`
	tree, bs := parseGo(t, src)
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1: %+v", len(types), types)
	}
	got := types[0]
	if got.Name != "Foo" {
		t.Errorf("Name = %q, want %q", got.Name, "Foo")
	}
	if got.Kind != "struct" {
		t.Errorf("Kind = %q, want %q", got.Kind, "struct")
	}
	if got.Line != 3 {
		t.Errorf("Line = %d, want 3", got.Line)
	}
	wantFields := []Field{
		{Name: "Name", Type: "string"},
		{Name: "Age", Type: "int"},
	}
	if !reflect.DeepEqual(got.Fields, wantFields) {
		t.Errorf("Fields = %#v, want %#v", got.Fields, wantFields)
	}
}

func TestExtractTypes_StructWithMultiNameField(t *testing.T) {
	src := `package x

type Person struct {
	Name, Email string
	Age         int
}
`
	tree, bs := parseGo(t, src)
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1", len(types))
	}
	wantFields := []Field{
		{Name: "Name", Type: "string"},
		{Name: "Email", Type: "string"},
		{Name: "Age", Type: "int"},
	}
	if !reflect.DeepEqual(types[0].Fields, wantFields) {
		t.Errorf("Fields = %#v, want %#v", types[0].Fields, wantFields)
	}
}

func TestExtractTypes_StructWithEmbeddedField(t *testing.T) {
	src := `package x

import "io"

type Foo struct {
	io.Reader
	Name string
}
`
	tree, bs := parseGo(t, src)
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1", len(types))
	}
	wantFields := []Field{
		{Name: "", Type: "io.Reader"},
		{Name: "Name", Type: "string"},
	}
	if !reflect.DeepEqual(types[0].Fields, wantFields) {
		t.Errorf("Fields = %#v, want %#v", types[0].Fields, wantFields)
	}
}

func TestExtractTypes_StructWithPointerAndSliceFields(t *testing.T) {
	src := `package x

type Foo struct {
	P *T
	S []int
}
`
	tree, bs := parseGo(t, src)
	types := ExtractTypes(tree, bs)
	wantFields := []Field{
		{Name: "P", Type: "*T"},
		{Name: "S", Type: "[]int"},
	}
	if !reflect.DeepEqual(types[0].Fields, wantFields) {
		t.Errorf("Fields = %#v, want %#v", types[0].Fields, wantFields)
	}
}

func TestExtractTypes_EmptyStruct(t *testing.T) {
	tree, bs := parseGo(t, "package x\ntype Foo struct{}\n")
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1", len(types))
	}
	if types[0].Kind != "struct" {
		t.Errorf("Kind = %q, want struct", types[0].Kind)
	}
	if len(types[0].Fields) != 0 {
		t.Errorf("Fields = %#v, want empty", types[0].Fields)
	}
}

func TestExtractTypes_Interface(t *testing.T) {
	src := `package x

type R interface {
	Read(p []byte) (n int, err error)
	Close() error
	Reset()
}
`
	tree, bs := parseGo(t, src)
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1", len(types))
	}
	got := types[0]
	if got.Name != "R" || got.Kind != "interface" {
		t.Errorf("Name=%q Kind=%q", got.Name, got.Kind)
	}
	wantMethods := []Method{
		{Name: "Read", Parameters: []string{"p []byte"}, ReturnTypes: []string{"n int", "err error"}},
		{Name: "Close", ReturnTypes: []string{"error"}},
		{Name: "Reset"},
	}
	if !reflect.DeepEqual(got.Methods, wantMethods) {
		t.Errorf("Methods = %#v, want %#v", got.Methods, wantMethods)
	}
}

func TestExtractTypes_EmptyInterface(t *testing.T) {
	tree, bs := parseGo(t, "package x\ntype Any interface{}\n")
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1", len(types))
	}
	if types[0].Kind != "interface" {
		t.Errorf("Kind = %q", types[0].Kind)
	}
	if len(types[0].Methods) != 0 {
		t.Errorf("Methods = %#v, want empty", types[0].Methods)
	}
}

func TestExtractTypes_Alias(t *testing.T) {
	tree, bs := parseGo(t, "package x\ntype X = string\n")
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1", len(types))
	}
	got := types[0]
	if got.Name != "X" || got.Kind != "alias" || got.Underlying != "string" {
		t.Errorf("got %+v", got)
	}
}

func TestExtractTypes_NamedType(t *testing.T) {
	tree, bs := parseGo(t, "package x\ntype Names []string\n")
	types := ExtractTypes(tree, bs)
	if len(types) != 1 {
		t.Fatalf("got %d types, want 1", len(types))
	}
	got := types[0]
	if got.Name != "Names" || got.Kind != "named" || got.Underlying != "[]string" {
		t.Errorf("got %+v", got)
	}
}

func TestExtractTypes_GroupedDeclaration(t *testing.T) {
	src := `package x

type (
	A struct {
		X int
	}
	B = string
	C int
)
`
	tree, bs := parseGo(t, src)
	types := ExtractTypes(tree, bs)
	if len(types) != 3 {
		t.Fatalf("got %d types, want 3: %+v", len(types), types)
	}
	if types[0].Name != "A" || types[0].Kind != "struct" {
		t.Errorf("[0] = %+v", types[0])
	}
	wantA := []Field{{Name: "X", Type: "int"}}
	if !reflect.DeepEqual(types[0].Fields, wantA) {
		t.Errorf("[0].Fields = %#v, want %#v", types[0].Fields, wantA)
	}
	if types[1].Name != "B" || types[1].Kind != "alias" || types[1].Underlying != "string" {
		t.Errorf("[1] = %+v", types[1])
	}
	if types[2].Name != "C" || types[2].Kind != "named" || types[2].Underlying != "int" {
		t.Errorf("[2] = %+v", types[2])
	}
}

func TestExtractTypes_MixedTopLevel(t *testing.T) {
	src := `package x

type Foo struct{ N int }

func F() {}

type Bar interface{ Do() }
`
	tree, bs := parseGo(t, src)
	types := ExtractTypes(tree, bs)
	if len(types) != 2 {
		t.Fatalf("got %d types, want 2: %+v", len(types), types)
	}
	if types[0].Name != "Foo" || types[0].Kind != "struct" {
		t.Errorf("[0] = %+v", types[0])
	}
	if types[1].Name != "Bar" || types[1].Kind != "interface" {
		t.Errorf("[1] = %+v", types[1])
	}
}
