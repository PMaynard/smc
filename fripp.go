package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/yaricom/goGraphML/graphml"
)

/*
	-. Parse IR Native
	6. Prase IR GraphML
	7. Output IR Native
	-. Lossy-Output IR GraphML
		- Loosing ExternalReferences.
*/

type FRIPP struct {
	data *etree.Document
}

type PlaybookProcess struct {
	XmiVers               string             `xml:"xmi:version,attr"`
	XmlnsXMI              string             `xml:"xmlns:xmi,attr"`
	XMLName               xml.Name           `xml:"FRIPP:PlaybookProcess"`
	XmlnsScheme           string             `xml:"xmlns:xsi,attr"`
	XmlnsFRIPP            string             `xml:"xmlns:FRIPP,attr"`
	Name                  string             `xml:"name,attr"`
	ArtifactInStateUsed   string             `xml:"artifactInStateUsed,attr"`
	ResultArtifactInState string             `xml:"resultArtifactInState,attr"`
	ResourceUsed          string             `xml:"resourceUsed,attr"`
	RelatedReferences     string             `xml:"relatedreferences,attr"`
	Artifact              Artifact           `xml:"artifact"`
	Process               []Process          `xml:"process"`
	Resource              []Resource         `xml:"resource"`
	Externalreferences    Externalreferences `xml:"externalreferences"`
	/* TODO: Externalreferences should be a slice. */
}

func (pb *PlaybookProcess) Init() {
	pb.XmiVers = "2.0"
	pb.XmlnsXMI = "http://www.omg.org/XMI"
	pb.XmlnsScheme = "http://www.w3.org/2001/XMLSchema-instance"
	pb.XmlnsFRIPP = "http://www.example.org/FRIPP"
}

type Artifact struct {
	Name  string `xml:"name,attr"`
	State State  `xml:"state"`
}

type State struct {
	ArtifactName          string                  `xml:"artifactName,attr"`
	Name                  string                  `xml:"name,attr"`
	Artifactstateinstance []Artifactstateinstance `xml:"artifactstateinstance"`
}

type Artifactstateinstance struct {
	UsedByActivity      string `xml:"usedByActivity,attr"`
	OriginatingActivity string `xml:"originatingActivity,attr"`
}

type Externalreferences struct {
	Name string `xml:"name,attr"`
}

type Process struct {
	Type                  string   `xml:"xsi:type,attr"`
	Name                  string   `xml:"name,attr"`
	XMLName               xml.Name `xml:"process"`
	ArtifactInStateUsed   string   `xml:"artifactInStateUsed,attr"`
	ResultArtifactInState string   `xml:"resultArtifactInState,attr"`
	ResourceUsed          string   `xml:"resourceUsed,attr"`
}

type Resource struct {
	Type string `xml:"xsi:type,attr"`
	Name string `xml:"name,attr"`
}

type nindex struct {
	name string
	id   string
}

type eindex struct {
	src string
	dst string
}

/*
PlaybookProcess{
	XmiVers:               "2.0",
	XmlnsXMI:              "http://www.omg.org/XMI",
	XmlnsScheme:           "http://www.w3.org/2001/XMLSchema-instance",
	XmlnsFRIPP:            "http://www.example.org/FRIPP",
	Name:                  "Test Name",
	ArtifactInStateUsed:   "//@artifact.0/@state.0/@artifactstateinstance.2",
	ResultArtifactInState: "//@artifact.0/@state.0/@artifactstateinstance.0",
	ResourceUsed:          "//@resource.2",
	RelatedReferences:     "//@externalreferences.0",
	Artifact: Artifact{
		Name: "undefined",
		State: State{
			ArtifactName: "undefined",
			Name:         "undefined",
			Artifactstateinstance: []Artifactstateinstance{
				{UsedByActivity: "//@process.0", OriginatingActivity: "/"},
				{UsedByActivity: "//@process.1", OriginatingActivity: "//@process.0"},
				{UsedByActivity: "/", OriginatingActivity: "//@process.1"},
			},
		},
	},
	Process: []Process{
		{
			Type:                  "FRIPP:PlaybookProcess",
			Name:                  "act1",
			ArtifactInStateUsed:   "//@artifact.0/@state.0/@artifactstateinstance.0",
			ResultArtifactInState: "//@artifact.0/@state.0/@artifactstateinstance.1",
			ResourceUsed:          "//@resource.0",
		},
		{
			Type:                  "FRIPP:PlaybookProcess",
			Name:                  "act2",
			ArtifactInStateUsed:   "//@artifact.0/@state.0/@artifactstateinstance.1",
			ResultArtifactInState: "//@artifact.0/@state.0/@artifactstateinstance.2",
			ResourceUsed:          "//@resource.1",
		},
	},
	Resource: []Resource{
		{Name: "actu1", Type: "FRIPP:Actuator"},
		{Name: "actu2", Type: "FRIPP:Actuator"},
		{Name: "actu0", Type: "FRIPP:Actuator"},
	},
}
res.Externalreferences = Externalreferences{Name: "google"}

*/

