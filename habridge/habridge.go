package main

import (
	"net/http"
	"fmt"
	"log"
	"github.com/pborges/huemulator"
	"strconv"
)
type SimpleDAO struct{}
func (d *SimpleDAO) GetLights() (lights []huemulator.Light) {
	lights = make([]huemulator.Light, 1)
	lights[0] = huemulator.Light{
		UUID:"d79951b3-78d7-4bdb-8168-d6fcab595160",
		Name:"Tester",
		OnFunc:func(l huemulator.Light) (ok bool) {
			fmt.Println("Turn on", l.Name)
			return true
		},
		OffFunc:func(l huemulator.Light) (ok bool) {
			fmt.Println("Turn off", l.Name)
			return false
		},
	}
	return
}

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[REQ]", r.RemoteAddr, r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

func main() {
	config := huemulator.Config{
		Hostname:"192.168.2.102",
		Port:10000,
		UDN:"f6543a06-800d-48ba-8d8f-bc2949eddc33", // any old guid will do so long as it doesnt change
		Protocol:"http",
	}

	go huemulator.UpnpResponder(config)
	r, _ := huemulator.NewRouter(config, new(SimpleDAO))
	fmt.Println(http.ListenAndServe(config.Hostname + ":" + strconv.Itoa(config.Port), RequestLogger(r)))
}