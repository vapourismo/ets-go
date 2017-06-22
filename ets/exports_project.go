// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"encoding/xml"
	"fmt"
)

// ProjectID is a project identifier.
type ProjectID string

// Project is a project.
type Project struct {
	ID   ProjectID
	Name string
}

// UnmarshalXML implements xml.Unmarshaler.
func (proj *Project) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
		return unmarshalProject11(d, start, proj)

	case schema13Namespace:
		return unmarshalProject13(d, start, proj)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", namespace)
	}
}
