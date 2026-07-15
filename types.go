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

// ForecastItem represents a single 3-hour weather prediction block in the forecast timeline.
type ForecastItem struct {
	Time       int64                `json:"dt"`
	Main       MainStats            `json:"main"`    // 🏆 Reused!
	Weather    []WeatherDescription `json:"weather"` // 🏆 Reused!
	Visibility int                  `json:"visibility"`
	Pop        float64              `json:"pop"` // Probability of precipitation
	TimeText   string               `json:"dt_txt"`
}

// ForecastCity represents the metadata about the city being forecasted.
type ForecastCity struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Coord      Coordinates `json:"coord"` // 🏆 Reused!
	Country    string      `json:"country"`
	Population int         `json:"population"`
	Timezone   int         `json:"timezone"`
	Sunrise    int64       `json:"sunrise"`
	Sunset     int64       `json:"sunset"`
}

// ForecastResponse represents the complete response from the 5-Day/3-Hour Forecast API.
type ForecastResponse struct {
	Cod     string         `json:"cod"`
	Message int            `json:"message"`
	Cnt     int            `json:"cnt"` // Number of timestamps returned (usually 40)
	List    []ForecastItem `json:"list"`
	City    ForecastCity   `json:"city"`
}
