package indicator

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/developerasun/SignalDash/server/models"
	colly "github.com/gocolly/colly/v2"
)

type IndicatorService struct {
	repo *IndicatorRepo
}

type Crawler interface {
	Init() *colly.Collector
	Scrape(c *colly.Collector) (string, error)
}

type DollarIndex struct{}

func (is *IndicatorService) CrawlAndInsert() error {
	di := DollarIndex{}
	crawler := di.Init()
	__dxy, sErr := di.Scrape(crawler)

	if sErr != nil {
		return sErr
	}

	_dxy := strings.Trim(__dxy, " ")
	log.Printf("CrawlAndInsert:_dxy:%s", _dxy)

	dxy, pfErr := strconv.ParseFloat(_dxy, 64)
	if pfErr != nil {
		return sErr
	}

	if cErr := is.repo.Create(&models.Indicator{
		Name:   "U.S. Dollar Index",
		Ticker: "DXY",
		Value:  dxy,
		Type:   "Fiat",
		Domain: "www.tradingview.com",
	}); cErr != nil {
		return cErr
	}

	return nil
}

func (di DollarIndex) Init() *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains("www.tradingview.com", "tradingview.com"),
		colly.UserAgent("Mozilla/5.0 (compatible; JakeBot/1.0)"),
		colly.IgnoreRobotsTxt(),
	)
}

func (di DollarIndex) Scrape(c *colly.Collector) (string, error) {
	var dxy string
	var hasError error = nil
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
		hasError = err
	})

	if vErr := c.Visit("https://www.tradingview.com/symbols/TVC-DXY/"); vErr != nil {
		hasError = vErr
		return "", hasError
	}

	// @dev wait for jobs to anchor values without error
	c.Wait()

	return dxy, hasError
}
