package main

import (
	"fmt"
	"log"
)

func main() {
	// IDepend2HashTree()
	SAND2FRIPPDM()
}

func IDepend2HashTree() {
	idm := IDependModel{}
	idm.ParseFile("data/dependencymodels/idepend/SCADA-DM-XML.xml")
}

func SAND2FRIPPDM() {
	st := SANDTree{}
	if err := st.ParseFile("data/ctrees/example.ctrees"); err != nil {
		log.Fatal(err)
	}

	var data DependencyModel
	data = data.Init()
	/* Add root */
	data.Paragon = Paragon{Description: st.Nodes[0].Desc, Probability: "1.0", Type: st.Nodes[0].Oper}

	// TODO: Implement tree walk method to create "data.Paragons" i.e. the children of the root and any of their children.

	// Output(GenXML(data), "test.xml")
	fmt.Println(GenXML(data))
}

// /*********************************************************************/
// func (dm *IDependModel) WalkPrint(prefix string, dependencies []int) {
// 	for _, d := range dependencies {
// 		dependency := dm.GetByID(d)

// 		name := strings.Replace(dependency.Name, "/", "\\", -1)
// 		name = strings.Replace(name, ",", "", -1)

// 		fmt.Printf("%s/%s\n", prefix, name)

// 		if len(dependency.Dependencies) != 0 {
// 			dm.WalkPrint(prefix+"/"+dependency.Name, dependency.Dependencies)
// 		}
// 	}
// }

/* TODO: Verify the XML output by reading it back in. */
// var dmtemp DependencyModel
// if err := xml.Unmarshal(out, &dmtemp); err != nil {
// 	panic(err)
// }
// fmt.Println(dmtemp)
