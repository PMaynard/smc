package main

import (
	"testing"
)

/**/
func TestDecide(t *testing.T) {

	tmp := t.TempDir()

	checks := []struct {
		format string
		in     string
		out    string
	}{
		/* FRIPP */
		{"fripp", FRIPPExample, tmp + "fripp.fripp"},
		{"fripp", FRIPPExample, tmp + "fripp.graphml"},
		{"fripp", FRIPPExampleGML, tmp + "frippGML.fripp"},
		{"fripp", FRIPPExampleGML, tmp + "frippGML.graphml"},

		/* DM */
		{"dm", DMSecMofExample, tmp + "secmof.dependencymodel"},
		{"dm", DMSecMofExample, tmp + "secmof.graphml"},
		{"dm", DMSecMofGraphML, tmp + "secmofGML.dependencymodel"},
		{"dm", DMSecMofGraphML, tmp + "secmofGML.graphml"},
		/* iDepend */
		{"idm", DMIDExample, tmp + "idepend.txt"},

		/* SAND */
		{"sand", SANDExample, tmp + "sand.graphml"},
		{"sand", SANDExampleGML, tmp + "sandGML.graphml"},

		/* Convert */
		{"sand", SANDExample, tmp + "sand.dependencymodel"},
		{"sand", SANDExample, tmp + "sand.fripp"},
	}

	for _, c := range checks {
		err := Decide(c.format, c.in, c.out)
		if err != nil {
			t.Error(err)
		}
	}
}

/**/
func TestDecideInvalid(t *testing.T) {
	tmp := t.TempDir()

	checks := []struct {
		format string
		in     string
		out    string
	}{
		{"nothin", FRIPPExample, ""},
		{"fripp", "missing", "fripp.fripp"},
		{"fripp", FRIPPExample, "fripp.fr"},
		{"dm", "missing", "fripp.fripp"},
		{"dm", DMSecMofExample, "fripp.fr"},

		{"idm", "missing", tmp + "idepend.txt"},
		{"idm", DMSecMofExample, "fripp.fr"},
		{"idm", "missing", tmp + "fripp.txt"},
		{"idm", "missing", tmp + "missing"},
		{"idm", "", tmp + ""},

		{"sand", "", tmp + ""},
	}

	for _, c := range checks {
		err := Decide(c.format, c.in, c.out)
		if err == nil {
			t.Error(err)
		}
	}
}
