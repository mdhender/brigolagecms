// brigolagecms/pkg/sessions/manager.go
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

// Package sessions implements a session manager.
package sessions

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Manager manages sessions
type Manager struct {
	sync.RWMutex
	data    map[string]*Session
	cookies struct {
		httpOnly bool
		secure   bool
	}
	ttl time.Duration // time-to-live for sessions
}

// New returns a new manager using the supplied duration as the time-to-live for sessions.
func New(duration time.Duration, httpOnly, secure bool) *Manager {
	m := Manager{
		data: make(map[string]*Session),
		ttl:  duration,
	}
	m.cookies.httpOnly = httpOnly
	m.cookies.secure = secure
	return &m
}

// AddCookie adds our session cookie to the response.
// If the session is not valid, we'll ask the client to delete the cached session cookie.
func (m *Manager) AddCookie(w http.ResponseWriter, s *Session) {
	cookie := http.Cookie{
		Name:     "id",
		Path:     "/",
		HttpOnly: m.cookies.httpOnly,
		// SameSite: http.SameSiteStrictMode,
		Secure: m.cookies.secure,
	}
	if s == nil || !s.IsActive() {
		cookie.MaxAge = -1
	} else {
		cookie.Value = s.id
		cookie.Expires = s.expiresAt
	}
	http.SetCookie(w, &cookie)
}

// AddSession adds our session to the request.
func (m *Manager) AddSession(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, contextKey("session"), s)
}

func (m *Manager) Get(id string) (*Session, bool) {
	m.RLock()
	s, ok := m.data[id]
	m.RUnlock()
	if !ok {
		s = &Session{}
	}
	return s, ok
}

// GetSession will search the request context and header for a session.
// It always returns a valid pointer, but it will point to an invalid
// session if the session is invalid. (That makes no sense.)
func (m *Manager) GetSession(r *http.Request) (*Session, bool) {
	// try to pull session from the request context
	if s, ok := GetSession(r); ok {
		return s, ok
	}

	// try to pull a session ID from cookie
	if cookie, err := r.Cookie("id"); err == nil {
		if s, ok := m.Get(cookie.Value); ok {
			return s, ok
		}
	}

	// no cookie, so try to pull the session ID from the bearer token in the header
	if headerAuthText := r.Header.Get("Authorization"); headerAuthText != "" {
		if authTokens := strings.SplitN(headerAuthText, " ", 2); len(authTokens) == 2 {
			if authType, bearerToken := authTokens[0], strings.TrimSpace(authTokens[1]); authType == "Bearer" {
				if s, ok := m.Get(bearerToken); ok {
					return s, ok
				}
			}
		}
	}

	// not in context, cookie, or bearer token
	return &Session{user: "anonymous"}, false
}

// InvalidateSession invalidates a session.
func (m *Manager) InvalidateSession(s *Session) {
	s.Invalidate()
}

// IsAuthenticated returns true only if the session is active and authenticated.
func (m *Manager) IsAuthenticated(id string) bool {
	s, ok := m.Get(id)
	return ok && s.IsAuthenticated()
}

// HasAll returns true only if the session is active and has been assigned all of the required roles.
func (m *Manager) HasAll(id string, roles ...string) bool {
	s, ok := m.Get(id)
	return ok && s.HasAll(roles...)
}

// HasAny returns true only if the session is active and has been assigned at least one of the required roles.
func (m *Manager) HasAny(id string, roles ...string) bool {
	s, ok := m.Get(id)
	return ok && s.HasAny(roles...)
}

// New returns a new session with expiration built in.
// The value for `admin` is ignored if the session is not authenticated.
func (m *Manager) New(user string, authenticated bool, roles ...string) *Session {
	s := &Session{
		id:            uuid.New().String(),
		user:          user,
		authenticated: authenticated,
		roles:         make(map[string]bool),
		createdAt:     time.Now(),
	}
	s.expiresAt = s.createdAt.Add(m.ttl)
	for _, role := range roles {
		s.roles[role] = true
	}
	m.Lock()
	m.data[s.id] = s
	m.Unlock()
	return s
}

// Purge removes expired sessions from the manager.
func (m *Manager) Purge() {
	m.Lock()
	for _, s := range m.data {
		if !s.IsActive() {
			delete(m.data, s.id)
		}
	}
	m.Unlock()
}

// User returns the user associated with the session, even if the session has expired.
func (m *Manager) User(id string) string {
	s, ok := m.Get(id)
	if !ok {
		return ""
	}
	return s.User()
}

// contextKey is the context key type for storing parameters in context.Context.
type contextKey string
