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
	t.mu.Lock()
	if t.DataInfoMap == nil {
		t.DataInfoMap = make(map[string]common.NCLinkDataInfoAPI)
	}
	for _, dataInfo := range t.ComponentMeta.DataInfo {
		dataInfoAPi, err := DataInfoInit(ctx, t.AdaptorID, t.DeviceID, t.ComponentID, dataInfo)
		if err != nil {
			logger.Error("数据项启动失败 component:%v \n data_info:%v \n err:%v \n", t, dataInfo, err)
			continue
		}
		t.DataInfoMap[dataInfo.DataItem.DataItemId] = dataInfoAPi
	}
	t.mu.Unlock()
	if err != nil {
		common.NCLinkComponentMap.Store(t.ComponentID, t)
	}
	return nil
}

func (t *NCLinkCommonComponent) GetDataInfoApi(dataItemID string) common.NCLinkDataInfoAPI {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.DataInfoMap[dataItemID]
}

func (t *NCLinkCommonComponent) UpdateMeta(ctx context.Context, meta *nclink.NCLinkComponent) (err error) {
	return
}

func (t *NCLinkCommonComponent) Shutdown() (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for dataItemID, dataInfo := range t.DataInfoMap {
		e := dataInfo.Shutdown()
		if err != nil {
			err = e
		}
		delete(t.DataInfoMap, dataItemID)
	}
	val, done := common.NCLinkComponentMap.LoadAndDelete(t.DeviceID)
	if done && val.(*NCLinkCommonComponent) != t {
		common.NCLinkComponentMap.LoadOrStore(t.DeviceID, val)
	}
	return
}