// brigolagecms/pkg/sessions/session.go
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
	"time"
)

// Session is
type Session struct {
	id            string
	user          string
	authenticated bool
	roles         map[string]bool
	createdAt     time.Time
	expiresAt     time.Time
}

// CreatedAt returns the timestamp when the session was created.
func (s Session) CreatedAt() time.Time {
	return s.createdAt
}

// HasAll returns true only if the session has been assigned all of the roles.
func (s Session) HasAll(roles ...string) bool {
	hasAll := len(roles) != 0
	for _, role := range roles {
		if hasAll = hasAll && s.roles[role]; !hasAll {
			return false
		}
	}
	return hasAll
}

// HasAny returns true only if the session has been assigned at least one of the roles.
func (s Session) HasAny(roles ...string) bool {
	for _, role := range roles {
		if s.roles[role] {
			return true
		}
	}
	return false
}

// Invalidate invalidates a session.
func (s *Session) Invalidate() {
	if s != nil {
		var zero time.Time
		s.expiresAt = zero
	}
}

// IsActive returns true only if the session has not expired.
func (s Session) IsActive() bool {
	return time.Now().Before(s.expiresAt)
}

// IsAuthenticated returns true only if the session has been authenticated.
func (s Session) IsAuthenticated() bool {
	return s.authenticated
}

// User returns the user.
func (s Session) User() string {
	return s.user
}
