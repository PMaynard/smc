package main

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	// fxml "github.com/m29h/xml"
	"fmt"
	"os"

	"github.com/yaricom/goGraphML/graphml"
)

// test := &DependencyModel{
// 	XmiVers: "2.0",
// 	XmlnsXMI: "http://www.omg.org/XMI",
// 	XmlnsDM: "http://www.example.org/dependencyModel",
// 	Paragon : Paragon{Description: "Company OK", Probability: "1.0", Type: "AND"},
// 	Paragons: []Paragon{
// 		{Description: "Personnel OK", Probability: "1.0", Type: "AND", Paragons: []Paragon{
// 			{Description: "Accepting Orders OK", Probability: "1.0", Type: "AND"},
// 		}},
// 		{Description: "Postage OK", Probability: "1.0", Type: "AND", Paragons: []Paragon{
// 			{Description: "Processing Orders OK", Probability: "1.0", Type: "AND"},
// 		}},
// 		{Description: "Services OK", Probability: "1.0", Type: "AND", Paragons: []Paragon{
// 			{Description: "AWS OK", Probability: "1.0", Type: "AND", Paragons: []Paragon{
// 				{Description: "HTTP OK", Probability: "1.0", Type: "AND"},
// 				{Description: "SMTP OK", Probability: "1.0", Type: "AND"},
// 			}},
// 			{Description: "Email OK", Probability: "1.0", Type: "AND", Paragons: []Paragon{
// 				{Description: "IMAP OK", Probability: "1.0", Type: "AND"},
// 			}},
// 			{Description: "Spreadsheet OK", Probability: "1.0", Type: "AND"},
// 		}},
// 	},
// }

type DependencyModel struct {
	XMLName xml.Name `xml:"http://www.example.org/dependencyModel Paragon"`
	// XMLName xml.Name `xml:"dependencyModel:Paragon"` // go1.22+ with patch for namesapce prefixes
	XmlnsXMI string `xml:"xmlns:xmi,attr"`
	XmiVers  string `xml:"xmi:version,attr"`
	XmlnsDM  string `xml:"xmlns:dependencyModel,attr"`
	Paragon  `xml:"paragon,attr"`
	Paragons []Paragon `xml:"paragon"`
	filename string
	name     string
}

type Paragon struct {
	XMLName     xml.Name  `xml:"paragon"`
	Description string    `xml:"description,attr"`
	Probability string    `xml:"probability,attr"`
	Type        string    `xml:"Type,attr,omitempty"`
	Paragons    []Paragon `xml:"paragon"`
}

/*********************************************************************/
func (dm *DependencyModel) Init(file string) {
	dm.XmiVers = "2.0"
	dm.XmlnsXMI = "http://www.omg.org/XMI"
	dm.XmlnsDM = "http://www.example.org/dependencyModel"
	dm.filename = file
	dm.name = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) + " DM"
}

/*********************************************************************/
func (dm *DependencyModel) ParseFile(file string) error {

	/* TODO: Read this data from the XML file */
	dm.Init(file)

	/* 1. Parse the file and create tree */
	rawdata, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	if err := xml.Unmarshal(rawdata, &dm); err != nil {
		return err
	}
	return nil
}

const (
	DM int = iota
	SAND
	Unknown
)

