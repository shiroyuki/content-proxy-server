package cps

import "log"
import yotsuba "github.com/shiroyuki/yotsuba-go"
import tori    "github.com/shiroyuki/tori-go"

type Application struct {
    CacheDriver    yotsuba.CacheDriver
    Cryptographer  yotsuba.Enigma
    FetcherService Fetcher
}

func (self *Application) HandleImage (h *tori.Handler) {
    var metadata  Metadata
    var kind      string
    var content   []byte
    var sourceUrl string
    var commonCacheKey      string
    var contentTypeCacheKey string
    var contentDataCacheKey string

    sourceUrl = self.Cryptographer.B64decode(h.Key("destination")[0])

    commonCacheKey = sourceUrl

    contentTypeCacheKey = "ct:"   + commonCacheKey
    contentDataCacheKey = "data:" + commonCacheKey

    kind    = string(self.CacheDriver.Load(contentTypeCacheKey))
    content = self.CacheDriver.Load(contentDataCacheKey)

    if &kind != nil && content != nil {
        self.respond(h, kind, &content)

        log.Println("Used in-memory cache.")
        log.Println("cps.WebCore.ServeHTTP: Responded HTTP 200")

        log.Println("D-2")

        return
    }

    metadata, content = self.FetcherService.Fetch(commonCacheKey)

    self.respond(h, metadata.Type, &content)

    self.CacheDriver.Save(contentTypeCacheKey, []byte(metadata.Type))
    self.CacheDriver.Save(contentDataCacheKey, content)
}

func (self *Application) respond(h *tori.Handler, kind string, content *[]byte) {
    h.SetContentType(kind)
    h.WriteByte(*content)
}
