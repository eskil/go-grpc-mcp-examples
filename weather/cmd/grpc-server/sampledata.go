package main

import (
	"errors"
	weatherv1 "weather/v1"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/proto"
)

// sampleDay holds one day of Copenhagen sample weather, transcribed from
// a historical/statistical forecast for May-July 2026. Temperatures are
// degrees Fahrenheit, as shown on the source.
type sampleDay struct {
	Date       *date.Date
	TempF      float32
	NightTempF float32
	Conditions string
}

// Static data scraped off https://world-weather.info/forecast/denmark/copenhagen/may-2026/
// That site has bot protection, so we can't scrape it in real time. Plus for a sample
// app, we don't want to be dependent on a live service.
var copenhagenSamples = []sampleDay{
	// May 2026
	{&date.Date{Year: 2026, Month: 5, Day: 1}, 68, 48, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 2}, 73, 52, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 3}, 63, 57, "cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 4}, 59, 54, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 5}, 55, 46, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 6}, 55, 45, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 7}, 57, 43, "cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 8}, 57, 43, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 9}, 59, 43, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 10}, 57, 46, "cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 11}, 48, 46, "rain"},
	{&date.Date{Year: 2026, Month: 5, Day: 12}, 54, 45, "cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 13}, 46, 46, "rain"},
	{&date.Date{Year: 2026, Month: 5, Day: 14}, 54, 41, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 15}, 55, 45, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 16}, 52, 46, "rain"},
	{&date.Date{Year: 2026, Month: 5, Day: 17}, 55, 45, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 18}, 61, 48, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 19}, 61, 46, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 20}, 59, 48, "cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 21}, 63, 52, "rain"},
	{&date.Date{Year: 2026, Month: 5, Day: 22}, 63, 52, "cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 23}, 73, 55, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 24}, 64, 57, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 25}, 70, 55, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 26}, 70, 59, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 27}, 61, 54, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 28}, 66, 50, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 29}, 73, 50, "sunny"},
	{&date.Date{Year: 2026, Month: 5, Day: 30}, 66, 59, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 5, Day: 31}, 63, 55, "partly cloudy"},

	// June 2026
	{&date.Date{Year: 2026, Month: 6, Day: 1}, 64, 55, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 2}, 68, 57, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 3}, 64, 59, "cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 4}, 63, 57, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 5}, 66, 57, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 6}, 64, 55, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 7}, 64, 57, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 8}, 64, 55, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 9}, 61, 57, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 10}, 63, 52, "sunny"},
	{&date.Date{Year: 2026, Month: 6, Day: 11}, 59, 52, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 12}, 61, 50, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 13}, 57, 54, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 14}, 61, 52, "rain"},
	{&date.Date{Year: 2026, Month: 6, Day: 15}, 61, 55, "sunny"},
	{&date.Date{Year: 2026, Month: 6, Day: 16}, 59, 54, "cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 17}, 63, 52, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 18}, 70, 61, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 19}, 72, 61, "sunny"},
	{&date.Date{Year: 2026, Month: 6, Day: 20}, 77, 64, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 21}, 68, 66, "cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 22}, 68, 57, "cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 23}, 72, 55, "sunny"},
	{&date.Date{Year: 2026, Month: 6, Day: 24}, 79, 61, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 25}, 75, 66, "sunny"},
	{&date.Date{Year: 2026, Month: 6, Day: 26}, 82, 63, "sunny"},
	{&date.Date{Year: 2026, Month: 6, Day: 27}, 90, 66, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 28}, 81, 72, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 29}, 73, 70, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 6, Day: 30}, 72, 59, "partly cloudy"},

	// July 2026
	{&date.Date{Year: 2026, Month: 7, Day: 1}, 72, 61, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 2}, 64, 59, "rain"},
	{&date.Date{Year: 2026, Month: 7, Day: 3}, 61, 57, "cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 4}, 66, 55, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 5}, 64, 59, "sunny"},
	{&date.Date{Year: 2026, Month: 7, Day: 6}, 57, 54, "rain"},
	{&date.Date{Year: 2026, Month: 7, Day: 7}, 59, 57, "rain"},
	{&date.Date{Year: 2026, Month: 7, Day: 8}, 68, 55, "sunny"},
	{&date.Date{Year: 2026, Month: 7, Day: 9}, 72, 57, "sunny"},
	{&date.Date{Year: 2026, Month: 7, Day: 10}, 73, 59, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 11}, 77, 63, "sunny"},
	{&date.Date{Year: 2026, Month: 7, Day: 12}, 81, 64, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 13}, 68, 64, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 14}, 75, 63, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 15}, 77, 61, "sunny"},
	{&date.Date{Year: 2026, Month: 7, Day: 16}, 77, 63, "sunny"},
	{&date.Date{Year: 2026, Month: 7, Day: 17}, 79, 63, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 18}, 70, 64, "cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 19}, 66, 59, "rain"},
	{&date.Date{Year: 2026, Month: 7, Day: 20}, 66, 57, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 21}, 68, 57, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 22}, 63, 55, "rain"},
	{&date.Date{Year: 2026, Month: 7, Day: 23}, 70, 59, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 24}, 68, 61, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 25}, 68, 61, "cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 26}, 70, 61, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 27}, 72, 63, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 28}, 73, 64, "partly cloudy"},
	{&date.Date{Year: 2026, Month: 7, Day: 29}, 66, 66, "rain"},
	{&date.Date{Year: 2026, Month: 7, Day: 30}, 66, 64, "rain"},
	{&date.Date{Year: 2026, Month: 7, Day: 31}, 73, 61, "sunny"},
}

func equalDate(a *date.Date, b *date.Date) bool {
	return a.Year == b.Year && a.Month == b.Month && a.Day == b.Day
}

func GetDataForDateRange(
	location string,
	temperatureUnit weatherv1.TemperatureUnit,
	dateRange *weatherv1.DateRange,
) (
	[]*weatherv1.WeatherInfo,
	error,
) {
	// We only have sample data for Copenhagen, DK. In a real application, this would query a database or call an external API.
	if location != "Copenhagen, DK" {
		return nil, errors.New("unknown location")
	}

	// Naively traverse the sample data and collect the days that fall within the requested date range.
	result := []*weatherv1.WeatherInfo{}
	collecting := false
	for _, day := range copenhagenSamples {
		if equalDate(day.Date, dateRange.Begin) {
			collecting = true
		}
		if collecting {
			result = append(result, &weatherv1.WeatherInfo{
				Date:           day.Date,
				HiTemperature:  proto.Float32(day.TempF),
				LowTemperature: proto.Float32(day.NightTempF),
				Conditions:     day.Conditions,
			})
		}
		if equalDate(day.Date, dateRange.End) {
			break
		}
	}

	// Convert temperatures to Celsius if requested. The sample data is in Fahrenheit.
	if temperatureUnit == weatherv1.TemperatureUnit_CELCIUS {
		for _, r := range result {
			r.HiTemperature = proto.Float32(float32(int32(10.0*((r.GetHiTemperature()-32)*5/9))) / 10.0)
			r.LowTemperature = proto.Float32(float32(int32(10.0*((r.GetLowTemperature()-32)*5/9))) / 10.0)
		}
	}

	return result, nil
}
