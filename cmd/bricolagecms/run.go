// brigolagecms/cmd/bricolagecms/run.go
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func run(cfg *config) error {
	if cfg == nil {
		return fmt.Errorf("assert(cfg != nil)")
	}
	if data, err := json.MarshalIndent(cfg, "  ", "  "); err != nil {
		return err
	} else {
		log.Printf("[config] %s\n", string(data))
	}

	if cfg.Server.Scheme == "https" {
		return fmt.Errorf("assert(http.Scheme != https)")
	}

	log.Printf("[server] serving files from %q\n", cfg.Server.Root)
	if err := os.Chdir(cfg.Server.Root); err != nil {
		return err
	}

	srv := &server{}
	srv.Addr = net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	srv.IdleTimeout = cfg.Server.Timeout.Idle
	srv.ReadTimeout = cfg.Server.Timeout.Read
	srv.WriteTimeout = cfg.Server.Timeout.Write
	srv.MaxHeaderBytes = 1 << 20

	srv.routes = http.DefaultServeMux
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("[server] listening on %s\n", srv.Addr)
	return srv.ListenAndServe()
}
