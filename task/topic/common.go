package topic

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/dubbo-go/common/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/mk1010/industry_adaptor/config"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/nclink/util"
	"github.com/mk1010/industry_adaptor/service"
)

func Init(ctx context.Context) {
	util.GoSafely(func() {
		NCLinkTopicSub(ctx, config.ConfInstance.AdaptorID, nclink.CommandTopic)
	}, nil)
}

func NCLinkTopicSub(ctx context.Context, adaptorID, topic string) {
	loop := true
	count := 0
	for ; loop; time.Sleep(3 * time.Second) {
		subClient, err := service.NCLinkClient.NCLinkSubscribe(ctx)
		if err != nil {
			err = fmt.Errorf("nclink topic:%s 订阅失败 %v", topic, err)
			logger.Error(err)
			continue
		}
		subTopic := &nclink.NCLinkTopicSub{
			Topic:     topic,
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
		// 各topic自行处理err
		switch topic {
		case nclink.CommandTopic:
			// 在这里进行loop配置
			count++
			// if count>100{
			// loop=false
			// }
			if err != nil {
				err = fmt.Errorf("nclink topic:%s 订阅失败 %v", topic, err)
				logger.Error(err)
				continue
			}
			NcLinkCommandTopic(ctx, adaptorID, subClient)
		default:
			logger.Error("无法订阅topic:", topic)
			return
		}
	}
}
