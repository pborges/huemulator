package huemulator
import (
	"log"
	"net"
	"strings"
	"html/template"
	"bytes"
)

const (
	upnp_multicast_address = "239.255.255.250:1900"
)

var responseTemplateText =
`HTTP/1.1 200 OK
CACHE-CONTROL: max-age=86400
EXT:
LOCATION: {{.Protocol}}://{{.Hostname}}:{{.Port}}/upnp/setup.xml
OPT: "http://schemas.upnp.org/upnp/1/0/"; ns=01
ST: urn:schemas-upnp-org:device:basic:1
USN: uuid:Socket-1_0-221438K0100073::urn:Belkin:device:**

`

func UpnpResponder(config Config) {
	responseTemplate, err := template.New("").Parse(responseTemplateText)

	log.Println("[UPNP] listening...")
	addr, err := net.ResolveUDPAddr("udp", upnp_multicast_address)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	l.SetReadBuffer(1024)

	for {
		b := make([]byte, 1024)
		n, src, err := l.ReadFromUDP(b)
		if err != nil {
			log.Fatal("[UPNP] ReadFromUDP failed:", err)
		}

		if strings.Contains(string(b[:n]), "MAN: \"ssdp:discover\"") {
			c, err := net.DialUDP("udp", nil, src)
			if err != nil {
				log.Fatal("[UPNP] DialUDP failed:", err)
			}

			log.Println("[UPNP] discovery request from", src)

			// For whatever reason I can't execute the template using c as the reader,
			// you HAVE to put it in a buffer first
			// possible timing issue?
			// don't believe me? try it
			b := &bytes.Buffer{}
			err = responseTemplate.Execute(b, config)
			if err != nil {
				log.Fatal("[UPNP] execute template failed:", err)
			}
			c.Write(b.Bytes())
		}
	}
}