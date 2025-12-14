package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/developerasun/SignalDash/server/dto"
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

func CreateExchangeRateDiff() (won float64, tether float64, err error) {
	responses, err := DoHttpGet([]string{
		"https://api.bithumb.com/v1/ticker?markets=KRW-USDT",
		"https://m.search.naver.com/p/csearch/content/qapirender.nhn?key=calculator&pkid=141&q=%ED%99%98%EC%9C%A8&where=m&u1=keb&u6=standardUnit&u7=0&u3=USD&u4=KRW&u8=down&u2=1",
	})
	if err != nil {
		return 0, 0, sderror.ErrInternalServer
	}

	var bithumbResp dto.ApiResponse[[]dto.BithumbApiItem]
	var naverResp dto.ApiResponse[dto.NaverApiItem]

	for _, v := range responses {
		if strings.Contains(string(v), "KRW-USDT") {
			err := json.Unmarshal(v, &bithumbResp.Data)
			if err != nil {
				return 0, 0, sderror.ErrInternalServer
			}
		}
		if strings.Contains(string(v), "pkid") {
			err := json.Unmarshal(v, &naverResp.Data)
			if err != nil {
				return 0, 0, sderror.ErrInternalServer
			}
		}
	}

	krwPurified := strings.ReplaceAll(naverResp.Data.Country[1].Value, ",", "")
	won, pErr := strconv.ParseFloat(krwPurified, 64)
	if pErr != nil {
		return 0, 0, sderror.ErrInternalServer
	}
	tether = bithumbResp.Data[0].OpeningPrice

	return won, tether, nil
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

func DoHttpGet(endpoints []string) (responses [][]byte, err error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	callback := func(endpoint string) (data []byte, err error) {
		defer wg.Done()

		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, sderror.ErrInternalServer
		}

		client := &http.Client{}
		res, dErr := client.Do(req)
		if dErr != nil {
			return nil, sderror.ErrInternalServer
		}

		data, rErr := io.ReadAll(res.Body)
		cErr := res.Body.Close()
		if cErr != nil {
			return nil, sderror.ErrInternalServer
		}

		if rErr != nil {
			return nil, sderror.ErrInternalServer
		}

		mu.Lock()
		responses = append(responses, data)
		mu.Unlock()
		return data, nil
	}

	if len(endpoints) >= 1 {
		wg.Add(len(endpoints))
		for _, v := range endpoints {
			go callback(v)
		}
		wg.Wait()
	} else {
		callback(endpoints[len(endpoints)-1])
	}

	return responses, nil
}
