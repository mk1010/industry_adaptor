package adaptor

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/nclink/util"
	"github.com/mk1010/industry_adaptor/service"
	"github.com/mk1010/industry_adaptor/task/common"
	"github.com/mk1010/industry_adaptor/task/device"
)

type NCLinkCommonAdaptor struct {
	AdaptorID string
	MetaData  *nclink.NCLinkAdaptor
	DeviceMap map[string]common.NCLinkDeviceAPI
	mu        sync.Mutex
}

func (ada *NCLinkCommonAdaptor) Start(ctx context.Context) (err error) {
	if len(ada.MetaData.DeviceId) <= 0 {
		err = fmt.Errorf("适配器%s下没有管理设备", ada.MetaData.AdaptorId)
		logger.Error(err)
		return
	}
	metaResp := new(nclink.NCLinkMetaDataResp)
	err = service.NCLinkClient.NCLinkGetMeta(ctx, &nclink.NCLinkMetaDataReq{
		DeviceId: ada.MetaData.DeviceId,
	}, metaResp)
	if err != nil || len(metaResp.Devices) <= 0 {
		err = fmt.Errorf("获取设备元数据失败 %v", err)
		logger.Error(err)
		return
	}

	util.GoSafely(func() {
		NcLinkCommandTopic(ctx, ada.AdaptorID)
	}, nil)
	ada.mu.Lock()
	defer ada.mu.Unlock()
	if ada.DeviceMap == nil {
		ada.DeviceMap = make(map[string]common.NCLinkDeviceAPI, len(ada.MetaData.DeviceId))
	}
	for _, deviceMeta := range metaResp.Devices {
		deviceAPI, err := device.DeviceInit(deviceMeta, ada.AdaptorID)
		if err != nil {
			return err
		}
		ada.DeviceMap[deviceMeta.DeviceId] = deviceAPI
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
			deviceAPI, e := device.DeviceInit(deviceMeta, ada.AdaptorID)
			if e != nil {
				err = e
				continue
			}
			ada.DeviceMap[deviceID] = deviceAPI
		}
	}
	ada.MetaData = meta
	return
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
	return
}
