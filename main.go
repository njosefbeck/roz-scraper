package main

import (
	"fmt"
	"strings"
	"time"
	"os"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/yosssi/gohtml"
)

func replaceSlashWithDash(str string) string {
	return strings.Replace(str, "/", "-", -1)
}

func replaceSpaceWithDash(str string) string {
	return strings.Replace(str, " ", "-", -1)
}

func replaceSpaceWithPeriod(str string) string {
	return strings.Replace(str, ".", "-", -1)
}

func formatString(str string) string {
	return replaceSlashWithDash(replaceSpaceWithDash(replaceSpaceWithPeriod(str)))
}

func removeAttrs(str string) string {
	attrs := `(\sstyle="[^"]*")|(\sclass="[^"]*")|(\sheight="[^"]*")|(\swidth="[^"]*")|(\sid="[^"]*")|(\salt="[^"]*")|(\sborder="[^"]*")|(\sname="[^"]*")|(\shref="[^"]*")|(\sspan="[^"]*")|(\salign="[^"]*")|(\sscope="[^"]*")|(\slang="[^"]*")|(\sgothic="[^"]*")|(\snanum="[^"]*")|(\snowrap="[^"]*")`
	re := regexp.MustCompile(attrs)
	return re.ReplaceAllString(str, "")
}

func removeDuplicates(str string) string {
	spans := `<span><span>|</span></span>|<span><span><span>|</span></span></span>|<span><span><span><span><span>|</span></span></span></span></span>|<section>\n<section>\n<section>\n<section>\n<section>\n|</section>\n</section>\n</section>\n</section>\n</section>\n|<section>\n<section>\n|</section>\n</section>\n`
	re := regexp.MustCompile(spans)
	return re.ReplaceAllString(str, "")
}

func removeEmpties(str string) string {
	re := regexp.MustCompile(`<p> </p>|<h2> </h2>|<col/>|<col />|<colgroup>|</colgroup>|<a></a>|<h1> </h1>|<span> </span>`)
	return re.ReplaceAllString(str, "")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
    // Set a delay between requests
    Delay: 1 * time.Second,
    // Add an additional random delay
    RandomDelay: 1 * time.Second,
	})

	pageCollector := c.Clone()

	c.OnHTML(".board_list.update table tbody tr", func(e *colly.HTMLElement) {
		number := e.ChildText("td:first-child")
		date := formatString(e.ChildText(".date"))
		category := formatString(e.ChildText(".icon"))
		title := formatString(e.ChildText(".title span"))
		url := "http://roz.gnjoy.com/news/update/" + e.ChildAttr(".title a[href]", "href")

		var pieces = []string{number, date, category, title}
		fileName := strings.Join(pieces, "_") + ".html"

		e.Request.Ctx.Put("fileName", fileName)

		pageCollector.Request("GET", url, nil, e.Request.Ctx, nil)
	})

	c.OnHTML(".pageing a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		e.Request.Visit(link)
	})

	pageCollector.OnHTML(".board_view.notice ul ", func(e *colly.HTMLElement) {
		selection, err := e.DOM.Html()
		check(err)

		stripped := removeDuplicates(removeEmpties(removeAttrs(selection)))
		formatted := gohtml.Format(stripped)

		fileName := e.Request.Ctx.Get("fileName")
		path := "/Users/nathanbeck/Sites/roz-scraper/html/" + fileName
		file, err := os.Create(path)
		check(err)

		defer file.Close()

		numBytes, err := file.WriteString(formatted)
		check(err)

		fmt.Println(numBytes, "bytes written successfully to", fileName)
	})

	c.Visit("http://roz.gnjoy.com/news/update/list.asp")
}
