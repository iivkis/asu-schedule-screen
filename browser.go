package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"

	"github.com/chromedp/chromedp"
)

const baseURL = "https://www.asu.ru/timetable"

var customizeDevJS string

var (
	browserCtx context.Context
	browserMX  sync.Mutex
)

func init() {
	//customize js
	{
		f, err := os.OpenFile("custom.js", os.O_CREATE|os.O_RDONLY, 0o777)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		codeJS, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		customizeDevJS = fmt.Sprintf("const f = () => {%s}; f();", string(codeJS))
	}

	//create browser context
	{
		opts := make([]chromedp.ExecAllocatorOption, 0)
		opts = append(opts, chromedp.DefaultExecAllocatorOptions[:]...)

		if p, ok := os.LookupEnv("GOOGLE_CHROME_SHIM"); ok {
			opts = append(opts, chromedp.ExecPath(p))
		}

		allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel2 := chromedp.NewContext(allocCtx)
		browserCtx = ctx

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		go func() {
			<-c
			cancel1()
			cancel2()
			os.Exit(0)
		}()
	}
}

func screenLink(link string) (buf []byte, err error) {
	browserMX.Lock()
	defer browserMX.Unlock()

	err = chromedp.Run(browserCtx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("%s/%s", baseURL, link)),
		chromedp.Evaluate(customizeDevJS, nil),
		chromedp.Screenshot("div.l-content-main", &buf, chromedp.NodeVisible),
	})

	if err != nil {
		return
	}

	err = os.WriteFile("last_screen.png", buf, 0o777)
	return buf, err
}
