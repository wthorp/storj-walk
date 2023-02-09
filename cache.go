package main

import (
	"encoding/binary"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Cache struct{ *leveldb.DB }

func NewCache() (*Cache, error) {
	db, err := leveldb.OpenFile("storj.db", nil)
	return &Cache{DB: db}, err
}

func (cache *Cache) SetCredentials(accessGrant, linksharingURL string) (err error) {
	return cache.Put([]byte("AG"), []byte(accessGrant+"|"+linksharingURL), nil)
}

func (cache *Cache) GetCredentials() (string, string) {
	credentials, _ := cache.Get([]byte("AG"), nil)
	parts := strings.SplitN(string(credentials), "|", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// file stuff
func (cache *Cache) HasFiles() bool {
	result := cache.indexSearch("F|", 0)
	return result != nil
}

func (cache *Cache) ListFiles(bucket, path string) (results []*FileInfo) {
	return cache.prefixSearch("F|" + bucket + "|" + path + "|")
}

func (cache *Cache) ListAllPaths(bucket string) map[string]struct{} {
	results := make(map[string]struct{}, 0)
	bytePrefix := util.BytesPrefix([]byte("F|" + bucket + "|"))
	iter := cache.NewIterator(bytePrefix, nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		if len(key) != 0 {
			results[strings.Split(string(iter.Key()), "|")[2]] = struct{}{}
		}
	}
	return results
}

func (cache *Cache) AddFile(bucket, path, name string, fi *FileInfo) (err error) {
	return cache.Put([]byte("F|"+bucket+"|"+path+"|"+name), fileInfoToBytes(fi), nil)
}

// other stuff
func (cache *Cache) indexSearch(prefix string, index int) (result *FileInfo) {
	bytePrefix := util.BytesPrefix([]byte(prefix))
	iter := cache.NewIterator(bytePrefix, nil)
	defer iter.Release()
	for x := 0; x <= index; x++ {
		iter.Next()
	}
	//TODO consider not supressing iter.Error()
	if iter.Value() == nil {
		return nil
	}
	return bytesToFileInfo(iter.Value())
}

func (cache *Cache) prefixSearch(prefix string) (results []*FileInfo) {
	bytePrefix := util.BytesPrefix([]byte(prefix))
	iter := cache.NewIterator(bytePrefix, nil)
	defer iter.Release()
	for iter.Next() {
		if iter.Value() == nil {
			return nil
		}
		results = append(results, bytesToFileInfo(iter.Value()))
	}
	//TODO consider not supressing iter.Error()
	return results
}

func bytesToFileInfo(b []byte) *FileInfo {
	fi := FileInfo{}
	fi.Modified = time.Unix(int64(binary.LittleEndian.Uint64(b[0:8])), 0)
	fi.Size = int64(binary.LittleEndian.Uint64(b[8:16]))
	fi.Name = string(b[16:])
	return &fi
}

func fileInfoToBytes(fi *FileInfo) []byte {
	b := []byte{}
	b = binary.LittleEndian.AppendUint64(b, uint64(fi.Modified.Unix()))
	b = binary.LittleEndian.AppendUint64(b, uint64(fi.Size))
	b = append(b, []byte(fi.Name)...)
	return b
}
