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

package main

import (
	"errors"

	"github.com/J-Siu/go-dtquery/dq"
	"github.com/J-Siu/go-ezlog"
	"github.com/J-Siu/go-is"
	"github.com/J-Siu/go-is/example/xfeed"
	"github.com/go-rod/rod"
)

func main() {
	ezlog.SetAllPrintln() // Setup ezlog print functions

	// Select log level
	ezlog.SetLogLevel(ezlog.ErrLevel)
	// ezlog.SetLogLevel(ezlog.DebugLevel)
	// ezlog.SetLogLevel(ezlog.TraceLevel)

	var xf *xfeed.XFeed

	// [is] start at [rod.Page] level
	page, err := getTab("localhost", 9222)
	if err == nil {
		property := is.Property{
			IInfoList: new(is.IInfoList),
			Page:      page,
			ScrollMax: 3, // number of time we will scroll
			UrlLoad:   true,
			UrlStr:    "https://x.com/home",
		}
		xf = xfeed.New(&property)
		xf.Run()
		err = xf.Err
	}
	if err == nil {
		xf.IInfoList.Print(is.PrintAll)
	}
	if err != nil {
		ezlog.Err(err.Error())
	}
}

// use [dq] to get devtools info, then set it up with [rod]
func getTab(host string, port int) (page *rod.Page, err error) {
	prefix := "GetTab"
	ezlog.Trace(prefix + ": Start")
	var (
		browser *rod.Browser
		pages   rod.Pages
	)

	// use [dq] to get devtools info
	devtools := dq.Get(host, port)
	err = devtools.Err

	// setup [rod]
	if err == nil {
		browser = rod.New().ControlURL(devtools.Ver.WsUrl)
		err = browser.Connect()
	}
	if err == nil {
		pages, err = browser.NoDefaultDevice().Pages()
	}
	if err == nil {
		page = pages.First()
		if page != nil {
			page.Activate()
		}
	}
	if err != nil {
		err = errors.New(prefix + ": Err: " + err.Error())
	}
	ezlog.Trace(prefix + ": End")
	return page, err
}
