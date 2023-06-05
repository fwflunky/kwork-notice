package parser

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"strconv"
	"time"
)

var runner *playwright.Playwright
var browser playwright.Browser
var context playwright.BrowserContext
var page playwright.Page

func StartBrowser() {
	_ = playwright.Install()
	runner, _ = playwright.Run()
	browser, _ = runner.Chromium.Launch()
	context, _ = browser.NewContext()
}

func LogIn(email, pass string) error {
	page, _ = context.NewPage()
	if _, err := page.Goto("https://kwork.ru/login", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		return err
	}
	for i := 0; i <= 10; i++ {
		if list, _ := page.QuerySelectorAll(".input-style.input-style--focus-blue.wMax.keep-placeholder"); len(list) == 2 {
			break
		} else if i == 10 {
			return fmt.Errorf("%s", "timeout while getting login form")
		}
		time.Sleep(1 * time.Second)
	}

	_, _ = page.Evaluate(`document.querySelectorAll(".input-style.input-style--focus-blue.wMax.keep-placeholder")[0].value = "` + email + `"`)
	_, _ = page.Evaluate(`document.querySelectorAll(".input-style.input-style--focus-blue.wMax.keep-placeholder")[0].dispatchEvent(new Event('input', {bubbles:true}));`)

	_, _ = page.Evaluate(`document.querySelectorAll(".input-style.input-style--focus-blue.wMax.keep-placeholder")[1].value = "` + pass + `"`)
	_, _ = page.Evaluate(`document.querySelectorAll(".input-style.input-style--focus-blue.wMax.keep-placeholder")[1].dispatchEvent(new Event('input', {bubbles:true}));`)

	isChanged := make(chan bool)
	page.Once("framenavigated", func(frame playwright.Frame) {
		if frame.URL() == "https://kwork.ru/seller" {
			isChanged <- true
		}
	})
	_, _ = page.Evaluate(`document.querySelectorAll(".kw-button--green")[4].click()`)

	select {
	case <-time.After(5 * time.Second):
		return fmt.Errorf("%s. %s", "timeout while waiting for seller page after login click", "Wrong password?")
	case <-isChanged:
		return nil
	}
}

func WaitForCards() bool {
	for i := 0; i <= 10; i++ {
		if list, _ := page.QuerySelectorAll(".card.want-card.js-want-container"); len(list) > 0 {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func OpenProjects() error {
	if _, err := page.Goto("https://kwork.ru/projects?a=1&page=1", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		return err
	}

	if !WaitForCards() {
		return fmt.Errorf("%s", "timeout while loading projects")
	}
	return nil
}

func GrabAllProjects() ([]Project, error) {
	var outs []Project

	out, _ := page.Evaluate(`document.querySelector(".mb10.pb20.t-align-c").childNodes[0].childNodes[0].childNodes[0].childNodes[document.querySelector(".mb10.pb20.t-align-c").childNodes[0].childNodes[0].childNodes[0].childNodes.length - 2].innerText`)
	maxPage, _ := strconv.Atoi(out.(string))

	for currentPage := 1; currentPage <= maxPage; currentPage++ {
		nextPage := "https://kwork.ru/projects?a=1&page=" + strconv.Itoa(currentPage)
		if page.URL() != nextPage {
			if _, err := page.Goto(nextPage, playwright.PageGotoOptions{
				WaitUntil: playwright.WaitUntilStateDomcontentloaded,
			}); err != nil {
				return nil, err
			}
		}
		if !WaitForCards() {
			return nil, fmt.Errorf("%s %d", "timeout while loading projects on page", currentPage)
		}
		cards, _ := page.QuerySelectorAll(".card.want-card.js-want-container")
		for current, c := range cards {
			var project Project
			project.WhatPageWas = currentPage
			project.ID, _ = c.GetAttribute("data-id")
			project.Link = "https://kwork.ru/projects/" + project.ID
			out, _ = page.Evaluate(`document.querySelectorAll(".card.want-card.js-want-container")[` + strconv.Itoa(current) + `].querySelector(".wants-card__header-title.first-letter.breakwords.pr250").childNodes[0].innerText`)
			if out != nil {
				project.Title = out.(string)
			}
			out, _ = page.Evaluate(`document.querySelectorAll(".card.want-card.js-want-container")[` + strconv.Itoa(current) + `].querySelector(".force-font.force-font--s12.mr8").childNodes[2].innerText`)
			if out != nil {
				project.ReportCount = out.(string)
			}
			outs = append(outs, project)
		}
	}
	return outs, nil
}
