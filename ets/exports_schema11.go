// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import "encoding/xml"

// unmarshalProject11 extracts project information from the current element.
func unmarshalProject11(d *xml.Decoder, start xml.StartElement, proj *Project) error {
	var doc struct {
		Project struct {
			ID                 string `xml:"Id,attr"`
			ProjectInformation struct {
				Name string `xml:"Name,attr"`
			}
		}
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	proj.ID = ProjectID(doc.Project.ID)
	proj.Name = doc.Project.ProjectInformation.Name

	return nil
}
