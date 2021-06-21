package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	DeviceReqTestStr = [...]string{
		`{                     
    "method": "report_info",            
    "clientToken": "123",                    
    "params": {            
        "imei":"358882046126540",     
        "mac":"548998153a5c",          
        "device_label": {                               
            "append_info": "other information"                
        }                
    }
}`,
		`{
  "method": "unbind_device",    
  "clientToken": "20a4ccfd-****-11e9-86c6-5254008a4f10",        
  "timestamp": 1212121221                    
}`,
		`{                    
   "method": "action",            
   "clientToken": "20a4ccfd-d308-****-86c6-5254008a4f10",                
   "actionId": "openDoor",                
   "timestamp": 1212121221,        
   "params": {                    
       "userid": "323343"            
   }
}`,
		`{
   "method":"event_post",
   "clientToken":"123",
   "version":"1.0",
   "eventId":"PowerAlarm",
   "type":"fault",
   "timestamp":1212121221,
   "params":{
       "Voltage":2.8,
       "Percent":20
   }
}`,
		`{
   "method": "get_status",
   "clientToken": "123",
   "type" : "report", 
   "showmeta": 0
}`,
		`{
   "method": "control",
   "clientToken": "123",    
   "params": {
       "power_switch": 1,
       "color": 1,
       "brightness": 66    
   }
}`,
		`{
   "method":"report",
   "clientToken":"123",
   "timestamp":1212121221,
   "params":{
       "power_switch":1,
       "color":1,
       "brightness":32
   }
}`,
	}
	DeviceRespTestStr = [...]string{
		`{        
   "method":"report_info_reply",    
   "clientToken":"123",    
   "code":0,    
   "status":"some message where error"
}`,
		`{            
   "method": "action_reply",        
   "clientToken": "20a4ccfd-d308-11e9-86c6-5254008a4f10",        
   "code": 0,            
   "status": "some message where error",        
   "response": {          
       "Code":  0            
   }
}`,
		`{
   "method": "event_reply",
   "clientToken": "123",
   "version": "1.0",
   "code": 0,
   "status": "some message where error",
   "data": {}
}`,
		`{
   "method": "get_status_reply",
   "code": 0,
   "clientToken": "123",
   "type": "report",
   "data": {
   "report": {
        "power_switch": 1,
        "color": 1,
        "brightness": 66    
   }
   }
}`,
		`{
   "method":"control_reply",
   "clientToken":"123",
   "code":0,
   "status":"some message where error"
}`,
		`{
   "method":"report_reply",
   "clientToken":"123",
   "code":0,
   "status":"some message where error"
}`,
	}
)

func TestDeviceReq(t *testing.T) {
	fmt.Println("TestDeviceReq")
	for _, dqStr := range DeviceReqTestStr {
		dq := DeviceReq{}
		json.Unmarshal([]byte(dqStr), &dq)
		fmt.Printf("src=%s|ans=%+v\n", dqStr, dq)
	}
}
func TestDeviceResp(t *testing.T) {
	fmt.Println("TestDeviceResp")
	for _, dqStr := range DeviceRespTestStr {
		dq := DeviceResp{}
		json.Unmarshal([]byte(dqStr), &dq)
		fmt.Printf("src=%s|ans=%+v\n", dqStr, dq)
	}
}
