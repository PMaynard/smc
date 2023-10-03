package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yaricom/goGraphML/graphml"
)

const (
	DESC = 0
	OPER = 1
)

type SANDTree struct {
	Nodes    []SANDNode
	Filename string
	Name     string
}

type SANDNode struct {
	ID     int
	Desc   string
	Oper   string
	Child  []int
	Indent int
	Parent int
}

/*********************************************************************/
func (st *SANDTree) ParseFile(file string) error {

	/* 1. Parse the whole file and create inital tree */
	st.Filename = file
	st.Name = strings.TrimSuffix(filepath.Base(st.Filename), filepath.Ext(st.Filename)) + " SAND"

	rawdata, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	/* For each line in the file create a Node. */
	for n, line := range strings.Split(string(rawdata), "\n") {
		res := SANDNode{}

		/* NOTE: IDs are based on their line number of the source file. */
		res.ID = n
		res.Parent = -1

		tmp := strings.Split(line, ":")
		res.Indent = strings.Count(tmp[DESC], "\t")
		res.Desc = strings.TrimSpace(tmp[DESC])

		/* If missing, the default operator is OR. */
		if len(tmp) != 2 {
			res.Oper = "OR"
		} else {
			res.Oper = strings.ToUpper(strings.TrimSpace(tmp[OPER]))
		}

		/* Don't add empty lines or comments 'i.e. starts with #' */
		if len(res.Desc) == 0 || strings.HasPrefix(res.Desc, "#") {
			continue
		}

		st.Nodes = append(st.Nodes, res)
	}

	/* 2. Expand Tree: Identify Parents and children. */
	prefix := []string{st.Nodes[0].Desc}
	for i, node := range st.Nodes {

		/* 2.1. Figure out the order. */
		/* For every node that's not root */
		if i != 0 {
			/* Prefix */
			cni := node.Indent          /* Curent Node Indent */
			pni := st.Nodes[i-1].Indent /* Past Node Indent */
			pl := len(prefix)           /* Prefix Length */

			// fmt.Println("\n", cni, pni, pl, cni, node.Desc)

			/*
			 * If Past node indent less than the current node indent,
			 * and the prefix length is the same as the current node
			 * indent. Then we need to append the prefix.
			 */
			if pni < cni && pl == cni {
				prefix = append(prefix, node.Desc)
			}

			/*
			 * If the prefix length is bigger than the current node
			 * indent, then we need to remove the diference in current
			 * node indent and the prefix length, then we need to
			 * append the prefix.
			 */
			if pl > cni {
				prefix = prefix[:len(prefix)-(pl-cni)]
				prefix = append(prefix, node.Desc)
			}
		}

		/* Debug prefix output */
		/* Usefull for Treemaps */
		// for i := range prefix {
		// 	fmt.Printf("%s/", prefix[i])
		// }
		// fmt.Printf("\n")

		/* 2.2. Assign Parents */
		if len(prefix) >= 2 {
			// fmt.Printf("%d %s parent %s \n", i, prefix, prefix[len(prefix)-2])
			st.Nodes[i].Parent = st.GetByName(prefix[len(prefix)-2]).ID
		}
	}

	/* 2.3. Assign Children */
	for _, child := range st.Nodes {
		if child.Parent != -1 {
			for ii, parent := range st.Nodes {
				if child.Parent == parent.ID {
					st.Nodes[ii].Child = append(parent.Child, child.ID)
				}
			}
		}
	}
	return nil
}

/*********************************************************************/
func (st *SANDTree) ParseFileGraphML(file string) error {

	/* 1. Parse the whole file and create inital tree */
	st.Filename = file
	st.Name = strings.TrimSuffix(filepath.Base(st.Filename), filepath.Ext(st.Filename)) + " SAND"

	rawdata, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rawdata.Close()

	/* GraphML */
	gml := graphml.NewGraphML(st.Name)
	err = gml.Decode(rawdata)
	if err != nil {
		return err
	}

	/* Make sure there is only one graph */
	if len(gml.Graphs) != 1 {
		return fmt.Errorf("%d graphs found. Needs 1.", len(gml.Graphs))
	}

	graph := gml.Graphs[0]

	/* 2. Get nodes */
	nodeIndex := make(map[string]int)
	for i, node := range graph.Nodes {

		if len(node.Data) != 1 {
			return fmt.Errorf("Node %s '%s' has more data '%d' than I wanted 1.", node.ID, node.Description, len(node.Data))
		}

		res := SANDNode{}
		res.ID = i
		res.Parent = -1
		res.Indent = -1
		res.Desc = node.Description
		res.Oper = node.Data[0].Value
		st.Nodes = append(st.Nodes, res)
		nodeIndex[node.ID] = i
	}

	/* 3. Get Edges */
	for _, edge := range graph.Edges {
		for i := range st.Nodes {
			/* Add Children */
			if nodeIndex[edge.Source] == i {
				st.Nodes[i].Child = append(st.Nodes[i].Child, nodeIndex[edge.Target])
			}

			/* Add Parent */
			if nodeIndex[edge.Target] == i {
				st.Nodes[i].Parent = nodeIndex[edge.Source]
			}
		}
	}
	return nil
}

/*********************************************************************/
func (st *SANDTree) GraphML() (*graphml.GraphML, error) {

	/* TODO: Make sure st is populated */

	/* GraphML */
	gml := graphml.NewGraphML(st.Name)

	/* Graph */
	attributes := make(map[string]interface{})
	attributes["SrcFormat"] = "SAND"

	graph, err := gml.AddGraph(st.Name, graphml.EdgeDirectionDirected, attributes)
	if err != nil {
		return nil, err
	}

	/* Add Nodes */
	nodeIndex := make(map[int]*graphml.Node)
	for i := range st.Nodes {
		attributes := make(map[string]interface{})
		attributes["Operator"] = st.Nodes[i].Oper

		node, err := graph.AddNode(attributes, st.Nodes[i].Desc)
		if err != nil {
			return nil, err
		}

		/* Keep new nodes mapping to old IDs */
		nodeIndex[st.Nodes[i].ID] = node
	}

	/* Add Edges */
	for _, node := range st.Nodes {
		if node.Parent != -1 {
			n1 := nodeIndex[node.Parent]
			n2 := nodeIndex[node.ID]

			attributes = make(map[string]interface{})
			attributes["Operator"] = st.Nodes[node.Parent].Oper

			_, err = graph.AddEdge(n1, n2, attributes, graphml.EdgeDirectionDefault, st.Nodes[node.Parent].Oper)
			if err != nil {
				return nil, err
			}
		}
	}
	return gml, nil
}

/*********************************************************************/
func (st *SANDTree) GetByName(desc string) SANDNode {
	for i := range st.Nodes {
		if st.Nodes[i].Desc == desc {
			return st.Nodes[i]
		}
	}
	return SANDNode{}
}
