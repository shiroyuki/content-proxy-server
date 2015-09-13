package cps

import (
    "log"
    "net/http"
)

type Server struct {
    Fetcher    Fetcher
    Internal   http.Server
    Compressed bool
}

func NewServer(
    address      string,
    cachePath    string,
    metadataPath string,
) Server {
    var memory CacheDriver

    memory  = CacheDriver(&InMemoryCacheDriver{})

    enigma  := Enigma{}
    fetcher := NewFetcher(enigma, cachePath, metadataPath, true)
    router  := NewWebCore(memory, enigma, fetcher, true)

    internalServer := http.Server{
        Addr:    address,
        Handler: &router,
    }

    app := Server{
        Fetcher:    fetcher,
        Internal:   internalServer,
        Compressed: true,
    }

    return app
}

func (self *Server) Listen() {
    log.Println("Bind the web service to:", self.Internal.Addr)
    log.Fatal("cps.server.Server.Listen/error:", self.Internal.ListenAndServe())
}
