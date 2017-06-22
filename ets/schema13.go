// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import "encoding/xml"

func unmarshalProjectInfo13(d *xml.Decoder, start xml.StartElement, pi *ProjectInfo) error {
	// Schema 11 and 13 are compatible for our purposes.
	return unmarshalProjectInfo11(d, start, pi)
}

func unmarshalProject13(d *xml.Decoder, start xml.StartElement, p *Project) error {
	return unmarshalProject11(d, start, p)
}
