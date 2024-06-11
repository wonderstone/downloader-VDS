package vds

import (
	"errors"

	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/data"
	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/hiscenter"
	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/rtcenter"
)

// init option type
type InitOption func(*conf)

// conf struct
type conf struct {
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	hiscenter.HisConn `yaml:"hiscenter"`
	rtcenter.RtConn   `yaml:"rtcenter"`
}

// option type
type Option func(*option)

// option struct
type option struct {
	StartDateTime string
	EndDateTime   string
	CandleType    string
	Fields        string
	Adj           string
	FType         string
	Year          string
	Count         int
	Page          int
	Pagesize      int
	Rate          int
	Day           int
	Director      int
	Indicator     string
	Filter        string
}

// 传输用户配置信息
func WithConf(config data.Conf) InitOption {
	return func(c *conf) {
		if config.User != "" {
			c.User = config.User
		}
		if config.Password != "" {
			c.Password = config.Password
		}
		if c.HisConn.Addr == "" {
			panic(errors.New("历史接口地址不能为空"))
		}
		if config.Rtcenter.Level1.Addr != "" {
			c.RtConn.L1Conn = &rtcenter.SubConn{
				Addr: config.Rtcenter.Level1.Addr,
			}
		}
		if config.Rtcenter.Level2.Addr != "" {
			c.RtConn.L2Conn = &rtcenter.SubConn{
				Addr: config.Rtcenter.Level2.Addr,
			}
		}
	}
}

// 可选参数
func WithStartDateTime(startDateTime string) Option {
	return func(o *option) {
		o.StartDateTime = startDateTime
	}
}

func WithEndDateTime(endDateTime string) Option {
	return func(o *option) {
		o.EndDateTime = endDateTime
	}
}

func WithCandleType(candleType string) Option {
	return func(o *option) {
		o.CandleType = candleType
	}
}

func WithFields(fields string) Option {
	return func(o *option) {
		o.Fields = fields
	}
}

func WithAdj(adj string) Option {
	return func(o *option) {
		o.Adj = adj
	}
}

func WithFType(fType string) Option {
	return func(o *option) {
		o.FType = fType
	}
}

func WithYear(year string) Option {
	return func(o *option) {
		o.Year = year
	}
}

func WithCount(count int) Option {
	return func(o *option) {
		o.Count = count
	}
}

func WithPage(page int) Option {
	return func(o *option) {
		o.Page = page
	}
}

func WithPagesize(pagesize int) Option {
	return func(o *option) {
		o.Pagesize = pagesize
	}
}

func WithRate(rate int) Option {
	return func(o *option) {
		o.Rate = rate
	}
}

func WithDay(day int) Option {
	return func(o *option) {
		o.Day = day
	}
}

func WithDirector(director int) Option {
	return func(o *option) {
		o.Director = director
	}
}

func WithIndicator(indicator string) Option {
	return func(o *option) {
		o.Indicator = indicator
	}
}

func WithFilter(filter string) Option {
	return func(o *option) {
		o.Filter = filter
	}
}
