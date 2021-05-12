package component

import (
	"context"
	"fmt"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/task/common"
)

func ComponentInit(ctx context.Context, componentMeta *nclink.NCLinkComponent, adaptorID, deviceID string) (common.NCLinkComponentAPI, error) {
	if componentMeta == nil || len(adaptorID) == 0 || len(deviceID) == 0 {
		err := fmt.Errorf("组件元数据为空或适配器ID或设备ID为空")
		logger.Error(err)
		return nil, err
	}
	switch componentMeta.ComponentType {
	case "nclink_common_component":
		{
			commonComponent := &NCLinkCommonComponent{
				AdaptorID:     adaptorID,
				DeviceID:      deviceID,
				ComponentID:   componentMeta.ComponentId,
				ComponentMeta: componentMeta,
			}
			err := commonComponent.Start(ctx)
			if err != nil {
				return nil, err
			}
			return commonComponent, nil
		}
	default:
		{
			err := fmt.Errorf("未知的组件类型元数据%+v", componentMeta)
			logger.Error(err)
			return nil, err
		}
	}
}
