/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/term"

	"github.com/whisper-project/client.go/api"
)

var (
	homeDir       = os.Getenv("HOME")
	prefsPath     = path.Join(homeDir, ".whisper")
	apiRoot       = "http://localhost:8080/api/console/v0"
	UserCancelled = errors.New("user cancelled")
)

type Prefs api.Prefs

func (p *Prefs) save() error {
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile(prefsPath, payload, 0o600)
}

func (p *Prefs) post() error {
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	url := apiRoot + "/api/prefs"
	resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		return nil
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		var prefs Prefs
		err = json.NewDecoder(resp.Body).Decode(&prefs)
		if err != nil {
			return err
		}
		*p = prefs
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("unexpected response %d: %s", resp.StatusCode, string(body))
}

func (p *Prefs) validate() error {
	for {
		err := p.post()
		if err != nil {
			return err
		}
		if p.ProfileId != "" {
			return nil
		}
		p.ProfileSecret, err = collectPassword(p.ProfileSecret != "")
		if err != nil {
			return err
		}
	}
}

func NewPrefs() *Prefs {
	return &Prefs{ClientId: uuid.NewString(), TypingOn: true}
}

func LoadPrefs() (*Prefs, error) {
	var prefs *Prefs
	created := false
	msg, err := os.ReadFile(prefsPath)
	if err == nil {
		err = json.Unmarshal(msg, &prefs)
		if err != nil {
			return nil, err
		}
	} else {
		prefs = NewPrefs()
		created = true
		prefs.ProfileEmail, err = collectEmail()
		if err != nil {
			return nil, err
		}
	}
	if err = prefs.validate(); err != nil {
		return nil, err
	}
	if err = prefs.save(); err != nil {
		return nil, err
	}
	if created {
		fmt.Println("A profile has been created for your email: " + prefs.ProfileEmail)
		fmt.Println("Your profile password is: " + prefs.ProfileSecret)
		fmt.Println("Please write it down if you want to use your profile on other devices.")
	}
	return prefs, nil
}

func collectEmail() (string, error) {
	re := regexp.MustCompile("^[-a-z0-9.]+@[-a-z0-9.]+$")
	fmt.Println("Welcome to Whisper! Let's get you set up on this device.")
	for {
		fmt.Print("Enter your email address (or 'help' for help): ")
		var email string
		_, err := fmt.Scanln(&email)
		if err != nil {
			return "", err
		}
		if email == "help" {
			fmt.Println("Every user's profile is tied to their email address.")
			fmt.Println("If you don't already have a profile, one will be created for you.")
			fmt.Println("If you already have a profile, you'll enter your password to use it on this device.")
			fmt.Println("If you don't want to continue now, enter 'exit' to quit.")
			continue
		}
		if email == "exit" {
			return "", UserCancelled
		}
		if re.MatchString(email) {
			return email, nil
		}
		fmt.Println("Sorry, that doesn't look like a valid email address.")
	}
}

func collectPassword(repeat bool) (string, error) {
	if repeat {
		fmt.Println("Sorry, wrong password. Please try again.")
	}
	for {
		fmt.Print("Enter your password (or 'help' for help): ")
		b, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		password := strings.TrimSpace(string(b))
		if password == "" {
			fmt.Println("Sorry, you must enter a password, or 'help' for help.")
			continue
		}
		if strings.ToLower(password) == "help" {
			fmt.Println("Enter the password you were shown when you registered your account.")
			fmt.Println("If you forgot your password, enter 'email' to have it mailed to you.")
			fmt.Println("If you don't want to continue now, enter 'exit' to quit.")
			continue
		}
		if strings.ToLower(password) == "exit" {
			return "", UserCancelled
		}
		if strings.ToLower(password) == "email" {
			fmt.Println("Sorry, this feature is not yet implemented.")
			continue
		}
		return password, nil
	}
}
