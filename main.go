package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"tus-server/storage"

	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

func main() {
	flag.Parse()
	storage.LoadConfig()

	var uploadPath string = flag.Args()[0]

	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory %s does not exist", uploadPath))
	}

	store := filestore.FileStore{
		Path: uploadPath,
	}

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	var basePath = "/files/"

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              basePath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})

	if err != nil {
		panic(fmt.Errorf("Unable to create handler: %s", err))
	}

	go func() {
		for {
			event := <-handler.CompleteUploads
			storage.ProcessFile(event)
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	http.Handle("/files/", http.StripPrefix("/files/", handler))

	for _, address := range getIP() {
		fmt.Println(fmt.Sprintf("Listening on: http://%s:%d%s", address, 8080, basePath))
	}

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Errorf("Unable to listen: %s", err))
	}
}

func getIP() (addressList []string) {
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ipv4 := ip.To4(); ipv4 != nil && !ip.IsLoopback() {
				addressList = append(addressList, ipv4.String())
			}
		}
	}

	return addressList
}
