package main

import (
	"encoding/xml"
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

type Paragon struct {
	XMLName     xml.Name `xml:"paragon"`
	Description string   `xml:"description,attr" json:" "`
	Probability string   `xml:"probability,attr"`
	Type        string   `xml:"type,attr,omitempty"`
	Paragons    []Paragon
}

type DependencyModel struct {
	XMLName  xml.Name `xml:"dependencyModel:Paragon"`
	XmiVers  string   `xml:"xmi:version,attr"`
	XmlnsXMI string   `xml:"xmlns:xmi,attr"`
	XmlnsDM  string   `xml:"xmlns:dependencyModel,attr"`
	Paragon  `xml:"paragon,attr" json:"Paragon"`
	Paragons []Paragon
}

/*********************************************************************/
func (dm DependencyModel) Init() DependencyModel {
	return DependencyModel{
		XmiVers:  "2.0",
		XmlnsXMI: "http://www.omg.org/XMI",
		XmlnsDM:  "http://www.example.org/dependencyModel",
	}
}

/*********************************************************************/
// func (p Paragon) String() string {
// 	return fmt.Sprintf("Paragon Description=%v, Probability=%v, Type=%v", p.Description, p.Probability, p.Type)
// }

/*********************************************************************/
/*********************************************************************/

/* Parsing DM from Text File */

// type InternelParse struct {
// 	id          int
// 	description string
// 	indent      int
// }

// func ParseTextDM(file string) DependencyModel {
// 	rawdata, err := os.ReadFile(file)
// 	if err != nil {
// 		panic(err)
// 	}

// 	/* TODO: Remove comments */

// 	/**/
// 	parsed := []InternelParse{}
// 	for id, l := range strings.Split(string(rawdata), "\n") {
// 		parsed = append(parsed, InternelParse{description: strings.TrimSpace(l), id: id, indent: strings.Count(l, "\t")})
// 	}

// 	/**/
// 	var data DependencyModel
// 	data = data.Init()
// 	/* Add root */
// 	data.Paragon = Paragon{Description: parsed[0].description, Probability: "1.0", Type: "AND"}

// 	for i := range parsed {
// 		fmt.Println(parsed[i])

// 		// if parsed[i+1].indent < parsed[i].indent {
// 		// 	fmt.Println(parsed[i+1].description, "under", parsed[i].description)
// 		// }

// 		/* TODO: Depth frist search for all children*/
// 	}

// 	return data

// }
