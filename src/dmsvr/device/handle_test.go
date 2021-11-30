package device_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/device"
	"testing"
)

var template string = `
{
  "version": "1.0",
  "properties": [
    {
      "id": "GPS_Info",
      "name": "GPS定位",
      "desc": "",
      "mode": "rw",
      "define": {
        "type": "struct",
        "specs": [
          {
            "id": "longtitude",
            "name": "GPS经度",
            "dataType": {
              "type": "float",
              "min": "-180",
              "max": "180",
              "start": "0",
              "step": "0.001",
              "unit": "度"
            }
          },
          {
            "id": "latitude",
            "name": "GPS纬度",
            "dataType": {
              "type": "float",
              "min": "-90",
              "max": "90",
              "start": "0",
              "step": "0.001",
              "unit": "度"
            }
          }
        ]
      },
      "required": false
    },
    {
      "id": "GPS_ExtInfo",
      "name": "GPS定位扩展",
      "desc": "",
      "mode": "rw",
      "define": {
        "type": "struct",
        "specs": [
          {
            "id": "latitude",
            "name": "纬度",
            "dataType": {
              "type": "float",
              "min": "-90",
              "max": "90",
              "start": "0",
              "step": "0.001",
              "unit": "度"
            }
          },
          {
            "id": "longtitude",
            "name": "经度",
            "dataType": {
              "type": "float",
              "min": "-180",
              "max": "180",
              "start": "0",
              "step": "0.001",
              "unit": "度"
            }
          },
          {
            "id": "altitude",
            "name": "海拔",
            "dataType": {
              "type": "float",
              "min": "-5000",
              "max": "99999",
              "start": "0",
              "step": "0.01",
              "unit": "m"
            }
          },
          {
            "id": "gps_speed",
            "name": "GPS速度",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "1000",
              "start": "0",
              "step": "1",
              "unit": "km/h"
            }
          },
          {
            "id": "direction",
            "name": "方向角",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "360",
              "start": "0",
              "step": "1",
              "unit": "度"
            }
          },
          {
            "id": "location_state",
            "name": "定位状态",
            "dataType": {
              "type": "bool",
              "mapping": {
                "0": "无效",
                "1": "有效"
              }
            }
          },
          {
            "id": "satellites",
            "name": "卫星数",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "9999999999999",
              "start": "0",
              "step": "1",
              "unit": ""
            }
          },
          {
            "id": "gps_time",
            "name": "GPS时间",
            "dataType": {
              "type": "timestamp"
            }
          },
          {
            "id": "collect_time",
            "name": "采集时间",
            "dataType": {
              "type": "timestamp"
            }
          }
        ]
      },
      "required": false
    },
    {
      "id": "Wifi_Info",
      "name": "wifi定位",
      "desc": "",
      "mode": "rw",
      "define": {
        "arrayInfo": {
          "type": "struct",
          "specs": [
            {
              "id": "Mac",
              "name": "mac地址",
              "dataType": {
                "type": "string",
                "min": "0",
                "max": "2048"
              }
            },
            {
              "id": "Rssi",
              "name": "信号强度",
              "dataType": {
                "type": "int",
                "min": "-1000",
                "max": "1000",
                "start": "0",
                "step": "1",
                "unit": ""
              }
            }
          ]
        },
        "type": "array"
      },
      "required": false
    },
    {
      "id": "Cell_Info",
      "name": "蜂窝定位",
      "desc": "LAC代码为基站小区号；cellId为基站 ID；signal为基站信号强度；采集时间为设备采集基站信息时间",
      "mode": "rw",
      "define": {
        "type": "struct",
        "specs": [
          {
            "id": "mcc",
            "name": "国家代码",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "999",
              "start": "460",
              "step": "1",
              "unit": ""
            }
          },
          {
            "id": "mnc",
            "name": "基站网络码",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "9999999",
              "start": "460",
              "step": "1",
              "unit": ""
            }
          },
          {
            "id": "lac",
            "name": "LAC代码",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "9999999",
              "start": "0",
              "step": "1",
              "unit": ""
            }
          },
          {
            "id": "cid",
            "name": "cellId",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "999999999",
              "start": "0",
              "step": "1",
              "unit": ""
            }
          },
          {
            "id": "rss",
            "name": "signal",
            "dataType": {
              "type": "int",
              "min": "-99999",
              "max": "99999",
              "start": "0",
              "step": "1",
              "unit": "dbm"
            }
          },
          {
            "id": "networkType",
            "name": "设备网络制式",
            "dataType": {
              "type": "enum",
              "mapping": {
                "1": "GSM",
                "2": "CDMA",
                "3": "WCDMA",
                "4": "TD_CDMA",
                "5": "LTE"
              }
            }
          },
          {
            "id": "collect_time",
            "name": "采集时间",
            "dataType": {
              "type": "timestamp"
            }
          }
        ]
      },
      "required": false
    },
    {
      "id": "ipaddr",
      "name": "IP地址",
      "desc": "",
      "mode": "r",
      "define": {
        "type": "string",
        "min": "0",
        "max": "64"
      },
      "required": false
    },
    {
      "id": "rssi",
      "name": "信号强度",
      "desc": "",
      "mode": "r",
      "define": {
        "type": "string",
        "min": "0",
        "max": "8"
      },
      "required": false
    },
    {
      "id": "imageUrl",
      "name": "图片地址",
      "desc": "用于传输存储图片地址",
      "mode": "rw",
      "define": {
        "type": "string",
        "min": "0",
        "max": "2048"
      },
      "required": false
    },
    {
      "id": "shuxing",
      "name": "属性",
      "desc": "描述",
      "mode": "rw",
      "define": {
        "type": "string",
        "min": "0",
        "max": "2048"
      },
      "required": false
    },
    {
      "id": "biashijigou",
      "name": "结构体属性",
      "desc": "",
      "mode": "rw",
      "define": {
        "type": "struct",
        "specs": [
          {
            "id": "fwe",
            "name": "dd",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "100",
              "start": "0",
              "step": "1",
              "unit": ""
            }
          },
          {
            "id": "ase",
            "name": "fe",
            "dataType": {
              "type": "int",
              "min": "0",
              "max": "100",
              "start": "0",
              "step": "1",
              "unit": ""
            }
          }
        ]
      },
      "required": false
    },
    {
      "id": "df",
      "name": "dd",
      "desc": "e",
      "mode": "rw",
      "define": {
        "arrayInfo": {
          "type": "int",
          "min": "4",
          "max": "100",
          "start": "4",
          "step": "1",
          "unit": "df"
        },
        "type": "array"
      },
      "required": false
    },
    {
      "id": "serfa",
      "name": "dfefawe",
      "desc": "dfawef",
      "mode": "rw",
      "define": {
        "type": "enum",
        "mapping": {
          "1": "fefeags",
          "4": "segfae"
        }
      },
      "required": false
    },
    {
      "id": "awerawe",
      "name": "dfwef",
      "desc": "",
      "mode": "rw",
      "define": {
        "type": "enum",
        "mapping": {
          "1": "1",
          "4": "测试"
        }
      },
      "required": false
    }
  ],
  "events": [
    {
      "id": "fesf",
      "name": "ddd",
      "desc": "",
      "type": "info",
      "params": [
        {
          "id": "se",
          "name": "dfef",
          "define": {
            "type": "bool",
            "mapping": {
              "0": "关",
              "1": "开"
            }
          }
        },
        {
          "id": "dfa",
          "name": "awefa",
          "define": {
            "type": "int",
            "min": "100",
            "max": "238",
            "start": "100",
            "step": "2",
            "unit": ""
          }
        }
      ],
      "required": false
    },
    {
      "id": "dfawe",
      "name": "fwefa",
      "desc": "",
      "type": "alert",
      "params": [
        {
          "id": "fe",
          "name": "se",
          "define": {
            "type": "bool",
            "mapping": {
              "0": "关",
              "1": "开"
            }
          }
        }
      ],
      "required": false
    },
    {
      "id": "gafa",
      "name": "dfawe",
      "desc": "",
      "type": "fault",
      "params": [
        {
          "id": "sera",
          "name": "fawe",
          "define": {
            "type": "bool",
            "mapping": {
              "0": "关",
              "1": "开"
            }
          }
        }
      ],
      "required": false
    }
  ],
  "actions": [
    {
      "id": "biaoshifu",
      "name": "功能名称",
      "desc": "描述",
      "input": [
        {
          "id": "asdfwe",
          "name": "dd",
          "define": {
            "type": "string",
            "min": "0",
            "max": "2048"
          }
        },
        {
          "id": "ee",
          "name": "ff",
          "define": {
            "type": "int",
            "min": "0",
            "max": "100",
            "start": "1",
            "step": "1",
            "unit": ""
          }
        }
      ],
      "output": [
        {
          "id": "se",
          "name": "fe",
          "define": {
            "type": "string",
            "min": "0",
            "max": "2048"
          }
        }
      ],
      "required": false
    }
  ],
  "profile": {
    "ProductId": "2SNTHBM6O7",
    "CategoryId": "303"
  }
}
`

