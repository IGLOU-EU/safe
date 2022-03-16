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

package safe_test

import (
	"testing"

	"github.com/IGLOU-EU/safe"
)

// Test a panic case
func TestPanic(t *testing.T) {
	// Create a new safe listener
	safe.Register(
		func(p interface{}) {
			t.Log("Domo arigato, Mr. Roboto. We got it !")
		},
	)
	safe.Listen()
	// Wait 1 second to make sure the panic handler is called
	defer func() { safe.Catch(recover()) }()

	// Make sure we panic
	panic("Ho no, Ganondorf make me Panic. Safe, help me to catch him!")
}

// Test catch panic
func TestCatch(t *testing.T) {
	defer func() { recover(); t.Log("This log proves the success of the test") }()

	// Create a new safe listener
	safe.Register(
		func(p interface{}) {
			t.Log("You are not supposed to see this !")
		},
	)

	// Close
	safe.Close()

	// Be panic
	safe.Catch(nil)
}

// Test listener
func TestListener(t *testing.T) {
	var success string

	// Call the listener without registering it
	safe.Register(nil)
	if err := safe.Listen(); err == nil {
		t.Error("Calling safe.Listen() without registering a handler should return an error")
	}

	// Create a new safe listener
	safe.Register(
		func(p interface{}) {
			success = "panic"
		},
	)
	safe.Listen()

	// Call the listener again
	if err := safe.Listen(); err == nil {
		t.Error("Calling safe.Listen() twice should return an error")
	}

	// List of test cases
	testCases := []struct {
		send     interface{}
		expected string
	}{
		{
			send:     nil,
			expected: "",
		},
		{
			send:     "",
			expected: "panic",
		},
		{
			send:     "Minsc and Boo are the best for kick Panic ass !",
			expected: "panic",
		},
		{
			send:     -42,
			expected: "panic",
		},
		{
			send:     42,
			expected: "panic",
		},
		{
			send:     true,
			expected: "panic",
		},
		{
			send:     false,
			expected: "panic",
		},
		{
			send:     []int{1, 2, 3},
			expected: "panic",
		},
		{
			send:     map[string]string{"foo": "bar"},
			expected: "panic",
		},
	}

	// Run test cases
	for _, testCase := range testCases {
		success = ""
		safe.Catch(testCase.send)
		if success != testCase.expected {
			t.Errorf("Expected %s, got %s", testCase.expected, success)
		}
	}

	// Close
	safe.Close()

	// Call the listener again
	// Expect continue to work normally
	safe.Close()
}
