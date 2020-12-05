// brigolagecms/cmd/brigolagecms/server.go
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

// Package main implements the BrigolageCMS server.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	cfg, err := newConfig("")
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	} else if cfg.FileName != "" {
		log.Printf("[main] loaded config from %q\n", cfg.FileName)
	}

	if err := run(cfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
}

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

	srv := &server{}
	srv.Addr = net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
	srv.IdleTimeout = cfg.Server.Timeout.Idle
	srv.ReadTimeout = cfg.Server.Timeout.Read
	srv.WriteTimeout = cfg.Server.Timeout.Write
	srv.MaxHeaderBytes = 1 << 20

	srv.routes = http.DefaultServeMux

	log.Printf("[server] listening on %s\n", srv.Addr)
	return srv.ListenAndServe()
}

type server struct {
	http.Server
	routes http.Handler
}
