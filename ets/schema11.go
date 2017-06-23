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
		ID         string `xml:"Id,attr"`
		Name       string `xml:"Name,attr"`
		Address    uint   `xml:"Address,attr"`
		ComObjects []struct {
			RefID         string `xml:"RefId,attr"`
			DatapointType string `xml:"DatapointType,attr"`
			Connectors    struct {
				Elements []struct {
					XMLName xml.Name
					RefID   string `xml:"GroupAddressRefId,attr"`
				} `xml:",any"`
			}
		} `xml:"ComObjectInstanceRefs>ComObjectInstanceRef"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	di.ID = DeviceInstanceID(doc.ID)
	di.Name = doc.Name
	di.Address = doc.Address
	di.ComObjects = make([]ComObjectInstanceRef, len(doc.ComObjects))

	for n, docComObj := range doc.ComObjects {
		comObj := ComObjectInstanceRef{
			RefID:         ComObjectRefID(docComObj.RefID),
			DatapointType: docComObj.DatapointType,
			Connectors:    make([]Connector, len(docComObj.Connectors.Elements)),
		}

		for m, docConnElem := range docComObj.Connectors.Elements {
			comObj.Connectors[m] = Connector{
				Receive: docConnElem.XMLName.Local == "Receive",
				RefID:   GroupAddressID(docConnElem.RefID),
			}
		}

		di.ComObjects[n] = comObj

	}

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

type groupRange11 GroupRange

func (gar *groupRange11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID           string `xml:"Id,attr"`
		Name         string `xml:"Name,attr"`
		RangeStart   uint   `xml:"RangeStart,attr"`
		RangeEnd     uint   `xml:"RangeEnd,attr"`
		GroupAddress []struct {
			ID      string `xml:"Id,attr"`
			Name    string `xml:"Name,attr"`
			Address uint   `xml:"Address,attr"`
		}
		GroupRange []groupRange11
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	gar.ID = GroupRangeID(doc.ID)
	gar.Name = doc.Name
	gar.RangeStart = doc.RangeStart
	gar.RangeEnd = doc.RangeEnd
	gar.Addresses = make([]GroupAddress, len(doc.GroupAddress))
	gar.SubRanges = make([]GroupRange, len(doc.GroupRange))

	for n, docGrpAddr := range doc.GroupAddress {
		gar.Addresses[n] = GroupAddress{
			ID:      GroupAddressID(docGrpAddr.ID),
			Name:    docGrpAddr.Name,
			Address: docGrpAddr.Address,
		}
	}

	for n, docGrpRange := range doc.GroupRange {
		gar.SubRanges[n] = GroupRange(docGrpRange)
	}

	return nil
}

type installation11 Installation

func (i *installation11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Name        string         `xml:"Name,attr"`
		Areas       []area11       `xml:"Topology>Area"`
		GroupRanges []groupRange11 `xml:"GroupAddresses>GroupRanges>GroupRange"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	i.Name = doc.Name
	i.Topology = make([]Area, len(doc.Areas))
	i.GroupAddresses = make([]GroupRange, len(doc.GroupRanges))

	for n, docArea := range doc.Areas {
		i.Topology[n] = Area(docArea)
	}

	for n, docGrpRange := range doc.GroupRanges {
		i.GroupAddresses[n] = GroupRange(docGrpRange)
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
