package adaptor

import (
	"context"
	"fmt"
	"io"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/service"
)

func AdaptorInit(ctx context.Context, adaMeta *nclink.NCLinkAdaptor) (err error) {
	if adaMeta == nil {
		err = fmt.Errorf("适配器元数据为空")
		logger.Error(err)
		return
	}
	switch adaMeta.AdaptorType {
	case "nclink_commmon_adaptor":
		{
			commonAda := &NCLinkCommonAdaptor{
				AdaptorID:   adaMeta.AdaptorId,
				AdaptorMeta: adaMeta,
			}
			err = commonAda.Start(context.Background())
			if err != nil {
				return err
			}
			return
		}
	default:
		{
			err = fmt.Errorf("未知的适配器类型元数据 %+v", adaMeta)
			logger.Error(err)
			return
		}
	}
}

func NcLinkCommandTopic(ctx context.Context, adaptorID string) {
	for {
		subClient, err := service.NCLinkClient.NCLinkSubscribe(ctx, &nclink.NCLinkTopicSub{
			Topic:     nclink.CommandTopic,
			AdaptorId: adaptorID,
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
			// do sth
			}
		}
	}
}
