package main

import "shiroyuki/iip"
//import "log"

func test_drive() {
    enigma  := iip.Enigma{}
    fetcher := iip.NewFetcher(enigma, "cache", "mcache", true)

    fetcher.Fetch("https://farm4.staticflickr.com/3930/15247727947_e3de85030a_k_d.jpg")
}

func start_service() {
    server := iip.NewServer("0.0.0.0:9500", "cache", "mcache")
    server.Listen()
}

func main() {
    start_service()
    //test_drive()
}
