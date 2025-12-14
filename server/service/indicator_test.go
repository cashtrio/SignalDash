package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/developerasun/SignalDash/server/config"
	"github.com/developerasun/SignalDash/server/dto"
	"github.com/developerasun/SignalDash/server/models"
	"github.com/developerasun/SignalDash/server/sderror"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func cleanup(t *testing.T) *gorm.DB {
	absPath, err := filepath.Abs("../config")
	require.NoError(t, err)

	env := config.NewEnvironment(absPath, "options").Instance
	testDB := env.GetString("server.database.test")

	dbPath, err := filepath.Abs("../" + testDB)
	require.NoError(t, err)

	_ = os.Remove(dbPath)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Indicator{})
	require.NoError(t, err)

	return db
}

func TestFindAndCreateDollarIndex(t *testing.T) {
	t.Skip()
	db := cleanup(t)
	_, fErr := FindLatestDollarIndex(db)

	// @dev should be clean
	require.ErrorIs(t, fErr, sderror.ErrEmptyStorage)

	dxy := "89.44"
	cErr := CreateDollarIndex(db, dxy)
	require.NoError(t, cErr)

	_, rfErr := FindLatestDollarIndex(db)
	require.NotErrorIs(t, rfErr, sderror.ErrEmptyStorage)
}

func TestNewHttp(t *testing.T) {
	responses, err := DoHttpGet([]string{
		"https://api.bithumb.com/v1/ticker?markets=KRW-USDT",
		"https://m.search.naver.com/p/csearch/content/qapirender.nhn?key=calculator&pkid=141&q=%ED%99%98%EC%9C%A8&where=m&u1=keb&u6=standardUnit&u7=0&u3=USD&u4=KRW&u8=down&u2=1",
	})
	require.NoError(t, err)

	var bithumbResp dto.ApiResponse[[]dto.BithumbApiItem]
	var naverResp dto.ApiResponse[dto.NaverApiItem]

	for _, v := range responses {
		if strings.Contains(string(v), "KRW-USDT") {
			err := json.Unmarshal(v, &bithumbResp.Data)
			if err != nil {
				require.Fail(t, "failed to unmarshal bithumb response")
			}
		}
		if strings.Contains(string(v), "pkid") {
			err := json.Unmarshal(v, &naverResp.Data)
			if err != nil {
				require.Fail(t, "failed to unmarshal naver response")
			}
		}
	}

	require.True(t, len(bithumbResp.Data) == 1)
	require.True(t, naverResp.Data.PKID == 141)
	t.Log("value: ", bithumbResp.Data[0].OpeningPrice)
}
