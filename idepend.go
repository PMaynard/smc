package main

import (
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
func (dm *IDependModel) ParseFile(file string) {
	rawdata, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := xml.Unmarshal(rawdata, &dm); err != nil {
		panic(err)
	}

	/* Sort the entities based on their ID. */
	sort.Slice(dm.Entities, func(i, j int) bool {
		return dm.Entities[i].ID < dm.Entities[j].ID
	})

	/* assume the root node is the one with the lowest ID */
	rootNode := dm.Entities[0]
	dm.WalkPrint(rootNode.Name, rootNode.Dependencies)
}

/*********************************************************************/
func (dm *IDependModel) GetByName(name string) IDependNode {
	for i := range dm.Entities {
		if dm.Entities[i].Name == name {
			return dm.Entities[i]
		}
	}
	return IDependNode{}
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

/*********************************************************************/
func (dm *IDependModel) WalkPrint(prefix string, dependencies []int) {
	for _, d := range dependencies {
		dependency := dm.GetByID(d)

		name := strings.Replace(dependency.Name, "/", "\\", -1)
		name = strings.Replace(name, ",", "", -1)

		fmt.Printf("%s/%s\n", prefix, name)

		if len(dependency.Dependencies) != 0 {
			dm.WalkPrint(prefix+"/"+dependency.Name, dependency.Dependencies)
		}
	}
}
