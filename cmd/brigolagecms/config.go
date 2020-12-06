// brigolagecms/cmd/brigolagecms/config.go
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
	"flag"
	"github.com/peterbourgon/ff/v3"
	"os"
	"time"
)

type config struct {
	Debug    bool
	FileName string
	Server   struct {
		Scheme  string
		Host    string
		Port    string
		Timeout struct {
			Idle  time.Duration
			Read  time.Duration
			Write time.Duration
		}
		Root string
	}
}

// newConfig returns a configuration.
// It accepts an optional configuration file name.
// If provided, the file must contain a valid JSON object.
//
// The command line overrides environment variables overides configuration file override default values.
func newConfig() (*config, error) {
	var cfg config
	cfg.Server.Scheme = "http"
	cfg.Server.Host = "localhost"
	cfg.Server.Port = "8080"
	cfg.Server.Timeout.Idle = 10 * time.Second
	cfg.Server.Timeout.Read = 5 * time.Second
	cfg.Server.Timeout.Write = 10 * time.Second
	cfg.Server.Root = "D:/GoLand/brigolagecms/public"

	var (
		fs                 = flag.NewFlagSet("brigolagecms", flag.ExitOnError)
		fileName           = fs.String("config", cfg.FileName, "config file (optional)")
		debug              = fs.Bool("debug", cfg.Debug, "log debug information (optional)")
		serverScheme       = fs.String("scheme", cfg.Server.Scheme, "http scheme, either 'http' or 'https'")
		serverHost         = fs.String("host", cfg.Server.Host, "host name (or IP) to listen on")
		serverPort         = fs.String("port", cfg.Server.Port, "port to listen on")
		serverTimeoutIdle  = fs.Duration("idle-timeout", cfg.Server.Timeout.Idle, "http idle timeout")
		serverTimeoutRead  = fs.Duration("read-timeout", cfg.Server.Timeout.Read, "http read timeout")
		serverTimeoutWrite = fs.Duration("write-timeout", cfg.Server.Timeout.Write, "http write timeout")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("BRIGOLAGECMS"), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ff.JSONParser)); err != nil {
		return nil, err
	}

	cfg.Debug = *debug
	cfg.FileName = *fileName
	cfg.Server.Scheme = *serverScheme
	cfg.Server.Host = *serverHost
	cfg.Server.Port = *serverPort
	cfg.Server.Timeout.Idle = *serverTimeoutIdle
	cfg.Server.Timeout.Read = *serverTimeoutRead
	cfg.Server.Timeout.Write = *serverTimeoutWrite

	return &cfg, nil
}
