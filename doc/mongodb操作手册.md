# MongoDB的表结构及操作语句

## 简介
在我们这个项目中mongodb用来存储用户设备产生的属性,操作及时间记录.
之所以采用mongodb则是因为设备的属性是多变,并且可以属性中包含属性的,而且属性的量会很大,这时候如果还是采用mysql的话就没办法进行拓展.

## 设备数据表结构设计
由于设备的数据包含属性,操作及事件,如果分开存储那么表会变得很大,并且难以维护,所以采用了以下的格式进行存储:
```json
{
	"property" : {
		"GPS_Info" : {
			"longtitude" : 3,
			"latitude" : 4
		},
		"biashijigou" : {
			"fwe" : 44,
			"ase" : 32
		}
	},
	"event" : {
		"dfa" : 123,
		"se" : 1
	},
    "action": {
        "dd": {
            "input": {
              "value": 123
            } ,
            "output": {
                "data": "val"
             } 
        }     
    },
	"timestamp" : 123123
}
```
其中property是设备属性,event是设备上报的事件,action是用户操作设备的记录,timestamp则是上报的时间,每次记录的只会有其中一个,并且记录的属性可能会有多个,需要查询的时候过滤了.

## 过滤语句
以属性为例过滤语句如下:
```mongojs
db.test4.find({"property.GPS_Info":{$ne:null }},{"property":1,"timestamp":1})
```
