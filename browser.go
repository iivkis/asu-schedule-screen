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

var browser struct {
	Ctx context.Context
	Mx  sync.Mutex
}

func init() {
	//load customize js
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
		opts := []chromedp.ExecAllocatorOption{}
		opts = append(opts, chromedp.DefaultExecAllocatorOptions[:]...)

		if p, ok := os.LookupEnv("GOOGLE_CHROME_SHIM"); ok {
			opts = append(opts, chromedp.ExecPath(p))
		}

		var cancels []context.CancelFunc
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		cancels = append(cancels, cancel)

		ctx, cancel = chromedp.NewContext(ctx)
		cancels = append(cancels, cancel)

		browser.Ctx = ctx

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		go func() {
			<-c
			for _, cancel := range cancels {
				cancel()
			}
			os.Exit(0)
		}()
	}
}

func screenLink(link string) (buf []byte, err error) {
	browser.Mx.Lock()
	defer browser.Mx.Unlock()

	err = chromedp.Run(browser.Ctx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("%s/%s", baseURL, link)),
		chromedp.Evaluate(customizeDevJS, nil),
		chromedp.Screenshot(".l-content-main", &buf, chromedp.NodeVisible),
	})

	if err != nil {
		return
	}

	err = os.WriteFile("last_screen.png", buf, 0o777)
	return buf, err
}
