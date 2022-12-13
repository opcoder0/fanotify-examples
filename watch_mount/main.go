package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opcoder0/fanotify"
	"golang.org/x/sys/unix"
)

func main() {
	var mountPoint string

	flag.StringVar(&mountPoint, "mount-path", "", "mount point path")
	flag.Parse()

	if mountPoint == "" {
		fmt.Println("missing mount path")
		os.Exit(1)
	}
	listener, err := fanotify.NewListener(mountPoint, true, fanotify.PermissionNone)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Listening to events for:", mountPoint)
	var eventTypes fanotify.EventType
	eventTypes = fanotify.FileAccessed |
		fanotify.FileOrDirectoryAccessed |
		fanotify.FileModified |
		fanotify.FileClosedAfterWrite |
		fanotify.FileClosedWithNoWrite |
		fanotify.FileOpened |
		fanotify.FileOrDirectoryOpened |
		fanotify.FileOpenedForExec
	fmt.Printf("Listening on the event mask %x\n", uint64(eventTypes))
	err = listener.WatchMount(eventTypes)
	if err != nil {
		fmt.Println("WatchMount:", err)
		os.Exit(1)
	}
	go listener.Start()
	for event := range listener.Events {
		fmt.Println(event)
		unix.Close(event.Fd)
	}
	listener.Stop()
}
