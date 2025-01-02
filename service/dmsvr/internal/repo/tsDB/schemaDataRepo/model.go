package schemaDataRepo

import "time"

type Event struct {
	ProductID  string    `gorm:"column:product_id;type:varchar(100);NOT NULL"`  // 产品id
	DeviceName string    `gorm:"column:device_name;type:varchar(100);NOT NULL"` // 设备名称
	Identifier string    `gorm:"column:identifier;type:varchar(100);NOT NULL"`  // 事件id
	Type       string    `gorm:"column:type;type:varchar(100);NOT NULL"`        // 事件内容
	Param      string    `gorm:"column:param;type:varchar(256);NOT NULL"`       // 时间戳
	Timestamp  time.Time `gorm:"column:ts;NOT NULL;"`                           // 操作时间
}

func (m *Event) TableName() string {
	return "dm_model_event"
}

type Property struct {
	ProductID  string    `gorm:"column:product_id;type:varchar(100);NOT NULL"`  // 产品id
	DeviceName string    `gorm:"column:device_name;type:varchar(100);NOT NULL"` // 设备名称
	Timestamp  time.Time `gorm:"column:ts;NOT NULL;"`                           // 操作时间
	Identifier string    `gorm:"column:identifier;type:varchar(100);NOT NULL"`  // 事件id
}

type PropertyString struct {
	Property
	Param string `gorm:"column:param;type:varchar(256);NOT NULL"` // 时间戳
}

func (m *PropertyString) TableName() string {
	return "dm_model_property_string"
}

type PropertyStringArray struct {
	Property
	Param string `gorm:"column:param;type:varchar(256);NOT NULL"` // 时间戳
	Pos   int64  `gorm:"column:pos;NOT NULL;"`
}

func (m *PropertyStringArray) TableName() string {
	return "dm_model_property_string_array"
}

type PropertyInt struct {
	Property
	Param int64 `gorm:"column:param;NOT NULL"` // 时间戳
}

func (m *PropertyInt) TableName() string {
	return "dm_model_property_int"
}

type PropertyIntArray struct {
	Property
	Param int64 `gorm:"column:param;NOT NULL"` // 时间戳
	Pos   int64 `gorm:"column:pos;NOT NULL;"`
}

func (m *PropertyIntArray) TableName() string {
	return "dm_model_property_int_array"
}

type PropertyFloat struct {
	Property
	Param float64 `gorm:"column:param;NOT NULL"` // 时间戳
}

func (m *PropertyFloat) TableName() string {
	return "dm_model_property_float"
}

type PropertyFloatArray struct {
	Property
	Param float64 `gorm:"column:param;NOT NULL"` // 时间戳
	Pos   int64   `gorm:"column:pos;NOT NULL;"`
}

func (m *PropertyFloatArray) TableName() string {
	return "dm_model_property_float_array"
}

type PropertyTimestamp struct {
	Property
	Param int64 `gorm:"column:param;NOT NULL"` // 时间戳
}

func (m *PropertyTimestamp) TableName() string {
	return "dm_model_property_timestamp"
}

type PropertyTimestampArray struct {
	Property
	Param int64 `gorm:"column:param;NOT NULL"` // 时间戳
	Pos   int64 `gorm:"column:pos;NOT NULL;"`
}

func (m *PropertyTimestampArray) TableName() string {
	return "dm_model_property_timestamp_array"
}

type PropertyEnum struct {
	Property
	Param int64 `gorm:"column:param;NOT NULL"` // 时间戳
}

func (m *PropertyEnum) TableName() string {
	return "dm_model_property_enum"
}

type PropertyEnumArray struct {
	Property
	Param int64 `gorm:"column:param;NOT NULL"` // 时间戳
	Pos   int64 `gorm:"column:pos;NOT NULL;"`
}

func (m *PropertyEnumArray) TableName() string {
	return "dm_model_property_enum_array"
}

type PropertyBool struct {
	Property
	Param bool `gorm:"column:param;NOT NULL"` // 时间戳
}

func (m *PropertyBool) TableName() string {
	return "dm_model_property_bool"
}

type PropertyBoolArray struct {
	Property
	Param bool  `gorm:"column:param;NOT NULL"` // 时间戳
	Pos   int64 `gorm:"column:pos;NOT NULL;"`
}

func (m *PropertyBoolArray) TableName() string {
	return "dm_model_property_bool_array"
}

type PropertyStruct struct {
	Property
	Param map[string]any `gorm:"column:param;type:json;serializer:json;NOT NULL"` // 时间戳
}

func (m *PropertyStruct) TableName() string {
	return "dm_model_property_struct"
}

type PropertyStructArray struct {
	Property
	Param map[string]any `gorm:"column:param;type:json;serializer:json;NOT NULL"` // 时间戳
	Pos   int64          `gorm:"column:pos;NOT NULL;"`
}

func (m *PropertyStructArray) TableName() string {
	return "dm_model_property_struct_array"
}
