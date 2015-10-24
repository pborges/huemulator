package huemulator


type LightsProvider interface {
	GetLights() ([]Light)
}

type Light struct {
	UUID    string
	Name    string
	OnFunc  func(Light) (ok bool)
	OffFunc func(Light) (ok bool)
}

type light struct {
	State            struct {
						 On        bool `json:"on"`
						 Bri       int `json:"bri"`
						 Hue       int `json:"hue"`
						 Sat       int `json:"sat"`
						 Effect    string `json:"effect"`
						 Ct        int `json:"ct"`
						 Alert     string `json:"alert"`
						 Colormode string `json:"colormode"`
						 Reachable bool `json:"reachable"`
						 XY        []float64 `json:"xy"`
					 } `json:"state"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	ModelId          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	UniqueId         string `json:"uniqueid"`
	SwVersion        string `json:"swversion"`
	PointSymbol      struct {
						 One   string `json:"1"`
						 Two   string `json:"2"`
						 Three string `json:"3"`
						 Four  string `json:"4"`
						 Five  string `json:"5"`
						 Six   string `json:"6"`
						 Seven string `json:"7"`
						 Eight string `json:"8"`
					 } `json:"pointsymbol"`
}

type lightsWrapper struct {
	Lights map[string]light `json:"lights"`
}

func wrapLights(lights []Light) (w lightsWrapper) {
	w.Lights = make(map[string]light)
	for _, v := range lights {
		l := light{
			Name:v.Name,
			UniqueId:v.UUID,
			Type:"Extended color light",
			ModelId:"LCT001",
			SwVersion:"65003148",
			ManufacturerName:"Philips",
		}
		l.State.Reachable = true
		w.Lights[v.UUID] = l
	}
	return
}
