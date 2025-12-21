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
	"errors"

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/go-rod/rod"
)

// IS processor structure
type Processor struct {
	*basestruct.Base
	*Property

	StateCurr *State
	StatePrev *State

	// -- Following 4 field func rarely need override

	// Load [UrlStr] into [Page]
	//
	// No override needed.
	LoadPage func() `json:"-"`

	// Determine whether the scroll loop should continue running
	//
	// No override needed.
	ScrollLoop func() bool `json:"-"`

	// Detect end of page, scroll no longer possible.
	//
	// No override needed.
	//
	// If elements are removed during [Run()], overload [V100_ExitScroll()] to do custom override.
	// As both of following checks can be flawed if elements are removed from page DOM.
	ScrollCalculation func() (scroll bool) `json:"-"`

	// Use [MustScrollIntoView] on [element]
	//
	// No override needed.
	ScrollElement func(element *rod.Element) `json:"-"`

	// --- Overload following field func as needed

	// Return the container element.
	//
	// build-in behavior is to return [property.Container]
	//
	// Override if needed
	V010_Container func() (container *rod.Element) `json:"-"`

	// Return collection of repeating elements within [property.Page] or [property.Container]
	//
	// build-in behavior is to return `nil`
	//
	// **Must override**
	V020_Elements func(container *rod.Element) *rod.Elements `json:"-"`

	// Extract information from [element] and put into an [IInfo] structure and return it.
	//
	// build-in behavior is to return `nil`
	//
	// **Must override**
	V030_ElementInfo func() (info IInfo) `json:"-"`

	// Determine [element] is a match or not base on [info]
	//
	// build-in behavior is to return (`true`, `""`)
	//
	// Override if needed
	V040_ElementMatch func() (matched bool, matchedStr string) `json:"-"`

	// Do some processing (eg, print, write to file, db, etc) if [element] is a match
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V050_ElementProcessMatched func() `json:"-"`

	// Do some processing (eg, print, write to file, db, etc) if [element] is not a match
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V060_ElementProcessUnmatch func() `json:"-"`

	// Do some processing (eg, print, write to file, db, etc) regardless [element] is a match or not
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V070_ElementProcess func() `json:"-"`

	// If current element is scrollable
	//
	// build-in behavior is to return `true``
	//
	// Override if needed
	V080_ElementScrollable func() bool `json:"-"`

	// Do some processing if required
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V090_ElementLoopEnd func() `json:"-"`

	// Do some processing if required
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V100_ScrollLoopEnd func() `json:"-"`
}

// Parameters:
//   - property *Property
//
// Returns:
//   - *Processor
func (t *Processor) New(property *Property) *Processor {
	t.Base = new(basestruct.Base)
	t.MyType = "is.Processor"
	prefix := t.MyType + ".New" + "(base)"

	if property == nil {
		t.Err = errors.New("is.New: property cannot be nil")
	} else if property.Page == nil {
		t.Err = errors.New("is.New: page/tab cannot be nil")
	} else {
		t.Property = property
		t.setFunc()
		t.Initialized = true
	}

	ezlog.Trace().N(prefix).M("Done").Out()
	return t
}

