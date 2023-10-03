package main

import (
	"bytes"
	"os"
	"testing"
)

/**/
var SANDExample = "testdata/sand/example.ctrees"
var SANDExampleGML = "testdata/sand/example.graphml"

/**/
func TestSANDParseFile(t *testing.T) {
	/* 1. Parse SAND from CTrees  */
	st := SANDTree{}
	if err := st.ParseFile(SANDExample); err != nil {
		t.Fatal(err)
	}

	/* 2. Check */
	got := len(st.Nodes)
	exp := 9
	if got != exp {
		t.Errorf("%s should be '%d' not '%d'", SANDExample, got, exp)
	}

	/* TODO: More validation of the loaded data. */
}

/*********************************************************************/
func TestSANDParseFileGraphML(t *testing.T) {
	st := SANDTree{}
	if err := st.ParseFileGraphML(SANDExampleGML); err != nil {
		t.Fatal(err)
	}

	/* 2. Check */
	got := len(st.Nodes)
	exp := 9
	if got != exp {
		t.Errorf("%s should be '%d' not '%d'", SANDExample, got, exp)
	}

	/* TODO: More validation of the loaded data. */
}

/*********************************************************************/
func TestConvertSANDtoGraphML(t *testing.T) {
	/* 1. Parse SAND from CTrees  */
	st := SANDTree{}
	if err := st.ParseFile(SANDExample); err != nil {
		t.Fatal(err)
	}

	/* 2. Convert from SAND to GraphML */

	gml, err := st.GraphML()
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = gml.Encode(&buf, true)
	if err != nil {
		t.Fatal(err)
	}

	/* 3. Load Golden File */
	expected, err := os.ReadFile(SANDExampleGML)
	if err != nil {
		t.Fatal(err)
	}

	/* 4. Compare */
	got := RemoveWS(buf.String())
	exp := RemoveWS(string(expected))
	if got != exp {
		t.Error("Generated GraphML does not match:", SANDExampleGML)
		t.Error("Expected:\n", exp)
		t.Error("Got:\n", got)
	}
}
