# makerbot-rpc

[![GoDoc](https://godoc.org/github.com/tjhorner/makerbot-rpc?status.svg)](https://godoc.org/github.com/tjhorner/makerbot-rpc) [![Build Status](https://travis-ci.org/tjhorner/makerbot-rpc.svg?branch=master)](https://travis-ci.org/tjhorner/makerbot-rpc)

![](https://repository-images.githubusercontent.com/183076033/f4f05c80-676d-11e9-98cc-b9f18c43b91f)

_[Go Gopher](https://www.thingiverse.com/thing:439511) by lclemente_

A Go client library to interact with your MakerBot printer via the network.

Documentation, examples and more at [GoDoc](https://godoc.org/github.com/tjhorner/makerbot-rpc).

Since this is currently mid-development, things will probably change _very, very often._ No stable API is guaranteed until the first stable version of this project.

## Example

```shell
go get github.com/tjhorner/makerbot-rpc
```

```golang
// WARNING: This example may fail to work at any time.
// This library is still in development. This example
// is only provided to give a sense of what the library can do.
// Errors are ignored for brevity.

client := makerbot.Client{}
defer client.Close()

// React when the printer's state changes
client.HandleStateChange(func(old, new *makerbot.PrinterMetadata) {
  if new.CurrentProcess == nil {
    return
  }

  // Log the thing the printer is currently doing
  log.Printf("Current process: %s, %v%% done\n", new.CurrentProcess.Name, new.CurrentProcess.Progress)
})

// Make initial TCP connection w/ printer
client.ConnectLocal("192.168.1.2", "9999") // most MakerBot printers listen on 9999

log.Printf("Connected to MakerBot printer: %s\n", client.Printer.MachineName)

// Authenticate with Thingiverse
client.AuthenticateWithThingiverse("my_access_token", "my_username")

log.Println("Queuing file for printing...")

// Print a file named `print.makerbot` in the same directory
client.PrintFileVerify("print.makerbot")

log.Println("Done! Bye bye.")
```

## Features and TODO

- [x] Connecting to printers (`ConnectLocal()`, `ConnectRemote()`)
- [x] Printer discovery via mDNS (`DiscoverPrinters()`)
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
- [x] Parse `.makerbot` print files along with their metadata, thumbnails, and toolpath (see `printfile` package)
- [ ] Get machine config (low priority; isn't very useful)
- [ ] Write tests
  - [ ] `makerbot` package (will need to make a mock MakerBot RPC server)
  - [ ] `jsonrpc` package
  - [x] `printfile` package
  - [ ] `reflector` package
- [ ] Write examples
- [ ] Better errors
- [ ] Fuzz the shizz out of thizz

## License

TBD, but probably MIT later.