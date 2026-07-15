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
