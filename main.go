/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"log"
)

func main() {
	log.SetFlags(0)
	prefs, err := LoadPrefs()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", prefs)
}
