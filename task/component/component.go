package component

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/task/common"
)

type NCLinkCommonComponent struct {
	AdaptorID     string
	DeviceID      string
	ComponentID   string
	ComponentMeta *nclink.NCLinkComponent
	mu            sync.Mutex
	DataInfoMap   map[string]common.NCLinkDataInfoAPI
}

func (t *NCLinkCommonComponent) Start(ctx context.Context) (err error) {
	if t.ComponentMeta == nil || t.AdaptorID == "" || t.DeviceID == "" {
		err = fmt.Errorf("未被标识的设备")
		logger.Error(err)
		return
	}
	t.ComponentID = t.ComponentMeta.ComponentId
	// 配置文件解析，但是没想好干嘛
	// configMap := make(map[string]interface{})
	// if len(t.ComponentMeta.Config) > 0 {
	// e := jsoniter.Unmarshal(t.ComponentMeta.Config, configMap)
	// if e != nil {
	// logger.Error("组件配置文件解析失败 %v", e)
	// }
	// }
	for _, dataInfo := range t.ComponentMeta.DataInfo {
		DataInfoInit(ctx, t.AdaptorID, t.DeviceID, t.ComponentID, dataInfo)
	}
	if err != nil {
		common.NCLinkComponentMap.Store(t.ComponentID, t)
	}
	return nil
}

func (t *NCLinkCommonComponent) GetDataInfoApi(ctx, dataInfoID string) common.NCLinkDataInfoAPI {
	return nil
}

func (t *NCLinkCommonComponent) UpdateMeta(ctx context.Context, meta *nclink.NCLinkComponent) (err error) {
	return
}

func (t *NCLinkCommonComponent) Shutdown() error {
	return nil
}
