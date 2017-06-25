// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"encoding/xml"
	"fmt"
	"io"
)

// ComObjectID is the ID of a communication object.
type ComObjectID string

// ComObject is a communication object.
type ComObject struct {
	ID                ComObjectID
	Name              string
	Text              string
	Description       string
	FunctionText      string
	ObjectSize        string
	DatapointType     string
	Priority          string
	ReadFlag          bool
	WriteFlag         bool
	CommunicationFlag bool
	TransmitFlag      bool
	UpdateFlag        bool
	ReadOnInitFlag    bool
}

// ApplicationProgramID is the ID of an application program.
type ApplicationProgramID string

// ApplicationProgram is an application program.
type ApplicationProgram struct {
	ID      ApplicationProgramID
	Name    string
	Version uint
	Objects []ComObject
}

// ManufacturerID is the ID of a manufacturer.
type ManufacturerID string

// ManufacturerData contains manufacturer-specific data.
type ManufacturerData struct {
	Manufacturer ManufacturerID
	Programs     []ApplicationProgram
}

// UnmarshalXML implements xml.Unmarshaler.
func (md *ManufacturerData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Decide which schema to use based on the value of the 'xmlns' attribute.
	ns := getNamespace(start)
	switch ns {
	case schema11Namespace, schema12Namespace, schema13Namespace:
		return d.DecodeElement((*manufacturerData11)(md), &start)

	default:
		return fmt.Errorf("Unexpected namespace '%s'", ns)
	}
}

// DecodeManufacturerData parses the contents of a manufacturer file.
func DecodeManufacturerData(r io.Reader) (*ManufacturerData, error) {
	md := &ManufacturerData{}
	if err := xml.NewDecoder(r).Decode(md); err != nil {
		return nil, err
	}

	return md, nil
}
