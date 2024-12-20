/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"fmt"
	"strings"
)

func processCommand(prefs *Prefs, cmd string, rest string) error {
	switch cmd {
	case "help":
		fmt.Println("> Available commands:")
		fmt.Println("> /help: show this help")
		fmt.Println("> /nwc <name>: new whisper conversation")
		fmt.Println("> /quit: exit the program")
		fmt.Println("> /wc: show whisper conversations")
	case "wc":
		m, err := getWhisperConversations(prefs.ProfileId)
		if err != nil {
			fmt.Println("> Error getting whisper conversations:", err)
			return nil
		}
		if len(m) == 0 {
			fmt.Println("> No whisper conversations found.")
			return nil
		}
		fmt.Println("> Whisper conversations (name: id):")
		for k, v := range m {
			fmt.Printf(">     %s: %s\n", k, v)
		}
	case "nwc":
		words := strings.Fields(rest)
		if len(words) == 0 {
			fmt.Println("> No name specified.\n> Usage: /nwc <name>")
			return nil
		} else if len(words) > 1 {
			fmt.Println("> Names cannot contain spaces.\n> Usage: /nwc <name>")
			return nil
		}
		id, err := newWhisperConversation(prefs.ProfileId, words[0])
		if err != nil {
			fmt.Println("> Error creating whisper conversation:", err)
			return nil
		}
		fmt.Println("> New whisper conversation (name: id):")
		fmt.Printf(">     %s: %s\n", words[0], id)
	default:
		fmt.Printf("> Unknown command: %s\n> Type '/help' for help\n", cmd)
	}
	return nil
}
