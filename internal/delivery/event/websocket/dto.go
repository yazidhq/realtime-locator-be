package websocket

type LocationMessage struct {
	Type      string  `json:"type"`
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CrossedRadiusMessage struct {
	Type    string `json:"type"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type RadiusArea struct {
	Radius    float64 `json:"radius"`
	CenterLat float64 `json:"center_lat"`
	CenterLon float64 `json:"center_lon"`
}