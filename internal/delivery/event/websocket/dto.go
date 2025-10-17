package websocket

type LocationMessage struct {
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type RadiusArea struct {
	Radius    float64 `json:"radius"`
	CenterLat float64 `json:"center_lat"`
	CenterLon float64 `json:"center_lon"`
}

type CrossedRadiusMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}