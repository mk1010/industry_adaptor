package adaptor

import (
	"context"
	"fmt"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
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
			err = commonAda.Start(ctx)
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
