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
	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/ezlog"
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
	// --
	Name string `json:"FuncName"` // current function/state name
	// --
	Elements      rod.Elements `json:"-"`             // Result of [Processor.V020_Elements()]
	ElementsCount int          `json:"ElementsCount"` // Number of elements in current iteration
	// --
	Element      *rod.Element `json:"Element"`      // Element being process
	ElementIndex int          `json:"ElementIndex"` // Index of element being process
	ElementInfo  IInfo        `json:"ElementInfo"`  // [Info] of [ElementLast]. Return from [Processor.V030_ElementInfo()]
	// --
	ElementScrollable bool `json:"ElementScrollable"` // update by V080_ElementScrollable
	// --
	ScrollableElement      *rod.Element `json:"ScrollableElement"`      // Last scrollable element
	ScrollableElementIndex int          `json:"ScrollableElementIndex"` // Index of element being process
	ScrollableElementInfo  IInfo        `json:"ScrollableElementInfo"`  // [Info] of [ElementScrollable].
	// --
	Scroll      bool `json:"Scroll"`      // True = to scroll. False = don't scroll.
	ScrollCount int  `json:"ScrollCount"` // Total number of times [Processor.ElementScroll()] called
	ScrollPage  bool `json:"ScrollLoop"`  // update by ScrollLoop
}

func (t *State) New(scrollCount int) *State {
	t.Base = new(basestruct.Base)
	t.MyType = "State"
	prefix := t.MyType + ".New"
	t.Initialized = true
	t.Scroll = true // 'Scroll' need to be init, as the default value is 'false'
	t.ScrollPage = true
	t.ScrollCount = scrollCount
	t.ScrollableElementInfo = nil
	ezlog.Debug().M(prefix).Out()
	return t
}
