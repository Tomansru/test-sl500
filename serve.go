package main

import (
	"git.tomans.ru/Tomansru/sl500-api"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	reader, err := sl500_api.NewConnection("COM5", sl500_api.Baud.Baud19200, false)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = reader.RfLight(sl500_api.ColorOff)

	_, err = reader.RfAntennaSta(sl500_api.AntennaOn)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = reader.RfLight(sl500_api.ColorRed)

	ch := time.Tick(100 * time.Millisecond)

	var resp, cardId, cardCapacity, blockData []byte

	for range ch {
		// Get card type and detect card
		resp, _ = reader.RfRequest(sl500_api.RequestAll)
		if len(resp) != 2 {
			_, _ = reader.RfLight(sl500_api.ColorRed)
			log.Printf("[ERROR 1]: Wrong length for rf request %v\n", resp)
			continue
		}

		// Anti collision check and get card capacity
		cardId, _ = reader.RfAnticoll()
		cardCapacity, _ = reader.RfSelect(cardId)
		if len(cardCapacity) == 0 {
			log.Printf("[ERROR 2]: Capacity length is wrong %v\n", cardCapacity)
			continue
		}

		// Set Authentication
		_, err = reader.RfM1Authentication2(sl500_api.AuthModeKeyA, 28, []byte{0x6a, 0xa9, 0xd7, 0x71, 0x3e, 0xb3})
		if err != nil {
			log.Println("[ERROR 3]:", err)
			continue
		}

		// Get Data
		blockData, err = reader.RfM1Read(28)
		if err != nil {
			log.Println("[ERROR 4]:", err)
			continue
		}
		if len(blockData) != 16 {
			log.Printf("[ERROR 5]: Wrong block data length %v\n", blockData)
			continue
		}

		_, _ = reader.RfLight(sl500_api.ColorGreen)
		//_, _ = reader.RfBeep(30)
		log.Printf("% x\n", blockData)
	}
}