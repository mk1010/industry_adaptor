package nclink

import "sync"

var NCLinkAdaptorMeta map[string]*NCLinkAdaptor

var NCLinkDeviceMeta map[string]*NCLinkDevice

var NCLinkComponentMeta map[string]*NCLinkComponent

var NCLinkDataItemMeta map[string]*NCLinkDataItem

var NCLinkSampleInfoMeta map[string]*NCLinkSampleInfo

var t sync.Map
