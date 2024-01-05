package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

/* TODO: Rewrite using Generics */

var helpmsg = "Usage: smc [fripp|dm|idm|sand] in.fripp out.graphml"

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Incorrect Agruments:", helpmsg)
		os.Exit(64)
	}

	if err := Decide(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*********************************************************************/
func Decide(format string, in string, out string) error {
	switch format {
	case "sand":
		temp, err := ReadSAND(in)
		if err != nil {
			return err
		}
		if err := WriteSAND(temp, out); err != nil {
			return err
		}
	case "dm":
		temp, err := ReadDM(in)
		if err != nil {
			return err
		}
		if err := WriteDM(temp, out); err != nil {
			return err
		}
	case "idm":
		if err := ReadWriteIDM(in, out); err != nil {
			return err
		}
	case "fripp":
		temp, err := ReadFRIPP(in)
		if err != nil {
			return err
		}
		if err := WriteFRIPP(temp, out); err != nil {
			return err
		}
	default:
		return errors.New("smc: Unsupported Source Format: " + helpmsg)
	}
	return nil
}

/*********************************************************************/
/**/
func ReadSAND(filename string) (SANDTree, error) {
	st := SANDTree{}
	switch filepath.Ext(filename) {
	case ".ctrees":
		if err := st.ParseFile(filename); err != nil {
			return st, err
		}
	case ".graphml":
		if err := st.ParseFileGraphML(filename); err != nil {
			return st, err
		}
	default:
		return st, errors.New("smc: Unsupported SAND Read Format.")
	}
	return st, nil
}

/**/
func WriteSAND(st SANDTree, filename string) error {
	switch filepath.Ext(filename) {
	case ".graphml":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}

		gml, err := st.GraphML()
		if err != nil {
			return err
		}

		if err := gml.Encode(f, true); err != nil {
			return err
		}
	case ".fripp":
		return ConvertSAND(st, filename)
	case ".dependencymodel":
		return ConvertSAND(st, filename)
	default:
		return errors.New("smc: Unsupported SAND Write Format.")
	}
	return nil
}

/*********************************************************************/

/**/
func ReadDM(filename string) (DependencyModel, error) {
	dm := DependencyModel{}
	switch filepath.Ext(filename) {
	case ".graphml":
		if err := dm.ParseFileGraphML(filename); err != nil {
			return dm, err
		}
	case ".dependencymodel":
		if err := dm.ParseFile(filename); err != nil {
			return dm, err
		}
	default:
		return dm, errors.New("smc: Unsupported DM Read Format.")
	}
	return dm, nil
}

/**/
func WriteDM(dm DependencyModel, filename string) error {
	switch filepath.Ext(filename) {
	case ".graphml":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}

		gml, err := dm.GraphML()
		if err != nil {
			return err
		}

		if err := gml.Encode(f, true); err != nil {
			return err
		}
	case ".dependencymodel":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = f.WriteString(dm.String())
		if err != nil {
			return err
		}
	default:
		return errors.New("smc: Unsupported DM Write Format.")
	}

	return nil
}

/*********************************************************************/
/**/
func ReadWriteIDM(in string, out string) error {
	if filepath.Ext(in) != ".xml" {
		return errors.New("smc: Unsupported Read/Write format for iDepend.")
	}

	if filepath.Ext(out) != ".txt" {
		return errors.New("smc: Unsupported Read/Write format for iDepend.")
	}

	idm := IDependModel{}
	if err := idm.ParseFile(in, out); err != nil {
		return err
	}
	return nil
}

/*********************************************************************/

/**/
func ReadFRIPP(filename string) (FRIPP, error) {
	ir := FRIPP{}
	switch filepath.Ext(filename) {
	case ".graphml":
		if err := ir.ParseFileGraphML(filename); err != nil {
			return ir, err
		}
	case ".fripp":
		if err := ir.ParseFile(filename); err != nil {
			return ir, err
		}
	default:
		return ir, errors.New("smc: Unsupported FRIPP Read Format.")
	}
	return ir, nil
}

/**/
func WriteFRIPP(ir FRIPP, filename string) error {
	switch filepath.Ext(filename) {
	case ".fripp":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = f.WriteString(ir.String())
		if err != nil {
			return err
		}
	case ".graphml":
		f, err := os.Create(filename)
		if err != nil {
			return err
		}

		gml, err := ir.GraphML()
		if err != nil {
			return err
		}

		if err := gml.Encode(f, true); err != nil {
			return err
		}
	default:
		return errors.New("smc: Unsupported FRIPP Write Format.")
	}

	return nil
}

/*********************************************************************/
/**/
func ConvertSAND(st SANDTree, out string) error {
	/* 1. Convert to GraphML */
	gml, err := st.GraphML()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = gml.Encode(&buf, true)
	if err != nil {
		return err
	}
	/* 2. Write GraphML to tempary file */
	f, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	if _, err := f.WriteString(buf.String()); err != nil {
		return err
	}

	/* 3. Read in the temp GraphML and output the desired format. */
	if filepath.Ext(out) == ".dependencymodel" {
		dm := DependencyModel{}
		if err := dm.ParseFileGraphML(f.Name()); err != nil {
			return err
		}

		if err := WriteDM(dm, out); err != nil {
			return nil
		}
	}
	if filepath.Ext(out) == ".fripp" {
		fripp := FRIPP{}
		if err := fripp.ParseFileGraphML(f.Name()); err != nil {
			return err
		}

		if err := WriteFRIPP(fripp, out); err != nil {
			return nil
		}
	}

	return nil
}
