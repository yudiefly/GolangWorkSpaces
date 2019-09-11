package main

import (
	"log"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)"),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*.douban.*", Parallelism: 5})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Fatal(err)
		}

		nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)
		for _, node := range nodes {
			url := htmlquery.FindOne(node, "./a/@href")
			title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)
			log.Println(strings.Split(htmlquery.InnerText(url), "/")[4], htmlquery.InnerText(title))
		}
	})

	c.OnHTML(".paginator a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.Visit("https://movie.douban.com/top250?start=0&filter=")
	c.Wait()
}
