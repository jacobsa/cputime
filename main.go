// Copyright 2015 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// cputime runs its arguments as a command, then prints CPU usage information
// to stderr. It is like the shell builtin 'time' command, except with stable
// formatting that includes a total across user and system CPU time.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func formatBytes(count int64) string {
	switch {
	case count < 1<<10:
		return fmt.Sprintf("%d bytes", count)

	case count < 1<<20:
		return fmt.Sprintf("%.2f KiB", float64(count)/(1<<10))

	case count < 1<<30:
		return fmt.Sprintf("%.2f MiB", float64(count)/(1<<20))

	default:
		fallthrough

	case count < 1<<40:
		return fmt.Sprintf("%.2f GiB", float64(count)/(1<<30))
	}
}

func main() {
	// Set up flags.
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s [flags...] command [args...]\n",
			os.Args[0])

		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check usage.
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Run the command.
	args := flag.Args()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)
	if err != nil {
		log.Fatalf("Run: %v", err)
		return
	}

	// Print usage info.
	rusage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
	user := time.Duration(rusage.Utime.Nano()) * time.Nanosecond
	sys := time.Duration(rusage.Stime.Nano()) * time.Nanosecond

	fmt.Fprintf(
		os.Stderr,
		"\n"+
			"Max RSS:   %v\n"+
			"Wall time: %v\n"+
			"\n"+
			"User  CPU: %v\n"+
			"Sys   CPU: %v\n"+
			"Total CPU: %v\n",
		formatBytes(rusage.Maxrss*1024),
		duration,
		user,
		sys,
		user+sys)
}
