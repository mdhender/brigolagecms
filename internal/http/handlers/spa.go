// brigolagecms/pkg/http/handlers/spa.go
//
// Copyright (c) 2020, Michael D Henderson.
// Copyright (c) 2002-2009 Kineticode, Inc. and others.
// Copyright (c) 2001-2002 About.com.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
//  3. Neither the name of the copyright holder nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package handlers

import (
	"log"
	"net/http"
	"os"
	"path"
)

// SPA is a static file server that works with single page applications (Vue, React, Svelte, etc).
//
// If a file is not found, it will return index.html from the root.
// We assume that the caller stripped the path prefix if needed.
//
//    http.StripPrefix("/static", handlers.SPA("/tmp"))
//
// Code based on http.ServeFile, but updated to refuse a directory listing.
//
func SPA(root string) http.HandlerFunc {
	root = path.Clean(root)
	log.Printf("[spa] root %q\n", root)

	indexFile := root + "/index.html"
	if stat, err := os.Stat(indexFile); err != nil {
		log.Printf("[spa] %q %+v\n", indexFile, err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else if stat.IsDir() {
		log.Printf("[spa] %q is a directory!", indexFile)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	} else if !stat.Mode().IsRegular() {
		log.Printf("[spa] %q is not a regular file!", indexFile)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
	log.Printf("[spa] index %q\n", indexFile)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		file := path.Clean(r.URL.Path)
		if file == "" || file == "/" || file == "." {
			file = "index.html"
		}
		file = path.Clean(root + "/" + file)

		stat, err := os.Stat(file)
		if err != nil {
			// with spa, we must serve the index file when we can't find the requested item
			file = indexFile
			if stat, err = os.Stat(file); err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
		}

		// we never serve directory indexes or irregular files
		if stat.IsDir() || !stat.Mode().IsRegular() {
			// we never serve irregular files
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// let the stdlib serve the file.
		rdr, err := os.Open(file)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer rdr.Close()
		http.ServeContent(w, r, file, stat.ModTime(), rdr)
	}
}
