package huemulator
import (
	"net/http"
	"text/template"
	"log"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"strconv"
)

var setupTemplateText =
`<?xml version="1.0"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
        <specVersion>
                <major>1</major>
                <minor>0</minor>
        </specVersion>
        <URLBase>{{.Protocol}}://{{.Hostname}}:{{.Port}}/</URLBase>
        <device>
			<deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>
			<friendlyName>Amazon-Echo-HA-Bridge ({{.Hostname}})</friendlyName>
			<manufacturer>Royal Philips Electronics</manufacturer>
			<modelName>Philips hue bridge 2012</modelName>
			<modelNumber>929000226503</modelNumber>
			<UDN>uuid:{{.UDN}}</UDN>
        </device>
</root>`

type Router struct {
	*httprouter.Router
	config        Config
	setupTemplate *template.Template
	lightsStatus  lightsWrapper
	lightLookup   map[string]Light
}

func NewRouter(config Config, dao LightsProvider) (m *Router, err error) {
	if err != nil {
		log.Fatalln("[WEB] executing template", err)
	}

	m = new(Router)
	m.setupTemplate, err = template.New("").Parse(setupTemplateText)
	m.config = config

	lights := dao.GetLights()

	m.lightLookup = make(map[string]Light)
	for _, l := range lights {
		m.lightLookup[l.UUID] = l
	}

	m.lightsStatus = wrapLights(lights)

	m.Router = httprouter.New()
	m.GET("/upnp/setup.xml", m.upnpSetup)

	m.GET("/api/:userId", m.lights)
	m.PUT("/api/:userId/lights", m.lightsList)
	m.PUT("/api/:userId/lights/:lightId/state", m.lightState)
	m.GET("/api/:userId/lights/:lightId", m.lightInfo)
	return
}

func (m *Router)upnpSetup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/xml")

	err := m.setupTemplate.Execute(w, m.config)
	if err != nil {
		log.Fatalln("[WEB] execute", err)
	}
}

func (m *Router)lights(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(m.lightsStatus)
	if err != nil {
		log.Fatalln("[WEB] Error encoding json", err)
	}
}

func (m *Router)lightInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(m.lightsStatus.Lights[p.ByName("lightId")])
	if err != nil {
		log.Fatalln("[WEB] Error encoding json", err)
	}
}

func (m *Router)lightsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	lightList := make(map[string]string)

	for _, l := range m.lightsStatus.Lights {
		lightList[l.UniqueId] = l.Name
	}

	err := json.NewEncoder(w).Encode(lightList)
	if err != nil {
		log.Fatalln("[WEB] Error encoding json", err)
	}
}

// this whole method needs to be reworked I think, it is not very good
func (m *Router)lightState(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	req := make(map[string]bool)
	json.NewDecoder(r.Body).Decode(&req)
	lightStatus := m.lightsStatus.Lights[p.ByName("lightId")]
	light := m.lightLookup[p.ByName("lightId")]

	state := false
	lightStatus.State.On = false
	lightStatus.State.XY = nil // this seems to be voodoo, if it is nil the echo says it could not turn on/off the device, useful...
	m.lightsStatus.Lights[lightStatus.UniqueId] = lightStatus

	if req["on"] {
		if state = light.OnFunc(light); state {
			lightStatus.State.On = true
			lightStatus.State.XY = []float64{0.4255, 0.3998}
			state = true
		}
	}else {
		if state = light.OffFunc(light); state {
			lightStatus.State.XY = []float64{0.4255, 0.3998}
		}
	}
	m.lightsStatus.Lights[lightStatus.UniqueId] = lightStatus

	// this is very ugly...
	w.Write([]byte("[{\"success\":{\"/lights/" + lightStatus.UniqueId + "/state/on\":" + strconv.FormatBool(state) + "}}]"))
}