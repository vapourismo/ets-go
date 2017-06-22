// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import "encoding/xml"

// unmarshalProject11 extracts project information from the current element.
func unmarshalProject13(d *xml.Decoder, start xml.StartElement, proj *Project) error {
	// Schema 11 and 13 are compatible for our purposes.
	return unmarshalProject11(d, start, proj)
}