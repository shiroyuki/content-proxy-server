package iip

import (
    "fmt"
    "bytes"
    "compress/gzip"
    "crypto/md5"
    "encoding/json"
    "hash"
    "io"
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

    actualPath   := filepath.Join(self.basePath, Key)
    content, err := ioutil.ReadFile(actualPath)

    log.Println("Reading from:", actualPath)

    if err != nil {
        return nil
    }

    if !self.Compressed {
        log.Println("Read (uncompressed)")
        return content
    }

    readingB := new(bytes.Buffer)
    writingB := new(bytes.Buffer)

    readingB.Write(content)

    r, _ := gzip.NewReader(readingB)

    defer r.Close()

    io.Copy(writingB, r)
    log.Println("Read (compressed)")

    //ioutil.WriteFile("sample.jpg", writingB.Bytes(), 0644)

    return writingB.Bytes()
}

func (self *GenericCacheDriver) Save(Key string, Content []byte) {
    self.initialize()

    actualPath := filepath.Join(self.basePath, Key)

    log.Println("Writing to:", actualPath)

    if !self.Compressed {
        ioutil.WriteFile(actualPath, Content, 0644)
        log.Println("Save (uncompressed)")
        return
    }

    // Compress the data into the buffer.
    b := new(bytes.Buffer)
    w := gzip.NewWriter(b)

    defer w.Close()

    w.Write(Content)

    ioutil.WriteFile(actualPath, b.Bytes(), 0644)
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

    content := self.FileStorage.Load(key)

    if content != nil {
        log.Println("Cache: Hit")
        return content, nil
    }

    content, contentType, contentLength, err := self.request(url)

    if err != nil {
        log.Println(err)
    }

    self.FileStorage.Save(key, content)
    self.saveMetadata(key, url, contentType, contentLength)

    log.Println("Cache: Missed")

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
