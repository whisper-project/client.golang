/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	apiRoot = "http://localhost:8080/api/console/v0"
)

func main() {
	log.SetFlags(0)
	prefs, err := LoadPrefs()
	if err != nil {
		log.Fatalf("> Couldn't load preferences: %v", err)
	}
	log.Printf("> Loaded profile ID %s for user %s", prefs.ProfileId, prefs.ProfileEmail)
	err = processInput(prefs)
	if err != nil {
		log.Fatal(err)
	}
}

func processInput(prefs *Prefs) error {
	input := bufio.NewScanner(os.Stdin)
	for {
		if !input.Scan() {
			return input.Err()
		}
		line := strings.TrimSpace(input.Text())
		if strings.HasPrefix(line, "/") {
			cmd, rest, _ := strings.Cut(line[1:], " ")
			if line == "/quit" || line == "/exit" {
				return nil
			}
			if cmd == "" {
				fmt.Println("> No command specified; type '/help' for help")
				continue
			}
			if err := processCommand(prefs, strings.ToLower(cmd), rest); err != nil {
				return err
			}
		} else {
			if err := processTyping(prefs, line); err != nil {
				return err
			}
		}
	}
}

func processTyping(prefs *Prefs, line string) error {
	return nil
}
