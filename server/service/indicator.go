package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/developerasun/SignalDash/server/models"
	"github.com/developerasun/SignalDash/server/sderror"
	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"
)

type indicator struct {
	crawler *colly.Collector
}

type Indicator interface {
	ScrapeDollarIndex() (dxy string, err error)
}

func NewIndicator(domains []string, botHeader string) Indicator {
	_crawler := NewCrawler(domains, botHeader)

	return &indicator{
		crawler: _crawler,
	}
}

func (i indicator) ScrapeDollarIndex() (string, error) {
	c := i.crawler

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

func FindLatestDollarIndex(db *gorm.DB) (*models.Indicator, error) {
	ctx := context.Background()
	indicator, lErr := gorm.G[models.Indicator](db).Last(ctx)
	if errors.Is(lErr, gorm.ErrRecordNotFound) {
		return nil, sderror.ErrEmptyStorage
	}

	return &indicator, nil
}

func CreateDollarIndex(db *gorm.DB, __dxy string) error {
	_dxy := strings.Trim(__dxy, " ")
	log.Printf("CreateDollarIndex:_dxy:%s", _dxy)

	dxy, pfErr := strconv.ParseFloat(_dxy, 64)
	if pfErr != nil {
		return pfErr
	}

	cErr := db.Create(&models.Indicator{
		Name:   "U.S. Dollar Index",
		Ticker: "DXY",
		Value:  dxy,
		Type:   "Fiat",
		Domain: "www.tradingview.com",
	}).Error

	if cErr != nil {
		log.Println("CreateDollarIndex: failed to create a new Dxy record")
		return sderror.ErrInternalServer
	}

	return nil
}

// ================================================================== //
// ============================== deps ============================== //
// ================================================================== //

func NewCrawler(domains []string, botHeader string) *colly.Collector {
	return colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.UserAgent(botHeader),
		colly.IgnoreRobotsTxt(),
	)
}
