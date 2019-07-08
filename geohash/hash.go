package geohash

import "bytes"

const (
	BASE32               = "0123456789bcdefghjkmnpqrstuvwxyz"
	MaxLatitude  float64 = 90
	MinLatitude  float64 = -90
	MaxLongitude float64 = 180
	MinLongitude float64 = -180
)

var (
	bits   = []int{16, 8, 4, 2, 1}
	base32 = []byte(BASE32)
)

type Box struct {
	MinLat, MaxLat float64 // 纬度
	MinLng, MaxLng float64 // 经度
}

func (b *Box) Width() float64 {
	return b.MaxLng - b.MinLng
}

func (b *Box) Height() float64 {
	return b.MaxLat - b.MinLat
}

func Encode(latitude, longitude float64, precision int) (string, *Box) {
	var geoHash bytes.Buffer
	var minLat, maxLat = MinLatitude, MaxLatitude
	var minLng, maxLng = MinLongitude, MaxLongitude
	var mid float64 = 0

	bit, ch, length, isEven := 0, 0, 0, true
	for length < precision {
		if isEven {
			if mid = (minLng + maxLng) / 2; mid < longitude {
				ch |= bits[bit]
				minLng = mid
			} else {
				maxLng = mid
			}
		} else {
			if mid = (minLat + maxLat) / 2; mid < latitude {
				ch |= bits[bit]
				minLat = mid
			} else {
				maxLat = mid
			}
		}

		isEven = !isEven
		if bit < 4 {
			bit++
		} else {
			geoHash.WriteByte(base32[ch])
			length, bit, ch = length+1, 0, 0
		}
	}

	b := &Box{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLng: minLng,
		MaxLng: maxLng,
	}

	return geoHash.String(), b
}

// 计算该点（latitude, longitude）在精度precision下的邻居 -- 周围8个区域+本身所在区域
// 返回这些区域的geohash值，总共9个
func GetNeighbors(latitude, longitude float64, precision int) []string {
	geohashs := make([]string, 9)

	// 本身
	geohash, b := Encode(latitude, longitude, precision)
	geohashs[0] = geohash

	// 上下左右
	geohashUp, _ := Encode((b.MinLat+b.MaxLat)/2+b.Height(), (b.MinLng+b.MaxLng)/2, precision)
	geohashDown, _ := Encode((b.MinLat+b.MaxLat)/2-b.Height(), (b.MinLng+b.MaxLng)/2, precision)
	geohashLeft, _ := Encode((b.MinLat+b.MaxLat)/2, (b.MinLng+b.MaxLng)/2-b.Width(), precision)
	geohashRight, _ := Encode((b.MinLat+b.MaxLat)/2, (b.MinLng+b.MaxLng)/2+b.Width(), precision)

	// 四个角
	geohashLeftUp, _ := Encode((b.MinLat+b.MaxLat)/2+b.Height(), (b.MinLng+b.MaxLng)/2-b.Width(), precision)
	geohashLeftDown, _ := Encode((b.MinLat+b.MaxLat)/2-b.Height(), (b.MinLng+b.MaxLng)/2-b.Width(), precision)
	geohashRightUp, _ := Encode((b.MinLat+b.MaxLat)/2+b.Height(), (b.MinLng+b.MaxLng)/2+b.Width(), precision)
	geohashRightDown, _ := Encode((b.MinLat+b.MaxLat)/2-b.Height(), (b.MinLng+b.MaxLng)/2+b.Width(), precision)

	geohashs[1], geohashs[2], geohashs[3], geohashs[4] = geohashUp, geohashDown, geohashLeft, geohashRight
	geohashs[5], geohashs[6], geohashs[7], geohashs[8] = geohashLeftUp, geohashLeftDown, geohashRightUp, geohashRightDown

	return geohashs
}