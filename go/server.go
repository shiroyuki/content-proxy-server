package main

import "shiroyuki/cps"

// func test_drive() {
//     enigma  := cps.Enigma{}
//     fetcher := cps.NewFetcher(enigma, "cache", "mcache", true)
//
//     fetcher.Fetch("https://farm4.staticflickr.com/3930/15247727947_e3de85030a_k_d.jpg")
// }

func main() {
    cps.StartService(
        "0.0.0.0:9500",
        "cache",
        "mcache",
        true,
    )
}
