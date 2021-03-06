// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at http://mozilla.org/MPL/2.0/.

package instance

import (
	"io"
	"strings"

	"github.com/ta2gch/iris/reader/tokenizer"
	"github.com/ta2gch/iris/runtime/ilos"
)

type Stream struct {
	Column *int
	Reader *tokenizer.Reader
	Writer io.Writer
}

func NewStream(r io.Reader, w io.Writer) ilos.Instance {
	if r == nil {
		return Stream{new(int), nil, w}
	}
	return Stream{new(int), tokenizer.NewReader(r), w}
}

func (Stream) Class() ilos.Class {
	return StreamClass
}

func (s Stream) Write(p []byte) (n int, err error) {
	i := strings.LastIndex(string(p), "\n")
	if i < 0 {
		*s.Column += len(p)
	} else {
		*s.Column = len(p[i+1:])
	}
	return s.Writer.Write(p)
}

func (s Stream) Read(p []byte) (n int, err error) {
	return s.Reader.Read(p)
}

func (Stream) String() string {
	return "#<INPUT-STREAM>"
}
