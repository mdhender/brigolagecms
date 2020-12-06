// brigolagecms/cmd/brigolagecms/run.go
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
	"fmt"
	"github.com/mdhender/brigolage/pkg/http/handlers"
	"github.com/mdhender/brigolage/pkg/services/version"
	"github.com/mdhender/brigolage/pkg/storage/memory"
	"log"
	"mime"
	"net/http"
)

func run(cfg *config) error {
	// go depends on the operating system to associate extensions with mime-types.
	// the default usually works, but some containers need a helping hand.
	for _, content := range []struct{ ext, typ string }{
		{".css", "text/css; charset=utf-8"},
		{".json", "application/json"},
	} {
		if err := mime.AddExtensionType(content.ext, content.typ); err != nil {
			return fmt.Errorf("adding mime type %q %w", content.ext, err)
		}
	}

	if cfg == nil {
		return fmt.Errorf("assert(cfg != nil)")
	}
	if cfg.Server.Scheme == "https" {
		return fmt.Errorf("assert(http.Scheme != https)")
	}

	ds, err := memory.New()
	if err != nil {
		return err
	}

	srv, err := newServer(cfg)
	if err != nil {
		return err
	}
	srv.routes(http.StripPrefix("/", handlers.SPA(cfg.Server.PublicRoot)), version.NewService(ds))

	log.Printf("[server] listening on %s\n", srv.Addr)
	return srv.ListenAndServe()
}
