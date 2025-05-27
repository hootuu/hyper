package amap

const (
	China = "100000"
)

type DistrictResponse struct {
	Status    string      `json:"status"`
	Info      string      `json:"info"`
	InfoCode  string      `json:"infocode"`
	Districts []*District `json:"districts"`
}

type District struct {
	Adcode    string      `json:"adcode"`
	Name      string      `json:"name"`
	Level     string      `json:"level"`
	Districts []*District `json:"districts"`
}
