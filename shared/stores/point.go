package stores

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/def"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
)

type Point struct {
	Longitude float64 `json:"longitude,range=[0:180]"` //经度
	Latitude  float64 `json:"latitude,range=[0:90]"`   //纬度
}

func (p Point) ToPo() def.Point {
	return def.Point{
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
	}
}
func ToPoint(p def.Point) Point {
	return Point{
		Longitude: p.Longitude,
		Latitude:  p.Latitude,
	}
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
	case string:
		va := value.(string)

		return p.parsePoint([]byte(va))
	default:
		return fmt.Errorf("failed to scan point: invalid type: %T", value)
	}
	return nil
}

//func (p Point) Value() (driver.Value, error) {
//	return []byte(fmt.Sprintf("ST_GeomFromText('POINT(%f %f)')", p.Longitude, p.Latitude)), nil
//}

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	switch dbType {
	case conf.Pgsql:
		return clause.Expr{
			//SQL:  "ST_PointFromText(?)",
			SQL: "ST_GeomFromText(ST_AsText(?),-1)::point", //如果你不知道 SRID 的值，可以使用 -1 来表示未知的空间参考系统。

			Vars: []interface{}{fmt.Sprintf("POINT(%f %f)", p.Longitude, p.Latitude)},
		}
	default:
		return clause.Expr{
			SQL:  "ST_PointFromText(?)",
			Vars: []interface{}{fmt.Sprintf("POINT(%f %f)", p.Longitude, p.Latitude)},
		}
	}
}
