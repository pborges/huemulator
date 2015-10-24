package huemulator

type Config struct {
	Hostname string
	Port     int
	UDN      string // guid, doesn't seem to matter what it is so long as it does not change
	Protocol string // http or https
}
