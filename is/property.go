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

// Package [is] is an infinite scroll processor using [go-rod/rod](https://github.com/go-rod/rod).
package is

import (
	"github.com/go-rod/rod"
)

type Property struct {
	// -- [rod] element

	Page      *rod.Page    `json:"Page,omitempty"`      // REQUIRED: Page element of [rod].
	Container *rod.Element `json:"Container,omitempty"` // The outer most rod.Element containing all repeating items

	// -- URL

	UrlCheck bool   `json:"UrlCheck,omitempty"` // Check [UrlStr] before loading
	UrlLoad  bool   `json:"UrlLoad,omitempty"`  // Control if [UrlStr] should be load at the beginning of [Run]
	UrlStr   string `json:"UrlStr,omitempty"`   // URL string used in [LoadPage]. Not use if [UrlLoad] = false

	// -- Flow control

	ScrollMax int `json:"ScrollMax,omitempty"` // Maximum time the page should be scrolled

	// -- Information collection

	IInfoList *IInfoList `json:"IInfoList,omitempty"` // Pointer of array of IInfo. If not nil, IInfo item will be added to the array
}
