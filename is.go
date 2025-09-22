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
	"strconv"

	"github.com/J-Siu/go-basestruct"
	"github.com/J-Siu/go-ezlog"
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

// IS processor structure
type Processor struct {
	*basestruct.Base
	*Property

	// -- Following 4 field func rarely need override

	// Load [UrlStr] into [Page]
	//
	// No override needed.
	LoadPage func() `json:"-"`

	// Determine whether the scroll loop should continue running
	//
	// No override needed.
	ScrollLoopBreak func(state *State) bool `json:"-"`

	// Detect end of page, scroll no longer possible.
	//
	// No override needed.
	//
	// If elements are removed during [Run()], overload [V100_ExitScroll()] to do custom override.
	// As both of following checks can be flawed if elements are removed from page DOM.
	ScrollCalculation func(state *State) (scroll bool) `json:"-"`

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
	V030_ElementInfo func(element *rod.Element, index int) (info IInfo) `json:"-"`

	// Determine [element] is a match or not base on [info]
	//
	// build-in behavior is to return (`true`, `""`)
	//
	// Override if needed
	V040_ElementMatch func(element *rod.Element, index int, info IInfo) (matched bool, matchedStr string) `json:"-"`

	// Do some processing (eg, print, write to file, db, etc) if [element] is a match
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V050_ElementProcessMatched func(element *rod.Element, index int, info IInfo) `json:"-"`

	// Do some processing (eg, print, write to file, db, etc) if [element] is not a match
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V060_ElementProcessUnmatch func(element *rod.Element, index int, info IInfo) `json:"-"`

	// Do some processing (eg, print, write to file, db, etc) regardless [element] is a match or not
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V070_ElementProcess func(element *rod.Element, index int, info IInfo) `json:"-"`

	// Determine if an element is scrollable
	//
	// build-in behavior is to return `true``
	//
	// Override if needed
	V080_ElementScrollable func(element *rod.Element, index int, info IInfo) bool `json:"-"`

	// Do some processing if required
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V090_ElementLoopEnd func(element *rod.Element, index int, info IInfo) `json:"-"`

	// Do some processing if required
	//
	// build-in behavior is to do nothing
	//
	// Override if needed
	V100_ScrollLoopEnd func(state *State) `json:"-"`
}

// Implement the default field functions
func (p *Processor) setFunc() {
	// -- Following 4 field func rarely need override
	{

		p.LoadPage = func() {
			prefix := p.MyType + ".LoadPage" + "(base)"
			ezlog.Trace(prefix + ": Start")
			if p.CheckErrInit(prefix) {
				if p.UrlLoad {
					ezlog.Debug(prefix + ": urlStr: " + p.UrlStr)
					p.Err = p.Page.Navigate(p.UrlStr)
					if p.Err == nil {
						ezlog.Trace(prefix + ": MustWaitDOMStable: Start")
						p.Page.MustWaitDOMStable()
						ezlog.Trace(prefix + ": MustWaitDOMStable: Start")
					}
				}
				if p.Err != nil {
					p.Err = errors.New(prefix + ": " + p.Err.Error())
					ezlog.Err(p.Err.Error())
				}
			}
			ezlog.Trace(prefix + ": End")
		}

		p.ScrollCalculation = func(state *State) (scroll bool) {
			prefix := p.MyType + ".ScrollCalculation (base)"
			scroll = state.Scroll
			/*
				"if (elementsCount == elementsCountLast)":
					will be triggered, if number of elements removed
							= number of new elements added after scroll

				"if (ElementLastScroll == ElementLast)":
					will be triggered, if all new elements added after scroll are removed
			*/
			if state.ElementLastScroll != nil && state.ElementLast != nil {
				if state.ElementLastScroll.Object.ObjectID == state.ElementLast.Object.ObjectID {
					scroll = false
				}
			} else if state.ElementLastScroll == state.ElementLast {
				scroll = false
			}
			ezlog.Trace(prefix + ": MustWaitDOMStable: Done")
			return scroll
		}

		p.ScrollElement = func(element *rod.Element) {
			prefix := p.MyType + ".ScrollElement (base)"
			ezlog.Trace(prefix + ": Start")
			if element != nil {
				element.MustScrollIntoView()
				ezlog.Trace(prefix + ": Scrolled")
				ezlog.Trace(prefix + ": MustWaitDOMStable: Start")
				p.Page.MustWaitDOMStable()
				ezlog.Trace(prefix + ": MustWaitDOMStable: End")
			}
			ezlog.Trace(prefix + ": End")
		}

		p.ScrollLoopBreak = func(state *State) bool {
			breakLoop := !(state.Scroll && (state.ScrollCount <= p.ScrollMax || p.ScrollMax < 0))
			if ezlog.LogLevel() == ezlog.TraceLevel {
				msg := "scrollMax: " + strconv.Itoa(p.ScrollMax) + "\n"
				msg += "breakLoop: " + "!(state.Scroll && (state.ScrollCount <= scrollMax || scrollMax < 0)) = " + strconv.FormatBool(breakLoop)
				ezlog.Trace(p.MyType + ".ScrollLoopBreak (base):")
				ezlog.TraceP(state.String())
				ezlog.TraceP(&msg)
			}
			return breakLoop
		}
	}
	// --- Overload following field func as needed
	{
		p.V010_Container = func() (container *rod.Element) {
			prefix := p.MyType + ".V010_Container" + "(base)"
			ezlog.Trace(prefix + ": Done")
			return p.Container
		}

		p.V020_Elements = func(container *rod.Element) *rod.Elements {
			prefix := p.MyType + ".V020_Elements" + "(base)"
			ezlog.Trace(prefix + ": Do nothing. Return 'nil'.")
			return nil
		}

		p.V030_ElementInfo = func(element *rod.Element, index int) (info IInfo) {
			prefix := p.MyType + ".V030_ElementInfo" + "(base)"
			ezlog.Trace(prefix + ": Do nothing. Return 'nil'.")
			return nil

		}
		p.V040_ElementMatch = func(element *rod.Element, index int, info IInfo) (matched bool, matchedStr string) {
			prefix := p.MyType + ".V040_ElementMatch" + "(base)"
			ezlog.Trace(prefix + ": return true, \"\"")
			return true, ""
		}

		p.V050_ElementProcessMatched = func(element *rod.Element, index int, info IInfo) {
			prefix := p.MyType + ".V050_ElementProcessMatched" + "(base)"
			ezlog.Trace(prefix + ": Do nothing")
		}

		p.V060_ElementProcessUnmatch = func(element *rod.Element, index int, info IInfo) {
			prefix := p.MyType + ".V060_ElementProcessUnmatch" + "(base)"
			ezlog.Trace(prefix + ": Do nothing")
		}

		p.V070_ElementProcess = func(element *rod.Element, index int, info IInfo) {
			prefix := p.MyType + ".V070_ElementProcess" + "(base)"
			ezlog.Trace(prefix + ": Do nothing")
		}

		p.V080_ElementScrollable = func(element *rod.Element, index int, info IInfo) bool {
			prefix := p.MyType + ".V080_ElementScrollable" + "(base)"
			ezlog.Trace(prefix + ": Do nothing. Return 'true'.")
			return true
		}

		p.V090_ElementLoopEnd = func(element *rod.Element, index int, info IInfo) {
			prefix := p.MyType + ".V090_ElementLoopEnd" + "(base)"
			ezlog.Trace(prefix + ": Do nothing")
		}

		p.V100_ScrollLoopEnd = func(state *State) {
			prefix := p.MyType + ".V100_ScrollLoopEnd" + "(base)"
			ezlog.Trace(prefix + ": Do nothing")
		}
	}
}

