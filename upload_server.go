package main

import (
	//"encoding/json"
	"github.com/JasonSoft/render"
	//"github.com/go-martini/martini"
	//"github.com/martini-contrib/binding"
	//"github.com/martini-contrib/sessions"
	//"io/ioutil"
	//"log"
	"net/http"
	"os"
	//"html/template"
	"path/filepath"
	//"strconv"
	//"fmt"
	//"time"
	"io"
)

type UploadServerOption struct {
	MaxFileSize int64
	Store       *Store
}

type File struct {
	Url string
}

type Result struct {
	Files []File
}

func (m *myClassic) UseUploadServer(option ...UploadServerOption) error {

	m.Post("/upload", func(render render.Render, r *http.Request, w http.ResponseWriter) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods",
			// "OPTIONS, HEAD, GET, POST, PUT, DELETE",
			"GET, POST, PUT, DELETE",
		)

		//get the multipart reader for the request.
		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := Result{}

		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}
			dst, err := os.Create(filepath.Join(_appDir, part.FileName()))
			defer dst.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			url := r.Host + "/" + part.FileName()
			file := File{Url: url}

			result.Files = append(result.Files, file)
		}

		render.JSON(200, result)
	})

	return nil
}
