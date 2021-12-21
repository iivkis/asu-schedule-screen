package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

const baseURL = "https://www.asu.ru/timetable"

var customize string

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

		customize = fmt.Sprintf("const f = () => {%s}; f();", string(codeJS))
	}
}

func screenLink(link string) (buf []byte, err error) {
	opts := make([]chromedp.ExecAllocatorOption, 0)
	opts = append(opts, chromedp.DefaultExecAllocatorOptions[:]...)

	fmt.Println(os.LookupEnv("GOOGLE_CHROME_SHIM"))

	if p, ok := os.LookupEnv("GOOGLE_CHROME_SHIM"); ok {
		opts = append(opts, chromedp.ExecPath(p))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("%s/%s", baseURL, link)),
		chromedp.Evaluate(customize, nil),
		chromedp.Screenshot("div.l-content-main", &buf, chromedp.NodeVisible),
	})

	if err != nil {
		return
	}

	// err = os.WriteFile("last_screen.png", buf, 0o777)

	return buf, err
}
