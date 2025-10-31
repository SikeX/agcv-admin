package initialize

import (
    "errors"
    "fmt"
    "net"
    "strings"

    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "go.uber.org/zap"
)

const (
    coapVersion             = 1
    coapTypeConfirmable     = 0
    coapTypeNonConfirmable  = 1
    coapTypeAcknowledgement = 2

    coapCodeGet              = 1
    coapCodeContent          = 69
    coapCodeNotFound         = 132
    coapCodeMethodNotAllowed = 133

    coapOptionURIPath = 11
)

type coapMessage struct {
    typ       int
    code      byte
    messageID uint16
    token     []byte
    options   map[uint16][][]byte
    payload   []byte
}

func CoapServer() {
    cfg := global.GVA_CONFIG.Coap

    global.GVA_COAP_SERVER = nil

    if !cfg.Enable {
        if global.GVA_LOG != nil {
            global.GVA_LOG.Info("CoAP server disabled, skip initialization")
        }
        return
    }

    host := cfg.Host
    if host == "" {
        host = "0.0.0.0"
    }

    port := cfg.Port
    if port == 0 {
        port = 5683
    }

    address := fmt.Sprintf("%s:%d", host, port)

    udpAddr, err := net.ResolveUDPAddr("udp", address)
    if err != nil {
        if global.GVA_LOG != nil {
            global.GVA_LOG.Error("resolve CoAP udp addr failed", zap.Error(err))
        }
        return
    }

    listener, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        if global.GVA_LOG != nil {
            global.GVA_LOG.Error("listen on CoAP udp failed", zap.Error(err))
        }
        return
    }

    global.GVA_COAP_SERVER = listener
    if global.GVA_LOG != nil {
        global.GVA_LOG.Info("CoAP server started", zap.String("address", address))
    }

    go serveCoap(listener)
}

func serveCoap(conn *net.UDPConn) {
    buffer := make([]byte, 1500)
    for {
        n, addr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            if !errors.Is(err, net.ErrClosed) && global.GVA_LOG != nil {
                global.GVA_LOG.Error("CoAP server read failed", zap.Error(err))
            }
            return
        }

        packet := append([]byte(nil), buffer[:n]...)
        msg, err := parseCoapMessage(packet)
        if err != nil {
            if global.GVA_LOG != nil {
                global.GVA_LOG.Warn("invalid CoAP message", zap.String("remote", addr.String()), zap.Error(err))
            }
            continue
        }

        responseType := coapTypeAcknowledgement
        if msg.typ != coapTypeConfirmable {
            responseType = coapTypeNonConfirmable
        }

        var (
            respCode = coapCodeNotFound
            payload  []byte
        )

        switch msg.code {
        case coapCodeGet:
            path := msg.path()
            if path == "/health" {
                respCode = coapCodeContent
                payload = []byte("ok")
            } else {
                respCode = coapCodeNotFound
            }
        default:
            respCode = coapCodeMethodNotAllowed
        }

        response := buildCoapResponse(msg, responseType, respCode, payload)
        if _, err = conn.WriteToUDP(response, addr); err != nil {
            if global.GVA_LOG != nil {
                global.GVA_LOG.Error("CoAP server write failed", zap.Error(err))
            }
        }
    }
}

func parseCoapMessage(data []byte) (coapMessage, error) {
    if len(data) < 4 {
        return coapMessage{}, errors.New("coap message too short")
    }

    version := int(data[0] >> 6)
    if version != coapVersion {
        return coapMessage{}, fmt.Errorf("unsupported coap version %d", version)
    }

    typ := int((data[0] >> 4) & 0x03)
    tkl := int(data[0] & 0x0F)
    if tkl > 8 {
        return coapMessage{}, errors.New("invalid token length")
    }

    if len(data) < 4+tkl {
        return coapMessage{}, errors.New("coap message truncated")
    }

    msg := coapMessage{
        typ:       typ,
        code:      data[1],
        messageID: uint16(data[2])<<8 | uint16(data[3]),
        options:   make(map[uint16][][]byte),
    }

    if tkl > 0 {
        msg.token = append([]byte(nil), data[4:4+tkl]...)
    }

    pos := 4 + tkl
    currentOpt := uint16(0)

    for pos < len(data) {
        if data[pos] == 0xFF {
            pos++
            if pos < len(data) {
                msg.payload = append([]byte(nil), data[pos:]...)
            }
            return msg, nil
        }

        deltaNibble := (data[pos] & 0xF0) >> 4
        lengthNibble := data[pos] & 0x0F
        pos++

        delta, err := readExtendedValue(deltaNibble, data, &pos)
        if err != nil {
            return coapMessage{}, err
        }

        length, err := readExtendedValue(lengthNibble, data, &pos)
        if err != nil {
            return coapMessage{}, err
        }

        if pos+int(length) > len(data) {
            return coapMessage{}, errors.New("coap option exceeds length")
        }

        currentOpt += delta
        value := append([]byte(nil), data[pos:pos+int(length)]...)
        msg.options[currentOpt] = append(msg.options[currentOpt], value)
        pos += int(length)
    }

    return msg, nil
}

func readExtendedValue(nibble byte, data []byte, pos *int) (uint16, error) {
    switch nibble {
    case 13:
        if *pos >= len(data) {
            return 0, errors.New("coap option truncated (extended 13)")
        }
        val := uint16(data[*pos]) + 13
        *pos += 1
        return val, nil
    case 14:
        if *pos+1 >= len(data) {
            return 0, errors.New("coap option truncated (extended 14)")
        }
        val := uint16(data[*pos])<<8 | uint16(data[*pos+1])
        *pos += 2
        return val + 269, nil
    case 15:
        return 0, errors.New("coap option nibble 15 reserved")
    default:
        return uint16(nibble), nil
    }
}

func buildCoapResponse(req coapMessage, respType int, code byte, payload []byte) []byte {
    tkl := len(req.token)
    if tkl > 8 {
        tkl = 8
    }

    header := make([]byte, 4+tkl)
    header[0] = byte((coapVersion << 6) | ((respType & 0x03) << 4) | byte(tkl))
    header[1] = code
    header[2] = byte(req.messageID >> 8)
    header[3] = byte(req.messageID)
    if tkl > 0 {
        copy(header[4:], req.token[:tkl])
    }

    if len(payload) == 0 {
        return header
    }

    resp := append(header, 0xFF)
    resp = append(resp, payload...)
    return resp
}

func (m coapMessage) path() string {
    segments := m.options[coapOptionURIPath]
    if len(segments) == 0 {
        return "/"
    }

    parts := make([]string, 0, len(segments))
    for _, seg := range segments {
        parts = append(parts, string(seg))
    }

    return "/" + strings.Join(parts, "/")
}
