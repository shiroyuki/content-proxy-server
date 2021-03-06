package cps

import (
    "fmt"
    "log"
    "net/http"
    "regexp"
    yotsuba "github.com/shiroyuki/yotsuba-go"
)

type WebCore struct { // implements http.Handler
    Cache      *yotsuba.CacheDriver
    Enigma     *yotsuba.Enigma
    Fetcher    *Fetcher
    Compressed bool
}

func NewWebCore(
    cache      *yotsuba.CacheDriver,
    enigma     *yotsuba.Enigma,
    fetcher    *Fetcher,
    compressed bool,
) *WebCore {
    return &WebCore{
        Cache:      cache,
        Enigma:     enigma,
        Fetcher:    fetcher,
        Compressed: compressed,
    }
}

func (self *WebCore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var metadata Metadata
    var kind     string
    var content  []byte
    var routePattern = regexp.MustCompile("^/p/(?P<url>[^/]+)")
    var requestPath  = r.URL.Path

    if !routePattern.MatchString(requestPath) {
        w.WriteHeader(404)
        log.Println("cps.WebCore.ServeHTTP: Responded HTTP 404")
        return // TODO HTTP 400
    }

    log.Println("Processing:", requestPath)

    rawActualUrl := URESearch(routePattern, requestPath)["url"]
    actualUrl    := self.Enigma.B64decode(rawActualUrl)

    log.Println("Interpreted:", actualUrl)

    commonCacheKey := actualUrl

    contentTypeCacheKey := "ct:"   + commonCacheKey
    contentDataCacheKey := "data:" + commonCacheKey

    kind    = string((*self.Cache).Load(contentTypeCacheKey))
    content = (*self.Cache).Load(contentDataCacheKey)

    if &kind != nil && content != nil {
        self.write(w, kind, content)
        log.Println("Used in-memory cache.")
        log.Println("cps.WebCore.ServeHTTP: Responded HTTP 200")

        return
    }

    metadata, content = self.Fetcher.Fetch(commonCacheKey)

    self.write(w, metadata.Type, content)

    (*self.Cache).Save(contentTypeCacheKey, []byte(metadata.Type))
    (*self.Cache).Save(contentDataCacheKey, content)
    log.Println("Used actual data.")
    log.Println("cps.WebCore.ServeHTTP: Responded HTTP 200")
}

func (self *WebCore) write(w http.ResponseWriter, kind string, content []byte) {
    w.Header().Set("Content-Type", kind)

    if !self.Compressed {
        w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
        w.Write(content)

        return
    }

    compressed := self.Enigma.Compress(content)

    w.Header().Set("Content-Encoding", "gzip")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", len(compressed)))
    w.Write(compressed)
}
