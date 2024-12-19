/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package api

type Prefs struct {
	ClientId      string `json:"clientId"`
	ProfileId     string `json:"profileId"`
	ProfileSecret string `json:"profileSecret"`
	ProfileEmail  string `json:"profileEmail"`
	TypingOn      bool   `json:"typingOn"`
	SpeakingOn    bool   `json:"speakingOn"`
}
