package common

import (
	"context"
	"sync"

	"github.com/mk1010/industry_adaptor/nclink"
)

// 存放公共变量

// 存在内存泄露 不使用  原功能通过反射实现
// var NCLinkAdaptorMeta = sync.Map{} // make(map[string]*nclink.NCLinkAdaptor)

// var NCLinkDeviceMeta = sync.Map{} // make(map[string]*nclink.NCLinkDevice)

// var NCLinkComponentMeta = sync.Map{} // make(map[string]*nclink.NCLinkComponent)

// var NCLinkDataItemMeta = make(map[string]*nclink.NCLinkDataItem)

// var NCLinkSampleInfoMeta = make(map[string]*nclink.NCLinkSampleInfo)

var NCLinkAdaptorMap = sync.Map{} // make(map[string]NCLinkAdaptorAPI)

var NCLinkDeviceMap = sync.Map{} // make(map[string]NCLinkDeviceAPI)

var NCLinkComponentMap = sync.Map{} // make(map[string]NCLinkComponentAPI)

var NClinkInstanceMap = sync.Map{} // make(map[string]NCLinkInstanceAPI)

type NCLinkAdaptorAPI interface {
	Start(ctx context.Context) (err error)
	UpdateMeta(ctx context.Context, meta *nclink.NCLinkAdaptor) error
	GetMeta() *nclink.NCLinkAdaptor
	Shutdown() error
}

type NCLinkDeviceAPI interface {
	Start(ctx context.Context) (err error)
	UpdateMeta(ctx context.Context, meta *nclink.NCLinkDevice) error
	GetMeta() *nclink.NCLinkDevice
	Shutdown() error
}

type NCLinkComponentAPI interface {
	Start(ctx context.Context) (err error)
	UpdateMeta(ctx context.Context, meta *nclink.NCLinkComponent) error
	GetMeta() *nclink.NCLinkComponent
	GetDataInfoApi(dataInfoID string) NCLinkDataInfoAPI
	Shutdown() error
}

type NCLinkDataInfoAPI interface {
	Start(ctx context.Context) (err error)
	SendData(data []byte) error
	UpdateMeta(ctx context.Context, meta *nclink.NCLinkDataInfo) error
	Shutdown() error
}

type NCLinkInstanceAPI interface {
	SendData(msg *nclink.NCLinkTopicMessage) error
	RecvRegister(deviceID, componentID, dataItemoID string, dataAPi NCLinkDataInfoAPI) error
	RecvUnRegister(deviceID, componentID, dataItemoID string, dataAPi NCLinkDataInfoAPI) error
}
