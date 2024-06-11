package vds

import (
	"errors"
	"os"
	"strconv"

	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/hiscenter"
	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/rtcenter"

	"github.com/go-gota/gota/dataframe"
	"gopkg.in/yaml.v3"
)

// 全局
var v vds

// vds interface
type Vds interface {
	GetBalance(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error)
	GetIncome(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error)
	GetCashflow(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error)
	GetCandle(symbol string, options ...Option) (*dataframe.DataFrame, error)
	GetAlpha101(symbol, factor string, options ...Option) (*dataframe.DataFrame, error)
	GetXbzq191alpha(symbol, factor string, options ...Option) (*dataframe.DataFrame, error)
	GetCalendar(market string, options ...Option) (*dataframe.DataFrame, error)
	GetIsTrading(market, date string) (int, error)
	GetLv1Snapshot(symbol string, options ...Option) (*dataframe.DataFrame, error)
	GetSecurity(symbol string, options ...Option) (*dataframe.DataFrame, error)
	GetLv2Snapshot(symbol string, options ...Option) (*dataframe.DataFrame, error)
	GetLv2Trade(symbol string, options ...Option) (*dataframe.DataFrame, error)
	GetLv2Entrust(symbol string, options ...Option) (*dataframe.DataFrame, error)
	Replay(symbol, startDateTime, endDateTime, replayType string, options ...Option) (<-chan *dataframe.DataFrame, error)
	DownloadCandleToPath(downloadPath, symbol, startDateTime, endDateTime string, options ...Option) error
	GetTimeline(symbol string, options ...Option) (*dataframe.DataFrame, error)
	SubscribeLv2(msgType, symbol string) (<-chan *dataframe.DataFrame, error)
	SubscribeLv1(msgType, symbol string) (<-chan *dataframe.DataFrame, error)
}

// vds struct
type vds struct {
	HisConn hiscenter.Conn
	L1Conn  rtcenter.Conn
	L2Conn  rtcenter.Conn
}

// 创建客户端
func Init(filePath string, options ...InitOption) error {
	// 默认配置
	conf := conf{}
	var err error

	//读取配置文件
	if filePath != "" {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		//序列化
		err = yaml.Unmarshal(file, &conf)
		if err != nil {
			return err
		}
	}

	// 重写配置信息
	for _, option := range options {
		option(&conf)
	}

	// 初始化hiscenter
	var hisConn hiscenter.Conn
	if conf.HisConn.Addr != "" {
		hisConn, err = hiscenter.NewConn(conf.HisConn)
		if err != nil {
			return err
		}
		// 登录hiscenter
		err = hisConn.Login(&hiscenter.LoginRequest{
			Uname: conf.User,
			Upwd:  conf.Password,
		})
		if err != nil {
			return err
		}
	}

	// 连接rtcenter
	rtConn, err := rtcenter.NewConn(conf.RtConn)
	if err != nil {
		return err
	}
	// 登录 L1
	if rtConn.L1Conn.Addr != "" {
		err = rtConn.L1Conn.Login(&rtcenter.LoginReq{
			User:     conf.User,
			Password: conf.Password,
		})
		if err != nil {
			return err
		}
	}

	// 登录 L2
	if rtConn.L2Conn.Addr != "" {
		err = rtConn.L2Conn.Login(&rtcenter.LoginReq{
			User:     conf.User,
			Password: conf.Password,
		})
		if err != nil {
			return err
		}
	}

	v = vds{
		HisConn: hisConn,
		L1Conn:  rtConn.L1Conn,
		L2Conn:  rtConn.L2Conn,
	}
	return err
}

// 是否已创建客户端
func IsInit() error {
	if v.HisConn == nil && v.L1Conn == nil && v.L2Conn == nil {
		return errors.New("客户端未初始化")
	}
	return nil
}

