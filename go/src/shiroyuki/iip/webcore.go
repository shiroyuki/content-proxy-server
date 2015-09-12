package iip

import (
    "fmt"
    "log"
    "net/http"
    "regexp"
)

type WebCore struct { // implements http.Handler
    Cache      CacheDriver
    Enigma     Enigma
    Fetcher    Fetcher
    Compressed bool
}

func NewWebCore(cache CacheDriver, enigma Enigma, fetcher Fetcher, compressed bool) WebCore {
    return WebCore{
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
    var routePattern = regexp.MustCompile("^/(?P<url>[^/]+)")

    if !routePattern.MatchString(r.URL.Path) {
        return // TODO HTTP 400
    }

    rawActualUrl := URESearch(routePattern, r.URL.Path)["url"]
    actualUrl    := self.Enigma.B64decode(rawActualUrl)

    log.Println("Processing:", r.URL)
    log.Println("Matched?:", routePattern.MatchString(r.URL.Path))
    log.Println("Interpreted:", actualUrl)

    commonCacheKey := actualUrl

    contentTypeCacheKey := "ct:"   + commonCacheKey
    contentDataCacheKey := "data:" + commonCacheKey

    kind    = string(self.Cache.Load(contentTypeCacheKey))
    content = self.Cache.Load(contentDataCacheKey)

    if &kind != nil && content != nil {
        self.write(w, kind, content)
        log.Println("Used in-memory cache.")

        return
    }

    metadata, content = self.Fetcher.Fetch(commonCacheKey)

    self.write(w, metadata.Type, content)

    self.Cache.Save(contentTypeCacheKey, []byte(metadata.Type))
    self.Cache.Save(contentDataCacheKey, content)
    log.Println("Used actual data.")
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
