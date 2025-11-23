package toolkit

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	colly "github.com/gocolly/colly/v2"
)

type Crawler interface {
	Init() *colly.Collector
	Scrape(c *colly.Collector) (string, error)
}

type DollarIndex struct{}

func (di DollarIndex) Init() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("www.tradingview.com", "tradingview.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; JakeBot/1.0)"),
		colly.IgnoreRobotsTxt(),
	)
}

func (di DollarIndex) Scrape(c *colly.Collector) (string, error) {
	var dxy string
	c.OnHTML("section[data-an-section-id=symbol-overview-page-section]", func(e *colly.HTMLElement) {
		substringToReplace := "The current value of U.S. Dollar Index is"
		expression := `(\d+\.\d+)`
		full := fmt.Sprintf("%s %s", substringToReplace, expression)
		target := regexp.MustCompile(full)

		match := target.FindString(e.Text)
		extracted := strings.Replace(match, substringToReplace, "", 1)
		log.Println("match: ", match, "extracted: ", extracted)

		dxy = extracted
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error:", err)
	})

	if err := c.Visit("https://www.tradingview.com/symbols/TVC-DXY/"); err != nil {
		return "", err
	}

	c.Wait()

	return dxy, nil
}

func CrawlDollarIndex() (string, error) {
	di := DollarIndex{}
	c := di.Init()
	return di.Scrape(c)
}

type DaiCoin struct{}

func (dc DaiCoin) Init() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("metamask.io"),
		colly.UserAgent("Mozilla/5.0 (compatible; JakeBot/1.0)"),
		colly.IgnoreRobotsTxt(),
	)
}

func (dc DaiCoin) Scrape(c *colly.Collector) (string, error) {
	// @dev register callbacks first. colly will run the callback recursively for the target query selector
	var prices []string
	c.OnHTML(`[class*="token-price_price"]`, func(h *colly.HTMLElement) {
		log.Println(h.Text)
		prices = append(prices, h.Text)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error:", err)
	})

	// @dev visit the target
	if err := c.Visit("https://metamask.io/price/dai"); err != nil {
		return "", err
	}

	// @dev wait for jobs to anchor values without error
	c.Wait()

	// @dev return result
	dai := prices[0]
	return dai, nil
}

func CrawlDaiPrice() (string, error) {
	dc := DaiCoin{}
	c := dc.Init()
	return dc.Scrape(c)
}
