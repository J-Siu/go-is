/*
Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package xfeed

import (
	"encoding/json"

	"github.com/J-Siu/go-ezlog"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type XInfo struct {
	is.InfoBase // MUST embed [is.InfoBase] to get [is.IInfo] interface

	User string `json:"user,omitempty"`
	Text string `json:"text,omitempty"`
}

func (xi *XInfo) String() string { return xi.User + ": " + xi.Text }

type XFeed struct {
	*is.Processor // MUST embed [*is.Processor]
}

// Override [is.Processor] field func
func (x *XFeed) override() {
	x.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := x.MyType + ".V020"
		ezlog.Trace(prefix + ": Start")
		var es rod.Elements
		tagName := "article"
		if element == nil {
			es = x.Page.MustElements(tagName)
		} else {
			es = element.MustElements(tagName)
		}
		ezlog.Trace(prefix + ": End")
		return &es
	}
	x.V030_ElementInfo = func(element *rod.Element, index int) is.IInfo {
		prefix := x.MyType + ".V030"
		ezlog.Trace(prefix + ": Start")
		ezlog.Trace(element.MustHTML())
		info := new(XInfo)
		var (
			err error
			e   *rod.Element
			tag string
		)

		// Username
		tag = "[data-testid='User-Name']"
		e, err = element.Element(tag)
		if err == nil && e != nil {
			tag = "a"
			e, err = e.Element(tag)
			if err == nil && e != nil {
				info.User = e.MustText()
			}
		}

		// Tweet text
		tag = "[data-testid='tweetText']"
		e, err = element.Element(tag)
		if err == nil && e != nil {
			info.Text = e.MustText()
		}
		ezlog.Debug(prefix + ": info:")
		ezlog.DebugP(MustToJsonStrP(info))

		ezlog.Trace(prefix + ": End")
		return info
	}
}

// Create and initialize XFeed
func New(property *is.Property) *XFeed {
	xf := new(XFeed)
	xf.Processor = is.New(property) // MUST use [is.New] to create and initialize the embedded [*is.Processor]
	xf.MyType = "xf"                // Optional: features of [basestruct.Base]
	xf.override()
	return xf
}

// helper function for printing/logging struct
func MustToJsonStrP(obj any) *string {
	prefix := "MustToJsonStrP"
	var str string
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		str = "{<MustToJsonStrP() failed>}"
		ezlog.Trace(prefix + ": Err: " + err.Error())
	} else {
		str = string(b)
	}
	return &str
}
