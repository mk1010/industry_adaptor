package task

import (
	"fmt"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
)

type NCLinkCommonAdaptor struct{}

func adaptorInit(adaMeta *nclink.NCLinkAdaptor) (err error) {
	if adaMeta == nil {
		err = fmt.Errorf("适配器元数据为空")
		logger.Error(err)
		return
	}
	switch adaMeta.AdaptorType {
	case "nclink_commmon_adaptor":
		{
			// commonAda := new(NCLinkCommonAdaptor)
		}
	default:
		{
			err = fmt.Errorf("未知的适配器元数据")
			logger.Error(err)
			return
		}
	}

	return
}
