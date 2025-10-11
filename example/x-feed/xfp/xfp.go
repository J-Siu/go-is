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

// xfp - X.com Feed Processor
// README (1)
package xfp

import (
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

// (1.1) Write a `info` struct
type XFeedInfo struct {
	is.InfoBase // (1.1) REQUIRED: embed [is.InfoBase] to get [is.IInfo] interface

	// Added fields
	User string `json:"User"`
	Text string `json:"Text"`
}

// (1.1) REQUIRED: embed `is.InfoBase` for `is.IInfo` interface
func (xi *XFeedInfo) String() string { return xi.User + ": " + xi.Text }

// (1.2) Write a `processor` struct
type XFeedProcessor struct {
	*is.Processor // (1.2) REQUIRED: embed `*is.Processor`
}

// (1.3) Write package/struct level `New` function. Must accept `*is.Property` as one of its arguments.
func (x *XFeedProcessor) New(
	property *is.Property, // (1.3) REQUIRED: `*is.Property` as one of its arguments
) *XFeedProcessor {
	x.Processor = is.New(property) // (1.3) REQUIRED: use [is.New] to create and initialize the embedded [*is.Processor]
	x.MyType = "xf"                // Optional: features of [basestruct.Base] embedded in [is.Processor]
	x.override()                   // (1.3) Override `is.Processor` field functions as needed
	return x
}

// (1.3) Override `is.Processor` field functions as needed
func (x *XFeedProcessor) override() {
	x.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := x.MyType + ".V020"
		ezlog.Trace().N(prefix).TxtStart().Out()
		var es rod.Elements
		tagName := "article"
		if element == nil {
			es = x.Page.MustElements(tagName)
		} else {
			es = element.MustElements(tagName)
		}
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return &es
	}
	x.V030_ElementInfo = func(element *rod.Element, index int) is.IInfo {
		prefix := x.MyType + ".V030"
		ezlog.Trace().N(prefix).TxtStart().Out()
		ezlog.Trace().M(element.MustHTML()).Out()
		info := new(XFeedInfo)
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
		ezlog.Debug().N(prefix).Nn("info").M(info).Out()

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return info
	}
}
