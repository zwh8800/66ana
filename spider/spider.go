package spider

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/proxy"
)

const (
	roomInfoUrl          = "http://open.douyucdn.cn/api/RoomApi/room/%d"
	openDouyuAddr        = "openbarrage.douyutv.com:8601"
	msgTypeSend   uint16 = 689
	msgTypeRecv   uint16 = 690
)

type giftData struct {
	Data struct {
		Gift []struct {
			Desc  string  `json:"desc"`
			Gx    float64 `json:"gx"`
			Himg  string  `json:"himg"`
			ID    string  `json:"id"`
			Intro string  `json:"intro"`
			Mimg  string  `json:"mimg"`
			Name  string  `json:"name"`
			Pc    float64 `json:"pc"`
			Type  string  `json:"type"`
		} `json:"gift"`
	} `json:"data"`
	Error int `json:"error"`
}

type SpiderStatus int

const (
	StatusRunning SpiderStatus = iota
	StatusClosed  SpiderStatus = iota
	StatusError   SpiderStatus = iota
)

type Spider struct {
	roomId     int
	giftMap    map[string]string
	conn       net.Conn
	msgChan    chan map[string]string
	status     SpiderStatus
	lastError  error
	httpClient *http.Client
}

func NewSpider(roomId int, dialer proxy.Dialer) (*Spider, error) {
	s := &Spider{
		roomId:  roomId,
		msgChan: make(chan map[string]string, 256),
		status:  StatusRunning,
	}
	var conn net.Conn
	var err error
	if dialer != nil {
		conn, err = dialer.Dial("tcp", openDouyuAddr)
		s.httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: dialer,
			},
		}
	} else {
		conn, err = net.Dial("tcp", openDouyuAddr)
		s.httpClient = http.DefaultClient
	}
	if err != nil {
		s.status = StatusError
		s.lastError = err
		return nil, err
	}
	s.conn = conn
	s.run()
	return s, nil
}

func (s *Spider) run() {
	if err := s.danmukuLogin(); err != nil {
		s.status = StatusError
		s.lastError = err
		return
	}
	if err := s.danmukuJoin(); err != nil {
		s.status = StatusError
		s.lastError = err
		return
	}
	go func() {
		for {
			if s.status != StatusRunning {
				return
			}

			if err := s.danmukuKeeplive(); err != nil {
				s.status = StatusError
				s.lastError = err
				return
			}
			time.Sleep(30 * time.Second)
		}
	}()

	go func() {
		defer func() {
			close(s.msgChan)
		}()
		for {
			if s.status != StatusRunning {
				return
			}

			message, err := s.danmukuReadAndPipe()
			if err != nil {
				s.status = StatusError
				s.lastError = err
				return
			} else {
				s.msgChan <- message
			}
		}
	}()
}

func (s *Spider) sendMsg(msg string) error {
	msgLen := len(msg) + 1 + 12
	buf := make([]byte, msgLen)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(msgLen-4))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(msgLen-4))
	binary.LittleEndian.PutUint16(buf[8:10], msgTypeSend)
	copy(buf[12:], []byte(msg))

	_, err := s.conn.Write(buf)
	return err
}

func (s *Spider) danmukuLogin() error {
	msg := fmt.Sprintf("type@=loginreq/roomid@=%d/", s.roomId)
	if err := s.sendMsg(msg); err != nil {
		return err
	}
	return nil
}

func (s *Spider) danmukuJoin() error {
	msg := fmt.Sprintf("type@=joingroup/rid@=%d/gid@=-9999/", s.roomId)
	if err := s.sendMsg(msg); err != nil {
		return err
	}
	return nil
}

func (s *Spider) danmukuKeeplive() error {
	msg := fmt.Sprintf("type@=keeplive/tick@=%d/", time.Now().Unix())
	if err := s.sendMsg(msg); err != nil {
		return err
	}
	return nil
}

func (s *Spider) readMessage() (string, error) {
	var (
		length      uint32
		length2     uint32
		messageType uint16
		unused      uint16
	)
	if err := binary.Read(s.conn, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	if err := binary.Read(s.conn, binary.LittleEndian, &length2); err != nil {
		return "", err
	}
	if length != length2 {
		return "", fmt.Errorf("length(%d) != length2(%d)\n", length, length2)
	}
	if err := binary.Read(s.conn, binary.LittleEndian, &messageType); err != nil {
		return "", err
	}
	if messageType != msgTypeRecv {
		return "", fmt.Errorf("messageData(%d) != typeRecv\n", messageType)
	}
	if err := binary.Read(s.conn, binary.LittleEndian, &unused); err != nil {
		return "", err
	}
	length = length - 8
	messageData := make([]byte, length)

	for i := 0; i < int(length); {
		n, err := s.conn.Read(messageData[i:])
		if err != nil {
			return "", err
		}
		i += n
	}

	return string(messageData), nil
}

func (s *Spider) danmukuReadAndPipe() (map[string]string, error) {
	msgStr, err := s.readMessage()
	if err != nil {
		return nil, err
	}
	message := parseMessage(msgStr)
	return message, nil
}

func (s *Spider) GetMessageChan() chan map[string]string {
	return s.msgChan
}

func (s *Spider) GetStatus() SpiderStatus {
	return s.status
}

func (s *Spider) GetLastError() error {
	return s.lastError
}

func (s *Spider) GetGiftMap() map[string]string {
	if s.giftMap == nil {
		resp, err := s.httpClient.Get(fmt.Sprintf(roomInfoUrl, s.roomId))
		if err != nil {
			return make(map[string]string)
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return make(map[string]string)
		}

		var giftData giftData
		if err := json.Unmarshal(data, &giftData); err != nil {
			return make(map[string]string)
		}

		if giftData.Error != 0 {
			return make(map[string]string)
		}

		s.giftMap = make(map[string]string, len(giftData.Data.Gift))

		for _, gift := range giftData.Data.Gift {
			s.giftMap[gift.ID] = gift.Name
		}
	}

	return s.giftMap
}

var msgRegex = regexp.MustCompile(`(.*?)@=(.*?)/`)

func parseMessage(message string) map[string]string {
	msg := make(map[string]string)

	submatchs := msgRegex.FindAllStringSubmatch(message, -1)

	for _, submatch := range submatchs {
		msg[submatch[1]] = submatch[2]
	}
	return msg
}
