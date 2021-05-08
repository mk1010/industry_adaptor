package task

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

var NCLinkAdaptoMap = make(map[string]NCLinkAdaptorAPI)

var NCLinkDeviceMap = make(map[string]NCLinkDeviceAPI)

var NCLinkComponentMap = make(map[string]NCLinkComponentAPI)

type NCLinkAdaptorAPI interface {
	Start(ctx context.Context) (err error)
	UpdateMeta(meta *nclink.NCLinkAdaptor) error
	SetID(string)
}

type NCLinkDeviceAPI interface {
	Start(ctx context.Context) (err error)
	UpdateMeta(meta *nclink.NCLinkDevice) error
	SetID(string)
}

type NCLinkComponentAPI interface {
	Start(ctx context.Context) (err error)
	UpdateMeta(meta *nclink.NCLinkComponent) error
	SetID(string)
}
