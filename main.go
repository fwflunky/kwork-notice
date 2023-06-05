package main

import (
	"kworknotice/parser"
	"kworknotice/vkbot"
	"log"
	"strconv"
	"time"
)

func main() {
	vkbot.Initialize("токен группы", 1234567)
	parser.StartBrowser()
	if err := parser.LogIn("email", "pass"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in successfully")
	for {
		if err := parser.OpenProjects(); err != nil {
			log.Fatal(err)
		}
		log.Println("Projects successfully opened")
		if outs, err := parser.GrabAllProjects(); err != nil {
			log.Fatal("Unable to grab projects:", err)
		} else {
			for _, nsp := range parser.GetOnlyNotSeenProjects(outs) {
				nsp.MarkAsSeen()
				vkbot.SendNotice("Новый кворк на бирже:\n\n" + nsp.Title + "\n" + nsp.Link + "\n" + nsp.ReportCount + "\nСтраница: " + strconv.Itoa(nsp.WhatPageWas))
				time.Sleep(1 * time.Second)
			}

		}

		time.Sleep(2 * time.Minute)
	}
}
