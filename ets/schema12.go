// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import "encoding/xml"

const schema12Namespace = "http://knx.org/xml/project/12"

func unmarshalProjectInfo12(d *xml.Decoder, start xml.StartElement, pi *ProjectInfo) error {
	// Schema 11 and 12 are compatible for our purposes.
	return unmarshalProjectInfo11(d, start, pi)
}

func unmarshalProject12(d *xml.Decoder, start xml.StartElement, p *Project) error {
	// Schema 11 and 12 are compatible for our purposes.
	return unmarshalProject11(d, start, p)
}
