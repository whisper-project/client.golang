/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getWhisperConversations(p *Prefs) (map[string]string, error) {
	path := fmt.Sprintf("/profiles/%s/whisper-conversations", p.ProfileId)
	resp, err := SendRequest(p, path, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, newServerError(resp.StatusCode, resp.Body)
	}
	result := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, newJsonError(err)
	}
	return result, nil
}

func newWhisperConversation(p *Prefs, name string) (string, error) {
	url := fmt.Sprintf("/profiles/%s/whisper-conversations", p.ProfileId)
	resp, err := SendRequest(p, url, "POST", nil, name)
	if err != nil {
		return "", newNetworkError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return "", fmt.Errorf("whisper conversation %q already exists", name)
	}
	if resp.StatusCode != http.StatusCreated {
		return "", newServerError(resp.StatusCode, resp.Body)
	}
	var result string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", newJsonError(err)
	}
	return result, nil
}
