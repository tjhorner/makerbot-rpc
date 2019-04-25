# makerbot-rpc

[![GoDoc](http://godoc.org/github.com/tjhorner/makerbot-rpc?status.svg)](http://godoc.org/github.com/tjhorner/makerbot-rpc)

A Go client library to interact with your MakerBot printer via the network.

Documentation, examples and more at [GoDoc](https://godoc.org/github.com/tjhorner/makerbot-rpc).

**This is currently in beta and does not support many functions that MakerBot printers make available.** Most notably, it does not yet support sending print files. Also, some responses are not yet modelled.

## Features and TODO

- [x] Connecting to printers (`Connect()`)
- [x] Authenticating with local printers via Thingiverse (`AuthenticateWithThingiverse()`)
- [ ] Authenticating with local printers via local authentication (pushing the knob)
- [ ] Authenticating with remote printers via MakerBot Reflector
- [x] Printer state updates (`HandleStateUpdate()`)
- [x] Load filament method (`LoadFilament()`)
- [x] Unload filament method (`UnloadFilament()`)
- [x] Cancel method (`Cancel()`)
- [x] Change machine name (`ChangeMachineName()`)
- [ ] Send print files
- [ ] Camera stream/snapshots
- [ ] Get machine config
- [ ] Write tests (will need to make a mock MakerBot RPC server)
- [ ] Write examples
- [ ] Better errors