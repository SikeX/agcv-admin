#!/usr/bin/env python3
"""
CoAP客户端测试脚本
用于测试 /agvc/data 接口
"""

import socket
import json
import struct

def build_coap_post(path, payload):
    """构建CoAP POST请求"""
    # CoAP Header
    version = 1
    msg_type = 0  # Confirmable
    token_length = 0
    code = 2  # POST
    message_id = 0x1234
    
    # Build header
    header_byte = (version << 6) | (msg_type << 4) | token_length
    header = struct.pack('!BBH', header_byte, code, message_id)
    
    # Build URI-Path options
    path_segments = path.strip('/').split('/')
    options = b''
    current_option = 0
    
    for segment in path_segments:
        segment_bytes = segment.encode('utf-8')
        delta = 11 - current_option  # URI-Path = 11
        length = len(segment_bytes)
        
        # Option header
        option_header = (delta << 4) | length
        options += struct.pack('!B', option_header)
        options += segment_bytes
        
        current_option = 11
    
    # Payload marker and payload
    payload_bytes = payload.encode('utf-8')
    message = header + options + b'\xFF' + payload_bytes
    
    return message

def test_agvc_data():
    """测试AGVC数据接口"""
    # 准备测试数据
    test_data = [
        {
            "psid": 1,
            "eqid": 1,
            "eqType": 2,
            "dataType": 2,
            "point": "2",
            "value": 32.32
        }
    ]
    
    payload = json.dumps(test_data)
    print(f"测试数据: {payload}")
    
    # 构建CoAP请求
    coap_message = build_coap_post("/agvc/data", payload)
    
    # 发送请求
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.settimeout(5)
    
    try:
        # 发送到CoAP服务器
        server_address = ('127.0.0.1', 5683)
        print(f"发送到 {server_address}...")
        sock.sendto(coap_message, server_address)
        
        # 接收响应
        data, server = sock.recvfrom(4096)
        print(f"收到响应: {len(data)} 字节")
        
        # 解析响应
        if len(data) >= 4:
            header_byte = data[0]
            code = data[1]
            message_id = struct.unpack('!H', data[2:4])[0]
            
            print(f"响应代码: {code}")
            print(f"消息ID: {message_id}")
            
            # 查找payload
            if b'\xFF' in data:
                payload_start = data.index(b'\xFF') + 1
                response_payload = data[payload_start:].decode('utf-8', errors='ignore')
                print(f"响应内容: {response_payload}")
            else:
                print("无响应内容")
                
    except socket.timeout:
        print("请求超时")
    except Exception as e:
        print(f"错误: {e}")
    finally:
        sock.close()

if __name__ == "__main__":
    test_agvc_data()
