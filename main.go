package main

// import gota

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/vds"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"

	// "gopkg.in/yaml.v3"
	"flag"

	"time"

	"github.com/emirpasic/gods/maps/hashbidimap"
	"github.com/spf13/viper"
)

type VDSInfo struct {
	Targets       string
	StartDateTime string
	EndDateTime   string
	CandleType    string
	Fields        string
	Adj           string
	dtt           string
	if1Day		  bool
}

func main() {
	var VDSConfigPath string
	flag.StringVar(&VDSConfigPath, "VDSConfigPath", "./config.yaml", "The path of the VDS config file")
	var VQTYamlPath string
	flag.StringVar(&VQTYamlPath, "VQTYamlPath", "./vqt-task/framework.yaml", "The path of the VQT yaml file")
	var outputCSVPath string
	flag.StringVar(&outputCSVPath, "outputCSVPath", "./output", "The path of the output csv file")

	// Initialize the VDS
	filepath := VDSConfigPath
	err := vds.Init(filepath)
	if err != nil {
		// To Che : make sure no connection error at Xi'an！！
		fmt.Println(err)
	}

	vqtyamlfilepath := VQTYamlPath

	vdsInfo, m := readYaml(vqtyamlfilepath, "normal")
	// fmt.Println(vdsInfo)
	// rsp, err := vds.GetCandle(vdsInfo.Targets,
	// 	vds.WithEndDateTime(vdsInfo.EndDateTime),
	// 	vds.WithStartDateTime(vdsInfo.StartDateTime),
	// 	vds.WithCandleType(vdsInfo.CandleType),
	// 	vds.WithFields("*"),
	// 	vds.WithAdj(vdsInfo.Adj))

	// fmt.Println(vdsInfo)
	// todo : delete the following lines
	tmpf, err := os.Open("tmp.csv")
	if err != nil {
		panic(err)
	}
	defer tmpf.Close()
	// fmt.Println(vdsInfo)

	rsp := dataframe.ReadCSV(tmpf)

	// todo : delete end

	// print the dttPlus1
	
	// fmt.Println(tmpdtt, tmpdttp1)

	// iter Output the csv
	for _, target := range m.Keys() {
		// filter the dataframe by symbol first 6 digits
		tmptarget := strings.Split(target.(string), ".")[0]
		tmpdf := filterBySymbol(&rsp, tmptarget)
		// tmpdf := filterBySymbol(&rsp, target.(string))
		sel := tmpdf.Select([]string{"Open", "Close", "High", "Low", "Volume", "Amount"})
		// write the dataframe to csv file
		Date := getDate(tmpdf)
		Time := getTime(tmpdf)
		// new a dataframe with Date and Time column with column name
		// add Date and Time column to the dataframe at leftmost
		selCols := dataframe.New(*Date, *Time)
		sel = selCols.CBind(sel)
		// check if outputCSVPath is a directory and already exists
		if _, err := os.Stat(outputCSVPath); os.IsNotExist(err) {
			// path/to/whatever does not exist
			os.Mkdir(outputCSVPath, 0755)
		}

		// todo :filter by dtt and dttp1
		// filter the dataframe by dtt and dttp1
		if vdsInfo.if1Day {
			tmpdtt, tmpdttp1 := dttPlus1(vdsInfo.dtt)
			filteredDF := filterByTime(&sel, tmpdtt, tmpdttp1)
			// write the dataframe to csv file
			// writeCSV(&sel, target.(string) + ".csv")
			writeCSV(filteredDF, outputCSVPath+"/"+target.(string)+".csv")
		} else {
			writeCSV(&sel, outputCSVPath+"/"+target.(string)+".csv")
		}

	}

	fmt.Println("VDSConfigPath: ", VDSConfigPath)
	fmt.Println("VQTYamlPath: ", VQTYamlPath)
	fmt.Println("outputCSVPath: ", outputCSVPath)
}

// func to read yaml file at filePath

func readYaml(filePath string, adj string) (VDSInfo, *hashbidimap.Map) {
	// use viper to read yaml file
	v := viper.New()
	v.SetConfigFile(filePath)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// - Step 1: read the framework part
	fw := v.GetStringMap("framework")
	// - Step 2: read the instrument part and give targets field
	targets := fw["instrument"].([]interface{})
	// build a hashbidimap to store the targets
	m := hashbidimap.New() // empty
	for _, target := range targets {
		m.Put(target.(string), label2vds(target.(string)))
	}
	// combine the m values to a string with ","
	var targetsCombined string
	for _, target := range m.Values() {
		targetsCombined = targetsCombined + target.(string) + ","
	}
	targetsCombined = strings.TrimSuffix(targetsCombined, ",")

	// deep copy the string targetsCombined to avoid the pointer problem
	dtCombined := targetsCombined

	// - Step 3: read the time part and give startDateTime and endDateTime field
	startDateTime := date2vds(strconv.Itoa(fw["begin"].(int)))
	endDateTime := date2vds(strconv.Itoa(fw["end"].(int)))
	// - Step 4: read the frequency part and give candleType field
	if1Day := false
	candleType := fw["frequency"].(string)
	var candleTypeFinal string
	// for VDS support the following 1min、5min、15min、30min、60min、120min
	switch candleType {
	case "1day":
		candleTypeFinal = "1min"
		if1Day = true
	case "1min", "5min", "15min", "30min", "60min", "120min":
		candleTypeFinal = candleType
	default:
		panic(fmt.Sprintf("candleType %s is not supported", candleType))
	}
	var dtt string
	// -Step 5 : read daily-trigger-time
	if candleType == "1day" {
		dtt = fw["daily-trigger-time"].(string)
		// fmt.Println(dtt)
	}
	return VDSInfo{
		Targets:       dtCombined,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		CandleType:    candleTypeFinal,
		Fields:        "*",
		Adj:           adj,
		dtt:           dtt,
		if1Day:        if1Day,
	}, m
}

