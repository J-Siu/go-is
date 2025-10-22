# IS - an Infinite Scroll processor

An infinite scroll processing package using [go-rod/rod](https://github.com/go-rod/rod).

- [Creating an IS app](#creating-an-is-app)
  - [Example](#example)
- [How to Use](#how-to-use)
  - [(1.1) Create Your Info Struct](#11-create-your-info-struct)
  - [(1.2) Create Your Processor Struct](#12-create-your-processor-struct)
  - [(1.3.3) Override Process Struct field functions](#133-override-process-struct-field-functions)
  - [(2.1) Property Struct](#21-property-struct)
  - [(2.4) Processing Flow inside Run()](#24-processing-flow-inside-run)
- [Logging](#logging)
- [Use What is Needed](#use-what-is-needed)
  - [Info and IInfoList](#info-and-iinfolist)
  - [The Element Functions](#the-element-functions)
  - [Change Log](#change-log)
  - [License](#license)

<!-- more -->

## Creating an IS app

A basic workflow of creating a infinite scroll processor with `IS`.

- (1) Create a package
  - [(1.1)](#11-create-your-info-struct) Write a `info` struct
    - REQUIRED: embed `is.InfoBase` for `is.IInfo` interface
    - Add additional field as needed
    - Override `String()` member function if needed
  - [(1.2)](#12-create-your-processor-struct) Write a `processor` struct
    - REQUIRED: embed `*is.Processor`
    - Add additional field as needed
  - (1.3) Write package/struct level `New` function
    - (1.3.1) REQUIRED: `*is.Property` as one of its arguments
    - (1.3.2) REQUIRED: use `is.New(property)` to initialize embedded `*is.Processor`
    - [(1.3.3)](#133-override-process-struct-field-functions) Override `is.Processor` field functions as needed
      - As a bare minimum, MUST override `V020_Elements` and `V030_ElementInfo`. Else `is.Run` will do nothing.
- (2) Write `main`
  - (2.1) Prepare a `is.Property` object, populate field as needed
    - REQUIRED: populate `Page` field (a `*rod.Page`, representing a browser tab)
    - Set `UrlLoad`, `true` to load page at `UrlStr`. (Default: `false`)
    - Set `UrlStr` to target site address. Not required if `UrlLoad` is `false`
  - (2.2) Allocate the `processor`
  - (2.3) Initialize the `processor` struct with the `property`
  - [(2.4)](#24-processing-flow-inside-run) Call `Run`
  - (2.5) Output result

### Example

- [x-feed](/example/x-feed/) in [example](/example/) - X.com feed processing using [IS]. With comment referencing workflow above.
- [yt-toolbox](https://github.com/J-Siu/yt-toolbox) - A more elaborate [IS] command line application.

## How to Use

### (1.1) Create Your Info Struct

The `info` struct and `IInfoList` provide a basic means to store and process information during `Run()`.

[xfp.go](/example/x-feed/xfp/xfp.go):

```go
// (1.1) Write a `info` struct
type XFeedInfo struct {
  is.InfoBase // (1.1) REQUIRED: embed [is.InfoBase] to get [is.IInfo] interface

  // Added fields
  User string `json:"user,omitempty"`
  Text string `json:"text,omitempty"`
}
```

The `is.InfoBase` implemented the `is.IInfo` interface functions:

Function|Description|Override Required
--|--|--
Matched() bool | Getter, return value of `matched`| No
MatchedStr() string | Getter, return value if `matchedStr` | No
SetMatched(matched bool) | Setter, set value of `matched`| No
SetMatchedStr(matchedStr string) | Setter, set value of `matchedStr`|No
String() string | Info struct to string| As needed

The `is.IInfo` allow info struct to be passed between the processor's `V*` field functions in `Run()`.

Add fields to the struct to store information.

### (1.2) Create Your Processor Struct

[xfp.go](/example/x-feed/xfp/xfp.go):

```go
// (1.2) Write a `processor` struct
type XFeedProcessor struct {
  *is.Processor // (1.2) REQUIRED: embed `*is.Processor`
}
```

### (1.3.3) Override Process Struct field functions

[xfp.go](/example/x-feed/xfp/xfp.go):

```go
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
```

`is.Processor` comes with 14 field functions:

Function|Description|Override Required
--|--|--
LoadPage func() | Load `UrlStr` | No
ScrollCalculation func(state *State) (scroll bool) | Detect end of page | No
ScrollElement func(element *rod.Element) | Use `rod.element.MustScrollIntoView` for scrolling | No
V010_Container func() (container *rod.Element) | Return a `container` element. (default: `Property.Container`) | As needed
V020_Elements func(container *rod.Element) *rod.Elements | Return collection of repeating elements in `container` from `V010_Container` (default: `nil`) | Yes
V030_ElementInfo func(element *rod.Element, index int) (info IInfo) |Extract information from `element`, and put them into an [IInfo] structure, and return it. (default: `nil) | Yes
V040_ElementMatch func(element *rod.Element, index int, info IInfo) (matched bool, matchedStr string)|Determine `element` is a match or not base on `info` (default: `true`, `""`)| As needed
V050_ElementProcessMatched func(element *rod.Element, index int, info IInfo)|Do some processing (eg, print, write to file, db, etc) if `element` is a match (default: do nothing)|As needed
V060_ElementProcessUnmatch func(element *rod.Element, index int, info IInfo)|Do some processing if `element` is not a match (default: do nothing)|As needed
V070_ElementProcess func(element *rod.Element, index int, info IInfo)|Do some processing regardless of `element` is a match or not (default: do nothing)|As needed
V080_ElementScrollable func(element *rod.Element, index int, info IInfo) bool|Determine if `element` is scrollable (default: true)|As needed (eg. `element` removed from DOM)
V090_ElementLoopEnd func(element *rod.Element, index int, info IInfo)|Do some processing if required (default: do nothing)|As needed
V100_ScrollLoopEnd func(state *State)|Do some processing if required (default: do nothing)|As needed

### (2.1) Property Struct

### (2.4) Processing Flow inside Run()

Following is pseudo code of `is.Processor.Run()`. Full code is [here](/is.go).

```go
Run() {
  state := new(State).New()
  LoadPage()
  Container = V010_Container()
  for {
    // -- SCROLL LOOP - START
    if ScrollLoopBreak(state) { break }
    ScrollElement(state.ElementLast)
    elements = V020_Elements(Container)
    for element(new ones after scroll) in elements {
      // -- ELEMENTS LOOP - END
      info := V030_ElementInfo(element, index)
      matched, matchedStr := V040_ElementMatch(element, index, info)
      if matched {
        V050_ElementProcessMatched(element, index, info)
      } else {
        V060_ElementProcessUnmatch(element, index, info)
      }
      V070_ElementProcess(element, index, info)
      if IInfoList != nil && info != nil { append(IInfoList, info) }
      V080_ElementScrollable(element, index, info) { update state }
      V090_ElementLoopEnd(element, index, info)
      // -- ELEMENTS LOOP - END
    }
    ScrollCalculation(state)
    V100_ScrollLoopEnd(state)
    // -- SCROLL LOOP - END
  }
}
```
## Logging

## Use What is Needed

### Info and IInfoList

The `info` struct and `IInfoList` provide a basic means to store and process information during `Run()`.

### The Element Functions

`V030_ElementInfo`, `V040_ElementMatch`, `V050_ElementProcessMatched`, `V060_ElementProcessUnmatch`, `V070_ElementProcess`

### License

The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
