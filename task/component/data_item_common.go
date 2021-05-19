package component

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/apache/dubbo-go/common/logger"
	"github.com/mk1010/industry_adaptor/nclink"
	"github.com/mk1010/industry_adaptor/task/common"
)

func DataInfoInit(ctx context.Context, adaptorID, deviceID, componentID string, dataInfoMeta *nclink.NCLinkDataInfo) (common.NCLinkDataInfoAPI, error) {
	if dataInfoMeta == nil || dataInfoMeta.DataItem == nil || dataInfoMeta.SampleInfo == nil {
		err := fmt.Errorf("组件元数据为空")
		logger.Error(err)
		return nil, err
	}
	switch dataInfoMeta.DataItem.DataItemType {
	case "nclink_common_data_item":
		{
			dataItem := &NCLinkCommonDataInfo{
				AdaptorID:   adaptorID,
				DeviceID:    deviceID,
				ComponentID: componentID,
				DataInfo:    dataInfoMeta,
				DataChan:    make(chan []byte, 10),
				done:        make(chan struct{}),
			}
			err := dataItem.Start(ctx)
			if err != nil {
				return nil, err
			}
			return dataItem, nil
		}
		// 花里胡哨的还是算了
	/*case "MethodRegisterExample":
	methReg := reflect.ValueOf(NCLinkComponetMethods{})
	callFunc := methReg.MethodByName(dataInfo.DataItem.DataItemType)
	if callFunc.Kind() == reflect.Func {
		msgChan := make(chan *nclink.NCLinkTopicMessage, 3)
		util.GoSafely(func() {
			callFunc.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(dataInfo), reflect.ValueOf(t.done), reflect.ValueOf(msgChan)})
		}, nil)
	} else {
		err = fmt.Errorf("组件数据项%s未找到可执行方法", dataInfo.DataItem.DataItemType)
		logger.Error(err)
		return
	}*/
	default:
		{
			err := fmt.Errorf("未知的数据项类型%s", dataInfoMeta.DataItem.DataItemType)
			logger.Error(err)
			return nil, err
		}
	}
}

var kindDealFuncMap = make(map[nclink.DataKind]func(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error))

func init() {
	kindDealFuncMap[nclink.DataKind_Int8] = getInt8
	kindDealFuncMap[nclink.DataKind_Int16] = getInt16
	kindDealFuncMap[nclink.DataKind_Int32] = getInt32
	kindDealFuncMap[nclink.DataKind_Int64] = getInt64
	kindDealFuncMap[nclink.DataKind_Uint8] = getUint8
	kindDealFuncMap[nclink.DataKind_Uint16] = getUint16
	kindDealFuncMap[nclink.DataKind_Uint32] = getUint32
	kindDealFuncMap[nclink.DataKind_Uint64] = getUint64
	kindDealFuncMap[nclink.DataKind_Float32] = getFloat32
	kindDealFuncMap[nclink.DataKind_Float64] = getFloat64
	kindDealFuncMap[nclink.DataKind_String] = getString
}

func GetFiledValueLittleEnd(kind nclink.DataKind, data *bytes.Buffer) (interface{}, error) {
	if f, ok := kindDealFuncMap[kind]; ok {
		return f(data, binary.LittleEndian)
	}
	err := fmt.Errorf("未找到执行方法 kind:%s", kind.String())
	logger.Error(err)
	return nil, err
}

func GetFiledValueBigEnd(kind nclink.DataKind, data *bytes.Buffer) (interface{}, error) {
	if f, ok := kindDealFuncMap[kind]; ok {
		return f(data, binary.BigEndian)
	}
	err := fmt.Errorf("未找到执行方法 kind:%s", kind.String())
	logger.Error(err)
	return nil, err
}

/*DataKind_Bool       DataKind = 1
DataKind_Int8       DataKind = 2
DataKind_Int16      DataKind = 3
DataKind_Int32      DataKind = 4
DataKind_Int64      DataKind = 5
DataKind_Uint       DataKind = 6
DataKind_Uint8      DataKind = 7
DataKind_Uint16     DataKind = 8
DataKind_Uint32     DataKind = 9
DataKind_Uint64     DataKind = 10
DataKind_Float32    DataKind = 11
DataKind_Float64    DataKind = 12
DataKind_Complex64  DataKind = 13
DataKind_Complex128 DataKind = 14
DataKind_String     DataKind = 15*/
func getInt8(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val int8
	err := binary.Read(data, b, &val)
	return val, err
}

func getInt16(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val int16
	err := binary.Read(data, b, &val)
	return val, err
}

func getInt32(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val int32
	err := binary.Read(data, b, &val)
	return val, err
}

func getInt64(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val int64
	err := binary.Read(data, b, &val)
	return val, err
}

func getUint8(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val uint8
	err := binary.Read(data, b, &val)
	return val, err
}

func getUint16(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val uint16
	err := binary.Read(data, b, &val)
	return val, err
}

func getUint32(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val uint32
	err := binary.Read(data, b, &val)
	return val, err
}

func getUint64(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val uint64
	err := binary.Read(data, b, &val)
	return val, err
}

func getFloat32(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val float32
	err := binary.Read(data, b, &val)
	return val, err
}

func getFloat64(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	var val float64
	err := binary.Read(data, b, &val)
	return val, err
}

func getString(data *bytes.Buffer, b binary.ByteOrder) (interface{}, error) {
	val, err := data.ReadString('\n')
	return val, err
}
