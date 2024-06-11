package hiscenter

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const TIMEOUT = 60 * time.Second
const MAX_RECV_MSG_SIZE = 1024 * 1024 * 20

type Conn interface {
	Login(in *LoginRequest) error
	CurrentUser(in *CurrentUserRequest) (*dataframe.DataFrame, error)
	Column(in *ColumnRequest) (*dataframe.DataFrame, error)
	FinanceCommon(in *FinanceCommonRequest) (*dataframe.DataFrame, error)
	Balance(in *BalanceRequest) (*dataframe.DataFrame, error)
	Income(in *IncomeRequest) (*dataframe.DataFrame, error)
	Cashflow(in *CashflowRequest) (*dataframe.DataFrame, error)
	PerformanceLetters(in *PerformanceLettersRequest) (*dataframe.DataFrame, error)
	PerformanceForecast(in *PerformanceForecastRequest) (*dataframe.DataFrame, error)
	FinIndicator(in *FinIndicatorRequest) (*dataframe.DataFrame, error)
	Candle(in *CandleRequest) (*dataframe.DataFrame, error)
	Alpha101(in *Alpha101Request) (*dataframe.DataFrame, error)
	Xbzq191Alpha(in *Xbzq191AlphaRequest) (*dataframe.DataFrame, error)
	Calendar(in *CalendarRequest) (*dataframe.DataFrame, error)
	IsTrading(in *IsTradingRequest) (string, error)
	Lv1Snapshot(in *Lv1SnapshotRequest) (*dataframe.DataFrame, error)
	Security(in *SecurityRequest) (*dataframe.DataFrame, error)
	Lv2Snapshot(in *Lv2SnapshotRequest) (*dataframe.DataFrame, error)
	Lv2Trade(in *Lv2TradeRequest) (*dataframe.DataFrame, error)
	Lv2Entrust(in *Lv2EntrustRequest) (*dataframe.DataFrame, error)
	Replay(in *ReplayRequest) (<-chan *dataframe.DataFrame, error)
	DownloadCandle(in *DownloadCandleRequest) (<-chan *dataframe.DataFrame, error)
	Timeline(in *TimelineRequest) (*dataframe.DataFrame, error)
	Bargain(in *BargainRequest) (*dataframe.DataFrame, error)
}

// hiscenter conn struct
type HisConn struct {
	Conn     *grpc.ClientConn
	Addr     string `yaml:"addr"`
	Ssl      bool   `yaml:"ssl"`
	Token    string
	Response *dataframe.DataFrame
}

// hiscenter response body data param struct
type Data struct {
	Fields []string        `json:"fields"`
	Items  [][]interface{} `json:"items"`
}

// hiscenter response body data param struct
type DetailData struct {
	Fields []string      `json:"fields"`
	Items  []interface{} `json:"items"`
}

// 创建新的连接
func NewConn(c HisConn) (Conn, error) {
	return &c, nil
}

// 登录
func (c *HisConn) Login(in *LoginRequest) error {
	var addr string = c.Addr
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	res, err := client.Login(ctx, in)
	if err != nil {
		return err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return errors.New(res.GetText())
	}
	c.Token = res.GetData()
	return nil
}

