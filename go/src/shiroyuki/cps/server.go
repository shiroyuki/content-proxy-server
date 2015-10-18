package cps

import yotsuba "github.com/shiroyuki/yotsuba-go"
import tori    "github.com/shiroyuki/tori-go"

func StartService(
    address      string,
    cachePath    string,
    metadataPath string,
    debugMode    bool,
) {
    cryptographer := yotsuba.Enigma{}

    memoryCacheDriver := yotsuba.NewInMemoryCacheDriver(&cryptographer, !debugMode)

    cacheDriver    := yotsuba.CacheDriver(memoryCacheDriver)
    fetcherService := NewFetcher(
        &cryptographer,
        cachePath,
        metadataPath,
        true,
    )

    app := Application{
        CacheDriver:    cacheDriver,
        Cryptographer:  cryptographer,
        FetcherService: fetcherService,
    }

    core := tori.NewSimpleCore()

    core.Router.DebugMode              = debugMode
    core.Router.PriorityList.DebugMode = debugMode

    core.Router.OnGet(
        "image-default",
        "/i/<destination>",
        app.HandleImage,
        true,
    )

    core.Router.OnGet(
        "image-default",
        "/i/<destination>/<spec>",
        app.HandleImage,
        true,
    )

    core.Listen(&address)
}
