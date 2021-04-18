// brigolagecms/pkg/sessions/handler.go
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

package sessions

import (
	"context"
	"log"
	"net/http"
	"time"
)

// GetSession will search the request context and header for a session.
// It always returns a valid pointer, but it will point to an invalid
// session if the session is invalid. (That makes no sense.)
func GetSession(r *http.Request) (*Session, bool) {
	s, ok := r.Context().Value(contextKey("session")).(*Session)
	if !ok {
		s = &Session{user: "anonymous"}
	}
	return s, ok
}

// Handle is our middleware adapter for managing sessions.
func (m *Manager) Handle(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// pull session from the request
		sess, ok := m.GetSession(r)
		log.Printf("[sessions] %s %s sess(%q)\n", r.Method, r.URL.Path, sess.user)

		if r.Method == "POST" && r.URL.Path == "/api/session/login" {
			if !ok {
				sess = m.New("mdhender", true, "person")
				log.Printf("[sessions] %s %s sess(%q, login)\n", r.Method, r.URL.Path, sess.user)
			}
			sess.expiresAt = time.Now().Add(m.ttl)
			log.Printf("[sessions] %s %s sess(%q, addTTL)\n", r.Method, r.URL.Path, sess.user)

			log.Printf("[sessions] %s %s sess(%q, addCookie)\n", r.Method, r.URL.Path, sess.user)
			m.AddCookie(w, sess)

			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == "POST" && r.URL.Path == "/api/session/logout" {
			if ok {
				log.Printf("[sessions] %s %s sess(%q, invalidate)\n", r.Method, r.URL.Path, sess.user)
				m.InvalidateSession(sess)
			}

			log.Printf("[sessions] %s %s sess(%q, delCookie)\n", r.Method, r.URL.Path, sess.user)
			m.AddCookie(w, nil)

			w.WriteHeader(http.StatusNoContent)
			return
		}

		if sess.id != "" {
			log.Printf("[sessions] %s %s sess(%q, addCookie)\n", r.Method, r.URL.Path, sess.user)
			m.AddCookie(w, sess)
		}

		// inject the session into the context and call the next handler
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey("session"), sess)))
	}
}
