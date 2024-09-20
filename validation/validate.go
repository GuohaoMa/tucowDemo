package validation

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/GuohaoMa/tucowDemo/model"
)

func Validate(filePath string) error {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening XML file: %v", err)
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)

	nodeCount := 0
	inNodes, inEdges := false, false
	edgeElementFound := false
	nodeIdMap := make(map[string]string)
	fromInEdgeCount, inInEdgeCount := 0, 0
	idInGraphCount, nameInGraphCount := 0, 0
	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Error decoding XML: %v", err)
		}

		switch elem := t.(type) {
		case xml.StartElement:
			if elem.Name.Local == "nodes" {
				inNodes = true
				if edgeElementFound {
					return errors.New("The <nodes> group must come before the <edges> group.")
				}
			}
			if elem.Name.Local == "edges" {
				inEdges = true
				edgeElementFound = true
			}
			if elem.Name.Local == "node" {
				if inNodes == true && inEdges == false {
					nodeCount += 1
					var node model.Node
					decoder.DecodeElement(&node, &elem)
					if _, ok := nodeIdMap[node.Identity]; ok {
						return errors.New("All nodes must have different <id> tags.")
					} else {
						nodeIdMap[node.Identity] = node.Name
					}
				}
				if inEdges == true && inNodes == false {
					if nodeCount == 0 {
						return errors.New("There must be at least one <node> in the <nodes> group")
					}
					var edge model.Edge
					decoder.DecodeElement(&edge, &elem)
					if edge.FromIdentity == "" {
						return errors.New("For every <edge>, there must be a single <from> tag")
					}
					if edge.ToIdentity == "" {
						return errors.New("For every <edge>, there must be a single <to> tag")
					}
					if _, ok := nodeIdMap[edge.FromIdentity]; !ok {
						return errors.New("From node of an edge must be predefined.")
					}
					if _, ok := nodeIdMap[edge.ToIdentity]; !ok {
						return errors.New("To node of an edge must be predefined.")
					}
					if edge.Cost < 0 {
						return errors.New("Cost of an edge must be non-negative.")
					}
				}
			}
			if elem.Name.Local == "from" {
				if inEdges == true {
					if fromInEdgeCount == 0 {
						fromInEdgeCount += 1
					} else {
						return errors.New("For every <edge>, there must be a single <from> tag")
					}
				}
			}
			if elem.Name.Local == "to" {
				if inEdges == true {
					if inInEdgeCount == 0 {
						inInEdgeCount += 1
					} else {
						return errors.New("For every <edge>, there must be a single <to> tag")
					}
				}
			}
			if elem.Name.Local == "id" {
				if inNodes == false && inEdges == false {
					if idInGraphCount == 0 {
						idInGraphCount += 1
					}
				}
			}
			if elem.Name.Local == "name" {
				if inNodes == false && inEdges == false {
					if nameInGraphCount == 0 {
						nameInGraphCount += 1
					}
				}
			}
		case xml.EndElement:
			if elem.Name.Local == "nodes" {
				inNodes = false
			}
			if elem.Name.Local == "edges" {
				inEdges = false
			}
			if elem.Name.Local == "node" && inEdges == true {
				fromInEdgeCount = 0
				inInEdgeCount = 0
			}
		}
	}

	if nodeCount == 0 {
		return errors.New("There must be at least one <node> in the <nodes> group")
	}
	if idInGraphCount == 0 {
		return errors.New("There must be an <id> in the <graph>")
	}
	if nameInGraphCount == 0 {
		return errors.New("There must be an <name> in the <graph>")
	}
	return nil
}
