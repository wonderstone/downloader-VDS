package vds

import (
	"fmt"
	"testing"

	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/data"
)

func TestNew(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}

	err = Init("", WithConf(data.Conf{
		User:     "admin",
		Password: "XBZQ_vds_2023@",
		Hiscenter: data.Hiscenter{
			Addr: "10.1.90.91:9905",
		},
		Rtcenter: data.Rtcenter{
			Level1: data.Level1{
				Addr: "10.1.90.91:9901",
			},
			Level2: data.Level2{
				Addr: "10.1.90.91:9904",
			},
		},
	}))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetCurrentUser(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%+v\n", rsp)
}

func TestGetColumn(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetColumn("balance")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%+v\n", rsp)
}

func TestGetFinanceCommon(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetFinanceCommon("balance", WithFields("*"), WithFilter("symbol=000001.SZ;end_date>=20200101;end_date<=20220101"), WithPage(1), WithPagesize(10))
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%+v\n", rsp)
}

func TestGetBalance(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetBalance("600000.SH", "20221024", "20241031", WithFields("*"))
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("%+v\n", rsp)
}

func TestGetIncome(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetIncome("000001.SZ", "20221024", "20231031", WithFields("*"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetCashflow(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetCashflow("000001.SZ", "20221024", "20231031", WithFields("*"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetPerformanceLetters(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetPerformanceLetters("000001.SZ", "20221024", "20231031", WithFields("*"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetPerformanceForecast(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetPerformanceForecast("000001.SZ", "20111024", "20231031", WithFields("*"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetFinIndicator(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetFinIndicator("000001.SZ", "20221024", "20231031", WithIndicator("*"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetCandle(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetCandle("600000.SH,601360.SH,603893.SH,601138.SH,002673.SZ,600096.SH,002100.SZ,002707.SZ,300059.SZ,002714.SZ", WithEndDateTime("20241228000000000"), WithCount(10000), WithCandleType("1min"), WithFields("*"), WithAdj("normal"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetAlpha101(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetAlpha101("600000.SH", "WQAlpha2", WithEndDateTime("20240724093600000"), WithCount(10), WithCandleType("1d"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetXbzq191alpha(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetXbzq191alpha("600000.SH", "XbzqAlpha1", WithEndDateTime("20240724093600000"), WithCount(10), WithCandleType("1d"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestCalendar(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetCalendar("SH", WithYear("2023"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestIsTrading(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetIsTrading("SH", "20230406")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetLv1Snapshot(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetLv1Snapshot("IF2312", WithStartDateTime("20231122150000000"), WithEndDateTime("20231123000000000"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetSecurity(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetSecurity("600000.SH", WithFields("*"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetLv2Snapshot(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetLv2Snapshot("600000.SH", WithStartDateTime("20240304000000000"), WithEndDateTime("20240305000000000"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetLv2Trade(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetLv2Trade("600000.SH", WithStartDateTime("20240304000000000"), WithEndDateTime("20240305000000000"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetLv2Entrust(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetLv2Entrust("600000.SH", WithStartDateTime("20240304000000000"), WithEndDateTime("20240305000000000"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestReplay(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	dataChan, err := Replay("600000.SH", "20240130000000000", "20240131000000000", "lv1snapshot", WithRate(-1))
	if err != nil {
		t.Error(err)
		return
	}
	for {
		for v := range dataChan {
			fmt.Println(v)
		}
	}
}

func TestDownloadCandleToPath(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	err = DownloadCandleToPath("./", "300300.SZ", "20240130000000000", "20240131000000000", WithCandleType("1min"))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetTimeline(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetTimeline("600000.SH", WithDay(1))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestGetBargain(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := GetBargain("600000.SH", WithStartDateTime("20240304000000000"), WithEndDateTime("20240905000000000"), WithDirector(1))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v\n", rsp)
}

func TestSubscribeLv2(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	dataChan, err := SubscribeLv2(data.Snapshot, "600000.SH")
	if err != nil {
		t.Error(err)
		return
	}
	for {
		for v := range dataChan {
			fmt.Println(v)
		}
	}
}

func TestSubscribeLv1(t *testing.T) {
	filepath := "../example.yaml"
	err := Init(filepath)
	if err != nil {
		t.Error(err)
		return
	}
	dataChan, err := SubscribeLv1(data.Bargain, "600000.SH")
	if err != nil {
		t.Error(err)
		return
	}
	for {
		for v := range dataChan {
			fmt.Println(v)
		}
	}
}