var propertyParamStr = [...]string{
	`{
    "GPS_Info": {
      "longtitude": 180,
      "latitude": 90
    }
  }`,
	`
{
    "GPS_Info": {
      "longtitude": 180,
      "latitude": 90
    },
    "GPS_ExtInfo": {
      "latitude": 0,
      "longtitude": 0,
      "altitude": 0,
      "gps_speed": 6,
      "direction": 5,
      "location_state": 1,
      "satellites": 3,
      "gps_time": 1624896182,
      "collect_time": 1624377600
    },
    "Cell_Info": {
      "mcc": 460,
      "mnc": 460,
      "lac": 0,
      "cid": 0,
      "rss": 0,
      "networkType": 3,
      "collect_time": 1623772800
    },
    "ipaddr": "awefra",
    "rssi": "fawega",
    "imageUrl": "feagaerga",
    "shuxing": "aghearfawef",
    "biashijigou": {
      "fwe": 4,
      "ase": 4
    },
    "serfa": 4,
    "awerawe": 4,
	"Wifi_Info":[
		{
			"Mac":"1231321dfa",
			"Rssi":123
		},
{
			"Mac":"4524asrgst",
			"Rssi":452
		},
{
			"Mac":"1231321dfafawe",
			"Rssi":214
		}
	]
  }
`,
}

