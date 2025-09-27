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
	"github.com/J-Siu/go-ezlog/v2"
	"github.com/J-Siu/go-is"
	"github.com/J-Siu/go-is/example/x-feed/xfp"
	"github.com/go-rod/rod"
)

// (2) Write `main`
func main() {
	ezlog.StrAny.IndentEnable(true)

	// Select log level
	ezlog.SetLogLevel(ezlog.ErrLevel)
	// ezlog.SetLogLevel(ezlog.DebugLevel)
	// ezlog.SetLogLevel(ezlog.TraceLevel)

	var x *xfp.XFeedProcessor

	page, err := getTab("localhost", 9222)
	if err == nil {
		// (2.1) Prepare a `is.Property` object, populate field as needed
		property := is.Property{
			Page:      page,              // (2.1) REQUIRED: populate `Page` field (a `*rod.Page`, representing a browser tab)
			IInfoList: new(is.IInfoList), // Initialize this to use build-in info array
			ScrollMax: 5,                 // number of time we will scroll, -1 for infinite (default: 0)
			UrlLoad:   true,
			UrlStr:    "https://x.com/home",
		}

		// (2.2) Allocate the `processor`
		x = new(xfp.XFeedProcessor)

		// (2.3) Initialize the `processor` struct with the `property`
		x.New(&property)

		// (2.4) Call `Run`
		x.Run()

		err = x.Err
	}
	if err == nil {
		// (2.5) Output result
		x.IInfoList.Print(is.PrintAll)
	}
	if err != nil {
		ezlog.Err().Msg(err).Out()
	}
}

// Helper function to connect to remote/running devtools with host and port only
//
// [dq] is not part of, [IS] package
func getTab(host string, port int) (page *rod.Page, err error) {
	prefix := "GetTab"
	ezlog.Trace().Name(prefix).Msg("Start").Out()
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
	ezlog.Trace().Name(prefix).Msg("End").Out()
	return page, err
}
