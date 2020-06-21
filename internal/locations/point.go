package locations

import (
"database/sql/driver"
"encoding/json"
"fmt"
"github.com/paulmach/orb"
"github.com/paulmach/orb/encoding/wkb"
)

//
type Point struct {
	orb.Point
}

// NewPoint creates new Point
func NewPoint(lng, lat float64) Point {
	return Point{Point: orb.Point{lng, lat}}
}

// Value enables serialization to SQL
func (p Point) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%g %g)", p.Lon(), p.Lat()), nil
}

// Scan enables deserialization from SQL
func (p *Point) Scan(src interface{}) error {
	return wkb.
		Scanner(&p.Point).
		Scan(src)
}

func (p Point) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%g, %g]", p.Lat(), p.Lon())), nil
}

func (p *Point) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &p.Point)
	if err != nil {
		return err
	}
	p.Point[0], p.Point[1] = p.Point[1], p.Point[0]
	return nil
}
