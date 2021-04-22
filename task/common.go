package task

import "github.com/mk1010/industry_adaptor/nclink"

// 存放公共变量

var NCLinkAdaptorMeta = make(map[string]*nclink.NCLinkAdaptor)

var NCLinkDeviceMeta = make(map[string]*nclink.NCLinkDevice)

var NCLinkComponentMeta = make(map[string]*nclink.NCLinkComponent)

var NCLinkDataItemMeta = make(map[string]*nclink.NCLinkDataItem)

var NCLinkSampleInfoMeta = make(map[string]*nclink.NCLinkSampleInfo)
