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

package is

import (
	"strconv"

	"github.com/J-Siu/go-ezlog"
)

type IInfoListPrintMode int8

const (
	PrintAll IInfoListPrintMode = iota
	PrintMatched
	PrintUnmatched
)

// Interface for info struct
type IInfo interface {
	Matched() bool                   // Get matched bool value
	MatchedStr() string              // Get matched string value
	SetMatched(matched bool)         // Set matched bool value
	SetMatchedStr(matchedStr string) // Set matched string value
	String() string                  // Info struct to string
}

// IInfo base struct to be embedded
//   - Only String() should be overloaded
type InfoBase struct {
	matched    bool
	matchedStr string
}

// Get matched bool value
func (s *InfoBase) Matched() bool { return s.matched }

// Get matched string value
func (s *InfoBase) MatchedStr() string { return s.matchedStr }

// Set matched bool value
func (s *InfoBase) SetMatched(matched bool) { s.matched = matched }

// Set matched string value
func (s *InfoBase) SetMatchedStr(matchedStr string) { s.matchedStr = matchedStr }

// Place holder only
func (s *InfoBase) String() string { return "String() placeholder!" }

type IInfoList []IInfo

func (list *IInfoList) Print(mode IInfoListPrintMode) {
	var mark string
	for c, info := range *list {
		mark = " [ ] "
		if info.Matched() {
			mark = " [X] "
		}
		switch mode {
		case PrintAll:
			ezlog.Msg(strconv.Itoa(c+1) + mark + info.String())
		case PrintMatched:
			if info.Matched() {
				ezlog.Msg(strconv.Itoa(c+1) + mark + info.String())
			}
		case PrintUnmatched:
			if !info.Matched() {
				ezlog.Msg(strconv.Itoa(c+1) + mark + info.String())
			}
		}
	}
}
