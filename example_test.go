package makerbot_test

import (
	"log"
	"os"
	"time"

	"github.com/tjhorner/makerbot-rpc"
)

func ExampleThisIsJustForTestingSorry() {
	pip, ok := os.LookupEnv("PRINTER_IP")
	if !ok {
		log.Fatalln("Please provide a printer IP in the form of a PRINTER_IP environment variable")
	}

	tvTok, ok := os.LookupEnv("THINGIVERSE_TOKEN")
	if !ok {
		log.Fatalln("Please provide a Thingiverse token in the form of a THINGIVERSE_TOKEN environment variable")
	}

	tvUsr, ok := os.LookupEnv("THINGIVERSE_USERNAME")
	if !ok {
		log.Fatalln("Please provide a Thingiverse username in the form of a THINGIVERSE_USERNAME environment variable")
	}

	client := makerbot.NewClient(pip)
	defer client.Close()

	err := client.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Connected to printer: %s\n", client.Printer.MachineName)

	client.HandleStateChange(func(old, new *makerbot.PrinterMetadata) {
		if old != nil {
			log.Printf("\nold: %+v\nnew: %+v\n", old.CurrentProcess, new.CurrentProcess)
		}
	})

	log.Println("Authenticating with printer...")

	err = client.AuthenticateWithThingiverse(tvTok, tvUsr)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting filament load process...")

	process, err := client.UnloadFilament(0)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v\n", string(*process))

	time.Sleep(time.Second * 5)

	cancellation, err := client.Cancel()
	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v\n", string(*cancellation))

	// Output:
}
