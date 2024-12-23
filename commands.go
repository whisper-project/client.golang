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
	rest = strings.TrimSpace(rest)
	switch cmd {
	case "help":
		fmt.Println("> Available commands:")
		fmt.Println("> /help: show this help")
		fmt.Println("> /reset: restart the program, removing preferences")
		fmt.Println("> /nwc <name>: new whisper conversation")
		fmt.Println("> /quit: exit the program")
		fmt.Println("> /wc: show whisper conversations")
	case "wc":
		m, err := getWhisperConversations(prefs)
		if err != nil {
			fmt.Println("> Error getting whisper conversations:", err)
			return nil
		}
		if len(m) == 0 {
			fmt.Println("> No whisper conversations found.")
			return nil
		}
		fmt.Println("> Whisper conversations:")
		for k, v := range m {
			fmt.Printf(">     %s: %s\n", k, v)
		}
	case "dwc":
		err := deleteWhisperConversation(prefs, rest)
		if err != nil {
			fmt.Println("> Error deleting whisper conversation:", err)
			return nil
		}
		fmt.Println("> Deleted whisper conversation:", rest)
	case "iwc":
		id, err := getWhisperConversationId(prefs, rest)
		if err != nil {
			fmt.Println("> Error getting whisper conversation:", err)
			return nil
		}
		fmt.Printf("> Whisper conversation %q has id %s\n", rest, id)
	case "nwc":
		id, err := newWhisperConversation(prefs, rest)
		if err != nil {
			fmt.Println("> Error creating whisper conversation:", err)
			return nil
		}
		fmt.Println("> New whisper conversation:")
		fmt.Printf(">     %s: %s\n", rest, id)
	default:
		fmt.Printf("> Unknown command: %s\n> Type '/help' for help\n", cmd)
	}
	return nil
}
