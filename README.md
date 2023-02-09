# Totally 90's Style | Storj Browser

In a world where Windows 95 rules over the PC desktop land, one man will travel through time... to bring back distributed object storage.

## What?!

Storj's distributed object storage platform is written in Go.  Unfortunately, a lot of Go programs are either cryptic command-line tools or web applications that require complex web-application stacks.  Significant community efforts have gone into bringing Storj to various programming languages, often for some end-user facing functionality.

So why not write simple native-desktop applications in Go?  One issue is the lack of a clear winning technology.  A [myriad of options](https://github.com/go-graphics/go-gui-projects/blob/master/README.md) exist, but which one to use?  [Gio](https://gioui.org/) is interesting, but lacks that native desktop feel.  [Qt](https://github.com/therecipe/qt) is a solid choice, with some potential licensing questions.  

At the end of the day, you get started with something.  [Walk](https://github.com/lxn/walk) ("Windows Application Library Kit") is a Windows only wrapper 
for the Win32 GUI Common Controls.  These controls date back to 1995, and should inspire a visual nostalgia for anyone who predates Google Docs.  While Android 
greatly dominates multi-platform OS statistics, Windows still holds a ~75% market share for desktops.  In many cases, prototyping or releasing 
via Windows-only applications is the easiest approach, rather than building anything too complex.

## Features

- multi-bucket file browsing of the Storj Network
- file metadata caching via LevelDB for a responsive GUI
- icons!  without them, the GUI looks sad
- uses Storj Link Sharing + WebView to preview some content types

## Unimplemented

- data virtualization:  EG: sorting, filtering, paging at the cache layer rather than in the GUI
- file upload / download
- icon caching (I'm not sure if this has value?)
- refreshing the cache
- changing credentials (just delete the database for now)
- a modern WebView... this features is based off of IE version ancient and doesn't really work