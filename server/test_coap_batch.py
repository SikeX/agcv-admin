#!/usr/bin/env python3
"""
CoAP批量数据测试脚本
模拟多个数据点的批量发送
"""

import socket
import json
import struct
import time
import random

def build_coap_post(path, payload):
    """构建CoAP POST请求"""
    version = 1
    msg_type = 0  # Confirmable
    token_length = 0
    code = 2  # POST
    message_id = random.randint(0, 65535)
    
    header_byte = (version << 6) | (msg_type << 4) | token_length
    header = struct.pack('!BBH', header_byte, code, message_id)
    
    path_segments = path.strip('/').split('/')
    options = b''
    current_option = 0
    
    for segment in path_segments:
        segment_bytes = segment.encode('utf-8')
        delta = 11 - current_option
        length = len(segment_bytes)
        option_header = (delta << 4) | length
        options += struct.pack('!B', option_header)
        options += segment_bytes
        current_option = 11
    
    payload_bytes = payload.encode('utf-8')
    message = header + options + b'\xFF' + payload_bytes
    
    return message

def send_coap_request(data, host='127.0.0.1', port=5683):
    """发送CoAP请求"""
    payload = json.dumps(data)
    coap_message = build_coap_post("/agvc/data", payload)
    
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.settimeout(5)
    
    try:
        sock.sendto(coap_message, (host, port))
        response, _ = sock.recvfrom(4096)
        
        if len(response) >= 4:
            code = response[1]
            if b'\xFF' in response:
                payload_start = response.index(b'\xFF') + 1
                response_payload = response[payload_start:].decode('utf-8', errors='ignore')
                return code, response_payload
            return code, None
    except Exception as e:
        print(f"错误: {e}")
        return None, None
    finally:
        sock.close()

def test_single_point():
    """测试单个数据点"""
    print("=" * 60)
    print("测试1: 单个数据点")
    print("=" * 60)
    
    data = [{
        "psid": 1,
        "eqid": 101,
        "eqType": 1,
        "dataType": 1,
        "point": "temperature",
        "value": 25.5
    }]
    
    print(f"发送数据: {json.dumps(data, indent=2)}")
    code, response = send_coap_request(data)
    print(f"响应代码: {code}")
    print(f"响应内容: {response}")
    print()

def test_multiple_points():
    """测试多个数据点"""
    print("=" * 60)
    print("测试2: 多个数据点")
    print("=" * 60)
    
    data = [
        {
            "psid": 1,
            "eqid": 101,
            "eqType": 1,
            "dataType": 1,
            "point": "temperature",
            "value": 25.5
        },
        {
            "psid": 1,
            "eqid": 101,
            "eqType": 1,
            "dataType": 2,
            "point": "humidity",
            "value": 65.8
        },
        {
            "psid": 1,
            "eqid": 102,
            "eqType": 2,
            "dataType": 1,
            "point": "pressure",
            "value": 101.3
        }
    ]
    
    print(f"发送数据: {json.dumps(data, indent=2)}")
    code, response = send_coap_request(data)
    print(f"响应代码: {code}")
    print(f"响应内容: {response}")
    print()

def test_different_stations():
    """测试不同电站的数据"""
    print("=" * 60)
    print("测试3: 不同电站的数据")
    print("=" * 60)
    
    data = [
        {
            "psid": 1,
            "eqid": 1,
            "eqType": 2,
            "dataType": 2,
            "point": "2",
            "value": 32.32
        },
        {
            "psid": 2,
            "eqid": 5,
            "eqType": 3,
            "dataType": 1,
            "point": "sensor_1",
            "value": 45.67
        },
        {
            "psid": 3,
            "eqid": 10,
            "eqType": 1,
            "dataType": 3,
            "point": "voltage",
            "value": 220.5
        }
    ]
    
    print(f"发送数据: {json.dumps(data, indent=2)}")
    code, response = send_coap_request(data)
    print(f"响应代码: {code}")
    print(f"响应内容: {response}")
    print()

def test_continuous_sending():
    """测试连续发送"""
    print("=" * 60)
    print("测试4: 连续发送 (10次)")
    print("=" * 60)
    
    success_count = 0
    fail_count = 0
    
    for i in range(10):
        data = [{
            "psid": 1,
            "eqid": 1,
            "eqType": 2,
            "dataType": 2,
            "point": "continuous_test",
            "value": round(random.uniform(20.0, 30.0), 2)
        }]
        
        code, response = send_coap_request(data)
        if code == 65:  # Created
            success_count += 1
            print(f"✓ 第{i+1}次发送成功: value={data[0]['value']}")
        else:
            fail_count += 1
            print(f"✗ 第{i+1}次发送失败: code={code}")
        
        time.sleep(0.1)  # 100ms间隔
    
    print()
    print(f"统计: 成功 {success_count}, 失败 {fail_count}")
    print()

def test_invalid_data():
    """测试无效数据"""
    print("=" * 60)
    print("测试5: 无效数据")
    print("=" * 60)
    
    # 测试空数组
    print("发送空数组...")
    code, response = send_coap_request([])
    print(f"响应代码: {code}")
    print(f"响应内容: {response}")
    print()
    
    # 测试无效JSON（通过字符串）
    print("注意: 无效JSON会在解析时失败")
    print()

def main():
    """主函数"""
    print("\n" + "=" * 60)
    print("CoAP AGVC Data 批量测试")
    print("=" * 60 + "\n")
    
    try:
        test_single_point()
        time.sleep(0.5)
        
        test_multiple_points()
        time.sleep(0.5)
        
        test_different_stations()
        time.sleep(0.5)
        
        test_continuous_sending()
        time.sleep(0.5)
        
        test_invalid_data()
        
        print("=" * 60)
        print("所有测试完成!")
        print("=" * 60)
        
    except KeyboardInterrupt:
        print("\n测试被中断")
    except Exception as e:
        print(f"\n测试出错: {e}")

if __name__ == "__main__":
    main()
