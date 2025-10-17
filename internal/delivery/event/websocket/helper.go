package websocket

import (
	"math"
)

// Fungsi untuk menghitung jarak antara dua koordinat (Haversine formula)
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // jari-jari bumi dalam meter

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c // hasil dalam meter
}