// Process the page
//
// No override needed.
func (t *Processor) Run() {
	prefix := t.MyType + ".Run" + "(base)"
	if t.CheckErrInit(prefix) {
		t.LoadPage()
	}
	if t.Err == nil {
		// Initial container
		t.Container = t.V010_Container()
		// Scroll Loop
		for t.ScrollLoop() {
			// -- SCROLL LOOP - START
			t.StatePrev = t.StateCurr
			if t.StatePrev != nil {
				if t.StatePrev.ScrollableElement != nil {
					t.ScrollElement(t.StatePrev.ScrollableElement)
				}
			} else {
				t.StatePrev = new(State).New(0)
			}
			t.StateCurr = new(State).New(t.StatePrev.ScrollCount)
			// -- Get elements
			t.StateCurr.ElementsCount = 0
			t.StateCurr.Elements = t.V020_Elements(t.Container)

			if t.StateCurr.Elements == nil {
				t.StateCurr.Scroll = false // no element, no scroll
			} else {
				t.StateCurr.ElementsCount = len(*t.StateCurr.Elements)
				ezlog.Trace().N(prefix).N("elements count").M(t.StateCurr.ElementsCount).Out()
				for index := t.StatePrev.ElementsCount; index < t.StateCurr.ElementsCount; index++ {
					// -- ELEMENTS LOOP - START
					t.StateCurr.Element = (*t.StateCurr.Elements)[index]
					t.StateCurr.ElementIndex = index
					t.StateCurr.ElementInfo = t.V030_ElementInfo()
					matched, matchedStr := t.V040_ElementMatch()
					if t.StateCurr.ElementInfo != nil {
						// Remove burden from package user
						t.StateCurr.ElementInfo.SetMatched(matched)
						t.StateCurr.ElementInfo.SetMatchedStr(matchedStr)
					}
					if matched {
						t.V050_ElementProcessMatched()
					} else {
						t.V060_ElementProcessUnmatch()
					}
					t.V070_ElementProcess()
					// info list
					if t.IInfoList != nil && t.StateCurr.ElementInfo != nil {
						*t.IInfoList = append(*t.IInfoList, t.StateCurr.ElementInfo)
					}
					if t.V080_ElementScrollable() {
						t.StateCurr.ScrollableElement = t.StateCurr.Element
						t.StateCurr.ScrollableElementIndex = t.StateCurr.ElementIndex
						t.StateCurr.ScrollableElementInfo = t.StateCurr.ElementInfo
					}
					t.V090_ElementLoopEnd()
					// -- ELEMENTS LOOP - END
				}
				t.StateCurr.Scroll = t.ScrollCalculation()
			}

			t.V100_ScrollLoopEnd()
			t.StateCurr.ScrollCount++
			// -- SCROLL LOOP - END
		}
	}
}

// Implement the default field functions
func (t *Processor) setFunc() {
	// -- Following 4 field func rarely need override
	t.LoadPage = t.base_LoadPage
	t.ScrollCalculation = t.base_ScrollCalculation
	t.ScrollElement = t.base_ScrollElement
	t.ScrollLoop = t.base_ScrollLoop
	// --- Overload following field func as needed
	t.V010_Container = t.base_V010_Container
	t.V020_Elements = t.base_V020_Elements
	t.V030_ElementInfo = t.base_V030_ElementInfo
	t.V040_ElementMatch = t.base_V040_ElementMatch
	t.V050_ElementProcessMatched = t.base_V050_ElementProcessMatched
	t.V060_ElementProcessUnmatch = t.base_V060_ElementProcessUnmatch
	t.V070_ElementProcess = t.base_V070_ElementProcess
	t.V080_ElementScrollable = t.base_V080_ElementScrollable
	t.V090_ElementLoopEnd = t.base_V090_ElementLoopEnd
	t.V100_ScrollLoopEnd = t.base_V100_ScrollLoopEnd
}

func (t *Processor) base_LoadPage() {
	prefix := t.MyType + ".LoadPage" + "(base)"
	ezlog.Trace().N(prefix).TxtStart().Out()
	if t.CheckErrInit(prefix) {
		if t.UrlLoad {
			ezlog.Debug().N(prefix).N("urlStr").M(t.UrlStr).Out()
			t.Err = t.Page.Navigate(t.UrlStr)
			if t.Err == nil {
				ezlog.Trace().N(prefix).N("MustWaitDOMStable").TxtStart().Out()
				t.Page.MustWaitDOMStable()
				ezlog.Trace().N(prefix).N("MustWaitDOMStable").TxtEnd().Out()
			}
		}
		if t.Err != nil {
			t.Err = errors.New(prefix + ": " + t.Err.Error())
			ezlog.Err().M(t.Err).Out()
		}
	}
	ezlog.Trace().N(prefix).TxtEnd().Out()
}

