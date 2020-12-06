// brigolagecms/cmd/brigolagecms/routes.go
//
// Copyright (c) 2020, Michael D Henderson.
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

package main

import (
	"errors"
	"github.com/mdhender/brigolage/pkg/http/jsonapi"
	"github.com/mdhender/brigolage/pkg/services/version"
	"github.com/mdhender/brigolage/pkg/way"
	"net/http"
)

func (s *server) routes(spa http.Handler, v version.Service) {
	router := way.NewRouter()

	router.Handle("GET", "/version", getVersion(v))

	router.Handle("POST", "/login", postLogin())
	router.Handle("POST", "/logout", postLogout())

	router.NotFound = spa

	s.Handler = router
}

// getVersion returns a handler for GET /version requests
func getVersion(svc version.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := svc.GetVersion()
		if err != nil {
			if errors.Is(err, version.ErrNotFound) {
				jsonapi.Error(w, r, http.StatusNotFound, err.Error())
				return
			}
			jsonapi.Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		jsonapi.Ok(w, r, http.StatusOK, data)
	}
}

// postLogin returns a handler for POST /login requests
func postLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonapi.Ok(w, r, http.StatusOK, true)
	}
}

// postLogout returns a handler for POST /logout requests
func postLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonapi.Ok(w, r, http.StatusOK, true)
	}
}
