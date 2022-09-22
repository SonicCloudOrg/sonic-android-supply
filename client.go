package go_android_supply

type adbClient struct {
	host    string
	port    int
	bin     string
	timeout int
}

func (c *adbClient) CreateTcpUsbBridge() {

}
