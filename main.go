/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"strings"
)

var (
	apiRoot = "http://localhost:8080/api/console/v0"
)

func main() {
	log.SetFlags(0)
	err := processInput()
	if err != nil {
		log.Fatal(err)
	}
}

func processInput() error {
	var prefs *Prefs
	var err error
	for {
		if prefs == nil {
			prefs, err = LoadPrefs()
			if err != nil {
				log.Fatalf("> Couldn't load preferences: %v", err)
			}
			log.Printf("> Loaded profile ID %s for user %s", prefs.ProfileId, prefs.ProfileEmail)
		}
		line, err := readLine(prefs.TypingOff)
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/") {
			cmd, rest, _ := strings.Cut(line[1:], " ")
			if line == "/reset" {
				if err := DeletePrefs(); err != nil {
					log.Fatalf("> Couldn't delete preferences: %v", err)
				}
				prefs = nil
				continue
			}
			if line == "/quit" || line == "/exit" {
				return nil
			}
			if cmd == "" {
				fmt.Println("> No command specified; type '/help' for help")
				continue
			}
			if err = processCommand(prefs, strings.ToLower(cmd), rest); err != nil {
				return err
			}
		} else {
			if err = processTyping(prefs, line); err != nil {
				return err
			}
		}
	}
}

func processTyping(prefs *Prefs, line string) error {
	return nil
}

func readLine(silent bool) (string, error) {
	if os.Stdin.Fd() != os.Stdout.Fd() && term.IsTerminal(int(os.Stdin.Fd())) {
		return readLineRaw(silent)
	} else {
		return readLineBuffered()
	}
}

func readLineRaw(silent bool) (string, error) {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", userCancelledError(err.Error())
	}
	defer term.Restore(int(os.Stdin.Fd()), state)
	isTyping := false
	var line []byte
	for {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			return "", userCancelledError(err.Error())
		}
		if b[0] == '\n' || b[0] == '\r' {
			fmt.Print("\r\n")
			if isTyping {
				stopTypingSound(silent)
			}
			return string(line), nil
		}
		if b[0] == 3 {
			fmt.Print("\r\n")
			if isTyping {
				stopTypingSound(silent)
			}
			return "", newUserCancelledError("Interrupt received!")
		}
		if b[0] == 127 || b[0] == 8 {
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Printf("\b \b")
				if isTyping && len(line) == 0 {
					stopTypingSound(silent)
					isTyping = false
				}
			}
		} else {
			line = append(line, b[0])
			if !isTyping && len(line) == 1 && line[0] != '/' {
				startTypingSound(silent)
				isTyping = true
			}
			fmt.Printf("%s", string(b))
		}
	}
}

func readLineBuffered() (string, error) {
	input := bufio.NewReader(os.Stdin)
	val, err := input.ReadString('\n')
	if err != nil {
		return "", userCancelledError(err.Error())
	}
	return strings.TrimSuffix(val, "\n"), nil
}

func startTypingSound(silent bool) {
	if !silent {
		fmt.Print("... typing sound starts ...\r\n")
	}
}

func stopTypingSound(silent bool) {
	if !silent {
		fmt.Print("... typing sound stops ...\r\n")
	}
}