// turn VQT label style "600664.XSHG.CS" to VDS style "600000.SH"
func label2vds(label string) string {
	labelComponents := strings.Split(label, ".")

	switch labelComponents[1] {
	case "XSHG":
		return labelComponents[0] + ".SH"
	case "XSHE":
		return labelComponents[0] + ".SZ"
	default:
		panic(fmt.Sprintf("label %s is not supported", label))
	}
}

// turn VQT Date style "20100102" to VDS style "20241228000000000" with same length and fill with 0
func date2vds(date string) string {
	return date + "000000000"
}

// go run main.go -VDSConfigPath=./config.yaml -VQTYamlPath=./vqt-task/framework.yaml -outputCSVPath=./output.csv

// filter the dataframe by symbol
func filterBySymbol(df *dataframe.DataFrame, symbol string) *dataframe.DataFrame {
	// filter the dataframe by symbol
	filter := df.Filter(dataframe.F{
		Colname:    "Symbol",
		Comparator: series.Eq,
		Comparando: symbol,
	})
	return &filter
}

// func to write dataframe to csv file with name
func writeCSV(df *dataframe.DataFrame, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	df.WriteCSV(f)
}

// get Date column in string format "20190102" from
// dataframe Time column which is sth like "20240308140100000"
func getDate(df *dataframe.DataFrame) *series.Series {
	time := df.Col("Time")
	date := make([]string, time.Len())
	for i := 0; i < time.Len(); i++ {
		date[i] = time.Elem(i).String()[:8]
	}
	tmp := series.Strings(date)
	// change the column name
	tmp.Name = "Date"
	return &tmp
}

// get Time column in string format "2019.01.02T14:50:00.000" from
// dataframe Time column which is sth like "20240308140100000"
func getTime(df *dataframe.DataFrame) *series.Series {
	time := df.Col("Time")
	date := make([]string, time.Len())
	for i := 0; i < time.Len(); i++ {
		date[i] = time.Elem(i).String()[:4] + "." + time.Elem(i).String()[4:6] + "." + time.Elem(i).String()[6:8] + "T" + time.Elem(i).String()[8:10] + ":" + time.Elem(i).String()[10:12] + ":" + time.Elem(i).String()[12:14] + ".000"
	}

	tmp := series.Strings(date)
	// change the column name
	tmp.Name = "Time"
	return &tmp
}

// change dtt from string to time and plus 1 minute
// and change them to string format, output like "145000000"
func dttPlus1(dtt string) (time.Time, time.Time) {
	// parse the dtt string to time
	dttTime, err := time.Parse("15:04", dtt)
	if err != nil {
		fmt.Println(err)
	}
	// add 1 minute
	dttTimep1 := dttTime.Add(time.Minute)
	// change the time to string
	// dttTimeStr := dttTime.Format("15:04")
	// dttTimep1Str := dttTimep1.Format("15:04")
	return dttTime, dttTimep1

}

// func to filter the dataframe by time labels at dtt , dttp1, and should be the latest at the sameday
func filterByTime(df *dataframe.DataFrame, dtt time.Time, dttp1 time.Time) *dataframe.DataFrame {
	// new a dataframe to store the result
	// set the header of the resultDF the same as the df
	// fmt.Println(df.Names())

	tmpSRecords := [][]string{}
	tmpSRecords = append(tmpSRecords, df.Names())

	closett, err := time.Parse("15:04", "15:00")
	if err != nil {
		fmt.Println(err)

	}
	// iter the df to check if the time is eq to dtt or dttp1 or the latest time at the same day
	for i := 0; i < df.Nrow(); i++ {
		// check if the time is equal to dtt or dttp1 or latest time at the same day
		rowTime := df.Col("Time").Elem(i).String()
		// fmt.Println(rowTime)
		//"2024.03.08T14:01:00.000"
		// change the rowTime to time
		formattedTime, _ := time.Parse("2006.01.02T15:04:05.000", rowTime)
		// check if the formattedTime hours and minutes is equal to dtt or dttp1
		if formattedTime.Hour() == dtt.Hour() && formattedTime.Minute() == dtt.Minute() {
			tmpSRecords = append(tmpSRecords, df.Records()[i+1])
		}
		if formattedTime.Hour() == dttp1.Hour() && formattedTime.Minute() == dttp1.Minute() {
			tmpSRecords = append(tmpSRecords, df.Records()[i+1])
		}
		// check if the time is the latest time at the same day
		if formattedTime.Hour() == closett.Hour() && formattedTime.Minute() == closett.Minute() {
			tmpSRecords = append(tmpSRecords, df.Records()[i+1])
		}
	}
	// load the tmpSRecords to the resultDF with header the same as the df
	resultDF := dataframe.LoadRecords(tmpSRecords)
	// set the header of the resultDF the same as the df
	return &resultDF

}
