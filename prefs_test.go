/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"testing"
)

func TestLoadPrefs(t *testing.T) {
	prefs, err := LoadPrefs()
	if err != nil {
		t.Fatal(err)
	}
	if prefs == nil {
		t.Fatal("prefs is nil")
	}
	t.Log(prefs)
}
