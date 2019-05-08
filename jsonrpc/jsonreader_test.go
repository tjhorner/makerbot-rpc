package jsonrpc_test

import (
	"crypto/rand"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/tjhorner/makerbot-rpc/jsonrpc"
)

func TestJSONReader(t *testing.T) {
	testJson := `[
		{
			"_id": "5cc4918e5ffc07a556e7a467",
			"index": 0,
			"guid": "ccaaa030-db12-4e41-901a-71db5eddd1a5",
			"isActive": false,
			"balance": "$2,952.59",
			"picture": "http://placehold.it/32x32",
			"age": 36,
			"eyeColor": "green",
			"name": {
				"first": "Rush",
				"last": "Best"
			},
			"company": "AQUASURE",
			"email": "rush.best@aquasure.us",
			"phone": "+1 (963) 588-3601",
			"address": "438 Murdock Court, Winchester, New Mexico, 1957",
			"about": "Proident exercitation et Lorem est. Deserunt occaecat culpa aute fugiat. Ea adipisicing culpa veniam est et qui anim ut sit tempor ut laboris est dolore. Cillum nulla incididunt eiusmod et cillum sit incididunt ullamco incididunt ex quis deserunt excepteur. Laboris anim esse duis duis Lorem in elit adipisicing laboris sit cupidatat esse tempor incididunt. Cillum elit laborum voluptate sint commodo exercitation laboris adipisicing non ipsum.",
			"registered": "Monday, August 1, 2016 11:54 PM",
			"latitude": "-2.77938",
			"longitude": "-177.303158",
			"tags": [
				"fugiat",
				"elit",
				"qui",
				"in",
				"officia"
			],
			"range": [
				0,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9
			],
			"friends": [
				{
					"id": 0,
					"name": "Dorothea Maxwell"
				},
				{
					"id": 1,
					"name": "Fry Blair"
				},
				{
					"id": 2,
					"name": "Jessie Ware"
				}
			],
			"greeting": "Hello, Rush! You have 5 unread messages.",
			"favoriteFruit": "banana"
		}
	]`

	var wg sync.WaitGroup

	wg.Add(1)
	reader := jsonrpc.NewJSONReader(func(data []byte) error {
		if string(data) != testJson {
			t.Errorf("JSONReader's json was incorrect, got: %s, want: [long test string]\n", string(data))
		}

		wg.Done()
		return nil
	})

	reader.Write([]byte(testJson))

	wg.Wait()
}

func TestJSONReader_GetRawData(t *testing.T) {
	randBytes := make([]byte, 32) // 32 random bytes
	rand.Read(randBytes)

	reader := jsonrpc.NewJSONReader(func(d []byte) error { return errors.New("") })

	go func() {
		reader.Write(randBytes)
	}()

	result := reader.GetRawData(32)

	if !reflect.DeepEqual(randBytes, result) {
		t.Errorf("JSONReader did not return randBytes\n")
	}
}
