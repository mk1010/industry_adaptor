package common

import (
	"context"
	"sync"

	"github.com/mk1010/industry_adaptor/nclink"
)

// 存放公共变量

var NCLinkAdaptorMeta = sync.Map{} // make(map[string]*nclink.NCLinkAdaptor)

var NCLinkDeviceMeta = sync.Map{} // make(map[string]*nclink.NCLinkDevice)

var NCLinkComponentMeta = sync.Map{} // make(map[string]*nclink.NCLinkComponent)

var NCLinkDataItemMeta = sync.Map{} // make(map[string]*nclink.NCLinkDataItem)

var NCLinkSampleInfoMeta = sync.Map{} // make(map[string]*nclink.NCLinkSampleInfo)

var NCLinkAdaptorMap = make(map[string]NCLinkAdaptorAPI)

var NCLinkDeviceMap = make(map[string]NCLinkDeviceAPI)

var NCLinkComponentMap = make(map[string]NCLinkComponentAPI)

type NCLinkAdaptorAPI interface {
	Start(ctx context.Context) (err error)
	UpdateMeta(ctx context.Context, meta *nclink.NCLinkAdaptor) error
	Shutdown() error
}

type NCLinkDeviceAPI interface {
	Start(ctx context.Context, id string, config interface{}) (err error)
	UpdateMeta(ctx context.Context, meta *nclink.NCLinkDevice) error
	Shutdown() error
}

type NCLinkComponentAPI interface {
	Start(ctx context.Context, id string, config interface{}) (err error)
	UpdateMeta(ctx context.Context, meta *nclink.NCLinkComponent) error
	Shutdown() error
}
