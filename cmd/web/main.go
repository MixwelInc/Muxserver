package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type application struct { //thats a pattern called Dependency Injection
	errorlog *log.Logger
	infolog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "Web address HTTP")
	flag.Parse()

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) //creating log for info messages
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{ //initializing app (bcz the handlers are now metods of struct)
		errorlog: errorlog,
		infolog:  infolog,
	}

	srv := &http.Server{ //initializing the server (the only change is redirecting errors to created errlog)
		Addr:     *addr,
		ErrorLog: errorlog,
		Handler:  app.routes(),
	}

	infolog.Printf("Запуск сервера на %s", *addr)
	err := srv.ListenAndServe()
	errorlog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat() // check index.html file in the path
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}
	return f, nil
}
