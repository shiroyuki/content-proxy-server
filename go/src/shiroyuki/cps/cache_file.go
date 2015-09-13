package cps

import "io/ioutil"
import "log"
import "os"
import "path/filepath"

// Cache Driver
type FileCacheDriver struct {
    CacheDriver

    Enigma      Enigma
    StoragePath string
    Compressed  bool
    basePath    string
}

func NewFileCacheDriver(enigma Enigma, storagePath string, compressed bool) FileCacheDriver {
    fcd := FileCacheDriver{
        Enigma:      enigma,
        StoragePath: storagePath,
        Compressed:  compressed,
    }

    fcd.Initialize()

    return fcd
}

func (self *FileCacheDriver) Initialize() {
    self.basePath, _ = filepath.Abs(self.StoragePath)
    os.MkdirAll(self.basePath, 0755)
}

func (self *FileCacheDriver) Load(key string) []byte {
    actualPath   := filepath.Join(self.basePath, key)
    content, err := ioutil.ReadFile(actualPath)

    if err != nil {
        return nil
    }

    if !self.Compressed {
        return content
    }

    return self.Enigma.Decompress(content)
}

func (self *FileCacheDriver) Save(key string, content []byte) {
    actualPath := filepath.Join(self.basePath, key)

    log.Println("Writing to:", actualPath)

    if !self.Compressed {
        ioutil.WriteFile(actualPath, content, 0644)

        return
    }

    compressed := self.Enigma.Compress(content)

    ioutil.WriteFile(actualPath, compressed, 0644)
}