// 当前登录用户
func (c *HisConn) CurrentUser(in *CurrentUserRequest) (*dataframe.DataFrame, error) {
	var addr string = c.Addr
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.CurrentUser(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDetailDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 财务表字段
func (c *HisConn) Column(in *ColumnRequest) (*dataframe.DataFrame, error) {
	var addr string = c.Addr
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Column(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 财务数据
func (c *HisConn) FinanceCommon(in *FinanceCommonRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.FinanceCommon(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 资产负债
func (c *HisConn) Balance(in *BalanceRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Balance(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 利润分配
func (c *HisConn) Income(in *IncomeRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Income(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 现金流
func (c *HisConn) Cashflow(in *CashflowRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Cashflow(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 业绩快报
func (c *HisConn) PerformanceLetters(in *PerformanceLettersRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.PerformanceLetters(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 业绩预告
func (c *HisConn) PerformanceForecast(in *PerformanceForecastRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.PerformanceForecast(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 业绩预告
func (c *HisConn) FinIndicator(in *FinIndicatorRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.FinIndicator(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// K线
func (c *HisConn) Candle(in *CandleRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Candle(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// Alpha101因子
func (c *HisConn) Alpha101(in *Alpha101Request) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Alpha101(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// Xbzq191alpha因子
func (c *HisConn) Xbzq191Alpha(in *Xbzq191AlphaRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Xbzq191Alpha(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 日历
func (c *HisConn) Calendar(in *CalendarRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Calendar(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 是否交易日
func (c *HisConn) IsTrading(in *IsTradingRequest) (string, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.IsTrading(ctx, in)
	if err != nil {
		return "", err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return "", errors.New(res.GetText())
	}
	return res.GetData(), nil
}

// Lv1快照
func (c *HisConn) Lv1Snapshot(in *Lv1SnapshotRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Lv1Snapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 码表
func (c *HisConn) Security(in *SecurityRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Security(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// Lv2快照
func (c *HisConn) Lv2Snapshot(in *Lv2SnapshotRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Lv2Snapshot(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// Lv2逐笔成交
func (c *HisConn) Lv2Trade(in *Lv2TradeRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Lv2Trade(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// Lv2逐笔委托
func (c *HisConn) Lv2Entrust(in *Lv2EntrustRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Lv2Entrust(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 回放
func (c *HisConn) Replay(in *ReplayRequest) (<-chan *dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	client := NewHiscenterClient(conn)
	ctx := context.TODO()
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Replay(ctx, in)
	if err != nil {
		return nil, err
	}
	ch := make(chan *dataframe.DataFrame)
	go RspData(res, ch)
	return ch, nil

}

// 下载K线
func (c *HisConn) DownloadCandle(in *DownloadCandleRequest) (<-chan *dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	client := NewHiscenterClient(conn)
	ctx := context.TODO()
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.DownloadCandle(ctx, in)
	if err != nil {
		return nil, err
	}
	ch := make(chan *dataframe.DataFrame)
	go RspData(res, ch)
	return ch, nil
}

// 分时
func (c *HisConn) Timeline(in *TimelineRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Timeline(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// 分笔明细
func (c *HisConn) Bargain(in *BargainRequest) (*dataframe.DataFrame, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MAX_RECV_MSG_SIZE)))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewHiscenterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", c.Token)
	res, err := client.Bargain(ctx, in)
	if err != nil {
		return nil, err
	}
	// 判断返回结果
	if res.GetCode() != 0 {
		return nil, errors.New(res.GetText())
	}
	c.Response, err = parseDataFrame(res.GetData())
	if err != nil {
		return nil, err
	}
	return c.Response, nil
}

// json转dataframe
func parseDataFrame(dataStr string) (*dataframe.DataFrame, error) {
	// 解析data参数
	var data Data
	// json解析
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, err
	}
	if len(data.Items) == 0 {
		return nil, nil
	}
	// 组装records数据
	var records [][]string
	var colTypes = make(map[string]series.Type)
	records = append(records, data.Fields)

	for key, item := range data.Items {
		var record []string
		for k, v := range item {
			if key == 0 {
				if _, ok := v.(string); ok {
					colTypes[data.Fields[k]] = series.String
				}
			}
			record = append(record, fmt.Sprint(v))
		}
		records = append(records, record)
	}
	df := dataframe.LoadRecords(
		records,
		dataframe.WithTypes(colTypes),
	)
	return &df, nil
}

// json转dataframe
func parseDetailDataFrame(dataStr string) (*dataframe.DataFrame, error) {
	// 解析data参数
	var data DetailData
	// json解析
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, err
	}
	if len(data.Items) == 0 {
		return nil, nil
	}
	// 组装records数据
	var records [][]string
	records = append(records, data.Fields)

	var record []string
	for _, item := range data.Items {
		record = append(record, fmt.Sprint(item))
	}
	records = append(records, record)
	// 转变为dataframe形式
	df := dataframe.LoadRecords(
		records,
	)
	return &df, nil
}

func RspData(res Hiscenter_ReplayClient, ch chan<- *dataframe.DataFrame) error {
	defer func() {
		close(ch)
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	for {
		rec, err := res.Recv()
		if err != nil {
			if io.EOF == err {
				return nil
			}
			panic(err)
		}
		if rec.GetCode() != 0 {
			panic(errors.New(rec.GetText()))
		}
		df, err := parseDataFrame(rec.GetData())
		if err != nil {
			panic(err)
		}
		ch <- df
	}
}