// Process the page
//
// No override needed.
func (p *Processor) Run() {
	prefix := p.MyType + ".Run (base)"

	if !p.CheckErrInit(prefix) {
		return
	}

	state := new(State).New()
	p.LoadPage()
	if p.Err != nil {
		return // LoadPage failed
	}
	p.Container = p.V010_Container()
	// -- Scroll Loop
	for {
		// -- SCROLL LOOP - START
		if p.ScrollLoopBreak(state) {
			break // exit Run()
		}
		p.ScrollElement(state.ElementLast)

		// -- Get elements
		elementsCount := 0
		state.ElementLastScroll = state.ElementLast
		state.Elements = p.V020_Elements(p.Container)

		if state.Elements == nil {
			state.Scroll = false // no element, no scroll
		} else {
			elementsCount = len(*state.Elements)
			ezlog.Debug(prefix + ": elementCount: " + strconv.Itoa(elementsCount))
			for index := state.ElementCountLast; index < elementsCount; index++ {
				// -- ELEMENTS LOOP - START
				element := (*state.Elements)[index]
				info := p.V030_ElementInfo(element, index)
				matched, matchedStr := p.V040_ElementMatch(element, index, info)
				if info != nil {
					// Remove burden from package user
					info.SetMatched(matched)
					info.SetMatchedStr(matchedStr)
				}
				if matched {
					p.V050_ElementProcessMatched(element, index, info)
				} else {
					p.V060_ElementProcessUnmatch(element, index, info)
				}
				p.V070_ElementProcess(element, index, info)
				if p.IInfoList != nil && info != nil {
					tmp := append(*p.IInfoList, info)
					p.IInfoList = &tmp
				}
				if p.V080_ElementScrollable(element, index, info) {
					state.ElementLast = element
					state.InfoLast = info
				}
				p.V090_ElementLoopEnd(element, index, info)
				// -- ELEMENTS LOOP - END
			}
			state.Scroll = p.ScrollCalculation(state)
			state.ElementCountLast = elementsCount
		}

		p.V100_ScrollLoopEnd(state)
		state.ScrollCount++
		// -- SCROLL LOOP - END
	}
}

// Parameters:
//   - property *Property
//
// Returns:
//   - *Processor
func New(property *Property) *Processor {
	p := new(Processor)
	p.Base = new(basestruct.Base)
	p.MyType = "is.Processor"

	if property == nil {
		p.Err = errors.New("is.New: property cannot be nil")
	}
	if p.Err == nil && property.Page == nil {
		p.Err = errors.New("is.New: page/tab cannot be nil")
	}
	if p.Err == nil {
		p.Property = property
		p.setFunc()
		p.Initialized = true
	}
	ezlog.Trace("is.New(): Done")
	return p
}
