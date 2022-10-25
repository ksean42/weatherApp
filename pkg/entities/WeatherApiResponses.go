package entities

//Forecast - response with forecast from weather api
type Forecast struct {
	CityID  int    `json:"city_id"`
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []List `json:"list"`
	City    struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

//List part of forecast weather api response
type List struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
		Humidity  int     `json:"humidity"`
		TempKf    float64 `json:"temp_kf"`
	} `json:"main"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Visibility int     `json:"visibility"`
	Pop        float64 `json:"pop"`
	Sys        struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	DtTxt string `json:"dt_txt"`
	Rain  struct {
		ThreeH float64 `json:"3h"`
	} `json:"rain,omitempty"`
}

//City - response with city list from weather api
type City struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LocalNames struct {
		Ar          string `json:"ar"`
		ASCII       string `json:"ascii"`
		Bg          string `json:"bg"`
		Ca          string `json:"ca"`
		De          string `json:"de"`
		El          string `json:"el"`
		En          string `json:"en"`
		Fa          string `json:"fa"`
		FeatureName string `json:"feature_name"`
		Fi          string `json:"fi"`
		Fr          string `json:"fr"`
		Gl          string `json:"gl"`
		He          string `json:"he"`
		Hi          string `json:"hi"`
		ID          string `json:"id"`
		It          string `json:"it"`
		Ja          string `json:"ja"`
		La          string `json:"la"`
		Lt          string `json:"lt"`
		Pt          string `json:"pt"`
		Ru          string `json:"ru"`
		Sr          string `json:"sr"`
		Th          string `json:"th"`
		Tr          string `json:"tr"`
		Vi          string `json:"vi"`
		Zu          string `json:"zu"`
	} `json:"local_names"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}
