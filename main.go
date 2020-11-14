package main

import (
	"fmt"
	"net"
	"net/http"
	"nookbox/config"
	"nookbox/shlink"
	"nookbox/storage"
	"os"

	"github.com/gorilla/mux"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

// Found this trick that makes half the work
// it is still called after init storage
var _ = func() error { config.Load(); return nil }()

func main() {
	if _, err := os.Stat(config.UploadDir); os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory %s does not exist", config.UploadDir))
	}

	store := filestore.FileStore{
		Path: config.UploadDir,
	}

	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	var basePath = "/files/"

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:      basePath,
		StoreComposer: composer,
	})

	if err != nil {
		panic(fmt.Errorf("Unable to create handler: %s", err))
	}

	http.Handle("/files/", http.StripPrefix("/files/", handler))

	r := mux.NewRouter()
	r.HandleFunc("/url/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vars := mux.Vars(r)
		name := storage.ProcessFile(vars["id"])
		response, _ := shlink.CreateLink(name)
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	})

	http.Handle("/", r)

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
