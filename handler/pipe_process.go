package handler

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/apache/dubbo-go/common/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/nclink/util"
	"github.com/mk1010/industry_adaptor/task/common"
)

type NCLinkCommonProcess struct {
	ProcessName   string
	route         map[string]map[string]common.NCLinkDataInfoAPI // map[comonentID]map[dataitemID]common.NCLinkDataInfoAPI
	routeMu       sync.Mutex
	WriteChan     chan *nclink.NCLinkTopicMessage
	CmdStdinPipe  io.WriteCloser
	CmdStdoutPipe io.ReadCloser
}

func (t *NCLinkCommonProcess) Start(ctx context.Context) error {
	if t.CmdStdinPipe == nil || t.CmdStdoutPipe == nil {
		err := fmt.Errorf("进程%s重定向管道为空", t.ProcessName)
		logger.Error(err)
		return err
	}
	if t.WriteChan == nil {
		t.WriteChan = make(chan *nclink.NCLinkTopicMessage, 10)
	}
	if t.route == nil {
		t.route = make(map[string]map[string]common.NCLinkDataInfoAPI)
	}
	util.GoSafely(func() {
		for {
			select {
			case msg, ok := <-t.WriteChan:
				{
					if !ok {
						t.Restart()
						continue
					}
					b, _ := jsoniter.Marshal(msg)
					err := binary.Write(t.CmdStdinPipe, binary.BigEndian, int32(len(b)))
					if err != nil {
						logger.Error("数据长度写入进程失败", err, string(b))
					}
					haveWrite := 0
					retries := 0
					totalLen := len(b)
					for haveWrite < totalLen && retries < 3 {
						n, err := t.CmdStdinPipe.Write(b)
						if err != nil {
							logger.Error("数据写入进程失败", err, string(b))
							retries++
						}
						haveWrite += n
						b = b[n:]
					}
					if retries >= 3 {
						t.Restart()
						continue
					}
				}
			}
		}
	}, nil)
	util.GoSafely(func() {
		var err error
		for {
			haveRead := 0
			retries := 0
			var size int32
			for haveRead != 4 && retries < 3 {
				buf := make([]byte, 4)
				haveRead, err = io.ReadFull(t.CmdStdoutPipe, buf)
				if err != nil {
					logger.Error("读取进程长度数据失败", err)
					retries++
				}
				binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &size)
			}
			haveRead = 0
			var byteBuf *bytes.Buffer
			for int32(haveRead) != size && retries < 3 {
				buf := make([]byte, size)
				haveRead, err = io.ReadFull(t.CmdStdoutPipe, buf)
				if err != nil {
					logger.Error("读取进程数据失败", err)
					retries++
				}
				byteBuf = bytes.NewBuffer(buf)
			}
			if retries >= 3 {
				t.Restart()
				continue
			}
			deviceID, err := byteBuf.ReadString('\n')
			if err != nil {
				continue
			}
			componentID, _ := byteBuf.ReadString('\n')
			dataInfoID, _ := byteBuf.ReadString('\n')
			data, err := ioutil.ReadAll(byteBuf)
			if err != nil {
				continue
			}
			var dataItemApi common.NCLinkDataInfoAPI
			if dataInfoMap, ok := t.route[componentID]; !ok {
				dataItemApi = t.SearchNCLinkDataInfo(deviceID, componentID, dataInfoID)
				if dataItemApi != nil {
					t.RecvRegister(deviceID, componentID, dataInfoID, dataItemApi)
				}
			} else {
				if dataItemApi, ok = dataInfoMap[dataInfoID]; !ok {
					dataItemApi = t.SearchNCLinkDataInfo(deviceID, componentID, dataInfoID)
					if dataItemApi != nil {
						t.RecvRegister(deviceID, componentID, dataInfoID, dataItemApi)
					}
				}
			}
			if dataItemApi == nil {
				logger.Errorf("device:%s component:%s dataItem:%s 数据未找到接收者 data:%v", deviceID, componentID, dataInfoID, data)
				continue
			}
			err = dataItemApi.SendData(data)
			if err != nil {
				logger.Errorf("device:%s component:%s dataItem:%s 接收者接受数据失败 data:%v err%v", deviceID, componentID,
					dataInfoID, data, err)
			}
		}
	}, nil)
	common.NClinkInstanceMap.Store(t.ProcessName, t)
	return nil
}

func (t *NCLinkCommonProcess) SendData(msg *nclink.NCLinkTopicMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	t.WriteChan <- msg
	return
}

func (t *NCLinkCommonProcess) RecvRegister(deviceID, componentID, dataInfoID string, dataAPi common.NCLinkDataInfoAPI) error {
	t.routeMu.Lock()
	defer t.routeMu.Unlock()
	var dataInfoMap map[string]common.NCLinkDataInfoAPI
	var ok bool
	if dataInfoMap, ok = t.route[componentID]; !ok {
		dataInfoMap = make(map[string]common.NCLinkDataInfoAPI)
		t.route[componentID] = dataInfoMap
	}
	dataInfoMap[dataInfoID] = dataAPi
	return nil
}

func (t *NCLinkCommonProcess) RecvUnRegister(deviceID, componentID, dataInfoID string, dataAPi common.NCLinkDataInfoAPI) error {
	t.routeMu.Lock()
	defer t.routeMu.Unlock()
	var dataInfoMap map[string]common.NCLinkDataInfoAPI
	var ok bool
	if dataInfoMap, ok = t.route[deviceID]; !ok {
		return nil
	}
	delete(dataInfoMap, dataInfoID)
	return nil
}

func (t *NCLinkCommonProcess) Restart() {
	in, out, err := StartPipe(t.ProcessName)
	if err != nil {
		logger.Error("重新启动进程失败 err", err)
	}
	oriIn := t.CmdStdinPipe
	oriOut := t.CmdStdoutPipe
	t.CmdStdinPipe = in
	t.CmdStdoutPipe = out
	oriIn.Close()
	oriOut.Close()
}

func (t *NCLinkCommonProcess) SearchNCLinkDataInfo(deviceID, componentID, dataInfoID string) common.NCLinkDataInfoAPI {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	component, ok := common.NCLinkComponentMap.Load(componentID)
	if ok {
		componentApi := component.(common.NCLinkComponentAPI) // may panic
		return componentApi.GetDataInfoApi(dataInfoID)
	}
	return nil
}
