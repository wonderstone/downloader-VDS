package rtcenter

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"gitlab.westresearch.west95582.com/department2/quantteam/vdsgo/data"

	"github.com/go-gota/gota/dataframe"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// rtcenter conn interface
type Conn interface {
	Login(login *LoginReq) error
	Subscribe(sub *SubReq) (<-chan *dataframe.DataFrame, error)
}

// rtcenter conn struct
type RtConn struct {
	L1Conn *SubConn `yaml:"level1"`
	L2Conn *SubConn `yaml:"level2"`
}

// rtcenter conn struct
type SubConn struct {
	Addr     string `yaml:"addr"`
	Conn     net.Conn
	Response <-chan *dataframe.DataFrame
}

// rtcenter login request struct
type LoginReq struct {
	MsgType  string `json:"MsgType"`
	User     string `json:"Name"`
	Password string `json:"PassWd"`
}

// rtcenter subscribe request struct
type SubReq struct {
	MsgType string `json:"MsgType"`
	Symbol  string `json:"Symbol"`
	//Fields  string `json:"Fields"`
}

// 创建新的连接
func NewConn(rtConn RtConn) (*RtConn, error) {
	c := RtConn{
		L1Conn: &SubConn{},
		L2Conn: &SubConn{},
	}
	// L1
	// 连接tcp
	if rtConn.L1Conn != nil && rtConn.L1Conn.Addr != "" {
		connL1, err := net.Dial("tcp", rtConn.L1Conn.Addr)
		if err != nil {
			return &c, err
		}
		c.L1Conn.Conn = connL1
		// 获取连接返回结果
		chL1 := make(chan *dataframe.DataFrame)
		c.L1Conn.Response = chL1
		go RspData(c.L1Conn.Conn, chL1)
	}

	// L2
	// 连接tcp
	if rtConn.L2Conn != nil && rtConn.L2Conn.Addr != "" {
		connL2, err := net.Dial("tcp", rtConn.L2Conn.Addr)
		if err != nil {
			return &c, err
		}
		c.L2Conn.Conn = connL2

		// 获取连接返回结果
		chL2 := make(chan *dataframe.DataFrame)
		c.L2Conn.Response = chL2
		go RspData(c.L2Conn.Conn, chL2)
	}

	return &c, nil
}

// 登录
func (c *SubConn) Login(login *LoginReq) error {

	// 组装json
	login.MsgType = data.Login
	loginStr, err := json.Marshal(login)
	if err != nil {
		return err
	}

	// 登录
	c.send(string(loginStr))

	return nil
}

// 订阅
func (c *SubConn) Subscribe(sub *SubReq) (<-chan *dataframe.DataFrame, error) {
	subStr, err := json.Marshal(sub)
	if err != nil {
		return nil, err
	}
	c.send(string(subStr))
	return c.Response, nil
}

// 发送请求
func (c *SubConn) send(str string) {
	var reqmap map[string]string
	json.Unmarshal([]byte(str), &reqmap)
	sub := &VDSReq{
		ReqMap: reqmap,
	}
	sData, err := proto.Marshal(sub)
	if err != nil {
		panic(err)
	}
	slen := len(sData)
	var b bytes.Buffer
	b.Write(IntToBytes(slen))
	b.Write([]byte(sData))
	c.Conn.Write([]byte(b.Bytes()))
}

// 获取返回信息
func RspData(newConn net.Conn, ch chan<- *dataframe.DataFrame) error {
	defer close(ch)
	for {
		//读消息头
		datalen := make([]byte, 4)
		_, err := io.ReadFull(newConn, datalen)
		if err != nil {
			panic(err)
		}
		// turn datalen into int
		dtlen := BytesToInt(datalen)
		buf := make([]byte, dtlen)
		_, err = io.ReadFull(newConn, buf)
		if err != nil {
			panic(err)
		}
		var s = VDSRsp{}
		err = proto.Unmarshal(buf, &s)
		if err != nil {
			panic(err)
		}
		// 判断数据返回
		code, ok := s.RspMap["Code"]
		if ok {
			var codeMsg Int32Msg
			code.GetValue().UnmarshalTo(&codeMsg)
			// 判断返回结果
			if codeMsg.GetData() != 0 {
				var errMsg StringMsg
				s.RspMap["Msg"].GetValue().UnmarshalTo(&errMsg)
				panic(errors.New(errMsg.GetData()))
			}
		}
		// 重新组装数据
		var records [][]string
		// 组装fields
		var fieldsStr StringMsg
		s.RspMap["Fields"].GetValue().UnmarshalTo(&fieldsStr)
		data := fieldsStr.GetData()
		if data == "" {
			continue
		}
		fields := strings.Split(data, ",")
		records = append(records, fields)
		var record []string
		//组装record
		for _, field := range fields {
			v := s.RspMap[field]
			// 获取类型
			msgType, err := protoregistry.GlobalTypes.FindMessageByURL(v.GetValue().GetTypeUrl())
			if err != nil {
				panic(err)
			}
			// 初始化类型
			msg := msgType.New().Interface()
			// 解析数据到固定类型
			v.GetValue().UnmarshalTo(msg)
			// 赋值到record
			record = append(record, fmt.Sprint(msg.ProtoReflect().Get(msg.ProtoReflect().Descriptor().Fields().ByName("data"))))
		}
		records = append(records, record)
		// 转换为dataframe格式
		df := dataframe.LoadRecords(records)
		ch <- &df
	}
}

//  turn an int32 into a byte array with 4 bytes
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//  turn a byte array with 4 bytes into an int32
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}