var eventParamStr = [...]string{
	`{
   "method":"event_post",
   "clientToken":"123",
   "version":"1.0",
   "eventId":"fesf",
   "type":"info",
   "timestamp":1212121221,
   "params":{
       "se":true,
       "dfa":120
   }
}`,
	`{
   "method":"event_post",
   "clientToken":"123",
   "version":"1.0",
   "eventId":"gafa",
   "type":"fault",
   "timestamp":1212121221,
   "params":{
       "sera":true
   }
}`,
}

var actionInParamStr = [...]string{
	`{                    
"method": "action",            
"clientToken": "20a4ccfd-d308-****-86c6-5254008a4f10",                
"actionId": "biaoshifu",                
"timestamp": 1212121221,        
"params": {                    
    "asdfwe": "323343",
	"ee":23
    }
}`,
	`{
"method": "action",            
"clientToken": "20a4ccfd-d308-****-86c6-5254008a4f10",                
"actionId": "biaoshifu",                
"timestamp": 1212121221,    
   "params":{
      "asdfwe": "4831",
	"ee":12
   }
}`,
}

var actionOutParamStr = [...]string{
	`{            
"method": "action_reply",        
"actionId": "biaoshifu",       
"clientToken": "20a4ccfd-d308-11e9-86c6-5254008a4f10",        
"code": 0,            
"status": "some message where error",        
"response": {          
    "se":  "afeafeag"           
     }
}`,
	`{            
"method": "action_reply",        
"clientToken": "20a4ccfd-d308-11e9-86c6-5254008a4f10",        
"actionId": "biaoshifu",       
"code": 0,            
"status": "some message where error",        
"response": {          
    "se":  "dfawe"            
     }
}`,
}

