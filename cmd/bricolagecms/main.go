// brigolagecms/cmd/bricolagecms/main.go
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

// Package main implements a static web service for our copy of the BricolageCMS website.
package main

import (
	"fmt"
	"log"
	"mime"
	"os"
)

func main() {
	// go depends on the operating system to associate extensions with mime-types.
	// the default usually works, but some containers need a helping hand.
	if err := mime.AddExtensionType(".css", "text/css; charset=utf-8"); err != nil {
		log.Fatal(err)
	}

	cfg, err := newConfig()
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