func (t *Processor) base_ScrollCalculation() (scroll bool) {
	prefix := t.MyType + ".ScrollCalculation" + "(base)"
	scroll = t.StateCurr.Scroll
	/*
		"if (elementsCount == elementsCountLast)":
			will be triggered, if number of elements removed
					= number of new elements added after scroll

		"if (ElementLastScroll == ElementLast)":
			will be triggered, if all new elements added after scroll are removed
	*/
	if t.StateCurr.ScrollableElement == nil {
		scroll = false
	} else if t.StatePrev != nil && t.StatePrev.ScrollableElement != nil && t.StateCurr.ScrollableElement != nil {
		if t.StatePrev.ScrollableElement.Object.ObjectID == t.StateCurr.ScrollableElement.Object.ObjectID {
			// prev scroll element == curr scroll element
			scroll = false
		}
	}
	ezlog.Trace().N(prefix).M("Done").Out()
	return scroll
}

func (t *Processor) base_ScrollElement(element *rod.Element) {
	prefix := t.MyType + ".ScrollElement" + "(base)"
	ezlog.Trace().N(prefix).TxtStart().Out()
	if element != nil {
		element.MustScrollIntoView()
		ezlog.Trace().N(prefix).M("Scrolled").Out()
		// ezlog.Trace().N(prefix).N("MustWaitDOMStable").TxtStart().Out()
		t.Page.MustWaitDOMStable()
		// ezlog.Trace().N(prefix).N("MustWaitDOMStable").TxtEnd().Out()
	}
	ezlog.Trace().N(prefix).TxtEnd().Out()
}

func (t *Processor) base_ScrollLoop() (scroll bool) {
	prefix := t.MyType + ".ScrollLoop" + "(base)"
	scroll = t.StateCurr == nil || (t.StateCurr.Scroll && (t.StateCurr.ScrollCount < t.ScrollMax || t.ScrollMax < 0))
	if ezlog.GetLogLevel() == ezlog.TRACE {
		ezlog.Trace().N(prefix).
			Ln("StateCurr").M(t.StateCurr).
			Ln("scrollMax").M(t.ScrollMax).
			Ln("scroll").N("t.StateCurr == nil || (t.StateCurr.Scroll && (t.StateCurr.ScrollCount < t.ScrollMax || t.ScrollMax < 0))").M(scroll).
			Out()
	}
	return
}

func (t *Processor) base_V010_Container() (container *rod.Element) {
	prefix := t.MyType + ".V010_Container" + "(base)"
	ezlog.Trace().N(prefix).M("Done").Out()
	return t.Container
}

func (t *Processor) base_V020_Elements(container *rod.Element) *rod.Elements {
	prefix := t.MyType + ".V020_Elements" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing. Return `nil`").Out()
	return nil
}

func (t *Processor) base_V030_ElementInfo() (info IInfo) {
	prefix := t.MyType + ".V030_ElementInfo" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing. Return `nil`").Out()
	return nil

}

func (t *Processor) base_V040_ElementMatch() (matched bool, matchedStr string) {
	prefix := t.MyType + ".V040_ElementMatch" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing. Return `true`,\"\"").Out()
	return true, ""
}

func (t *Processor) base_V050_ElementProcessMatched() {
	prefix := t.MyType + ".V050_ElementProcessMatched" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing").Out()
}

func (t *Processor) base_V060_ElementProcessUnmatch() {
	prefix := t.MyType + ".V060_ElementProcessUnmatch" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing").Out()
}

func (t *Processor) base_V070_ElementProcess() {
	prefix := t.MyType + ".V070_ElementProcess" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing").Out()
}

func (t *Processor) base_V080_ElementScrollable() bool {
	prefix := t.MyType + ".V080_ElementScrollable" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing. Return `true`").Out()
	return true
}

func (t *Processor) base_V090_ElementLoopEnd() {
	prefix := t.MyType + ".V090_ElementLoopEnd" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing").Out()
}

func (t *Processor) base_V100_ScrollLoopEnd() {
	prefix := t.MyType + ".V100_ScrollLoopEnd" + "(base)"
	ezlog.Trace().N(prefix).M("Do nothing").Out()
}
