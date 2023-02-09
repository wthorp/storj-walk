package main

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type FileInfo struct {
	Name     string
	Size     int64
	Modified time.Time
}

type FileInfoModel struct {
	walk.SortedReflectTableModelBase
	items []*FileInfo
}

var _ walk.ReflectTableModel = new(FileInfoModel)

func NewFileInfoModel() *FileInfoModel {
	return new(FileInfoModel)
}

func (m *FileInfoModel) Items() interface{} {
	return m.items
}

func (m *FileInfoModel) Image(row int) interface{} {
	//TODO: should this be cache?
	hI := hIconForFilePath(m.items[row].Name)
	if hI != 0 {
		ic, err := walk.NewIconFromHICONForDPI(hI, mainWindow.DPI())
		if err != nil {
			fmt.Printf("%+v", err)
		}
		return ic
	}
	return nil
}

func hIconForFilePath(filePath string) win.HICON {
	var shfi win.SHFILEINFO
	flags := uint32(win.SHGFI_USEFILEATTRIBUTES | win.SHGSI_ICON | win.SHGFI_SMALLICON)
	if strings.HasSuffix(filePath, "/") {
		win.SHGetFileInfo(syscall.StringToUTF16Ptr(filePath), 0x80|0x10, &shfi, uint32(unsafe.Sizeof(shfi)), flags)
	} else {
		win.SHGetFileInfo(syscall.StringToUTF16Ptr(filePath), 0x80, &shfi, uint32(unsafe.Sizeof(shfi)), flags)
	}
	return shfi.HIcon
}
