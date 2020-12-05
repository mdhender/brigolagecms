package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
)

func main() {
	publicRoot := "D:/GoLand/brigolagecms/cmd/about/public"
	if val := os.Getenv("ABOUT_PUBLIC_ROOT"); val != "" {
		publicRoot = val
		fmt.Printf("[about] %-30s == %q\n", "ABOUT_PUBLIC_ROOT", val)
	}
	if err := os.Chdir(publicRoot); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
	fmt.Printf("[about] serving files from %q\n", publicRoot)
	http.Handle("/", http.FileServer(http.Dir(".")))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
}

func init() {
	// this is stupid, really stupid.
	// go depends on the operating system to associate extensions with mime-types.
	// the default works on most, but not all, environments.
	// i'm hitting it with this hammer to make everyone happy.
	if err := mime.AddExtensionType(".css", "text/css; charset=utf-8"); err != nil {
		log.Fatal(err)
	}
}
