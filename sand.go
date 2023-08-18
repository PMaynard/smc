package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	DESC = 0
	OPER = 1
)

type SANDTree struct {
	Nodes []SANDNode
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
func (st *SANDTree) GetByName(desc string) SANDNode {
	for i := range st.Nodes {
		if st.Nodes[i].Desc == desc {
			return st.Nodes[i]
		}
	}
	return SANDNode{}
}

/*********************************************************************/
func (st *SANDTree) GetByID(id int) SANDNode {
	for i := range st.Nodes {
		if st.Nodes[i].ID == id {
			return st.Nodes[i]
		}
	}
	return SANDNode{}
}

/*********************************************************************/
func (st *SANDTree) Print() {
	for _, n := range st.Nodes {
		fmt.Printf("%d:\"%s(%s)\" Parent(%s) Child(%d)\n", n.ID, n.Desc, n.Oper, st.GetByID(n.Parent).Desc, n.Child)
	}
}
