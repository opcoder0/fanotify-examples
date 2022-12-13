package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opcoder0/fanotify"
)

func main() {
	var listenPath string

	flag.StringVar(&listenPath, "listen-path", "", "path to watch events")
	flag.Parse()

	if listenPath == "" {
		fmt.Println("missing listen path")
		os.Exit(1)
	}
	mountPoint := "/"
	listener, err := fanotify.NewListener(mountPoint, false, fanotify.PermissionNone)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Listening to events for:", listenPath)
	var eventTypes fanotify.EventType
	eventTypes =
		fanotify.FileAccessed |
			fanotify.FileOrDirectoryAccessed |
			fanotify.FileModified |
			fanotify.FileOpenedForExec |
			fanotify.FileAttribChanged |
			fanotify.FileOrDirectoryAttribChanged |
			fanotify.FileCreated |
			fanotify.FileOrDirectoryCreated |
			fanotify.FileDeleted |
			fanotify.FileOrDirectoryDeleted |
			fanotify.WatchedFileDeleted |
			fanotify.WatchedFileOrDirectoryDeleted |
			fanotify.FileMovedFrom |
			fanotify.FileOrDirectoryMovedFrom |
			fanotify.FileMovedTo |
			fanotify.FileOrDirectoryMovedTo |
			fanotify.WatchedFileMoved |
			fanotify.WatchedFileOrDirectoryMoved
	err = listener.AddWatch(listenPath, eventTypes)
	if err != nil {
		fmt.Println("MarkMount:", err)
		os.Exit(1)
	}
	go listener.Start()
	for event := range listener.Events {
		fmt.Println(event)
	}
	listener.Stop()
}
