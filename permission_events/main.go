package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/opcoder0/fanotify"
	"golang.org/x/sys/unix"
)

func main() {
	var path string

	flag.StringVar(&path, "path", "", "point path")
	flag.Parse()

	if path == "" {
		fmt.Println("missing path")
		os.Exit(1)
	}
	listener, err := fanotify.NewListener(path, false, fanotify.PostContent)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Listening to events for:", path)
	err = listener.AddWatch(path, fanotify.FileOpenToExecutePermission|fanotify.FileClosedAfterWrite)
	if err != nil {
		fmt.Println("AddWatch:", err)
		os.Exit(1)
	}
	go listener.Start()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for event := range listener.Events {
			fmt.Println(event)
			unix.Close(event.Fd)
		}
	}()

	go func() {
		theList := map[string]bool{
			"/home/saikiran/tmp/a.sh":       true,
			"/home/saikiran/tmp/b-file.txt": false,
		}
		defer wg.Done()
		for event := range listener.PermissionEvents {
			v, ok := theList[event.Path]
			if !ok || v {
				fmt.Println("Allowed:", event)
				listener.Allow(event)
			} else {
				fmt.Println("Denied:", event)
				listener.Deny(event)
			}
			unix.Close(event.Fd)
		}
	}()
	wg.Wait()
	listener.Stop()
}
