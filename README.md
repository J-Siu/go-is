# IS - Infinite Scroll

An infinite scroll processing package using [go-rod/rod](https://github.com/go-rod/rod).

## Assumption and Concept

## Implementation

### Property Structure

### IInfo and IInfoList Interfaces

### Processor Structure

## How to Use

### The V* field func

The build-in field functions.

#### Responsibilities of the Element Func

Field Function|Responsibility
--|--
`V030_ElementInfo` | Extract information from [element] and put into an [IInfo] structure and return it.
`V040_ElementMatch` | Determine [element] is a match or not base on [info]
`V050_ElementProcessMatched`| Do some processing (eg, print, write to file, db, etc) if [element] is a match
`V060_ElementProcessUnmatch`| Do some processing (eg, print, write to file, db, etc) if [element] is not a match
`V070_ElementProcess`| Do some processing (eg, print, write to file, db, etc) regardless [element] is a match or not

### Use What Is Needed

### Change Log

- v1.0.0
  - Initial commit

### License

The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
