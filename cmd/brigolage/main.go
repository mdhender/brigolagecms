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

// Package main implements the BrigolageCMS server.
package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	log.Println("######################################################################")
	log.Println("[main] harvesting entropy")
	start := time.Now()
	var entropyBits int64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &entropyBits); err != nil {
		log.Fatal(err)
	}
	rand.Seed(entropyBits)
	fmt.Printf("[main] harvested entropy: %s\n", time.Since(start))

	log.Println("######################################################################")
	log.Println("[main] setting up configuration")
	cfg, err := newConfig()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	} else if cfg.FileName != "" {
		log.Printf("[main] loaded config from %q\n", cfg.FileName)
	} else if data, err := json.MarshalIndent(cfg, "  ", "  "); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	} else {
		log.Printf("[main] config %s\n", string(data))
	}

	log.Println("######################################################################")
	log.Println("[main] starting application")
	if err := run(cfg); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
}
