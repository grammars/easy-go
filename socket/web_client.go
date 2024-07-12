package socket

type WebClient struct {
	Host        string
	Port        int
	Name        string
	Monitor     *Monitor
	PrintDetail bool
}
