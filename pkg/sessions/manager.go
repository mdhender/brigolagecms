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
	"github.com/google/uuid"
	"sync"
	"time"
)

// Manager manages sessions
type Manager struct {
	sync.RWMutex
	data map[string]*Session
	ttl  time.Duration // time-to-live for sessions
}

// New returns a new manager using the supplied duration as the time-to-live for sessions.
func New(duration time.Duration) *Manager {
	m := Manager{
		data: make(map[string]*Session),
		ttl:  duration,
	}
	return &m
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

// IsAuthenticated returns true only if the session is active and authenticated.
func (m *Manager) IsAuthenticated(id string) bool {
	m.RLock()
	s, ok := m.data[id]
	m.RUnlock()
	return ok && s.IsAuthenticated()
}

// HasAll returns true only if the session is active and has been assigned all of the required roles.
func (m *Manager) HasAll(id string, roles ...string) bool {
	m.RLock()
	s, ok := m.data[id]
	m.RUnlock()
	return ok && s.HasAll(roles...)
}

// HasAny returns true only if the session is active and has been assigned at least one of the required roles.
func (m *Manager) HasAny(id string, roles ...string) bool {
	m.RLock()
	s, ok := m.data[id]
	m.RUnlock()
	return ok && s.HasAny(roles...)
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
	m.RLock()
	s, ok := m.data[id]
	m.RUnlock()
	if !ok {
		return ""
	}
	return s.User()
}
