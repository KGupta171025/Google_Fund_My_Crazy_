package spatial

import (
	"fmt"
	"math"
	"metro-synapse/backend/internal/models"
	"strings"
)

const (
	base32Alphabet = "0123456789bcdefghjkmnpqrstuvwxyz"
	minLat         = -90.0
	maxLat         = 90.0
	minLon         = -180.0
	maxLon         = 180.0
)

var base32Map = func() map[byte]int {
	m := make(map[byte]int)
	for i := 0; i < len(base32Alphabet); i++ {
		m[base32Alphabet[i]] = i
	}
	return m
}()

// EncodeGeohash encodes latitude and longitude into a geohash string with specified precision
func EncodeGeohash(lat, lon float64, precision int) string {
	if precision <= 0 || precision > 12 {
		precision = 9
	}

	latInterval := [2]float64{minLat, maxLat}
	lonInterval := [2]float64{minLon, maxLon}

	var geohash strings.Builder
	geohash.Grow(precision)

	isEven := true
	bit := 0
	ch := 0

	for geohash.Len() < precision {
		var mid float64
		if isEven {
			mid = (lonInterval[0] + lonInterval[1]) / 2.0
			if lon >= mid {
				ch |= (1 << (4 - bit))
				lonInterval[0] = mid
			} else {
				lonInterval[1] = mid
			}
		} else {
			mid = (latInterval[0] + latInterval[1]) / 2.0
			if lat >= mid {
				ch |= (1 << (4 - bit))
				latInterval[0] = mid
			} else {
				latInterval[1] = mid
			}
		}

		isEven = !isEven
		if bit < 4 {
			bit++
		} else {
			geohash.WriteByte(base32Alphabet[ch])
			bit = 0
			ch = 0
		}
	}

	return geohash.String()
}

// DecodeGeohash decodes a geohash string into latitude and longitude bounds and midpoint
func DecodeGeohash(geohash string) (models.Location, models.BoundingBox, error) {
	geohash = strings.ToLower(geohash)
	latInterval := [2]float64{minLat, maxLat}
	lonInterval := [2]float64{minLon, maxLon}
	isEven := true

	for i := 0; i < len(geohash); i++ {
		c := geohash[i]
		val, ok := base32Map[c]
		if !ok {
			return models.Location{}, models.BoundingBox{}, fmt.Errorf("invalid geohash character: %c", c)
		}

		for mask := 16; mask > 0; mask >>= 1 {
			if isEven {
				mid := (lonInterval[0] + lonInterval[1]) / 2.0
				if (val & mask) != 0 {
					lonInterval[0] = mid
				} else {
					lonInterval[1] = mid
				}
			} else {
				mid := (latInterval[0] + latInterval[1]) / 2.0
				if (val & mask) != 0 {
					latInterval[0] = mid
				} else {
					latInterval[1] = mid
				}
			}
			isEven = !isEven
		}
	}

	bbox := models.BoundingBox{
		MinLat: latInterval[0],
		MaxLat: latInterval[1],
		MinLon: lonInterval[0],
		MaxLon: lonInterval[1],
	}

	midLoc := models.Location{
		Lat: (latInterval[0] + latInterval[1]) / 2.0,
		Lon: (lonInterval[0] + lonInterval[1]) / 2.0,
	}

	return midLoc, bbox, nil
}

// HaversineDistanceMeters calculates exact great-circle distance between two coordinates in meters
func HaversineDistanceMeters(loc1, loc2 models.Location) float64 {
	const earthRadiusMeters = 6371000.0

	dLat := (loc2.Lat - loc1.Lat) * (math.Pi / 180.0)
	dLon := (loc2.Lon - loc1.Lon) * (math.Pi / 180.0)

	lat1Rad := loc1.Lat * (math.Pi / 180.0)
	lat2Rad := loc2.Lat * (math.Pi / 180.0)

	a := math.Sin(dLat/2.0)*math.Sin(dLat/2.0) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2.0)*math.Sin(dLon/2.0)

	c := 2.0 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	return earthRadiusMeters * c
}
