package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chromedp/chromedp"
)

const baseURL = "https://www.asu.ru/timetable"

var customize string

func init() {
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

func screenLink(link string) (buf []byte, err error) {
	ctx, cancel := chromedp.NewContext(context.Background())
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
