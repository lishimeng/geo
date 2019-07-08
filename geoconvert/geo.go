package geoconvert

import "math"

func RMCtoWGS84(lon float64, lat float64) (float64, float64) {
	return rmc2wgs(lon), rmc2wgs(lat)
}

//BD09toGCJ02
func BD09toGCJ02(lon, lat float64) (float64, float64) {
	x := lon - 0.0065
	y := lat - 0.006

	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*XPI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*XPI)

	gLon := z * math.Cos(theta)
	gLat := z * math.Sin(theta)

	return gLon, gLat
}

//GCJ02toBD09
func GCJ02toBD09(lon, lat float64) (float64, float64) {
	z := math.Sqrt(lon*lon+lat*lat) + 0.00002*math.Sin(lat*XPI)
	theta := math.Atan2(lat, lon) + 0.000003*math.Cos(lon*XPI)

	bdLon := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006

	return bdLon, bdLat
}

//WGS84toGCJ02
func WGS84toGCJ02(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}

	mgLon, mgLat := delta(lon, lat)

	return mgLon, mgLat
}

//GCJ02toWGS84
func GCJ02toWGS84(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}

	mgLon, mgLat := delta(lon, lat)

	return lon*2 - mgLon, lat*2 - mgLat
}

//BD09toWGS84
func BD09toWGS84(lon, lat float64) (float64, float64) {
	lon, lat = BD09toGCJ02(lon, lat)
	return GCJ02toWGS84(lon, lat)
}

//WGS84toBD09
func WGS84toBD09(lon, lat float64) (float64, float64) {
	lon, lat = WGS84toGCJ02(lon, lat)
	return GCJ02toBD09(lon, lat)
}