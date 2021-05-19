package topic

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/apache/dubbo-go/common/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/mk1010/industry_adaptor/config"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/nclink/util"
	"github.com/mk1010/industry_adaptor/service"
	"github.com/mk1010/industry_adaptor/task"
	"github.com/mk1010/industry_adaptor/task/common"
)

func Init(ctx context.Context) {
	util.GoSafely(func() {
		NcLinkCommandTopic(ctx, config.ConfInstance.AdaptorID)
	}, nil)
}

func NcLinkCommandTopic(ctx context.Context, adaptorID string) {
	for {
		subClient, err := service.NCLinkClient.NCLinkSubscribe(ctx)
		if err != nil {
			err = fmt.Errorf("nclink命令topic订阅失败 %v", err)
			logger.Error(err)
			continue
		}
		subTopic := &nclink.NCLinkTopicSub{
			Topic:     nclink.CommandTopic,
			AdaptorId: adaptorID,
		}
		payload, _ := jsoniter.Marshal(subTopic)
		err = subClient.Send(&nclink.NCLinkTopicMessage{
			MessageId:   util.GetUuid(),
			MessageKind: int32(nclink.NclinkCommandMessageKind_Subscribe),
			Payload: &nclink.NCLinkPayloads{
				UnixTimeMs: util.TimeToUnixMs(time.Now()),
				Payload:    payload,
			},
		})
		if err != nil {
			err = fmt.Errorf("nclink命令topic订阅失败 %v", err)
			logger.Error(err)
			continue
		}
		for {
			msg, err := subClient.Recv()
			if err != nil {
				if err == io.EOF {
					logger.Error("nclink命令通道被关闭，可能是被重新订阅")
					return
				}
				logger.Error("nclink命令通道错误：", err)
				break
			}
			switch msg.MessageKind {
			case int32(nclink.NclinkCommandMessageKind_UpdateMeta):
				{
					resp := new(nclink.NCLinkMetaDataResp)
					err := jsoniter.Unmarshal(msg.Payload.Payload, resp)
					if err != nil {
						logger.Errorf("topic msg 解析失败 msg:%v err:%v", msg, err)
					}
					if len(resp.Adaptors) > 0 {
						for _, adaptor := range resp.Adaptors {
							val, ok := common.NCLinkAdaptorMap.Load(adaptor.AdaptorId)
							if ok {
								adaApi, ok := val.(common.NCLinkAdaptorAPI)
								if ok {
									adaApi.UpdateMeta(ctx, adaptor)
								}
							}
						}
						for _, device := range resp.Devices {
							val, ok := common.NCLinkAdaptorMap.Load(device.DeviceId)
							if ok {
								deviceApi, ok := val.(common.NCLinkDeviceAPI)
								if ok {
									deviceApi.UpdateMeta(ctx, device)
								}
							}
						}
						for _, component := range resp.Components {
							val, ok := common.NCLinkAdaptorMap.Load(component.ComponentId)
							if ok {
								componentApi, ok := val.(common.NCLinkComponentAPI)
								if ok {
									componentApi.UpdateMeta(ctx, component)
								}
							}
						}
					}
				}
			case int32(nclink.NclinkCommandMessageKind_GetMeta):
				{
					req := new(nclink.NCLinkMetaDataReq)
					err := jsoniter.Unmarshal(msg.Payload.Payload, req)
					if err != nil {
						logger.Errorf("topic msg 解析失败 msg:%v err:%v", msg, err)
					}
					resp := new(nclink.NCLinkMetaDataResp)
					for _, id := range req.AdaptorId {
						val, ok := common.NCLinkAdaptorMap.Load(id)
						if ok {
							adaApi, ok := val.(common.NCLinkAdaptorAPI)
							if ok {
								resp.Adaptors = append(resp.Adaptors, adaApi.GetMeta())
							}
						}
					}
					for _, id := range req.DeviceId {
						val, ok := common.NCLinkDeviceMap.Load(id)
						if ok {
							deviceApi, ok := val.(common.NCLinkDeviceAPI)
							if ok {
								resp.Devices = append(resp.Devices, deviceApi.GetMeta())
							}
						}
					}
					for _, id := range req.ComponentId {
						val, ok := common.NCLinkDeviceMap.Load(id)
						if ok {
							componentApi, ok := val.(common.NCLinkComponentAPI)
							if ok {
								resp.Components = append(resp.Components, componentApi.GetMeta())
							}
						}
					}
					payload, _ := jsoniter.Marshal(resp)
					subClient.Send(&nclink.NCLinkTopicMessage{
						MessageId:   msg.MessageId,
						MessageKind: nclink.NCLinkMsgResp,
						Payload: &nclink.NCLinkPayloads{
							UnixTimeMs: util.TimeToUnixMs(time.Now()),
							Payload:    payload,
						},
					})
				}
			case int32(nclink.NclinkCommandMessageKind_Shutdown):
				{
					req := new(nclink.NCLinkMetaDataReq)
					err := jsoniter.Unmarshal(msg.Payload.Payload, req)
					if err != nil {
						logger.Errorf("topic msg 解析失败 msg:%v err:%v", msg, err)
					}
					for _, id := range req.AdaptorId {
						val, ok := common.NCLinkAdaptorMap.Load(id)
						if ok {
							adaApi, ok := val.(common.NCLinkAdaptorAPI)
							if ok {
								adaApi.Shutdown()
							}
						}
					}
					for _, id := range req.DeviceId {
						val, ok := common.NCLinkDeviceMap.Load(id)
						if ok {
							deviceApi, ok := val.(common.NCLinkDeviceAPI)
							if ok {
								deviceApi.Shutdown()
							}
						}
					}
					for _, id := range req.ComponentId {
						val, ok := common.NCLinkDeviceMap.Load(id)
						if ok {
							componentApi, ok := val.(common.NCLinkComponentAPI)
							if ok {
								componentApi.Shutdown()
							}
						}
					}
				}
			case int32(nclink.NclinkCommandMessageKind_Restart):
				{
					val, ok := common.NCLinkAdaptorMap.Load(config.ConfInstance.AdaptorID)
					if ok {
						adaApi, ok := val.(common.NCLinkAdaptorAPI)
						if ok {
							adaApi.Shutdown()
						}
					}
					if err := task.Init(ctx); err != nil {
						subClient.CloseSend()
						return
					}
				}
			}
		}
	}
}
