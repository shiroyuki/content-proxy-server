package cps

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "log"
    "strconv"
)

// Data Fetcher
type Fetcher struct {
    FileStorage   FileCacheDriver
    MetadataRepo  FileCacheDriver
    Cryptographer Enigma
}

func NewFetcher(
    enigma           Enigma,
    cachePath        string,
    metadataPath     string,
    forceCompression bool,
) Fetcher {
    contentDriver  := NewFileCacheDriver(enigma, cachePath,    true)
    metadataDriver := NewFileCacheDriver(enigma, metadataPath, false)

    fetcher := Fetcher{
        FileStorage:   contentDriver,
        MetadataRepo:  metadataDriver,
        Cryptographer: enigma,
    }

    return fetcher
}

func (self *Fetcher) Fetch(url string) (Metadata, []byte) {
    var metadata *Metadata
    var content  []byte

    key := self.Cryptographer.HashString(url)

    log.Printf("Fetching the data from: %s\n", url)
    log.Printf("Cache Key: %s\n", key)

    content  = self.FileStorage.Load(key)
    metadata = self.loadMetadata(key)

    if content != nil && metadata != nil {
        log.Println("Cache: Hit")

        return *metadata, content
    }

    content, contentType, contentLength, err := self.request(url)

    if err != nil {
        log.Println(err)
    }

    self.FileStorage.Save(key, content)

    metadata = self.createMetadata(url, contentType, contentLength)

    self.saveMetadata(key, metadata)
    log.Println("Cache: Missed")

    return *metadata, content
}

func (self *Fetcher) request(url string) ([]byte, string, uint64, error) {
    resp, err := http.Get(url)

    if err != nil {
        return nil, "", 0, err
    }

    defer resp.Body.Close()

    contentType      := resp.Header.Get("content-type")
    contentLength, _ := strconv.ParseUint(resp.Header.Get("content-length"), 10, 64)

    content, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return nil, "", 0, err
    }

    return content, contentType, contentLength, nil
}

func (self *Fetcher) createMetadata(
    Url  string,
    Type string,
    Size uint64,
) *Metadata {
    return &Metadata{
        Url:  Url,
        Type: Type,
        Size: Size,
    }
}

func (self *Fetcher) loadMetadata(key string) *Metadata {
    var metadata *Metadata

    rawData := self.MetadataRepo.Load(key)
    err     := json.Unmarshal(rawData, &metadata)

    if err != nil {
        log.Fatal("cps.fetcher.Fetcher.loadMetadata/error:", err)
    }

    return metadata
}

func (self *Fetcher) saveMetadata(
    key      string,
    metadata *Metadata,
) {
    encoded, err := json.Marshal(metadata)

    if err != nil {
        log.Fatal("cps.fetcher.Fetcher.saveMetadata/error:", err)
    }

    self.MetadataRepo.Save(key, encoded)
}
