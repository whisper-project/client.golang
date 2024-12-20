/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func getWhisperConversations(pId string) (map[string]string, error) {
	url := fmt.Sprintf("%s/profiles/%s/whisper-conversations", apiRoot, pId)
	resp, err := http.Get(url)
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

func newWhisperConversation(pId, name string) (string, error) {
	url := fmt.Sprintf("%s/profiles/%s/whisper-conversations", apiRoot, pId)
	body, err := json.Marshal(name)
	if err != nil {
		return "", newJsonError(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", newNetworkError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", newServerError(resp.StatusCode, resp.Body)
	}
	var result string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", newJsonError(err)
	}
	return result, nil
}
