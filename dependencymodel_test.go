package main

import (
	"bytes"
	"os"
	"testing"
)

/**/
var DMSecMofExample = "testdata/dependencymodels/secmof/example.dependencymodel"
var DMSecMofGraphML = "testdata/dependencymodels/secmof/example.graphml"

/**/
func TestDMParseFile(t *testing.T) {

}

/**/
func TestDMParseFileGraphML(t *testing.T) {

}

/**/
func TestDMConvertDMtoGraphML(t *testing.T) {
	/* 1. Parse DM from SecMof  */
	dm := DependencyModel{}
	if err := dm.ParseFile(DMSecMofExample); err != nil {
		t.Fatal(err)
	}

	/* 2. Convert from DM to GraphML */
	var buf bytes.Buffer
	gml, err := dm.GraphML()
	if err != nil {
		t.Fatal(err)
	}
	if err := gml.Encode(&buf, true); err != nil {
		t.Fatal(err)
	}

	/* 3. Load Golden File */
	expected, err := os.ReadFile(DMSecMofGraphML)
	if err != nil {
		t.Fatal(err)
	}

	/* 4. Compare */
	got := RemoveWS(buf.String())
	exp := RemoveWS(string(expected))

	if got != exp {
		t.Errorf("Generated GraphML from '%s' does not match '%s':", DMSecMofExample, DMSecMofGraphML)
		t.Error("Expected:\n", string(exp))
		t.Error("Got:\n", buf.String())
	}
}

/**/
func TestDMConvertGraphMLtoDM(t *testing.T) {
	/* 1. Parse DM from Graphml  */
	dm := DependencyModel{}
	if err := dm.ParseFileGraphML(DMSecMofGraphML); err != nil {
		t.Fatal(err)
	}

	/* 2. Convert from DM to GraphML */
	var buf bytes.Buffer
	gml, err := dm.GraphML()
	if err != nil {
		t.Fatal(err)
	}
	if err := gml.Encode(&buf, true); err != nil {
		t.Fatal(err)
	}

	/* TODO:
	 *  Due to the way the GraphML is generated the order is not preserved.
	 *  An in depth comparison needs to be performed.
	 */

	/* 3. Load Golden File */
	// expected, err := os.ReadFile(DMSecMofGraphML)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// /* 4. Compare */
	// got := len(RemoveWS(buf.String()))
	// exp := len(RemoveWS(string(expected)))

	// if got != exp {
	// 	t.Error("Generated GraphML does not match:", DMSecMofGraphML)
	// 	t.Error("Expected:\n", exp)
	// 	t.Error("Got:\n", got)
	// }
}
