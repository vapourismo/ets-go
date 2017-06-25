// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import "encoding/xml"

const schema11Namespace = "http://knx.org/xml/project/11"

type projectInfo11 ProjectInfo

func (pi *projectInfo11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Project struct {
			ID                 string `xml:"Id,attr"`
			ProjectInformation struct {
				Name string `xml:",attr"`
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
		Name       string `xml:",attr"`
		Address    uint   `xml:",attr"`
		ComObjects []struct {
			RefID         string `xml:"RefId,attr"`
			DatapointType string `xml:",attr"`
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
		Name           string `xml:",attr"`
		Address        uint   `xml:",attr"`
		DeviceInstance []deviceInstance11
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	l.ID = LineID(doc.ID)
	l.Name = doc.Name
	l.Address = doc.Address
	l.Devices = make([]DeviceInstance, len(doc.DeviceInstance))

	for n, docDeviceInstance := range doc.DeviceInstance {
		l.Devices[n] = DeviceInstance(docDeviceInstance)
	}

	return nil
}

type area11 Area

func (a *area11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID      string `xml:"Id,attr"`
		Name    string `xml:",attr"`
		Address uint   `xml:",attr"`
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
		Name         string `xml:",attr"`
		RangeStart   uint   `xml:",attr"`
		RangeEnd     uint   `xml:",attr"`
		GroupAddress []struct {
			ID      string `xml:"Id,attr"`
			Name    string `xml:",attr"`
			Address uint   `xml:",attr"`
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
		Name        string         `xml:",attr"`
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

type project11 Project

func (p *project11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

type comObject11 ComObject

func (co *comObject11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID                string `xml:"Id,attr"`
		Name              string `xml:",attr"`
		Text              string `xml:",attr"`
		Description       string `xml:",attr"`
		FunctionText      string `xml:",attr"`
		ObjectSize        string `xml:",attr"`
		DatapointType     string `xml:",attr"`
		Priority          string `xml:",attr"`
		ReadFlag          string `xml:",attr"`
		WriteFlag         string `xml:",attr"`
		CommunicationFlag string `xml:",attr"`
		TransmitFlag      string `xml:",attr"`
		UpdateFlag        string `xml:",attr"`
		ReadOnInitFlag    string `xml:",attr"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	co.ID = ComObjectID(doc.ID)
	co.Name = doc.Name
	co.Text = doc.Text
	co.Description = doc.Description
	co.FunctionText = doc.FunctionText
	co.ObjectSize = doc.ObjectSize
	co.DatapointType = doc.DatapointType
	co.Priority = doc.Priority
	co.ReadFlag = doc.ReadFlag == "Enabled"
	co.WriteFlag = doc.WriteFlag == "Enabled"
	co.CommunicationFlag = doc.CommunicationFlag == "Enabled"
	co.TransmitFlag = doc.TransmitFlag == "Enabled"
	co.UpdateFlag = doc.UpdateFlag == "Enabled"
	co.ReadOnInitFlag = doc.ReadOnInitFlag == "Enabled"

	return nil
}

type comObjectRef11 ComObjectRef

func (cor *comObjectRef11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID                string  `xml:"Id,attr"`
		RefID             string  `xml:"RefId,attr"`
		Name              *string `xml:",attr"`
		Text              *string `xml:",attr"`
		Description       *string `xml:",attr"`
		FunctionText      *string `xml:",attr"`
		ObjectSize        *string `xml:",attr"`
		DatapointType     *string `xml:",attr"`
		Priority          *string `xml:",attr"`
		ReadFlag          *string `xml:",attr"`
		WriteFlag         *string `xml:",attr"`
		CommunicationFlag *string `xml:",attr"`
		TransmitFlag      *string `xml:",attr"`
		UpdateFlag        *string `xml:",attr"`
		ReadOnInitFlag    *string `xml:",attr"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	cor.ID = ComObjectRefID(doc.ID)
	cor.RefID = ComObjectID(doc.RefID)
	cor.Name = doc.Name
	cor.Text = doc.Text
	cor.Description = doc.Description
	cor.FunctionText = doc.FunctionText
	cor.ObjectSize = doc.ObjectSize
	cor.DatapointType = doc.DatapointType
	cor.Priority = doc.Priority

	if doc.ReadFlag != nil {
		cor.ReadFlag = new(bool)
		*cor.ReadFlag = *doc.ReadFlag == "Enabled"
	}

	if doc.WriteFlag != nil {
		cor.WriteFlag = new(bool)
		*cor.WriteFlag = *doc.WriteFlag == "Enabled"
	}

	if doc.CommunicationFlag != nil {
		cor.CommunicationFlag = new(bool)
		*cor.CommunicationFlag = *doc.CommunicationFlag == "Enabled"
	}

	if doc.TransmitFlag != nil {
		cor.TransmitFlag = new(bool)
		*cor.TransmitFlag = *doc.TransmitFlag == "Enabled"
	}

	if doc.UpdateFlag != nil {
		cor.UpdateFlag = new(bool)
		*cor.UpdateFlag = *doc.UpdateFlag == "Enabled"
	}

	if doc.ReadOnInitFlag != nil {
		cor.ReadOnInitFlag = new(bool)
		*cor.ReadOnInitFlag = *doc.ReadOnInitFlag == "Enabled"
	}

	return nil
}

type applicationProgram11 ApplicationProgram

func (ap *applicationProgram11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		ID      string `xml:"Id,attr"`
		Name    string `xml:",attr"`
		Version uint   `xml:"ApplicationVersion,attr"`
		Static  struct {
			Objects    []comObject11    `xml:"ComObjectTable>ComObject"`
			ObjectRefs []comObjectRef11 `xml:"ComObjectRefs>ComObjectRef"`
		}
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	ap.ID = ApplicationProgramID(doc.ID)
	ap.Name = doc.Name
	ap.Version = doc.Version
	ap.Objects = make([]ComObject, len(doc.Static.Objects))
	ap.ObjectRefs = make([]ComObjectRef, len(doc.Static.ObjectRefs))

	for n, docComObj := range doc.Static.Objects {
		ap.Objects[n] = ComObject(docComObj)
	}

	for n, docComObjRef := range doc.Static.ObjectRefs {
		ap.ObjectRefs[n] = ComObjectRef(docComObjRef)
	}

	return nil
}

type manufacturerData11 ManufacturerData

func (md *manufacturerData11) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var doc struct {
		Manufacturer struct {
			ID       string                 `xml:"RefId,attr"`
			Programs []applicationProgram11 `xml:"ApplicationPrograms>ApplicationProgram"`
		} `xml:"ManufacturerData>Manufacturer"`
	}

	if err := d.DecodeElement(&doc, &start); err != nil {
		return err
	}

	md.Manufacturer = ManufacturerID(doc.Manufacturer.ID)
	md.Programs = make([]ApplicationProgram, len(doc.Manufacturer.Programs))

	for n, docProg := range doc.Manufacturer.Programs {
		md.Programs[n] = ApplicationProgram(docProg)
	}

	return nil
}
