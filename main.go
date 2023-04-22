package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	function "example.com/function"
)

var (
	acceptingConnections int32
)

const defaultTimeout = 10 * time.Second

// It reads two environment variables, "read_timeout" and "write_timeout,"
// to set the read and write timeouts for the server. If these environment
// variables are not set, defaultTimeout value will be used for both.
func main() {
	fmt.Println("Jai Shree Ram")
	readTimeout := parseIntOrDurationValue(os.Getenv("read_timeout"), defaultTimeout)
	writeTimeout := parseIntOrDurationValue(os.Getenv("write_timeout"), defaultTimeout)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8082),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}
	fmt.Println("line : 35")
	http.HandleFunc("/", function.Handle)
	fmt.Println("line : 37")
	listenUntilShutdown(s, writeTimeout)
	fmt.Println("line : 39")
}

// The function starts a goroutine that listens for the SIGTERM signal
// from the operating system. When the signal is received, the function
// logs a message indicating that the server is shutting down and waits
// for the specified shutdown timeout duration using the time.Tick function.
func listenUntilShutdown(s *http.Server, shutdownTimeout time.Duration) {
	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)

		<-sig

		log.Printf("[entrypoint] SIGTERM received.. shutting down server in %s\n", shutdownTimeout.String())

		<-time.Tick(shutdownTimeout)

		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("[entrypoint] Error in Shutdown: %v", err)
		}

		log.Printf("[entrypoint] No new connections allowed. Exiting in: %s\n", shutdownTimeout.String())

		<-time.Tick(shutdownTimeout)

		close(idleConnsClosed)
	}()

	// Run the HTTP server in a separate go-routine.
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("[entrypoint] Error ListenAndServe: %v", err)
			close(idleConnsClosed)
		}
	}()
	// Finally, the function sets a global variable acceptingConnections to 1
	// using the atomic.StoreInt32 function and waits for the idleConnsClosed
	// channel to be closed. This indicates that all idle connections have
	// been closed, and the function can safely exit.
	atomic.StoreInt32(&acceptingConnections, 1)

	<-idleConnsClosed
}

func parseIntOrDurationValue(val string, fallback time.Duration) time.Duration {
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return time.Duration(parsedVal) * time.Second
		}
	}

	duration, durationErr := time.ParseDuration(val)
	if durationErr != nil {
		return fallback
	}
	return duration
}
