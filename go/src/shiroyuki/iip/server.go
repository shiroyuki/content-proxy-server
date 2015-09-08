package iip

import (
    "log"
    "net/http"
)

type Handler struct {
    Fetcher Fetcher
}

func (self *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    content, _ := self.Fetcher.Fetch("https://farm4.staticflickr.com/3930/15247727947_e3de85030a_k_d.jpg")

    w.Write(content)
}

type Server struct {
    Fetcher  Fetcher
    Internal http.Server
}

func NewServer(address string, fetcher Fetcher) Server {
    internalServer := http.Server{
        Addr:    address,
        Handler: &Handler{
            Fetcher: fetcher,
        },
    }

    app := Server{
        Fetcher:  fetcher,
        Internal: internalServer,
    }

    return app
}

func (self *Server) Listen() {
    log.Fatal(self.Internal.ListenAndServe())
}
