package locations

type Location struct {
	ID       string   `json:"id"`
	GeoPoint GeoPoint `json:"geoPoint"`
	MetaData MetaData `json:"metaData"`
}

type MetaData struct {
	LocationName string `json:"locationName"`
	LocationType string `json:"locationType"`
}

type GeoPoint struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type LocationType string

const (
	City    LocationType = "city"
	Town    LocationType = "Town"
	Station LocationType = "Station"
	Airport LocationType = "Airport"
)

func (l LocationType) String() string {
	return string(l)
}
func toPoint(g *GeoPoint) Point {
	if g != nil {
		return NewPoint(g.Longitude, g.Latitude)
	}
	return Point{}
}

func toGeoPoint(p Point) *GeoPoint {
	if p.Lat() != 0 && p.Lon() != 0 {
		return &GeoPoint{
			Longitude: p.Lon(),
			Latitude:  p.Lat(),
		}
	}
	return nil
}

type LocationStoreModel struct {
	ID           string       `db:"loc_id"`
	Point        Point        `db:"point"`
	LocationName string       `db:"loc_name"`
	LocationType LocationType `db:"loc_type"`
}
