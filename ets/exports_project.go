// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"encoding/xml"
	"fmt"
	"io"
)

// ProjectID is a project identifier.
type ProjectID string

// ProjectInfo is a project.
type ProjectInfo struct {
	ID   ProjectID
	Name string
}

// UnmarshalXML implements xml.Unmarshaler.
func (proj *ProjectInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Find 'xmlns' attribute.
	var namespace string
	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" {
			namespace = attr.Value
			break
		}
	}

	// Decide which schema to use based on the value of the 'xmlns' attribute.
	switch namespace {
	case schema11Namespace:
		return unmarshalProjectInfo11(d, start, proj)

	case schema13Namespace:
		return unmarshalProjectInfo13(d, start, proj)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", namespace)
	}
}

// DecodeProjectInfo parses the contents of project file.
func DecodeProjectInfo(r io.Reader) (*ProjectInfo, error) {
	proj := &ProjectInfo{}
	if err := xml.NewDecoder(r).Decode(proj); err != nil {
		return nil, err
	}

	return proj, nil
}
