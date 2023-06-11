package utils

import (
	"github.com/i-Things/things/shared/def"
	m "math"
	"time"
)

const rad = m.Pi / 180

// time conversions

const (
	daySec = 60 * 60 * 24
	j1970  = 2440588.0
	j2000  = 2451545.0
)

func toJulian(t time.Time) float64 {
	return float64(t.Unix())/daySec - 0.5 + j1970
}
func fromJulian(j float64) time.Time {
	return time.Unix(int64((j+0.5-j1970)*daySec), 0)
}
func toDays(t time.Time) float64 {
	return toJulian(t) - j2000
}

// general utilities for celestial body position

const e = rad * 23.4397

func rightAscension(l, b float64) float64 {
	return m.Atan2(m.Sin(l)*m.Cos(e)-m.Tan(b)*m.Sin(e), m.Cos(l))
}
func declination(l, b float64) float64 {
	return m.Asin(m.Sin(b)*m.Cos(e) + m.Cos(b)*m.Sin(e)*m.Sin(l))
}
func azimuth(H, phi, dec float64) float64 {
	return m.Atan2(m.Sin(H), m.Cos(H)*m.Sin(phi)-m.Tan(dec)*m.Cos(phi))
}
func altitude(H, phi, dec float64) float64 {
	return m.Asin(m.Sin(phi)*m.Sin(dec) + m.Cos(phi)*m.Cos(dec)*m.Cos(H))
}
func siderealTime(d, lw float64) float64 {
	return rad*(280.16+360.9856235*d) - lw
}

// general sun calculations

func solarMeanAnomaly(d float64) float64 {
	return rad * (357.5291 + 0.98560028*d)
}
func eclipticLongitude(ma float64) float64 {
	c := rad * (1.9148*m.Sin(ma) + 0.02*m.Sin(2*ma) + 0.0003*m.Sin(3*ma)) // equation of center
	p := rad * 102.9372                                                   // perihelion of the Earth
	return ma + c + p + m.Pi
}
func sunCoords(d float64) (float64, float64) {
	l := eclipticLongitude(solarMeanAnomaly(d))
	return declination(l, 0), rightAscension(l, 0)
}

// returns sun's azimuth and altitude given time and latitude/longitude

func SunPosition(t time.Time, lat, lng float64) (float64, float64) {
	lw := rad * -lng
	phi := rad * lat
	d := toDays(t)
	dec, ra := sunCoords(d)
	h := siderealTime(d, lw) - ra

	return azimuth(h, phi, dec), altitude(h, phi, dec)
}

// calculations for sun times

const j0 = 0.0009

func julianCycle(d, lw float64) float64 {
	return m.Floor(d - j0 - lw/(2.0*m.Pi) + 0.5)
}
func approxTransit(ht, lw, n float64) float64 {
	return j0 + (ht+lw)/(2.0*m.Pi) + n
}
func solarTransitJ(ds, ma, l float64) float64 {
	return j2000 + ds + 0.0053*m.Sin(ma) - 0.0069*m.Sin(2*l)
}
func hourAngle(h, phi, d float64) float64 {
	return m.Acos((m.Sin(h) - m.Sin(phi)*m.Sin(d)) / (m.Cos(phi) * m.Cos(d)))
}

// returns set time for the given sun altitude
func getSetJ(h, lw, phi, dec, n, m, l float64) float64 {
	w := hourAngle(h, phi, dec)
	a := approxTransit(w, lw, n)
	return solarTransitJ(a, m, l)
}

// sun times configuration

type SunAngle struct {
	angle    float64
	riseName string
	setName  string
}

var sunAngles = [...]SunAngle{
	SunAngle{-0.833, "sunrise", "sunset"},
	SunAngle{-0.3, "sunriseEnd", "sunsetStart"},
	SunAngle{-6.0, "dawn", "dusk"},
	SunAngle{-12.0, "nauticalDawn", "nauticalDusk"},
	SunAngle{-18.0, "nightEnd", "night"},
	SunAngle{6.0, "goldenHourEnd", "goldenHour"},
}

// 太阳升起的时候
func SunRiseTime(t time.Time, point def.Point) time.Time {
	return SunTimes(t, point)["sunrise"]
}

// 太阳落下的时候
func SunSetTime(t time.Time, point def.Point) time.Time {
	return SunTimes(t, point)["sunset"]
}

// calculates sun times for a given date and latitude/longitude
func SunTimes(t time.Time, point def.Point) map[string]time.Time {
	earthPos := PositionToEarth(point)
	lw := rad * -earthPos.Longitude
	phi := rad * earthPos.Latitude

	d := toDays(t)
	n := julianCycle(d, lw)
	ds := approxTransit(0, lw, n)

	ma := solarMeanAnomaly(ds)
	l := eclipticLongitude(ma)
	dec := declination(l, 0)

	jNoon := solarTransitJ(ds, ma, l)

	times := map[string]time.Time{
		"solarNoon": fromJulian(jNoon),
		"nadir":     fromJulian(jNoon - 0.5),
	}

	for _, sunAngle := range sunAngles {
		jSet := getSetJ(sunAngle.angle*rad, lw, phi, dec, n, ma, l)

		times[sunAngle.riseName] = fromJulian(jNoon - (jSet - jNoon))
		times[sunAngle.setName] = fromJulian(jSet)
	}

	return times
}

// moon calculations, based on http://aa.quae.nl/en/reken/hemelpositie.html formulas

func moonCoords(d float64) (float64, float64, float64) { // geocentric ecliptic coordinates of the moon
	el := rad * (218.316 + 13.176396*d) // ecliptic longitude
	ma := rad * (134.963 + 13.064993*d) // mean anomaly
	f := rad * (93.272 + 13.229350*d)   // mean distance

	l := rad*6.289*m.Sin(ma) + el // longitude
	b := rad * 5.128 * m.Sin(f)   // latitude

	dist := 385001 - 20905*m.Cos(ma) // distance to the moon in km

	return declination(l, b), rightAscension(l, b), dist
}

func MoonPosition(t time.Time, lat, lng float64) (float64, float64, float64) {
	lw := rad * -lng
	phi := rad * lat
	d := toDays(t)

	dec, ra, dist := moonCoords(d)
	ha := siderealTime(d, lw) - ra
	h := altitude(ha, phi, dec)

	// altitude correction for refraction
	h = h + rad*0.017/m.Tan(h+rad*10.26/(h+rad*5.10))
	return azimuth(ha, phi, dec), h, dist
}

// example:
// azimuth, altitude := SunPosition(time.Now(), 50.5, 30.5)
// times := SunTimes(time.Now(), 50.5, 30.5)
