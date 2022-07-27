package schema

import (
	"fmt"
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

func TestTempate(t *testing.T) {
	fmt.Println("TestTempate")
	T, err := NewSchema([]byte(template))
	if err != nil {
		t.Error(err)
	}
	if len(T.Properties) != len(T.Property) {
		t.Fail()
	}
	for i := 0; i < len(T.Properties); i++ {
		if &T.Properties[i] != T.Property[T.Properties[i].ID] {
			t.Fail()
		}
	}
	if len(T.Events) != len(T.Event) {
		t.Fail()
	}
	for i := 0; i < len(T.Events); i++ {
		if &T.Events[i] != T.Event[T.Events[i].ID] {
			t.Fail()
		}
	}
	if len(T.Actions) != len(T.Action) {
		t.Fail()
	}
	for i := 0; i < len(T.Actions); i++ {
		if &T.Actions[i] != T.Action[T.Actions[i].ID] {
			t.Fail()
		}
	}
}
