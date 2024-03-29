package fyers

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dragonzurfer/fyersgo/api"
	"github.com/dragonzurfer/fyersgo/utils"
)

func (c *client) GetQuote(symbols []string) ([]api.DataQuote, error) {
	if resp, err := c.invoke(utils.GET, c.toUriData(Data, QuoteUrl, QueryParam, SymbolsQueryParam, strings.Join(symbols, ",")), nil); err != nil {
		return nil, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var quoteResp []api.DataQuote
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "d.#.v")), &quoteResp); err != nil {
				return nil, err
			} else {
				return quoteResp, nil
			}
		} else {
			return nil, fmt.Errorf("failed to get quote for symbols %v. %v", symbols, utils.GetJsonValueAtPath(resp, "errmsg"))
		}
	}
}

type DomainHistoricalData struct {
	Status  string      `json:"s,omitempty" yaml:"s,omitempty"`
	Candles [][]float64 `json:"candles,omitempty" yaml:"candles,omitempty"`
}

// GetHistoricalData gives us candles of past data based on given time frame.
func (c *client) GetHistoricalData(symbol string, resolution api.Resolution, startDate, endDate time.Time) (api.HistoricalData, error) {
	queryParam := fmt.Sprintf("symbol=%s&date_format=%d&resolution=%s&cont_flag=1", symbol, 0, resolution)
	if !startDate.IsZero() {
		queryParam = fmt.Sprintf("%s&range_from=%d", queryParam, startDate.Unix())
	} else {
		queryParam = fmt.Sprintf("%s&range_from=2021-01-01", queryParam)
	}
	if !endDate.IsZero() {
		queryParam = fmt.Sprintf("%s&range_to=%d", queryParam, endDate.Unix())
	} else {
		queryParam = fmt.Sprintf("%s&range_to=2021-01-02", queryParam)
	}
	if resp, err := c.invoke(utils.GET, c.toUriData(DataV3, HistoricalDataApiUrl, QueryParam, queryParam), nil); err != nil {

		return api.HistoricalData{}, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var r DomainHistoricalData
			if json.Unmarshal(resp, &r); err != nil {
				return api.HistoricalData{}, err
			} else {
				response := api.HistoricalData{Symbol: symbol, Candles: make([]api.Candle, 0, 1)}
				for _, c := range r.Candles {
					response.Candles = append(response.Candles,
						api.Candle{
							Timestamp:    utils.ToIstTimeFromEpoch(int64(c[0])),
							OpenValue:    float32(c[1]),
							HighestValue: float32(c[2]),
							LowestValue:  float32(c[3]),
							CloseValue:   float32(c[4]),
							Volume:       int64(c[5]),
						},
					)
				}
				return response, nil
			}
		} else {
			return api.HistoricalData{}, fmt.Errorf("failed to get hostorical data for symbol %v. %v", symbol, utils.GetJsonValueAtPath(resp, "message"))
		}
	}
}

func (c *client) GetMarketDepth(symbol string) (map[string]api.MarketDepth, error) {
	if resp, err := c.invoke(utils.GET, c.toUriData(Data, MarketDepthApiUrl, QueryParam, SymbolQueryParam, symbol, AndQueryParam, OhlcvQueryParam), nil); err != nil {
		return nil, err
	} else {
		if utils.IsSuccessResponse(resp) {
			var response map[string]api.MarketDepth
			if json.Unmarshal([]byte(utils.GetJsonValueAtPath(resp, "d")), &response); err != nil {
				return nil, err
			} else {
				return response, nil
			}
		} else {
			return nil, fmt.Errorf("failed to get market depth for symbol %v. %v", symbol, utils.GetJsonValueAtPath(resp, "errmsg"))
		}
	}
}
