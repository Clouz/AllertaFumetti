package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type item struct {
	productItemName               string
	productItemAttribute          string
	productItemAttributeSecondary string
	specialPrice                  string
	productItemInner              string
}

func (i item) String() string {
	return fmt.Sprintf("%v | %v | %v | %v | %v", i.productItemName, i.productItemAttribute, i.productItemAttributeSecondary, i.specialPrice, i.productItemInner)
}

func main() {
	//stories := []item{}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36"),
	)

	// Find and visit all links
	//<div class="product details product-item-details">
	//<div class="product-item-info type1"
	c.OnHTML("div.product-item-info.type1", func(e *colly.HTMLElement) {
		//fmt.Println(e.Text)
		temp := item{}
		temp.productItemName = e.ChildText("h3.product-item-name")
		e.ForEach("div.product-item-attribute", func(i int, el *colly.HTMLElement) {
			if i == 0 {
				temp.productItemAttribute = el.Text
			} else {
				temp.productItemAttributeSecondary = el.Text
			}
		})
		//temp.productItemAttribute = e.ChildText("div.product-item-attribute")
		//temp.productItemAttributeSecondary = e.ChildText("div.product-item-attribute-secondary")
		temp.specialPrice = e.ChildText("span.special-price span.price")
		temp.productItemInner = e.ChildText("div.product-item-inner")
		//stories = append(stories, temp)
		fmt.Println(temp)

	})

	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("----Visiting", r.URL)
	})

	c.Visit("https://www.panini.it/shp_ita_it/fumetti-libri-riviste.html?p=1&product_list_limit=36")
}
