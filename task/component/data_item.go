package component

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/dubbo-go/common/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/nclink/util"
	"github.com/mk1010/industry_adaptor/service"
)

// 不同组件具有相同数据项类型，且可以对调用方法进行复用时，建议注册到该对象上。
type NCLinkCommonDataInfo struct {
	AdaptorID      string
	DeviceID       string
	ComponentID    string
	DataInfo       *nclink.NCLinkDataInfo
	DataChan       chan []byte
	sendTimeTicker *time.Ticker
	done           chan struct{}
	dataPayloads   []*nclink.NCLinkPayloads
	mu             sync.Mutex
	once           sync.Once
}

func (n *NCLinkCommonDataInfo) Start(ctx context.Context) error {
	// 这里不做参数校验了 相信DataInfoInit
	util.GoSafely(func() {
		n.listen(ctx)
	}, nil)
	return nil
}

func (n *NCLinkCommonDataInfo) listen(ctx context.Context) {
	n.sendTimeTicker = time.NewTicker(time.Duration(n.DataInfo.SampleInfo.UploadPeriod) * time.Millisecond)
	util.GoSafely(
		func() {
			n.sendData(ctx)
		}, nil)
	for {
		select {
		case data, ok := <-n.DataChan:
			{
				if !ok {
					logger.Errorf("收集数据通道被关闭  datainfo:%v", n)
					n.Shutdown()
					return
				}
				now := util.TimeToUnixMs(time.Now())
				originData := data
				dataMap := make(map[string]interface{})
				// 这里根据设备发送数据不同进行变换 这里示例 设备发送过来数据为json格式的

				//err:=jsoniter.Unmarshal(data, dataMap)
				//if err!=nil{
				//	logger.Errorf("数据解析失败  datainfo:%v bytes:%v",n,data)
				//	continue
				//}
				//这里是将其视为大头端数据流示例
				var val interface{}
				var err error
				byteBuffer := bytes.NewBuffer(data)
				for _, item := range n.DataInfo.DataItem.Items {
					val, err = GetFiledValueBigEnd(item.Kind, byteBuffer)
					if err != nil {
						logger.Errorf("数据解析失败 origin data:%v item:%v", originData, item)
						continue
					}
					dataMap[item.FiledName] = val
				}
				var payload []byte
				_ = jsoniter.Unmarshal(payload, dataMap)
				nclinkPayload := &nclink.NCLinkPayloads{
					UnixTimeMs: now,
					Payload:    payload,
				}
				n.mu.Lock()
				n.dataPayloads = append(n.dataPayloads, nclinkPayload)
				n.mu.Unlock()
			}
		case <-n.done:
			{
				return
			}
		}
	}
}

func (n *NCLinkCommonDataInfo) sendData(ctx context.Context) {
	var dataPayloads []*nclink.NCLinkPayloads
	for {
		select {
		case <-n.sendTimeTicker.C:
			{
				n.mu.Lock()
				dataPayloads = n.dataPayloads
				n.dataPayloads = nil
				n.mu.Unlock()
				resp := new(nclink.NCLinkBaseResp)
				in := &nclink.NCLinkDataMessage{
					DataId:      util.GetUuid(),
					DeviceId:    n.DeviceID,
					ComponentId: n.ComponentID,
					DataItemId:  n.DataInfo.DataItem.DataItemId,
					Payloads:    dataPayloads,
				}
				err := service.NCLinkClient.NCLinkSendData(ctx, in, resp)
				if err != nil || resp.StatusCode != nclink.StatusOk {
					logger.Errorf("生成数据发送失败 err=%v resp=%v 数据项元数据%+v\n 发送数据%+v", err, resp, n, in)
				}
			}
		case <-n.done:
			{
				return
			}
		}
	}
}

func (n *NCLinkCommonDataInfo) SendData(data []byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	n.DataChan <- data
	return
}

func (n *NCLinkCommonDataInfo) UpdateMeta(ctx context.Context, meta *nclink.NCLinkDataInfo) error {
	if meta == nil {
		return n.Shutdown()
	}
	if n.DataInfo.SampleInfo.SampleInfoId != meta.SampleInfo.SampleInfoId {
		n.sendTimeTicker.Reset(time.Duration(meta.SampleInfo.UploadPeriod) * time.Millisecond)
		// todo
		// n.SendMsgChan
	}
	n.DataInfo = meta
	return nil
}

func (n *NCLinkCommonDataInfo) Shutdown() (err error) {
	n.once.Do(func() {
		close(n.done)
		n.sendTimeTicker.Stop()
	})
	return nil
}
