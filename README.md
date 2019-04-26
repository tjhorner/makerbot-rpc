# makerbot-rpc

[![GoDoc](https://godoc.org/github.com/tjhorner/makerbot-rpc?status.svg)](https://godoc.org/github.com/tjhorner/makerbot-rpc)

A Go client library to interact with your MakerBot printer via the network.

Documentation, examples and more at [GoDoc](https://godoc.org/github.com/tjhorner/makerbot-rpc).

**This is currently in beta and does not support many functions that MakerBot printers make available.** ~~Most notably, it does not yet support sending print files.~~ Also, some responses are not yet modelled.

## Features and TODO

- [x] Connecting to printers (`ConnectLocal()`, `ConnectRemote()`)
- [ ] Printer discovery via mDNS
- [x] Authenticating with local printers via Thingiverse (`AuthenticateWithThingiverse()`)
- [ ] Authenticating with local printers via local authentication (pushing the knob)
- [x] Authenticating with remote printers via MakerBot Reflector (`ConnectRemote()`)
- [x] Printer state updates (`HandleStateUpdate()`)
- [x] Load filament method (`LoadFilament()`)
- [x] Unload filament method (`UnloadFilament()`)
- [x] Cancel method (`Cancel()`)
- [x] Change machine name (`ChangeMachineName()`)
- [x] Send print files (`Print()`, `PrintFile()`)
- [x] Camera stream/snapshots (`HandleCameraFrame()`, `GetCameraFrame()`)
- [ ] Get machine config (low priority; isn't very useful)
- [ ] Write tests (will need to make a mock MakerBot RPC server)
- [ ] Write examples
- [ ] Better errors