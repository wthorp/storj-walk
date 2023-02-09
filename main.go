package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lxn/walk"
	w "github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

var mainWindow *w.MainWindow

func main() {
	var accessGrant, linkURL string
	var splitter *w.Splitter
	var treeView *walk.TreeView
	var tableView *walk.TableView
	var webView *walk.WebView

	ctx := context.Background()

	m := NewModel()

	if err := (d.MainWindow{
		AssignTo: &mainWindow,
		Title:    "Totally 90's Style | Storj Browser",
		MinSize:  d.Size{Width: 600, Height: 400},
		Size:     d.Size{Width: 1024, Height: 640},
		Layout:   d.HBox{MarginsZero: true},
		Children: []d.Widget{
			d.HSplitter{
				AssignTo: &splitter,
				Children: []d.Widget{
					d.TreeView{
						AssignTo:             &treeView,
						Model:                m,
						OnCurrentItemChanged: func() { m.TreeOnCurrentItemChanged(*treeView) },
						ContextMenuItems: []d.MenuItem{
							d.Action{Text: "Refresh", OnTriggered: showAboutBoxAction_Triggered},
							d.Action{Text: "Config", OnTriggered: showAboutBoxAction_Triggered},
						},
					},
					d.TableView{
						AssignTo:      &tableView,
						Model:         m.TableModel,
						StretchFactor: 2,
						Columns: []d.TableViewColumn{
							{DataMember: "Name", Width: 192},
							{DataMember: "Size", Format: "%d", Alignment: d.AlignFar, Width: 64},
							{DataMember: "Modified", Format: "2006-01-02 15:04:05", Width: 120},
						},
						OnCurrentIndexChanged: func() { m.TableOnCurrentIndexChanged(*treeView, *tableView, *webView, linkURL) },
					},
					d.WebView{AssignTo: &webView, StretchFactor: 2},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	splitter.SetFixed(treeView, true)
	splitter.SetFixed(tableView, true)

	// open leveldb cache
	var err error
	cache, err := NewCache()
	if err != nil {
		return
	}
	defer cache.Close()
	accessGrant, linkURL = cache.GetCredentials()
	if accessGrant == "" {
		if cmd, err := RunAccessGrantDialog(mainWindow, &accessGrant, &linkURL); err != nil {
			log.Print(err)
		} else if cmd == walk.DlgCmdOK {
			cache.SetCredentials(accessGrant, linkURL)
		}
	}
	m.Init(ctx, accessGrant, cache)
	if err != nil {
		showError()
	}
	defer m.Close()

	mainWindow.Run()
}

func showAboutBoxAction_Triggered() {
	walk.MsgBox(mainWindow, "About", "Walk Actions Example", walk.MsgBoxIconInformation)
}

func showError(a ...any) {
	w.MsgBox(mainWindow, "Error", fmt.Sprint(a...), w.MsgBoxOK|w.MsgBoxIconError)
}

func RunAccessGrantDialog(owner walk.Form, access, linkURL *string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var accessBox *walk.LineEdit
	var linkURLBox *walk.LineEdit

	return d.Dialog{
		AssignTo:      &dlg,
		Title:         "Enter Storj Access Grant",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       d.Size{Width: 300, Height: 80},
		Layout:        d.VBox{},
		Children: []d.Widget{
			d.Composite{
				Layout: d.HBox{},
				Children: []d.Widget{
					d.Label{Text: "Access Grant"},
					d.LineEdit{AssignTo: &accessBox, Text: access},
				},
			},
			d.Composite{
				Layout: d.HBox{},
				Children: []d.Widget{
					d.Label{Text: "Linksharing URL"},
					d.LineEdit{AssignTo: &linkURLBox, Text: linkURL},
				},
			},
			d.Composite{
				Layout: d.HBox{},
				Children: []d.Widget{
					d.HSpacer{},
					d.PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							*access = accessBox.Text()
							*linkURL = linkURLBox.Text()
							dlg.Accept()
						},
					},
					d.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
