package evhub

import (
	"io"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type EvhubClient struct {
	netConn           net.Conn
	serverAddtr       string
	connTime          int64
	headerSlice       []byte
	secondHeaderBytes []byte
	lastMsgTime       int64
}

func NewEvhubClient(connAddr string) *EvhubClient {
	evc := &EvhubClient{
		netConn:           nil,
		serverAddtr:       connAddr,
		connTime:          0,
		headerSlice:       make([]byte, MsgHeadSize),
		secondHeaderBytes: make([]byte, SecondHeaderSize),
		lastMsgTime:       0,
	}
	return evc
}

func (evc *EvhubClient) connect() error {
	newConn, err := net.Dial("tcp", evc.serverAddtr)
	if err != nil {
		logrus.Errorf("evhubclient connect server(%s) error:%s", evc.serverAddtr, err.Error())
		return err
	}
	evc.netConn = newConn
	evc.connTime = time.Now().Unix()
	//evc.lastMsgTime = evc.connTime
	return nil
}
func (evc *EvhubClient) SendMsg(msg *NetMessage) error {
	wData := msg.Encode()
	var err error = nil
	var tryCount int32 = 0
	if evc.netConn == nil {
		err = evc.connect()
		if err != nil {
			return err
		}
	}
	for {
		_, err = evc.netConn.Write(wData)
		if err != nil {
			logrus.Errorf("EvhubClient send data to server failed, %s", err)
			evc.netConn.Close()
			if tryCount < 3 {
				errConn := evc.connect()
				if errConn == nil {
					// 重连成功,继续发送
					tryCount++
					time.Sleep(time.Millisecond * 1)
					continue
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	if err == nil {
		evc.lastMsgTime = time.Now().Unix()
	}

	return err
}

func (evc *EvhubClient) RecvMsg(timeoutMs int64) (*NetMessage, error) {
	evc.netConn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeoutMs)))
	_, err := io.ReadFull(evc.netConn, evc.headerSlice)
	if err != nil {
		return nil, err
	}

	hd := &NetMsgHead{}
	hd.Decode(evc.headerSlice)
	var hd2 *NetMsgSecondHead = nil
	if hd.HasSecond > 0 {
		_, err := io.ReadFull(evc.netConn, evc.secondHeaderBytes)
		if err != nil {
			logrus.Errorf("EvhubClient read failed, %s", err)
			return nil, err
		}
		hd2 = &NetMsgSecondHead{}
		hd2.Decode(evc.secondHeaderBytes)
	}

	messageBody := make([]byte, hd.Len-MsgHeadSize)
	_, err = io.ReadFull(evc.netConn, messageBody)
	if err != nil {
		logrus.Errorf("read msgbody failed, %s", err.Error())
		return nil, err
	}
	netMsg := NewNetMessage(hd, messageBody)
	netMsg.SecondHead = hd2
	netMsg.RecvTime = time.Now().UnixNano() / 1000000
	return netMsg, nil
}
func (evc *EvhubClient) Close() {
	if evc.netConn != nil {
		evc.netConn.Close()
		evc.netConn = nil
	}
}
