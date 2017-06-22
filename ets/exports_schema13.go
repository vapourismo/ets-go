// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import "encoding/xml"

// unmarshalProject11 extracts project information from the current element.
func unmarshalProjectInfo13(d *xml.Decoder, start xml.StartElement, proj *ProjectInfo) error {
	// Schema 11 and 13 are compatible for our purposes.
	return unmarshalProjectInfo11(d, start, proj)
}
