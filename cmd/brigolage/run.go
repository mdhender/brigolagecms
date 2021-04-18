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
	"fmt"
	"github.com/mdhender/brigolage/internal/http/handlers"
	"github.com/mdhender/brigolage/internal/services/version"
	"github.com/mdhender/brigolage/internal/sessions"
	"github.com/mdhender/brigolage/internal/storage/memory"
	"log"
	"mime"
	"net/http"
	"time"
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
	sm := sessions.New(15*time.Hour, true, cfg.Server.Scheme == "https")
	spa := http.StripPrefix("/", handlers.SPA(cfg.Server.PublicRoot))

	srv, err := newServer(cfg)
	if err != nil {
		return err
	}

	srv.Handler = sm.Handle(srv.routes(sm, spa, version.NewService(ds)))

	log.Printf("[server] listening on %s\n", srv.Addr)
	return srv.ListenAndServe()
}
