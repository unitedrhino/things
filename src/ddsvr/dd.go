package main

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/i-Things/things/src/ddsvr/dddirect"
	"os"
)

func main() {
	test()
	os.Exit(0)
	dddirect.NewDd()
}

var script string = `
var COMMAND_REPORT = 0x00; //属性上报。
var COMMAND_SET = 0x01; //属性设置。
var COMMAND_REPORT_REPLY = 0x02; //上报数据返回结果。
var COMMAND_SET_REPLY = 0x03; //属性设置设备返回结果。
var COMMAD_UNKOWN = 0xff;    //未知的命令。
var ALINK_PROP_REPORT_METHOD = 'thing.event.property.post'; //物联网平台Topic，设备上传属性数据到云端。
var ALINK_PROP_SET_METHOD = 'thing.service.property.set'; //物联网平台Topic，云端下发属性控制指令到设备端。
var ALINK_PROP_SET_REPLY_METHOD = 'thing.service.property.set'; //物联网平台Topic，设备上报属性设置的结果到云端。
var SELF_DEFINE_TOPIC_UPDATE_FLAG = '/user/update'  //自定义Topic：/user/update。
var SELF_DEFINE_TOPIC_ERROR_FLAG = '/user/update/error' //自定义Topic：/user/update/error。
/*
示例数据：
设备上报属性数据：
传入参数：
    0x000000000100320100000000
输出结果：
    {"method":"thing.event.property.post","id":"1","params":{"prop_float":0,"prop_int16":50,"prop_bool":1},"version":"1.0"}

属性设置的返回结果：
传入参数：
    0x0300223344c8
输出结果：
    {"code":"200","data":{},"id":"2241348","version":"1.0"}
*/
function rawDataToProtocol(bytes) {
    var uint8Array = new Uint8Array(bytes.length);
    for (var i = 0; i < bytes.length; i++) {
        uint8Array[i] = bytes[i] & 0xff;
    }
    var dataView = new DataView(uint8Array.buffer, 0);
    var jsonMap = new Object();
    var fHead = uint8Array[0]; // command
    if (fHead == COMMAND_REPORT) {
        jsonMap['method'] = ALINK_PROP_REPORT_METHOD; //ALink JSON格式，属性上报topic。
        jsonMap['version'] = '1.0'; //ALink JSON格式，协议版本号固定字段。
        jsonMap['id'] = '' + dataView.getInt32(1); //ALink JSON格式，标示该次请求id值。
        var params = {};
        params['prop_int16'] = dataView.getInt16(5); //对应产品属性中prop_int16。
        params['prop_bool'] = uint8Array[7]; //对应产品属性中prop_bool。
        params['prop_float'] = dataView.getFloat32(8); //对应产品属性中prop_float。
        jsonMap['params'] = params; //ALink JSON格式，params标准字段。
    } else if(fHead == COMMAND_SET_REPLY) {
        jsonMap['version'] = '1.0'; //ALink JSON格式，协议版本号固定字段。
        jsonMap['id'] = '' + dataView.getInt32(1); //ALink JSON格式，标示该次请求id值。
        jsonMap['code'] = ''+ dataView.getUint8(5);
        jsonMap['data'] = {};
    }

    return jsonMap;
}
/*
示例数据：
云端下发属性设置指令：
传入参数：
    {"method":"thing.service.property.set","id":"12345","version":"1.0","params":{"prop_float":123.452, "prop_int16":333, "prop_bool":1}}
输出结果：
    0x0100003039014d0142f6e76d

设备上报的返回结果：
传入数据：
    {"method":"thing.event.property.post","id":"12345","version":"1.0","code":200,"data":{}}
输出结果：
    0x0200003039c8
*/
function protocolToRawData(json) {
    var method = json['method'];
    var id = json['id'];
    var version = json['version'];
    var payloadArray = [];
    if (method == ALINK_PROP_SET_METHOD) //属性设置。
    {
        var params = json['params'];
        var prop_float = params['prop_float'];
        var prop_int16 = params['prop_int16'];
        var prop_bool = params['prop_bool'];
        //按照自定义协议格式拼接 rawData。
        payloadArray = payloadArray.concat(buffer_uint8(COMMAND_SET)); //command字段。
        payloadArray = payloadArray.concat(buffer_int32(parseInt(id))); //ALink JSON格式 'id'。
        payloadArray = payloadArray.concat(buffer_int16(prop_int16)); //属性'prop_int16'的值。
        payloadArray = payloadArray.concat(buffer_uint8(prop_bool)); //属性'prop_bool'的值。
        payloadArray = payloadArray.concat(buffer_float32(prop_float)); //属性'prop_float'的值。
    } else if (method ==  ALINK_PROP_REPORT_METHOD) { //设备上报数据返回结果。
        var code = json['code'];
        payloadArray = payloadArray.concat(buffer_uint8(COMMAND_REPORT_REPLY)); //command字段。
        payloadArray = payloadArray.concat(buffer_int32(parseInt(id))); //ALink JSON格式'id'。
        payloadArray = payloadArray.concat(buffer_uint8(code));
    } else { //未知命令，对于这些命令不做处理。
        var code = json['code'];
        payloadArray = payloadArray.concat(buffer_uint8(COMMAD_UNKOWN)); //command字段。
        payloadArray = payloadArray.concat(buffer_int32(parseInt(id))); //ALink JSON格式'id'。
        payloadArray = payloadArray.concat(buffer_uint8(code));
    }
    return payloadArray;
}

/*
  示例数据
  自定义Topic：
     /user/update，上报数据。
  输入参数：
     topic:/{productKey}/{deviceName}/user/update
     bytes: 0x000000000100320100000000
  输出参数：
  {
     "prop_float": 0,
     "prop_int16": 50,
     "prop_bool": 1,
     "topic": "/{productKey}/{deviceName}/user/update"
   }
 */
function transformPayload(topic, bytes) {
    var uint8Array = new Uint8Array(bytes.length);
    for (var i = 0; i < bytes.length; i++) {
        uint8Array[i] = bytes[i] & 0xff;
    }
    var dataView = new DataView(uint8Array.buffer, 0);
    var jsonMap = {};

    if(topic.includes(SELF_DEFINE_TOPIC_ERROR_FLAG)) {
        jsonMap['topic'] = topic;
        jsonMap['errorCode'] = dataView.getInt8(0)
    } else if (topic.includes(SELF_DEFINE_TOPIC_UPDATE_FLAG)) {
        jsonMap['topic'] = topic;
        jsonMap['prop_int16'] = dataView.getInt16(5);
        jsonMap['prop_bool'] = uint8Array[7];
        jsonMap['prop_float'] = dataView.getFloat32(8);
    }

    return jsonMap;
}

//以下是部分辅助函数。
function buffer_uint8(value) {
    var uint8Array = new Uint8Array(1);
    var dv = new DataView(uint8Array.buffer, 0);
    dv.setUint8(0, value);
    return [].slice.call(uint8Array);
}
function buffer_int16(value) {
    var uint8Array = new Uint8Array(2);
    var dv = new DataView(uint8Array.buffer, 0);
    dv.setInt16(0, value);
    return [].slice.call(uint8Array);
}
function buffer_int32(value) {
    var uint8Array = new Uint8Array(4);
    var dv = new DataView(uint8Array.buffer, 0);
    dv.setInt32(0, value);
    return [].slice.call(uint8Array);
}
function buffer_float32(value) {
    var uint8Array = new Uint8Array(4);
    var dv = new DataView(uint8Array.buffer, 0);
    dv.setFloat32(0, value);
    return [].slice.call(uint8Array);
}
`

func test() {
	vm := goja.New()
	_, err := vm.RunString(script)
	if err != nil {
		panic(err)
	}
	sum, ok := goja.AssertFunction(vm.Get("transformPayload"))
	if !ok {
		panic("Not a function")
	}
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x32, 0x01, 0x00, 0x00, 0x00, 0x00}
	res, err := sum(goja.Undefined(), vm.ToValue("/{productKey}/{deviceName}/user/update"), vm.ToValue(data))
	if err != nil {
		panic(err)
	}
	str, err := res.ToObject(nil).MarshalJSON()
	fmt.Println(str, err, res)
	// Output: 42
}
