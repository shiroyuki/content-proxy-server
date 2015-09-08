package main

import "shiroyuki/iip"

func bootstrap() iip.Fetcher {
    contentDriver := iip.GenericCacheDriver{
        StoragePath: "cache",
        Compressed:  true,
    }

    metadataDriver := iip.GenericCacheDriver{
        StoragePath: "mcache",
        Compressed:  false,
    }

    enigma := iip.Enigma{}

    fetcher := iip.Fetcher{
        FileStorage:   contentDriver,
        MetadataRepo:  metadataDriver,
        Cryptographer: enigma,
    }

    return fetcher
}

func test_drive() {
    fetcher := bootstrap()

    fetcher.Fetch("https://farm4.staticflickr.com/3930/15247727947_e3de85030a_k_d.jpg")
}

func start_service() {
    fetcher := bootstrap()

    server := iip.NewServer("0.0.0.0:9500", fetcher)
    server.Listen()
}

func main() {
    //start_service()
    test_drive()
}
