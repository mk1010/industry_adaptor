package device

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/service"
	"github.com/mk1010/industry_adaptor/task/common"
	"github.com/mk1010/industry_adaptor/task/component"
)

type NCLinkCommonDevice struct {
	AdaptorID    string
	DeviceID     string
	DeviceMeta   *nclink.NCLinkDevice
	ComponentMap map[string]common.NCLinkComponentAPI
	mu           sync.Mutex
}

func (t *NCLinkCommonDevice) Start(ctx context.Context) (err error) {
	if t.AdaptorID == "" || t.DeviceMeta == nil {
		err = fmt.Errorf("未被标识的设备")
		logger.Error(err)
		return
	}
	t.DeviceID = t.DeviceMeta.DeviceId
	if len(t.DeviceMeta.ComponentId) <= 0 {
		err = fmt.Errorf("设备%s下没有管理组件", t.DeviceMeta.DeviceId)
		logger.Error(err)
		return
	}
	metaResp := new(nclink.NCLinkMetaDataResp)
	err = service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
		ComponentId: t.DeviceMeta.ComponentId,
	}, metaResp)
	if err != nil || len(metaResp.Components) <= 0 {
		err = fmt.Errorf("获取组件元数据失败 %v", err)
		logger.Error(err)
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ComponentMap == nil {
		t.ComponentMap = make(map[string]common.NCLinkComponentAPI, len(t.DeviceMeta.ComponentId))
	}
	for _, componentMeta := range metaResp.Components {
		componentAPI, err := component.ComponentInit(ctx, componentMeta, t.AdaptorID, t.DeviceID)
		if err != nil {
			return err
		}
		t.ComponentMap[componentMeta.ComponentId] = componentAPI
	}
	if err != nil {
		common.NCLinkDeviceMap.Store(t.DeviceID, t)
	}
	return nil
}

func (t *NCLinkCommonDevice) UpdateMeta(ctx context.Context, meta *nclink.NCLinkDevice) (err error) {
	if meta == nil {
		return t.Shutdown()
	}
	metaResp := new(nclink.NCLinkMetaDataResp)
	err = service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
		ComponentId: meta.ComponentId,
	}, metaResp)
	if err != nil || len(metaResp.Components) <= 0 {
		err = fmt.Errorf("获取组件元数据失败 %v", err)
		logger.Error(err)
		return
	}
	componentMetaMap := make(map[string]*nclink.NCLinkComponent, len(metaResp.Components))
	for _, componentMeta := range metaResp.Components {
		componentMetaMap[componentMeta.ComponentId] = componentMeta
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for componentID, componentAPI := range t.ComponentMap {
		if _, ok := componentMetaMap[componentID]; !ok {
			e := componentAPI.Shutdown()
			if e != nil {
				err = e
			}
			delete(t.ComponentMap, componentID)
		}
	}
	for componentID, componentMeta := range componentMetaMap {
		if _, ok := t.ComponentMap[componentID]; !ok {
			componentAPI, e := component.ComponentInit(ctx, componentMeta, t.AdaptorID, t.DeviceID)
			if e != nil {
				err = e
				continue
			}
			t.ComponentMap[componentID] = componentAPI
		}
	}
	t.DeviceMeta = meta
	return
}

func (t *NCLinkCommonDevice) GetMeta() *nclink.NCLinkDevice {
	return t.DeviceMeta
}

func (t *NCLinkCommonDevice) Shutdown() (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for componentID, component := range t.ComponentMap {
		e := component.Shutdown()
		if err != nil {
			err = e
		}
		delete(t.ComponentMap, componentID)
	}
	// 这里是高危操作 有必要注释下
	// 首先无条件删除当前存在的键值对，并把当前键值对返回
	// 如果当前键值对存在，那么进行一次类型转换比较是不是删除了当前指针
	// 如果不是 那么按一定策略写回去，由于两个原子操作一起做就不是原子操作
	// 先检查当前这个键值对是否被写入了，如果写入了那么就应该是被写入的值
	// 如果没被写入 值应该是之前被删除的值
	val, done := common.NCLinkDeviceMap.LoadAndDelete(t.DeviceID)
	if done {
		if v, ok := val.(*NCLinkCommonDevice); ok && v != t {
			common.NCLinkDeviceMap.LoadOrStore(t.DeviceID, val)
		}
	}
	return
}