/*********************************************************************/
func (ir *FRIPP) ParseFile(file string) error {
	ir.data = etree.NewDocument()
	if err := ir.data.ReadFromFile(file); err != nil {
		return err
	}
	return nil
}

/*********************************************************************/
func (ir *FRIPP) ParseFileGraphML(file string) error {
	/* 1. Parse GraphML File. */
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

	/* 2. Create the Playbook. */

	res := PlaybookProcess{}
	res.Init()

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

	res.Name = graph.Description

	/* 2.1 PlaybookProccess is store in the graph Data. */
	for _, a := range graph.Data {
		switch Src {
		case SAND:
		default:
			if a.Key == gml.GetKey("artifactInStateUsed", graphml.KeyForGraph).ID {
				res.ArtifactInStateUsed = a.Value
			}
			if a.Key == gml.GetKey("resultArtifactInState", graphml.KeyForGraph).ID {
				res.ResultArtifactInState = a.Value
			}
			if a.Key == gml.GetKey("resourceUsed", graphml.KeyForGraph).ID {
				res.ResourceUsed = a.Value
			}
			if a.Key == gml.GetKey("resourceUsed-mapped", graphml.KeyForGraph).ID {
				/* TODO: SecMof maps the resource based on the order they appear in the
				 *		 output file. In this case we need to figure
				 *		 out how to keep the ordering, when GraphML does not.
				 */
				res.Resource = append(res.Resource, Resource{Name: a.Value, Type: "FRIPP:Actuator"})
			}
			if a.Key == gml.GetKey("relatedreferences", graphml.KeyForGraph).ID {
				res.RelatedReferences = a.Value
			}
			if a.Key == gml.GetKey("externalreferences", graphml.KeyForGraph).ID {
				res.Externalreferences.Name = a.Value
			}
		}
	}

	/* 2.2 Each Node is a Proccess */
	nodeIndex := make(map[string]nindex)
	for id, node := range graph.Nodes {

		switch Src {
		case SAND:
			newProcess := Process{
				Name: node.Description,
				Type: "FRIPP:PlaybookProcess"}
			res.Process = append(res.Process, newProcess)
			nodeIndex[node.ID] = nindex{node.Description, fmt.Sprintf("%d", id)}
		default:

			/* Remove the GraphML Specific nodes. */
			if node.Description != "Start" && node.Description != "End" {

				newProcess := Process{
					Name: node.Description,
					Type: "FRIPP:PlaybookProcess"}

				for _, a := range node.Data {
					if a.Key == gml.GetKey("artifactInStateUsed", graphml.KeyForNode).ID {
						newProcess.ArtifactInStateUsed = a.Value
					}
					if a.Key == gml.GetKey("resultArtifactInState", graphml.KeyForNode).ID {
						newProcess.ResultArtifactInState = a.Value
					}
					if a.Key == gml.GetKey("resourceUsed", graphml.KeyForNode).ID {
						newProcess.ResourceUsed = a.Value
					}
				}

				res.Process = append(res.Process, newProcess)
				res.Resource = append(res.Resource, Resource{Name: node.Data[0].Value, Type: "FRIPP:Actuator"})
			}
		}
	}

	/* 2.3 Each Edge is an "artifact in state" */

	/* TODO: Parse the Names instead of defaulting to undefined. */
	res.Artifact = Artifact{Name: "undefined", State: State{ArtifactName: "undefined", Name: "undefined"}}

	edgeIndex := make(map[int]eindex)
	for _, edge := range graph.Edges {
		switch Src {
		case SAND:
			usedby := fmt.Sprintf("%s%s", "//@process.", nodeIndex[edge.Target].id)
			originating := fmt.Sprintf("%s%s", "//@process.", nodeIndex[edge.Source].id)
			res.Artifact.State.Artifactstateinstance = append(res.Artifact.State.Artifactstateinstance, Artifactstateinstance{UsedByActivity: usedby, OriginatingActivity: originating})
			edgeIndex[len(res.Artifact.State.Artifactstateinstance)-1] = eindex{nodeIndex[edge.Source].name, nodeIndex[edge.Target].name}
		default:
			res.Artifact.State.Artifactstateinstance = append(res.Artifact.State.Artifactstateinstance, Artifactstateinstance{UsedByActivity: edge.Data[0].Value, OriginatingActivity: edge.Data[1].Value})
		}
	}

	if Src == SAND {
		for id, index := range edgeIndex {
			for p := range res.Process {
				if res.Process[p].Name == index.dst {
					res.Process[p].ResultArtifactInState = fmt.Sprintf("%s%d %s", "//@artifact.0/@state.0/@artifactstateinstance.", id, res.Process[p].ResultArtifactInState)
				}

				if res.Process[p].Name == index.src {
					res.Process[p].ArtifactInStateUsed = fmt.Sprintf("%s%d %s", "//@artifact.0/@state.0/@artifactstateinstance.", id, res.Process[p].ArtifactInStateUsed)
				}
			}
		}
	}

	/* 3. Restructure into XML. */
	xmlres, err := xml.MarshalIndent(res, " ", " ")
	if err != nil {
		return err
	}
	ir.data = etree.NewDocument()
	ir.data.ReadFromBytes(xmlres)
	return nil
}

