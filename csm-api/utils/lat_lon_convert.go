package utils

import "math"

func LatLonToXY(lat, lon float64) (int, int) {
	const (
		Re    = 6371.00877 // 지구 반경(km)
		grid  = 5.0        // 격자 간격(km)
		slat1 = 30.0       // 표준 위도1(degree)
		slat2 = 60.0       // 표준 위도2(degree)
		olon  = 126.0      // 기준점 경도(degree)
		olat  = 38.0       // 기준점 위도(degree)
		xo    = 43         // 기준점 X좌표(GRID)
		yo    = 136        // 기준점 Y좌표(GRID)
	)

	degToRad := func(deg float64) float64 {
		return deg * math.Pi / 180.0
	}

	re := Re / grid
	radSlat1 := degToRad(slat1)
	radSlat2 := degToRad(slat2)
	radOlon := degToRad(olon)
	radOlat := degToRad(olat)

	sn := math.Tan(math.Pi*0.25+radSlat2*0.5) / math.Tan(math.Pi*0.25+radSlat1*0.5)
	sn = math.Log(math.Cos(radSlat1)/math.Cos(radSlat2)) / math.Log(sn)
	sf := math.Tan(math.Pi*0.25 + radSlat1*0.5)
	sf = math.Pow(sf, sn) * math.Cos(radSlat1) / sn
	ro := math.Tan(math.Pi*0.25 + radOlat*0.5)
	ro = re * sf / math.Pow(ro, sn)

	radLat := degToRad(lat)
	radLon := degToRad(lon)

	ra := math.Tan(math.Pi*0.25 + radLat*0.5)
	ra = re * sf / math.Pow(ra, sn)
	theta := radLon - radOlon
	if theta > math.Pi {
		theta -= 2.0 * math.Pi
	}
	if theta < -math.Pi {
		theta += 2.0 * math.Pi
	}
	theta *= sn

	x := ra*math.Sin(theta) + float64(xo)
	y := ro - ra*math.Cos(theta) + float64(yo)

	return int(math.Floor(x + 0.5)), int(math.Floor(y + 0.5))
}
