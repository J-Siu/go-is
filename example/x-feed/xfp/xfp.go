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
	"github.com/J-Siu/go-is/v3/is"
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
func (t *XFeedProcessor) New(
	property *is.Property, // (1.3) REQUIRED: `*is.Property` as one of its arguments
) *XFeedProcessor {
	t.Processor = is.New(property) // (1.3) REQUIRED: use [is.New] to create and initialize the embedded [*is.Processor]
	t.MyType = "xf"                // Optional: features of [basestruct.Base] embedded in [is.Processor]
	t.override()                   // (1.3) Override `is.Processor` field functions as needed
	return t
}

// (1.3) Override `is.Processor` field functions as needed
func (t *XFeedProcessor) override() {
	t.V020_Elements = t.override_V20
	t.V030_ElementInfo = t.override_V30
}

func (t *XFeedProcessor) override_V20() {
	prefix := t.MyType + ".V020"
	t.StateCurr.Name = prefix
	var es rod.Elements
	tagName := "article"
	if t.StateCurr.Element == nil {
		es = t.Page.MustElements(tagName)
	} else {
		es = t.StateCurr.Element.MustElements(tagName)
	}
	t.StateCurr.Elements = es
}

func (t *XFeedProcessor) override_V30() {
	prefix := t.MyType + ".V030"
	t.StateCurr.Name = prefix
	ezlog.Trace().M(t.StateCurr.Element.MustHTML()).Out()
	info := new(XFeedInfo)
	var (
		err error
		e   *rod.Element
		tag string
	)

	// Username
	tag = "[data-testid='User-Name']"
	e, err = t.StateCurr.Element.Element(tag)
	if err == nil && e != nil {
		tag = "a"
		e, err = e.Element(tag)
		if err == nil && e != nil {
			info.User = e.MustText()
		}
	}

	// Tweet text
	tag = "[data-testid='tweetText']"
	e, err = t.StateCurr.Element.Element(tag)
	if err == nil && e != nil {
		info.Text = e.MustText()
	}
	ezlog.Debug().N(prefix).N("info").Lm(info).Out()
	t.StateCurr.ElementInfo = info
}
