/*
 * Copyright 2022 Iglou.eu
 * Copyright 2022 Adrien Kara
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package safe

import (
	"errors"
	"sync"
)

// Handler defines the function to call when a Panic is catched
// Take recovery interface as parameter
type Handler func(interface{})

var (
	// running is used to know if the listener is running
	running bool

	// mut is used to lock the package when accessing to the running flag
	mut sync.Mutex

	// handler is used to store the panic handler
	handler Handler

	// pending is the channel where the recover is sent
	pending = make(chan interface{})

	// done channel is used to notify listener to quit
	done = make(chan struct{})
)

// Register initializes the safe package
// It must be called before any other function
// If already called, it will clear the last setings
func Register(Panic Handler) bool {

	// Check if the package is already running
	// If it is close the listener
	mut.Lock()
	isRun := running
	mut.Unlock()

	if isRun {
		Close()
	}

	// Set the panic handler
	mut.Lock()
	handler = Panic
	mut.Unlock()

	return true
}

// Listen start Panic listener
// It must be called after Register
func Listen() error {
	mut.Lock()
	defer mut.Unlock()

	// Check if Register has been called
	if handler == nil {
		return errors.New("safe: panic handler not set")
	}

	// Check if the package is already running
	if running {
		return errors.New("safe: package already running, close it first")
	}

	// Create a goroutine to listen to the panic and signal
	go func() {
		for {
			p := <-pending

			// If the channel received nil continue to listen
			if p == nil {
				done <- struct{}{}
				continue
			}

			// If the channel received 'Q' rune quit the goroutine
			if p == 'Q' {
				done <- struct{}{}
				return
			}

			handler(p)
			done <- struct{}{}
		}
	}()

	// Set the running flag to true
	running = true

	return nil
}

// Catch catch a panic and send it to the pending channel and wait until the listener is done
// It must be called at the beginning of a function that we want to catch a panic
// defer func(){safe.Catch(recover())}()
func Catch(p interface{}) {
	mut.Lock()
	defer mut.Unlock()

	// Check if the listener is running
	if !running {
		panic("safe: package not running, call safe.Register and safe.Listen() first")
	}

	pending <- p
	Done()
}

// Done is used to lock until listener have done
func Done() {
	<-done
}

// Close is used to close the listening goroutine
func Close() {
	mut.Lock()
	defer mut.Unlock()

	// Check if the listener is running
	if !running {
		return
	}

	pending <- 'Q'
	Done()

	running = false
}
