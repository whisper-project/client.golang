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
	"net/url"
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
	path := fmt.Sprintf("/profiles/%s/whisper-conversations", p.ProfileId)
	resp, err := SendRequest(p, path, "POST", nil, name)
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

func deleteWhisperConversation(p *Prefs, name string) error {
	path := fmt.Sprintf("/profiles/%s/whisper-conversations/%s", p.ProfileId, url.PathEscape(name))
	resp, err := SendRequest(p, path, "DELETE", nil, nil)
	if err != nil {
		return newNetworkError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("whisper conversation %q not found", name)
	}
	if resp.StatusCode != http.StatusNoContent {
		return newServerError(resp.StatusCode, resp.Body)
	}
	return nil
}

func getWhisperConversationId(p *Prefs, name string) (string, error) {
	path := fmt.Sprintf("/profiles/%s/whisper-conversations/%s", p.ProfileId, url.PathEscape(name))
	resp, err := SendRequest(p, path, "GET", nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", newServerError(resp.StatusCode, resp.Body)
	}
	var id string
	if err := json.NewDecoder(resp.Body).Decode(&id); err != nil {
		return "", newJsonError(err)
	}
	return id, nil
}