// 获取当前登录用户
func GetCurrentUser() (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	df, err := v.HisConn.CurrentUser(&hiscenter.CurrentUserRequest{})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取财务表字段
func GetColumn(name string) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	df, err := v.HisConn.Column(&hiscenter.ColumnRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取财务数据
func GetFinanceCommon(name string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.FinanceCommon(&hiscenter.FinanceCommonRequest{
		Name:     name,
		Fields:   opt.Fields,
		Filter:   opt.Filter,
		Page:     int32(opt.Page),
		Pagesize: int32(opt.Pagesize),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取资产负债
func GetBalance(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Balance(&hiscenter.BalanceRequest{
		Symbol:    symbol,
		StartDate: startDate,
		EndDate:   endDate,
		Fields:    opt.Fields,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取利润分配
func GetIncome(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Income(&hiscenter.IncomeRequest{
		Symbol:    symbol,
		StartDate: startDate,
		EndDate:   endDate,
		Fields:    opt.Fields,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取现⾦流
func GetCashflow(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Cashflow(&hiscenter.CashflowRequest{
		Symbol:    symbol,
		StartDate: startDate,
		EndDate:   endDate,
		Fields:    opt.Fields,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取业绩快报
func GetPerformanceLetters(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.PerformanceLetters(&hiscenter.PerformanceLettersRequest{
		Symbol:    symbol,
		StartDate: startDate,
		EndDate:   endDate,
		Fields:    opt.Fields,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取业绩预告
func GetPerformanceForecast(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.PerformanceForecast(&hiscenter.PerformanceForecastRequest{
		Symbol:    symbol,
		StartDate: startDate,
		EndDate:   endDate,
		Fields:    opt.Fields,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取公司主要财务分析指标(新会计准则)
func GetFinIndicator(symbol, startDate, endDate string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.FinIndicator(&hiscenter.FinIndicatorRequest{
		Symbol:    symbol,
		StartDate: startDate,
		EndDate:   endDate,
		Indicator: opt.Indicator,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取K线
func GetCandle(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Candle(&hiscenter.CandleRequest{
		Symbol:        symbol,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Count:         int32(opt.Count),
		Fields:        opt.Fields,
		CandleType:    opt.CandleType,
		Adj:           opt.Adj,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取Alpha101因子
func GetAlpha101(symbol, factor string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Alpha101(&hiscenter.Alpha101Request{
		Symbol:        symbol,
		Factor:        factor,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Type:          opt.FType,
		Count:         int32(opt.Count),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取Xbzq191alpha因子
func GetXbzq191alpha(symbol, factor string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Xbzq191Alpha(&hiscenter.Xbzq191AlphaRequest{
		Symbol:        symbol,
		Factor:        factor,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Type:          opt.FType,
		Count:         int32(opt.Count),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取日历
func GetCalendar(market string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Calendar(&hiscenter.CalendarRequest{
		Market: market,
		Year:   opt.Year,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取是否交易日
func GetIsTrading(market, date string) (int, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return 0, err
	}

	rspStr, err := v.HisConn.IsTrading(&hiscenter.IsTradingRequest{
		Market: market,
		Date:   date,
	})
	if err != nil {
		return 0, err
	}
	rsp, err := strconv.Atoi(rspStr)
	if err != nil {
		return 0, err
	}
	return rsp, nil
}

// 获取Lv1快照
func GetLv1Snapshot(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Lv1Snapshot(&hiscenter.Lv1SnapshotRequest{
		Symbol:        symbol,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Fields:        opt.Fields,
		Page:          int32(opt.Page),
		Pagesize:      int32(opt.Pagesize),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取码表
func GetSecurity(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Security(&hiscenter.SecurityRequest{
		Symbol: symbol,
		Fields: opt.Fields,
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取Lv2快照
func GetLv2Snapshot(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Lv2Snapshot(&hiscenter.Lv2SnapshotRequest{
		Symbol:        symbol,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Fields:        opt.Fields,
		Page:          int32(opt.Page),
		Pagesize:      int32(opt.Pagesize),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取Lv2逐笔成交
func GetLv2Trade(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Lv2Trade(&hiscenter.Lv2TradeRequest{
		Symbol:        symbol,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Fields:        opt.Fields,
		Page:          int32(opt.Page),
		Pagesize:      int32(opt.Pagesize),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取Lv2逐笔委托
func GetLv2Entrust(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Lv2Entrust(&hiscenter.Lv2EntrustRequest{
		Symbol:        symbol,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Fields:        opt.Fields,
		Page:          int32(opt.Page),
		Pagesize:      int32(opt.Pagesize),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 回放
func Replay(symbol, startDateTime, endDateTime, replayType string, options ...Option) (<-chan *dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	ch, err := v.HisConn.Replay(&hiscenter.ReplayRequest{
		Symbol:        symbol,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		Type:          replayType,
		Rate:          int32(opt.Rate),
	})
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// 下载K线
func DownloadCandleToPath(downloadPath, symbol, startDateTime, endDateTime string, options ...Option) error {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	ch, err := v.HisConn.DownloadCandle(&hiscenter.DownloadCandleRequest{
		Symbol:        symbol,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		CandleType:    opt.CandleType,
		Fields:        "*",
		Adj:           "normal",
	})
	if err != nil {
		return err
	}

	for df := range ch {
		groups := df.GroupBy("FileName").GetGroups()
		for fileName, symbolDf := range groups {
			symbolDf = symbolDf.Drop("FileName")
			path := downloadPath + "/" + fileName
			_, err := os.Stat(path)
			var writeHeader dataframe.WriteOption
			if err != nil {
				writeHeader = dataframe.WriteHeader(true)
			} else {
				writeHeader = dataframe.WriteHeader(false)
			}
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			err = symbolDf.WriteCSV(file, writeHeader)
			if err != nil {
				return err
			}
			file.Close()
		}
	}
	return nil
}

// 获取分时
func GetTimeline(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Timeline(&hiscenter.TimelineRequest{
		Symbol: symbol,
		Day:    int32(opt.Day),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 获取分笔明细
func GetBargain(symbol string, options ...Option) (*dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}

	// 处理非必填参数
	opt := &option{}
	for _, do := range options {
		do(opt)
	}

	df, err := v.HisConn.Bargain(&hiscenter.BargainRequest{
		Symbol:        symbol,
		StartDateTime: opt.StartDateTime,
		EndDateTime:   opt.EndDateTime,
		Director:      int32(opt.Director),
		Fields:        opt.Fields,
		Page:          int32(opt.Page),
		Pagesize:      int32(opt.Pagesize),
	})
	if err != nil {
		return nil, err
	}
	return df, nil
}

// 发送订阅L2请求
func SubscribeLv2(msgType, symbol string) (<-chan *dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}
	// 订阅信息
	ch, err := v.L2Conn.Subscribe(&rtcenter.SubReq{
		MsgType: msgType,
		Symbol:  symbol,
	})
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// 发送订阅Lv1请求
func SubscribeLv1(msgType, symbol string) (<-chan *dataframe.DataFrame, error) {
	// 判断是否初始化客户端
	err := IsInit()
	if err != nil {
		return nil, err
	}
	// 订阅信息
	ch, err := v.L1Conn.Subscribe(&rtcenter.SubReq{
		MsgType: msgType,
		Symbol:  symbol,
	})
	if err != nil {
		return nil, err
	}
	return ch, nil
}