/*********************************************************************/
func (ir *FRIPP) GraphML() (*graphml.GraphML, error) {
	gml := graphml.NewGraphML("")
	attributes := make(map[string]interface{})

	/* Look up the PlaybookProcess's Name */
	playbookName := ""
	for _, a := range ir.data.FindElement("/PlaybookProcess").Attr {
		if a.Key == "name" {
			playbookName = a.Value
		}

		if a.Key == "artifactInStateUsed" {
			attributes["artifactInStateUsed"] = a.Value
		}

		if a.Key == "resultArtifactInState" {
			attributes["resultArtifactInState"] = a.Value
		}
		if a.Key == "resourceUsed" {
			attributes["resourceUsed"] = a.Value
			attributes["resourceUsed-mapped"] = ir.MapPath("/PlaybookProcess/resource", "name", ExtractInt(a.Value))
		}
		if a.Key == "relatedreferences" {
			attributes["relatedreferences"] = a.Value
		}
	}

	for _, a := range ir.data.FindElement("/PlaybookProcess/externalreferences").Attr {
		if a.Key == "name" {
			attributes["externalreferences"] = a.Value
		}
	}

	graph, err := gml.AddGraph(playbookName, graphml.EdgeDirectionDirected, attributes)
	if err != nil {
		return nil, err
	}

	/* Nodes */
	nodeIndex := make(map[int]*graphml.Node)
	for id, node := range ir.data.FindElements("/PlaybookProcess/process") {
		name := ""
		attributes := make(map[string]interface{})

		for _, a := range node.Attr {
			if a.Key == "name" {
				name = a.Value
			}
			if a.Key == "artifactInStateUsed" {
				attributes["artifactInStateUsed"] = a.Value
			}
			if a.Key == "resultArtifactInState" {
				attributes["resultArtifactInState"] = a.Value
			}
			if a.Key == "resourceUsed" {
				/* TODO: Rename to align with ProcessPlaybook.*/
				attributes["ActuatorName"] = ir.MapPath("/PlaybookProcess/resource", "name", ExtractInt(a.Value))
				attributes["resourceUsed"] = a.Value
			}
		}

		nodeIndex[id], err = graph.AddNode(attributes, name)
		if err != nil {
			return nil, err
		}
	}

	/*
	* Add the fake start and end nodes at the end so we can
	* still map the nodes to the correct FRIPP IDs
	 */

	/*
	* TODO: Maybe auto gen the start and end from the root
	* 		PlaybookProcess like we do for edges
	 */
	attributes = make(map[string]interface{})
	start, err := graph.AddNode(attributes, "Start")
	if err != nil {
		return nil, err
	}
	end, err := graph.AddNode(attributes, "End")
	if err != nil {
		return nil, err
	}

	/* Edges */
	id := 0
	for _, node := range ir.data.FindElements("/") {
		for _, a := range node.Attr {
			if a.Key == "resultArtifactInState" {
				n := strings.Split(a.Value, ".")
				lookup := n[len(n)-1:][0]

				src := ir.MapPath("/PlaybookProcess/artifact/state/artifactstateinstance", "originatingActivity", Str2Int(lookup))
				dst := ir.MapPath("/PlaybookProcess/artifact/state/artifactstateinstance", "usedByActivity", Str2Int(lookup))

				if src == "" || dst == "" {
					return nil, err
				}

				/* Set default to Start/End */
				n1 := start
				n2 := end

				if src != "/" {
					n1 = nodeIndex[ExtractInt(src)]
				}

				if dst != "/" {
					n2 = nodeIndex[ExtractInt(dst)]
				}

				attributes = make(map[string]interface{})
				attributes["usedByActivity"] = src
				attributes["originatingActivity"] = dst

				_, err = graph.AddEdge(n1, n2, attributes, graphml.EdgeDirectionDefault, "")
				if err != nil {
					return nil, err
				}
				id = id + 1
			}
		}
	}

	return gml, nil
}

/*********************************************************************/
func (ir FRIPP) String() string {
	str, err := ir.data.WriteToString()
	if err != nil {
		log.Fatal(err)
	}

	return str
}

/*********************************************************************/
func Str2Int(from string) int {
	n, err := strconv.Atoi(from)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

/*********************************************************************/
func ExtractInt(from string) int {
	return Str2Int(strings.Split(from, ".")[1])
}

/*********************************************************************/
func (ir *FRIPP) MapPath(path string, attribute string, id int) string {
	res := ir.data.FindElements(path)
	for _, aa := range res[id].Attr {
		if aa.Key == attribute {
			return aa.Value
		}
	}
	return ""
}
