package store

import (
	"context"
	"encoding/binary"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
)

type Point struct {
	Longitude float64 `json:"longitude,range=[0:180]"` //经度
	Latitude  float64 `json:"latitude,range=[0:90]"`   //纬度
}

func (p *Point) parsePoint(binaryData []byte) error {
	if len(binaryData) != 25 {
		return nil
	}
	longitudeBytes := binaryData[len(binaryData)-16 : len(binaryData)-8]
	latitudeBytes := binaryData[len(binaryData)-8:]
	var encode binary.ByteOrder = binary.LittleEndian
	if binaryData[4] != 1 {
		encode = binary.BigEndian
	}
	longitude := math.Float64frombits(encode.Uint64(longitudeBytes))
	latitude := math.Float64frombits(encode.Uint64(latitudeBytes))
	p.Longitude = longitude
	p.Latitude = latitude
	return nil
}
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("failed to scan point: value is nil")
	}
	switch value.(type) {
	case []byte:
		va := value.([]byte)
		return p.parsePoint(va)
	default:
		return fmt.Errorf("failed to scan point: invalid type: %T", value)
	}
	return nil
}

//func (p Point) Value() (driver.Value, error) {
//	return []byte(fmt.Sprintf("ST_GeomFromText('POINT(%f %f)')", p.Longitude, p.Latitude)), nil
//}

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%f %f)", p.Longitude, p.Latitude)},
	}
}
