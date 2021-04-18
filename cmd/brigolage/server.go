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
	"net"
	"net/http"
)

type server struct {
	http.Server
}

// newServer returns an initialized server.
// the main change from the default server is that we override the default timeouts.
// see the following sources for an explanation of why:
//   https://blog.cloudflare.com/exposing-go-on-the-internet/
//   https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
//   https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
func newServer(cfg *config, options ...func(*server) error) (*server, error) {
	srv := &server{}
	srv.Addr = net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	srv.IdleTimeout = cfg.Server.Timeout.Idle
	srv.ReadTimeout = cfg.Server.Timeout.Read
	srv.WriteTimeout = cfg.Server.Timeout.Write
	srv.MaxHeaderBytes = 1 << 20

	// allow caller to override the default values
	for _, option := range options {
		if err := option(srv); err != nil {
			return nil, err
		}
	}

	return srv, nil
}
