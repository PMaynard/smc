package main

import (
	"os"
	"testing"
)

var DMIDExample = "testdata/dependencymodels/idepend/SCADA-DM-XML.xml"
var DMIDExampleText = "testdata/dependencymodels/idepend/SCADA-DM-XML.txt"

/*********************************************************************/
func TestIDependParseFile(t *testing.T) {
	/* 1. Parse DM from iDepend  */
	idm := IDependModel{}

	out := t.TempDir() + "idepend"
	if err := idm.ParseFile(DMIDExample, out); err != nil {
		t.Fatal(err)
	}

	/* 2. Load Golden File */
	exp, err := os.ReadFile(DMIDExampleText)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}

	/* 3. Check */
	if string(exp) != string(got) {
		t.Error("Generated Output does not match:", DMIDExampleText, out)
		// t.Error("Expected:\n", string(exp))
		// t.Error("Got:\n", string(got))
	}

}
