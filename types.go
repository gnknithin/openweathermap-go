package openweathermap

// Coordinates represents the geographical location block.
type Coordinates struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

// WeatherDescription represents the short summary of weather conditions.
type WeatherDescription struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// MainStats represents the core temperature and pressure metrics.
type MainStats struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

// CurrentWeatherResponse represents the complete payload returned by the Current Weather API.
type CurrentWeatherResponse struct {
	Coord   Coordinates          `json:"coord"`
	Weather []WeatherDescription `json:"weather"`
	Main    MainStats            `json:"main"`
	Name    string               `json:"name"`
	Cod     int                  `json:"cod"`
}

// GeocodeLocation represents a geographical location entry returned by the Geocoding API.
type GeocodeLocation struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Country   string  `json:"country"`
	State     string  `json:"state,omitempty"` // omitempty handles states not applicable outside certain countries
}

// PointInTimeWeather represents the common meteorological metrics returned for a specific timestamp.
type PointInTimeWeather struct {
	Time       int64                `json:"dt"`
	Sunrise    int64                `json:"sunrise,omitempty"`
	Sunset     int64                `json:"sunset,omitempty"`
	Temp       float64              `json:"temp"`
	FeelsLike  float64              `json:"feels_like"`
	Pressure   int                  `json:"pressure"`
	Humidity   int                  `json:"humidity"`
	DewPoint   float64              `json:"dew_point"`
	UVI        float64              `json:"uvi"`
	Clouds     int                  `json:"clouds"`
	Visibility int                  `json:"visibility"`
	WindSpeed  float64              `json:"wind_speed"`
	WindDeg    int                  `json:"wind_deg"`
	Weather    []WeatherDescription `json:"weather"` // 🏆 Reusing our existing structure!
}

// OneCallResponse represents the massive payload returned by the One Call API.
type OneCallResponse struct {
	Latitude       float64              `json:"lat"`
	Longitude      float64              `json:"lon"`
	Timezone       string               `json:"timezone"`
	TimezoneOffset int                  `json:"timezone_offset"`
	Current        PointInTimeWeather   `json:"current"`
	Hourly         []PointInTimeWeather `json:"hourly"`
}
