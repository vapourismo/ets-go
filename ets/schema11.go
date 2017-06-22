// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import "encoding/xml"

func unmarshalProjectInfo11(d *xml.Decoder, start xml.StartElement, pi *ProjectInfo) error {
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

	pi.ID = ProjectID(doc.Project.ID)
	pi.Name = doc.Project.ProjectInformation.Name

	return nil
}

type deviceInstance11 DeviceInstance

func (di *deviceInstance11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID      string `xml:"Id,attr"`
		Name    string `xml:"Name,attr"`
		Address uint   `xml:"Address,attr"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	di.ID = DeviceInstanceID(doc.ID)
	di.Name = doc.Name
	di.Address = doc.Address

	return nil
}

type line11 Line

func (l *line11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID             string `xml:"Id,attr"`
		Name           string `xml:"Name,attr"`
		Address        uint   `xml:"Address,attr"`
		DeviceInstance []deviceInstance11
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	l.ID = LineID(doc.ID)
	l.Name = doc.Name
	l.Address = doc.Address
	l.DevicesInstances = make([]DeviceInstance, len(doc.DeviceInstance))

	for n, docDeviceInstance := range doc.DeviceInstance {
		l.DevicesInstances[n] = DeviceInstance(docDeviceInstance)
	}

	return nil
}

type area11 Area

func (a *area11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID      string `xml:"Id,attr"`
		Name    string `xml:"Name,attr"`
		Address uint   `xml:"Address,attr"`
		Line    []line11
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	a.ID = AreaID(doc.ID)
	a.Name = doc.Name
	a.Address = doc.Address
	a.Lines = make([]Line, len(doc.Line))

	for n, docLine := range doc.Line {
		a.Lines[n] = Line(docLine)
	}

	return nil
}

type installation11 Installation

func (i *installation11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Name     string   `xml:"Name,attr"`
		Topology []area11 `xml:"Topology>Area"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	i.Name = doc.Name
	i.Topology = make([]Area, len(doc.Topology))

	for n, docArea := range doc.Topology {
		i.Topology[n] = Area(docArea)
	}

	return nil
}

func unmarshalProject11(d *xml.Decoder, start xml.StartElement, p *Project) error {
	var doc struct {
		Project struct {
			ID            string           `xml:"Id,attr"`
			Installations []installation11 `xml:"Installations>Installation"`
		}
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	p.ID = ProjectID(doc.Project.ID)
	p.Installations = make([]Installation, len(doc.Project.Installations))

	for i, docInst := range doc.Project.Installations {
		p.Installations[i] = Installation(docInst)
	}

	return nil
}
