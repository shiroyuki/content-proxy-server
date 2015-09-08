package iip

import (
    "fmt"
    "archive/zip"
    "bytes"
    "crypto/md5"
    "encoding/json"
    "hash"
    "io/ioutil"
    "net/http"
    "os"
    "log"
    "path/filepath"
    "strconv"
)

// Enigma - Wrapper to Cryptographer and Hasher
type Enigma struct {
    h hash.Hash
}

func (self *Enigma) Hasher() hash.Hash {
    if self.h == nil {
        self.h = md5.New()
    }

    self.h.Reset()

    return self.h
}

func (self *Enigma) Hash(Content []byte) string {
    sum := self.Hasher().Sum(Content)

    return fmt.Sprintf("%x", sum)
}

func (self *Enigma) HashString(Content string) string {
    data := []byte(Content)

    return self.Hash(data)
}

// Cache Driver
type GenericCacheDriver struct {
    StoragePath string
    Compressed  bool
    initialized bool
    basePath    string
}

func (self *GenericCacheDriver) initialize() {
    if self.initialized {
        return
    }

    self.basePath, _ = filepath.Abs(self.StoragePath)
    os.MkdirAll(self.basePath, 0755)
}

func (self *GenericCacheDriver) Load(Key string) []byte {
    self.initialize()

    actualPath := filepath.Join(self.basePath, Key)
    content, err := ioutil.ReadFile(actualPath)

    log.Println("Reading from:", actualPath)

    if err != nil {
        return nil
    }

    return content
}

func (self *GenericCacheDriver) Save(Key string, Content []byte) {
    self.initialize()

    actualPath := filepath.Join(self.basePath, Key)

    log.Println("Writing to:", actualPath)

    if (!self.Compressed) {
        ioutil.WriteFile(actualPath, Content, 0644)
        log.Println("Save (uncompressed)")
        return
    }

    log.Println("Compression enabled")

    // Compress the data into the buffer.
    ob    := new(bytes.Buffer)
    w     := zip.NewWriter(ob)
    ib, _ := w.Create("image")

    ib.Write(Content)
    w.Close()

    ioutil.WriteFile(actualPath, ob.Bytes(), 0644)
    log.Println("Save (compressed)")
}

// Data Fetcher
type Fetcher struct {
    FileStorage   GenericCacheDriver
    MetadataRepo  GenericCacheDriver
    Cryptographer Enigma
}

func (self *Fetcher) Fetch(url string) ([]byte, error) {
    key := self.Cryptographer.HashString(url)

    log.Printf("Fetching the data from: %s\n", url)
    log.Printf("Cache Key: %s\n", key)

    content, contentType, contentLength, err := self.request(url)

    if err != nil {
        log.Println(err)
    }

    self.FileStorage.Save(key, content)
    self.saveMetadata(key, url, contentType, contentLength)

    return content, nil
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

func (self *Fetcher) saveMetadata(
    Key  string,
    Url  string,
    Type string,
    Size uint64,
) {
    metadata := Metadata{
        Url:  Url,
        Type: Type,
        Size: Size,
    }

    encoded, err := json.Marshal(metadata)

    if err != nil {
        log.Fatal(err)
    }

    self.MetadataRepo.Save(Key, encoded)
}
