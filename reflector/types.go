package reflector

import "net"

// CallPrinterResponse represents a response from the
// Client.CallPrinter method
type CallPrinterResponse struct {
	Call struct {
		ID         string `json:"id"`
		Relay      string `json:"relay"`
		ClientCode string `json:"client_code"`
	} `json:"call"`
}

// RelayAddr resolves the relay's address
func (r *CallPrinterResponse) RelayAddr() (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", r.Call.Relay)
}
