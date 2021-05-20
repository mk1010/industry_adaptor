package adaptor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/service"
	"github.com/mk1010/industry_adaptor/task/common"
	"github.com/mk1010/industry_adaptor/task/device"
)

type NCLinkCommonAdaptor struct {
	AdaptorID   string
	AdaptorMeta *nclink.NCLinkAdaptor
	DeviceMap   map[string]common.NCLinkDeviceAPI
	mu          sync.Mutex
}

func (ada *NCLinkCommonAdaptor) Start(ctx context.Context) (err error) {
	if ada.AdaptorMeta == nil {
		err = fmt.Errorf("适配器%s下没有适配器元数据", ada.AdaptorMeta.AdaptorId)
		logger.Error(err)
		return
	}
	ada.AdaptorID = ada.AdaptorMeta.AdaptorId
	if len(ada.AdaptorMeta.DeviceId) <= 0 {
		err = fmt.Errorf("适配器%s下没有管理设备", ada.AdaptorMeta.AdaptorId)
		logger.Error(err)
		return
	}
	metaResp := new(nclink.NCLinkMetaDataResp)
	err = service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
		DeviceId: ada.AdaptorMeta.DeviceId,
	}, metaResp)
	if err != nil || len(metaResp.Devices) <= 0 {
		err = fmt.Errorf("获取设备元数据失败 %v", err)
		logger.Error(err)
		return
	}
	ada.mu.Lock()
	defer ada.mu.Unlock()
	if ada.DeviceMap == nil {
		ada.DeviceMap = make(map[string]common.NCLinkDeviceAPI, len(ada.AdaptorMeta.DeviceId))
	}
	for _, deviceMeta := range metaResp.Devices {
		deviceAPI, err := device.DeviceInit(ctx, deviceMeta, ada.AdaptorID)
		if err != nil {
			continue
		}
		ada.DeviceMap[deviceMeta.DeviceId] = deviceAPI
	}
	if err != nil {
		common.NCLinkAdaptorMap.Store(ada.AdaptorID, ada)
	}
	return nil
}

func (ada *NCLinkCommonAdaptor) UpdateMeta(ctx context.Context, meta *nclink.NCLinkAdaptor) (err error) {
	if meta == nil {
		return ada.Shutdown()
	}
	metaResp := new(nclink.NCLinkMetaDataResp)
	err = service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
		DeviceId: meta.DeviceId,
	}, metaResp)
	if err != nil || len(metaResp.Devices) <= 0 {
		err = fmt.Errorf("获取设备元数据失败 %v", err)
		logger.Error(err)
		return
	}
	deviceMetaMap := make(map[string]*nclink.NCLinkDevice, len(metaResp.Devices))
	for _, deviceMeta := range metaResp.Devices {
		deviceMetaMap[deviceMeta.DeviceId] = deviceMeta
	}
	ada.mu.Lock()
	defer ada.mu.Unlock()
	for deviceID, deviceAPI := range ada.DeviceMap {
		if _, ok := deviceMetaMap[deviceID]; !ok {
			e := deviceAPI.Shutdown()
			if e != nil {
				err = e
			}
			delete(ada.DeviceMap, deviceID)
		}
	}
	for deviceID, deviceMeta := range deviceMetaMap {
		if _, ok := ada.DeviceMap[deviceID]; !ok {
			deviceAPI, e := device.DeviceInit(ctx, deviceMeta, ada.AdaptorID)
			if e != nil {
				err = e
				continue
			}
			ada.DeviceMap[deviceID] = deviceAPI
		}
	}
	ada.AdaptorMeta = meta
	return
}

func (ada *NCLinkCommonAdaptor) GetMeta() *nclink.NCLinkAdaptor {
	return ada.AdaptorMeta
}

func (ada *NCLinkCommonAdaptor) Shutdown() (err error) {
	ada.mu.Lock()
	defer ada.mu.Unlock()
	for _, device := range ada.DeviceMap {
		e := device.Shutdown()
		if err != nil {
			err = e
		}
	}
	// 这里是高危操作 有必要注释下
	// 首先无条件删除当前存在的键值对，并把当前键值对返回
	// 如果当前键值对存在，那么进行一次类型转换比较是不是删除了当前指针
	// 如果不是 那么按一定策略写回去，由于两个原子操作一起做就不是原子操作
	// 先检查当前这个键值对是否被写入了，如果写入了那么就应该是被写入的值
	// 如果没被写入 值应该是之前被删除的值
	val, done := common.NCLinkAdaptorMap.LoadAndDelete(ada.AdaptorID)
	if done {
		if v, ok := val.(*NCLinkCommonAdaptor); ok && v != ada {
			common.NCLinkAdaptorMap.LoadOrStore(ada.AdaptorID, val)
		}
	}
	return
}
