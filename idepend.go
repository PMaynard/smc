package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strings"
)

type IDependModel struct {
	ID       int           `xml:"id"`
	Name     string        `xml:"name"`
	Desc     string        `xml:"description"`
	OrgID    string        `xml:"organisation_id"`
	Checksum string        `xml:"checksum"`
	Entities []IDependNode `xml:"entities>entity"`
}

type IDependNode struct {
	ID           int    `xml:"id"`
	Name         string `xml:"name"`
	OrgID        string `xml:"organisation_id"`
	Dependencies []int  `xml:"dependencies>id"`
}

/*********************************************************************/
func (dm *IDependModel) ParseFile(in string, out string) error {
	rawdata, err := os.ReadFile(in)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(rawdata, &dm); err != nil {
		return err
	}

	/* Sort the entities based on their ID. */
	sort.Slice(dm.Entities, func(i, j int) bool {
		return dm.Entities[i].ID < dm.Entities[j].ID
	})

	/* Output file */
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	/* assume the root node is the one with the lowest ID */
	rootNode := dm.Entities[0]
	dm.WalkPrint(w, rootNode.Name, rootNode.Dependencies)
	return nil
}

/*********************************************************************/
func (dm *IDependModel) WalkPrint(w *bufio.Writer, prefix string, dependencies []int) {
	for _, d := range dependencies {
		dependency := dm.GetByID(d)

		name := strings.Replace(dependency.Name, "/", "\\", -1)
		name = strings.Replace(name, ",", "", -1)

		_, err := w.WriteString(fmt.Sprintf("%s/%s\n", prefix, name))
		if err != nil {
			return
		}

		if len(dependency.Dependencies) != 0 {
			dm.WalkPrint(w, prefix+"/"+dependency.Name, dependency.Dependencies)
		}
	}
}

/*********************************************************************/
func (dm *IDependModel) GetByID(id int) IDependNode {
	for i := range dm.Entities {
		if dm.Entities[i].ID == id {
			return dm.Entities[i]
		}
	}
	return IDependNode{}
}
