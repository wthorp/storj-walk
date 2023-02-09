package main

import (
	"context"
	"strings"

	"github.com/lxn/walk"
	"storj.io/uplink"
)

type Model struct {
	walk.TreeModelBase
	cache *Cache
	storj *uplink.Project
	roots []*Path

	TableModel *FileInfoModel
}

var _ walk.TreeModel = new(Model)

func NewModel() *Model {
	m := &Model{}
	m.TableModel = NewFileInfoModel()
	return m
}

func (m *Model) Init(ctx context.Context, accessGrant string, cache *Cache) error {
	m.cache = cache
	// Parse access grant from string
	access, err := uplink.ParseAccess(accessGrant)
	if err != nil {
		return err
	}
	// Open up the Project we will be working with.
	m.storj, err = uplink.OpenProject(ctx, access)
	if err != nil {
		return err
	}

	// check for any existing cached data
	hasCache := m.cache.HasFiles()

	// add the files and path
	buckets := m.storj.ListBuckets(ctx, nil)
	for buckets.Next() {
		// add bucket to roots
		bucket := buckets.Item().Name
		bucketV := NewPath(bucket, nil)
		m.roots = append(m.roots, bucketV)
		// add all files and paths
		if !hasCache {
			objects := m.storj.ListObjects(ctx, bucket, &uplink.ListObjectsOptions{Recursive: true, System: true, Custom: true})
			for objects.Next() {
				path, name := GetPathAndName(objects.Item().Key)
				m.cache.AddFile(bucket, path, name, &FileInfo{Name: name, Size: objects.Item().System.ContentLength, Modified: objects.Item().System.Created})
				addAllPaths(bucketV, path)
			}
		} else {
			for path, _ := range m.cache.ListAllPaths(bucket) {
				addAllPaths(bucketV, path)
			}
		}
	}
	m.PublishItemsReset(nil)

	return nil
}

func addAllPaths(parent *Path, path string) {
	if path == "" {
		return
	}
	parts := strings.Split(path, "/")
	for _, part := range parts {
		var pathV *Path
		for _, kid := range parent.children {
			if kid.name == part {
				pathV = kid
			}
		}
		if pathV == nil {
			pathV = NewPath(part, parent)
			parent.children = append(parent.children, pathV)
		}
		parent = pathV
	}
}

func GetPathAndName(key string) (path, name string) {
	i := strings.LastIndex(key, "/")
	if i == -1 {
		return "", key
	}
	path = key[0:i]
	name = key[i+1:]
	return path, name
}

// event handlers
func (m *Model) TreeOnCurrentItemChanged(treeView walk.TreeView) {
	dir := treeView.CurrentItem().(*Path)
	if err := m.UpdateTableView(dir.Path()); err != nil {
		walk.MsgBox(mainWindow, "Error", err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
	}
}

func (m *Model) UpdateTableView(bucket, path string) error {
	m.TableModel.items = m.cache.ListFiles(bucket, path)
	m.TableModel.PublishRowsReset()
	return nil
}

func (m *Model) TableOnCurrentIndexChanged(treeView walk.TreeView, tableView walk.TableView, webView walk.WebView, linkURL string) {
	if index := tableView.CurrentIndex(); index > -1 {
		name := m.TableModel.items[index].Name
		dir := treeView.CurrentItem().(*Path)
		bucket, path := dir.Path()
		url := linkURL + bucket + "/" + path + "/" + name
		webView.SetURL(url)
	} else {
		webView.SetURL("about:blank")
	}
}

// boring stuff
func (m *Model) Close() {
	m.storj.Close()
}

func (*Model) LazyPopulation() bool {
	return false
}

func (m *Model) RootCount() int {
	return len(m.roots)
}

func (m *Model) RootAt(index int) walk.TreeItem {
	return m.roots[index]
}