func TestVerifyPropertyParam(t *testing.T) {
	fmt.Println("TestVerifyPropertyParam")
	T, err := device.NewTemplate([]byte(template))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range propertyParamStr {
		var dq device.DeviceReq
		err := utils.Unmarshal([]byte(v), &dq.Params)
		if err != nil {
			t.Fatal(err)
		}
		out, err := T.VerifyReqParam(dq, device.PROPERTY)
		if err != nil {
			t.Fatal(err)
		}
		{
			p, _ := json.Marshal(out)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParam=%s\n", str.String())
		}
		{
			val := make(map[string]interface{}, len(out))
			for _, v := range out {
				val[v.ID] = v.ToVal()
			}
			p, _ := json.Marshal(val)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParamTomap=%s\n", str.String())
		}

	}

}

func TestVerifyEventParam(t *testing.T) {
	fmt.Println("TestVerifyEventParam")
	T, err := device.NewTemplate([]byte(template))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range eventParamStr {
		var dq device.DeviceReq
		err := utils.Unmarshal([]byte(v), &dq)
		if err != nil {
			t.Fatal(err)
		}
		out, err := T.VerifyReqParam(dq, device.EVENT)
		if err != nil {
			t.Fatal(err)
		}
		{
			p, _ := json.Marshal(out)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParam=%s\n", str.String())
		}
		{
			val := make(map[string]interface{}, len(out))
			for _, v := range out {
				val[v.ID] = v.ToVal()
			}
			p, _ := json.Marshal(val)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParamTomap=%s\n", str.String())
		}
	}
}

func TestVerifyActionInParam(t *testing.T) {
	fmt.Println("TestVerifyActionInParam")
	T, err := device.NewTemplate([]byte(template))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range actionInParamStr {
		var dq device.DeviceReq
		err := utils.Unmarshal([]byte(v), &dq)
		if err != nil {
			t.Fatal(err)
		}
		out, err := T.VerifyReqParam(dq, device.ACTION_INPUT)
		if err != nil {
			t.Fatal(err)
		}
		{
			p, _ := json.Marshal(out)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParam=%s\n", str.String())
		}
		{
			val := make(map[string]interface{}, len(out))
			for _, v := range out {
				val[v.ID] = v.ToVal()
			}
			p, _ := json.Marshal(val)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParamTomap=%s\n", str.String())
		}
	}
}

func TestVerifyActionOutParam(t *testing.T) {
	fmt.Println("TestVerifyActionOutParam")
	T, err := device.NewTemplate([]byte(template))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range actionOutParamStr {
		var dq device.DeviceResp
		err := utils.Unmarshal([]byte(v), &dq)
		if err != nil {
			t.Fatal(err)
		}
		out, err := T.VerifyRespParam(dq, "biaoshifu", device.ACTION_OUTPUT)
		if err != nil {
			t.Fatal(err)
		}
		{
			p, _ := json.Marshal(out)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParam=%s\n", str.String())
		}
		{
			val := make(map[string]interface{}, len(out))
			for _, v := range out {
				val[v.ID] = v.ToVal()
			}
			p, _ := json.Marshal(val)
			var str bytes.Buffer
			_ = json.Indent(&str, []byte(p), "", "    ")
			fmt.Printf("getParamTomap=%s\n", str.String())
		}
	}
}
