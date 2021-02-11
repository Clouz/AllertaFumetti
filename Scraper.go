package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type item struct {
	productItemName               string
	productItemAttribute          string
	productItemAttributeSecondary string
	specialPrice                  string
	productItemInner              string
	productItemPhoto              *url.URL
	productItemHref               *url.URL
}

func (i item) String() string {
	return fmt.Sprintf("%v | %v | %v | %v | %v | %v | %v", i.productItemName, i.productItemAttribute, i.productItemAttributeSecondary, i.specialPrice, i.productItemInner, i.productItemHref, i.productItemPhoto)
}

func main() {

	itemLimit := 500
	itemCount := 0

	stories := []item{}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36"),
		// colly.URLFilters(
		// 	regexp.MustCompile("https://www.panini\\.it/shp_ita_it/fumetti-libri-riviste\\.html?p=/(|e.+)$"),
		// 	regexp.MustCompile("http://httpbin\\.org/h.+"),
		// ),
	)

	// Find and visit all links
	//<div class="product-item-info type1"
	c.OnHTML("div.product-item-info.type1", func(e *colly.HTMLElement) {
		temp := item{}
		temp.productItemName = e.ChildText("h3.product-item-name")
		e.ForEach("div.product-item-attribute", func(i int, el *colly.HTMLElement) {
			if i == 0 {
				temp.productItemAttribute = el.Text
			} else {
				temp.productItemAttributeSecondary = el.Text
			}
		})
		temp.specialPrice = e.ChildText("span.special-price span.price")
		temp.productItemInner = e.ChildText("div.product-item-inner")
		temp.productItemHref, _ = url.Parse(e.ChildAttr("div.product-item-photo a", "href"))
		temp.productItemPhoto, _ = url.Parse(e.ChildAttr("div.product-item-photo img.product-image-photo", "data-src"))

		itemCount++

		stories = append(stories, temp)
		fmt.Println(itemCount, ":\t", temp)

	})

	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		if (itemLimit >= itemCount) || (itemLimit == 0) {
			link := e.Attr("href")
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("----Visiting", r.URL)
	})

	c.Limit(&colly.LimitRule{
		RandomDelay: 5 * time.Second,
	})

	c.Visit("https://www.panini.it/shp_ita_it/fumetti-libri-riviste.html?p=1&product_list_limit=36")

	toCSV(stories)

}

func toCSV(i []item) {
	records := [][]string{}

	for _, v := range i {

		xx := []string{
			v.productItemName,
			v.productItemAttribute,
			v.productItemAttributeSecondary,
			v.specialPrice,
			v.productItemInner,
			v.productItemPhoto.String(),
			v.productItemHref.String(),
		}

		records = append(records, xx)
	}

	f, _ := os.Create("export.csv")
	w := csv.NewWriter(f)
	w.Comma = ';'

	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}
