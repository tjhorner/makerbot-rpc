package reflector

import "net"

type CallPrinterResponse struct {
	Call struct {
		ID         string `json:"id"`
		Relay      string `json:"relay"`
		ClientCode string `json:"client_code"`
	} `json:"call"`
}

func (r *CallPrinterResponse) RelayAddr() (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", r.Call.Relay)
}