/*********************************************************************/
func (dm *DependencyModel) ParseFileGraphML(file string) error {
	/* 1. Parse the whole file and create inital tree */
	rawdata, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rawdata.Close()

	/* GraphML */
	gml := graphml.NewGraphML("")
	err = gml.Decode(rawdata)
	if err != nil {
		return err
	}

	/* Make sure there is only one graph */
	if len(gml.Graphs) != 1 {
		return fmt.Errorf("%d graphs found. Needs 1.", len(gml.Graphs))
	}

	graph := gml.Graphs[0]

	/* Check if the SrcFormat is defined. */
	Src := Unknown
	if gml.GetKey("SrcFormat", graphml.KeyForGraph) != nil {
		for _, a := range graph.Data {
			if a.Key == gml.GetKey("SrcFormat", graphml.KeyForGraph).ID {
				if a.Value == "SAND" {
					Src = SAND
				}
			}
		}
	}

	/* 2. Get nodes */
	nodeIndex := make(map[string]Paragon)
	for _, node := range graph.Nodes {
		switch Src {
		case SAND:
			/* TODO: DM does not have the notion of Sequental AND. */
			operator := "AND"
			if gml.GetKey("Operator", graphml.KeyForNode) != nil {
				for _, a := range node.Data {
					if a.Key == gml.GetKey("Operator", graphml.KeyForNode).ID {
						if a.Value == "OR" {
							operator = "OR"
						}
					}
				}
			}

			nodeIndex[node.ID] = Paragon{Description: node.Description, Probability: "1", Type: operator}
		default:
			if len(node.Data) != 2 {
				return fmt.Errorf("Node %s '%s' has more data '%d' than I wanted 2.", node.ID, node.Description, len(node.Data))
			}
			nodeIndex[node.ID] = Paragon{Description: node.Description, Probability: node.Data[0].Value, Type: node.Data[1].Value}
		}
	}

	/* 3. Map the Edges */

	/* Assume the first node is the root node. */

	edges := make(map[string]*graphml.Edge)
	for _, edge := range graph.Edges {
		edges[edge.ID] = edge
	}

	/* Intialise the Dependency Model */
	dm.Init(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)) + " DM")

	/* Create the root node, and its children. */
	dm.Paragon = nodeIndex[graph.Nodes[0].ID]
	for _, edge := range edges {
		if nodeIndex[edge.Source].Description == dm.Paragon.Description {
			dm.Paragons = append(dm.Paragons, nodeIndex[edge.Target])
			delete(edges, edge.ID)
		}
	}

	/* While there are no edges left, for each child, call addWalk. */
	for len(edges) != 0 {
		for i := range dm.Paragons {
			dm.addWalk(&dm.Paragons[i], edges, nodeIndex)
		}
	}

	return nil
}

/*********************************************************************/
func (dm *DependencyModel) addWalk(parent *Paragon, edges map[string]*graphml.Edge, nodeIndex map[string]Paragon) {

	/* for each edge see if it has an edge source */
	for _, edge := range edges {
		if nodeIndex[edge.Source].Description == parent.Description {
			/* if there is an edge source, add the child */
			parent.Paragons = append(parent.Paragons, nodeIndex[edge.Target])

			/* remove the edge from the abanonded child list */
			delete(edges, edge.ID)
		}
	}

	/* for all the children see if they have children. */
	for i := range parent.Paragons {
		dm.addWalk(&parent.Paragons[i], edges, nodeIndex)
	}

}

/*********************************************************************/
func (dm *DependencyModel) GraphML() (*graphml.GraphML, error) {
	gml := graphml.NewGraphML(dm.name)
	attributes := make(map[string]interface{})

	graph, err := gml.AddGraph(dm.name, graphml.EdgeDirectionDirected, attributes)
	if err != nil {
		return gml, err
	}

	nodeIndex := make(map[string]*graphml.Node)
	for _, edges := range dm.GetInternal() {

		/* Add Nodes */
		node := edges[len(edges)-1]

		attributes := make(map[string]interface{})
		attributes["Probability"] = node.Probability
		attributes["Type"] = node.Type

		nodeIndex[node.Description], err = graph.AddNode(attributes, node.Description)
		if err != nil {
			return gml, err
		}

		/*Add Edges */
		if len(edges) > 1 {
			parent := edges[len(edges)-2]

			n1 := nodeIndex[parent.Description]
			n2 := nodeIndex[node.Description]

			attributes = make(map[string]interface{})
			attributes["Type"] = parent.Type

			_, err = graph.AddEdge(n1, n2, attributes, graphml.EdgeDirectionDefault, parent.Type)
			if err != nil {
				return gml, err
			}
		}
	}
	return gml, nil
}

/*********************************************************************/
func (dm DependencyModel) GetInternal() [][]Paragon {
	res := [][]Paragon{}
	res = append(res, []Paragon{dm.Paragon})
	dm.interalWalk([]Paragon{dm.Paragon}, dm.Paragons, &res)

	return res
}

/*********************************************************************/
func (dm DependencyModel) interalWalk(prefix []Paragon, children []Paragon, res *[][]Paragon) {
	for _, child := range children {

		line := []Paragon{}
		for _, p := range prefix {
			line = append(line, p)
		}
		line = append(line, child)
		// fmt.Println("-", line)
		*res = append(*res, line)

		if len(child.Paragons) != 0 {
			dm.interalWalk(append(prefix, child), child.Paragons, res)
		}
	}
}

/*********************************************************************/
func (dm DependencyModel) String() string {
	out, err := xml.MarshalIndent(dm, " ", "  ")
	if err != nil {
		panic(err)
	}
	return xml.Header + string(out)
}
