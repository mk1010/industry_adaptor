package device

import (
	"context"
	"fmt"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/task/common"
)

func DeviceInit(ctx context.Context, deviceMeta *nclink.NCLinkDevice, adaptorID string) (common.NCLinkDeviceAPI, error) {
	if deviceMeta == nil || adaptorID == "" {
		err := fmt.Errorf("代理器元数据为空或适配器ID为空")
		logger.Error(err)
		return nil, err
	}
	switch deviceMeta.DeviceType {
	case "nclink_common_device":
		{
			commonDevice := &NCLinkCommonDevice{
				AdaptorID:  adaptorID,
				DeviceID:   deviceMeta.DeviceId,
				DeviceMeta: deviceMeta,
			}
			err := commonDevice.Start(ctx)
			if err != nil {
				return nil, err
			}
			return commonDevice, nil
		}
	default:
		{
			err := fmt.Errorf("未知的设备类型元数据%+v", deviceMeta)
			logger.Error(err)
			return nil, err
		}
	}
}
