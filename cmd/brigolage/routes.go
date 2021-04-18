/*
 * brigolagecms - a reimagining of Bricolage CMS
 *
 * Copyright (c) 2021, Michael D Henderson.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *   1. Redistributions of source code must retain the above copyright notice, this
 *      list of conditions and the following disclaimer.
 *   2. Redistributions in binary form must reproduce the above copyright notice,
 *      this list of conditions and the following disclaimer in the documentation
 *      and/or other materials provided with the distribution.
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"errors"
	"github.com/mdhender/brigolage/internal/http/jsonapi"
	"github.com/mdhender/brigolage/internal/services/version"
	"github.com/mdhender/brigolage/internal/sessions"
	"github.com/mdhender/brigolage/internal/way"
	"net/http"
)

func (s *server) routes(sm *sessions.Manager, spa http.Handler, v version.Service) http.Handler {
	router := way.NewRouter()

	router.Handle("GET", "/api/version", getVersion(v))

	router.NotFound = spa

	return router
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
func postLogin(m *sessions.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//m.AddCookie(w, m.New("mdhender", true, "person"))
		jsonapi.Ok(w, r, http.StatusOK, true)
	}
}

// postLogout returns a handler for POST /logout requests
func postLogout(m *sessions.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s, ok := m.GetSession(r); ok {
			m.InvalidateSession(s)
		}
		jsonapi.Ok(w, r, http.StatusOK, true)
	}
}
