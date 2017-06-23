// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"encoding/xml"
	"fmt"
	"io"
)

func getNamespace(start xml.StartElement) (ns string) {
	for _, attr := range start.Attr {
		if attr.Name.Local == "xmlns" {
			ns = attr.Value
			break
		}
	}

	return
}

// ProjectID is a project identifier.
type ProjectID string

// ProjectInfo contains project information.
type ProjectInfo struct {
	ID   ProjectID
	Name string
}

// UnmarshalXML implements xml.Unmarshaler.
func (pi *ProjectInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Decide which schema to use based on the value of the 'xmlns' attribute.
	ns := getNamespace(start)
	switch ns {
	case schema11Namespace:
		return unmarshalProjectInfo11(d, start, pi)

	case schema13Namespace:
		return unmarshalProjectInfo13(d, start, pi)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", ns)
	}
}

// DecodeProjectInfo parses the contents of project info file.
func DecodeProjectInfo(r io.Reader) (*ProjectInfo, error) {
	info := &ProjectInfo{}
	if err := xml.NewDecoder(r).Decode(info); err != nil {
		return nil, err
	}

	return info, nil
}

// Connector is a connection to a group address.
type Connector struct {
	Receive bool
	RefID   GroupAddressID
}

// ComObjectRefID is the ID of a communication object reference.
type ComObjectRefID string

// ComObjectInstanceRef connects a communication object reference with zero or more group addresses.
type ComObjectInstanceRef struct {
	RefID         ComObjectRefID
	DatapointType string
	Connectors    []Connector
}

// DeviceInstanceID is the ID of a device instance.
type DeviceInstanceID string

// DeviceInstance is a device instance.
type DeviceInstance struct {
	ID         DeviceInstanceID
	Name       string
	Address    uint
	ComObjects []ComObjectInstanceRef
}

// LineID is the ID of a line.
type LineID string

// Line is a line.
type Line struct {
	ID               LineID
	Name             string
	Address          uint
	DevicesInstances []DeviceInstance
}

// AreaID is the ID of an area.
type AreaID string

// Area is an area.
type Area struct {
	ID      AreaID
	Name    string
	Address uint
	Lines   []Line
}

// GroupAddressID is the ID of a group address.
type GroupAddressID string

// GroupAddress is a group address.
type GroupAddress struct {
	ID      GroupAddressID
	Name    string
	Address uint
}

// GroupRangeID is the ID of a group range.
type GroupRangeID string

// GroupRange is a range of group addresses.
type GroupRange struct {
	ID         GroupRangeID
	Name       string
	RangeStart uint
	RangeEnd   uint
	Addresses  []GroupAddress
	SubRanges  []GroupRange
}

// Installation is an installation within a project.
type Installation struct {
	Name           string
	Topology       []Area
	GroupAddresses []GroupRange
}

// Project contains an entire project.
type Project struct {
	ID            ProjectID
	Installations []Installation
}

// UnmarshalXML implements xml.Unmarshaler.
func (p *Project) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Decide which schema to use based on the value of the 'xmlns' attribute.
	ns := getNamespace(start)
	switch ns {
	case schema11Namespace:
		return unmarshalProject11(d, start, p)

	case schema13Namespace:
		return unmarshalProject13(d, start, p)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", ns)
	}
}

// DecodeProject parses the contents of a project file.
func DecodeProject(r io.Reader) (*Project, error) {
	proj := &Project{}
	if err := xml.NewDecoder(r).Decode(proj); err != nil {
		return nil, err
	}

	return proj, nil
}
