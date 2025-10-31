package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// AgvcDataItem AGVC数据项
type AgvcDataItem struct {
	Psid     int     `json:"psid"`
	Eqid     int     `json:"eqid"`
	EqType   int     `json:"eqType"`
	DataType int     `json:"dataType"`
	Point    string  `json:"point"`
	Value    float64 `json:"value"`
}

const (
	coapVersion             = 1
	coapTypeConfirmable     = 0
	coapCodePost            = 2
	coapCodeCreated         = 65
	coapCodeContent         = 69
	coapCodeBadRequest      = 128
	coapCodeNotFound        = 132
	coapCodeMethodNotAllowed = 133
	coapCodeInternalServerError = 160
	coapOptionURIPath       = 11
)

func buildCoapPOST(path string, payload []byte) []byte {
	var buf bytes.Buffer

	// Header
	messageID := uint16(0x1234)
	tkl := 0
	headerByte := byte(coapVersion<<6) | byte(coapTypeConfirmable<<4) | byte(tkl)

	buf.WriteByte(headerByte)
	buf.WriteByte(coapCodePost)
	binary.Write(&buf, binary.BigEndian, messageID)

	// URI-Path options
	segments := bytes.Split([]byte(path[1:]), []byte("/")) // Remove leading /
	currentOpt := uint16(0)

	for _, segment := range segments {
		delta := coapOptionURIPath - currentOpt
		length := len(segment)

		optionHeader := byte(delta<<4) | byte(length)
		buf.WriteByte(optionHeader)
		buf.Write(segment)

		currentOpt = coapOptionURIPath
	}

	// Payload marker and payload
	if len(payload) > 0 {
		buf.WriteByte(0xFF)
		buf.Write(payload)
	}

	return buf.Bytes()
}

func parseCoapResponse(data []byte) (code byte, payload []byte, err error) {
	if len(data) < 4 {
		return 0, nil, fmt.Errorf("response too short")
	}

	code = data[1]

	// Find payload marker
	for i := 4; i < len(data); i++ {
		if data[i] == 0xFF {
			payload = data[i+1:]
			break
		}
	}

	return code, payload, nil
}

func sendCoapRequest(host string, port int, path string, data interface{}) error {
	// Serialize data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	fmt.Printf("发送数据: %s\n", string(jsonData))

	// Build CoAP message
	coapMsg := buildCoapPOST(path, jsonData)

	// Create UDP connection
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Set timeout
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Send request
	_, err = conn.Write(coapMsg)
	if err != nil {
		return fmt.Errorf("failed to send: %v", err)
	}

	// Receive response
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return fmt.Errorf("failed to receive: %v", err)
	}

	// Parse response
	code, payload, err := parseCoapResponse(buf[:n])
	if err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	fmt.Printf("响应代码: %d\n", code)
	if len(payload) > 0 {
		fmt.Printf("响应内容: %s\n", string(payload))
	}

	// Check response code
	switch code {
	case coapCodeCreated:
		fmt.Println("✓ 数据发送成功!")
	case coapCodeBadRequest:
		fmt.Println("✗ 请求错误")
	case coapCodeNotFound:
		fmt.Println("✗ 路径未找到")
	case coapCodeMethodNotAllowed:
		fmt.Println("✗ 方法不允许")
	case coapCodeInternalServerError:
		fmt.Println("✗ 服务器内部错误")
	default:
		fmt.Printf("? 未知响应代码: %d\n", code)
	}

	return nil
}

func main() {
	host := "127.0.0.1"
	port := 5683
	path := "/agvc/data"

	fmt.Println("CoAP客户端测试")
	fmt.Println("================")
	fmt.Println()

	// Test 1: Single data point
	fmt.Println("测试1: 单个数据点")
	fmt.Println("------------------")
	singleData := []AgvcDataItem{
		{
			Psid:     1,
			Eqid:     1,
			EqType:   2,
			DataType: 2,
			Point:    "2",
			Value:    32.32,
		},
	}
	if err := sendCoapRequest(host, port, path, singleData); err != nil {
		fmt.Printf("错误: %v\n", err)
	}
	fmt.Println()

	time.Sleep(500 * time.Millisecond)

	// Test 2: Multiple data points
	fmt.Println("测试2: 多个数据点")
	fmt.Println("------------------")
	multipleData := []AgvcDataItem{
		{Psid: 1, Eqid: 101, EqType: 1, DataType: 1, Point: "temperature", Value: 25.5},
		{Psid: 1, Eqid: 101, EqType: 1, DataType: 2, Point: "humidity", Value: 65.8},
		{Psid: 1, Eqid: 102, EqType: 2, DataType: 1, Point: "pressure", Value: 101.3},
	}
	if err := sendCoapRequest(host, port, path, multipleData); err != nil {
		fmt.Printf("错误: %v\n", err)
	}
	fmt.Println()

	time.Sleep(500 * time.Millisecond)

	// Test 3: Different stations
	fmt.Println("测试3: 不同电站")
	fmt.Println("------------------")
	stationData := []AgvcDataItem{
		{Psid: 1, Eqid: 1, EqType: 2, DataType: 2, Point: "2", Value: 32.32},
		{Psid: 2, Eqid: 5, EqType: 3, DataType: 1, Point: "sensor_1", Value: 45.67},
		{Psid: 3, Eqid: 10, EqType: 1, DataType: 3, Point: "voltage", Value: 220.5},
	}
	if err := sendCoapRequest(host, port, path, stationData); err != nil {
		fmt.Printf("错误: %v\n", err)
	}
	fmt.Println()

	fmt.Println("所有测试完成!")
}
