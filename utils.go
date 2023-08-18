package main

import (
	"encoding/xml"
	"os"
)

/*********************************************************************/
func GenXML(input any) string {
	out, err := xml.MarshalIndent(input, " ", "  ")
	if err != nil {
		panic(err)
	}
	return xml.Header + string(out)
}

/*********************************************************************/
func Output(content string, file string) error {
	err := os.WriteFile(file, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
