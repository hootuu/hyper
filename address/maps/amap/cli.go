package amap

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/spf13/cast"
	"resty.dev/v3"
	"time"
)

type Client struct {
	cli *resty.Client
}

func NewClient(key string) *Client {
	cli := resty.New()
	cli.SetBaseURL("https://restapi.amap.com/v3/").
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetTimeout(60 * time.Second).
		AddRequestMiddleware(func(cli *resty.Client, req *resty.Request) error {
			req.SetQueryParam("key", key)
			return nil
		})
	return &Client{cli: cli}
}

func (p *Client) District(adCode string, page int, sub int) ([]*District, error) {
	var resp DistrictResponse
	err := retry.Do(func() error {
		_, err := p.cli.R().
			SetQueryParam("keywords", adCode).
			SetQueryParam("page", cast.ToString(page)).
			SetQueryParam("subdistrict", cast.ToString(sub)).
			SetResult(&resp).
			Get("config/district")
		if err != nil {
			return err
		}
		if resp.Status != "1" {
			return fmt.Errorf("amap err: status=%s:[%s]%s", resp.Status, resp.InfoCode, resp.Info)
		}
		return nil
	},
		retry.Attempts(5),
		retry.Delay(1500*time.Millisecond),
	)
	if err != nil {
		return nil, err
	}
	return resp.Districts, nil
}
