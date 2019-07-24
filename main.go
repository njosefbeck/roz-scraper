package main

import (
	"fmt"
	"strings"
	"github.com/gocolly/colly"
)

func formatString(str string) string {
	var withoutForwardSlash = strings.Replace(str, "/", "-", -1)
	var withoutSpaces = strings.Replace(withoutForwardSlash, " ", "-", -1)
	var withoutPeriods = strings.Replace(withoutSpaces, ".", "-", -1)

	return withoutPeriods
}

func main() {
	c := colly.NewCollector(
		colly.CacheDir("./roz_updates"),
	)

	//pageCollector := c.Clone()

	c.OnHTML(".board_list.update table tbody tr", func(e *colly.HTMLElement) {
		number := e.ChildText("td:first-child")
		date := formatString(e.ChildText(".date"))
		category := formatString(e.ChildText(".icon"))
		title := formatString(e.ChildText(".title span"))

		var pieces = []string{number, date, category, title}

		fileName := strings.Join(pieces, "_") + ".html"

		fmt.Println(fileName)
	})

	c.Visit("http://roz.gnjoy.com/news/update/list.asp")
}
