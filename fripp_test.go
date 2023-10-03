package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

/**/
var FRIPPExample = "testdata/fripp/example.fripp"
var FRIPPExampleGML = "testdata/fripp/example.graphml"

/* Remove all white space */
/* TODO: This is not be the best way to comapre XML. */
func RemoveWS(src string) string {
	rtn := strings.ReplaceAll(src, " ", "")
	return strings.ReplaceAll(rtn, "\n", "")
}

/**/
func check(t *testing.T, ir FRIPP) {
	checks := []struct {
		path string
		len  int
	}{
		{path: "/PlaybookProcess", len: 1},
		{path: "/PlaybookProcess[@name='Test Name']", len: 1},
		{path: "/PlaybookProcess[@artifactInStateUsed='//@artifact.0/@state.0/@artifactstateinstance.2']", len: 1},
		{path: "/PlaybookProcess[@resultArtifactInState='//@artifact.0/@state.0/@artifactstateinstance.0']", len: 1},
		{path: "/PlaybookProcess[@resourceUsed='//@resource.2']", len: 1},
		{path: "/PlaybookProcess[@relatedreferences='//@externalreferences.0']", len: 1},

		{path: "/PlaybookProcess/artifact/state/artifactstateinstance", len: 3},
		{path: "/PlaybookProcess/artifact/state/artifactstateinstance[@usedByActivity='//@process.0']", len: 1},
		{path: "/PlaybookProcess/artifact/state/artifactstateinstance[@originatingActivity='//@process.0']", len: 1},

		{path: "/PlaybookProcess/process", len: 2},
		{path: "/PlaybookProcess/process[@xsi:type='FRIPP:PlaybookProcess']", len: 2},
		{path: "/PlaybookProcess/process[@artifactInStateUsed='//@artifact.0/@state.0/@artifactstateinstance.0']", len: 1},
		{path: "/PlaybookProcess/process[@artifactInStateUsed='//@artifact.0/@state.0/@artifactstateinstance.1']", len: 1},
		{path: "/PlaybookProcess/process[@resultArtifactInState='//@artifact.0/@state.0/@artifactstateinstance.1']", len: 1},
		{path: "/PlaybookProcess/process[@resultArtifactInState='//@artifact.0/@state.0/@artifactstateinstance.2']", len: 1},

		{path: "/PlaybookProcess/resource[@xsi:type='FRIPP:Actuator']", len: 3},

		{path: "/PlaybookProcess/externalreferences[@name='google']", len: 1},
	}

	for _, check := range checks {
		res := ir.data.FindElements(check.path)
		found := len(res)
		if found != check.len {
			t.Errorf("%s should be '%d' not '%d'", check.path, check.len, found)
			// t.Log(ir.data.WriteToString())
		}
	}
}

/**/
func TestFRIPPParseFile(t *testing.T) {
	/* 1. Parse FRIPP from FRIPP  */
	ir := FRIPP{}
	if err := ir.ParseFile(FRIPPExample); err != nil {
		t.Fatal(err)
	}

	/* 2. Check the main elements exist. */
	check(t, ir)
}

/**/
func TestFRIPPParseFileGraphML(t *testing.T) {
	ir := FRIPP{}
	if err := ir.ParseFileGraphML(FRIPPExampleGML); err != nil {
		t.Fatal(err)
	}

	/* 2. Check the main elements exist. */
	check(t, ir)
}

/**/
func TestFRIPPConvertFRIPPtoGraphML(t *testing.T) {
	/* 1. Parse FRIPP from FRIPP  */
	ir := FRIPP{}
	if err := ir.ParseFile(FRIPPExample); err != nil {
		t.Fatal(err)
	}

	/* 2. Convert from FRIPP to GraphML */
	var buf bytes.Buffer
	gml, err := ir.GraphML()
	if err != nil {
		t.Fatal(err)
	}
	if err := gml.Encode(&buf, true); err != nil {
		t.Fatal(err)
	}

	/* 3. Load Golden File */
	expected, err := os.ReadFile(FRIPPExampleGML)
	if err != nil {
		t.Fatal(err)
	}

	/* 4. Compare */
	got := RemoveWS(buf.String())
	exp := RemoveWS(string(expected))
	if got != exp {
		t.Error("Generated GraphML does not match:", FRIPPExampleGML)
		t.Error("Expected:\n", exp)
		t.Error("Got:\n", got)
	}
}

/**/
func TestFRIPPConvertGraphMLtoFRIPP(t *testing.T) {
	/* 1. Parse FRIPP from GraphML  */
	fromFRIPP := FRIPP{}
	if err := fromFRIPP.ParseFileGraphML(FRIPPExampleGML); err != nil {
		t.Fatal(err)
	}

	/* 2. Output to FRIPP */
	toFRIPPFile := t.TempDir() + "FRIPP"
	f, err := os.Create(toFRIPPFile)
	if err != nil {
		t.Fatal(err)
	}

	n, err := f.WriteString(fromFRIPP.String())
	if err != nil {
		t.Fatal(err)
	}

	if n != 1456 {
		t.Fatal("1456 Bytes should have been written.", n)
	}

	/* 4. Parse FRIPP */
	newIR := FRIPP{}
	if err := newIR.ParseFile(toFRIPPFile); err != nil {
		t.Fatal(err)
	}

	/* 3. Compare */
	check(t, newIR)
}
