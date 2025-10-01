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

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/go-rod/rod"
)

// # [State]
//
// Used at the beginning of [Processor.Run()] scroll loop for [breakLoop] calculation,
// and [Processor.ElementScroll()] for scrolling.
//
// At the bottom of 'Run()' scroll loop, it is passed into [Processor.V100_ScrollLoopEnd()] for customized scroll calculation.
type State struct {
	*basestruct.Base

	Elements          *rod.Elements `json:"elements,omitempty"`            // Result of [Processor.V020_Elements()]
	ElementLast       *rod.Element  `json:"element_last,omitempty"`        // Last element of the previous scroll loop iteration
	ElementLastScroll *rod.Element  `json:"element_last_scroll,omitempty"` // Element used for previous scroll (not necessary last loop iteration)
	ElementCountLast  int           `json:"element_count_last,omitempty"`  // Number of elements of previous loop iteration
	InfoLast          IInfo         `json:"info_last,omitempty"`           // [Info] of [ElementLast]. Return from [Processor.V030_ElementInfo()]
	Scroll            bool          `json:"scroll,omitempty"`              // Used by [breakLoop]. True = to scroll. False = don't scroll.
	ScrollCount       int           `json:"scroll_count,omitempty"`        // Total number of times [Processor.ElementScroll()] called
}

func (s *State) New() *State {
	s.Base = new(basestruct.Base)
	s.MyType = "RunState"
	s.Initialized = true

	// 'Scroll' need to be init, as the default value is 'false'
	s.Scroll = true

	return s
}

// This should only be used at Trace level log
func (s *State) String() *string {
	var str string
	if s.Elements != nil {
		str += "esCount:" + strconv.Itoa(len(*s.Elements)) + "\n"
	}
	if s.ElementLast != nil {
		str += "eLast:" + string(s.ElementLast.Object.ObjectID) + "\n"
	}
	if s.ElementLastScroll != nil {
		str += "eLastScroll:" + string(s.ElementLastScroll.Object.ObjectID) + "\n"
	}
	if s.InfoLast != nil {
		str += "ListItemLast(matched):" + strconv.FormatBool(s.InfoLast.Matched()) + "\n"
	}
	str += "eCountLast:" + strconv.Itoa(s.ElementCountLast) + "\n"
	str += "Scroll:" + strconv.FormatBool(s.Scroll)
	return &str
}
