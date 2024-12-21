/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/whisper-project/client.golang/api"
)

var (
	homeDir   = os.Getenv("HOME")
	prefsPath = path.Join(homeDir, ".whisper")
)

type Prefs api.Prefs

func (p *Prefs) save() error {
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile(prefsPath, payload, 0o600)
}

func (p *Prefs) post() (bool, error) {
	pnp := *p
	pnp.ProfileSecret = ""
	pnp.ProfileEmail = makeSha1(p.ProfileEmail)
	resp, err := SendRequest(p, "/preferences", "POST", nil, &pnp)
	if err != nil {
		return false, newNetworkError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		// we couldn't authorize against the profile, so get a new password
		return false, nil
	}
	if resp.StatusCode == http.StatusUnauthorized {
		// we now have the profile, but are being challenged for the password
		var prefs Prefs
		err = json.NewDecoder(resp.Body).Decode(&prefs)
		if err != nil {
			return false, newJsonError(err)
		}
		p.ProfileId = prefs.ProfileId
		return false, nil
	}
	if resp.StatusCode == http.StatusNoContent {
		// profile is complete and correct
		return true, nil
	}
	if resp.StatusCode == http.StatusCreated {
		// server created a new profile and returned the password
		var prefs Prefs
		err = json.NewDecoder(resp.Body).Decode(&prefs)
		if err != nil {
			return false, newJsonError(err)
		}
		p.ProfileId = prefs.ProfileId
		p.ProfileSecret = prefs.ProfileSecret
		return true, nil
	}
	return false, newServerError(resp.StatusCode, resp.Body)
}

func (p *Prefs) validate() (bool, error) {
	var authorized bool
	var err error
	for {
		if p.ProfileEmail == "" {
			err = p.collectEmail()
			if err != nil {
				return false, err
			}
		}
		preAuthSecret := p.ProfileSecret
		authorized, err = p.post()
		if err != nil {
			return false, err
		}
		if authorized {
			return preAuthSecret == "", nil
		}
		// profile exists, but we need the password
		err = p.collectPassword()
		if err != nil {
			return false, err
		}
	}
}

func (p *Prefs) requestEmail() (bool, error) {
	resp, err := SendRequest(p, "/request-email", "POST", nil, p.ProfileEmail)
	if err != nil {
		return false, newNetworkError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusNoContent {
		return false, newServerError(resp.StatusCode, resp.Body)
	}
	return true, nil
}

func (p *Prefs) collectEmail() error {
	re := regexp.MustCompile("^[-a-z0-9.]+@[-a-z0-9.]+$")
	p.ProfileSecret = ""
	if p.ProfileEmail == "" {
		fmt.Println("> Welcome to Whisper! Let's get you set up on this device.")
	}
	for {
		fmt.Print("> Enter your email address (or 'help' for help): ")
		email, err := readLine()
		if err != nil {
			return userCancelledError("Interrupt received!")
		}
		email = strings.TrimSpace(strings.ToLower(email))
		if email == "" {
			fmt.Println("> Sorry, you must enter an email, or 'help' for help.")
			continue
		}
		if email == "help" {
			fmt.Println("> Every user's profile is tied to their email address.")
			fmt.Println("> If you don't already have a profile, one will be created for you.")
			fmt.Println("> If you already have a profile, you'll enter your password to use it on this device.")
			fmt.Println("> If you don't want to continue now, enter 'quit' to quit.")
			continue
		}
		if email == "quit" || email == "exit" {
			return newUserCancelledError("Quit!")
		}
		if re.MatchString(email) {
			p.ProfileEmail = email
			return nil
		}
		fmt.Println("> Sorry, that doesn't look like a valid email address.")
	}
}

func (p *Prefs) collectPassword() error {
	if p.ProfileSecret != "" {
		fmt.Println("> Sorry, wrong password. Please try again.")
	}
	p.ProfileSecret = ""
	for {
		fmt.Print("> Enter your password (or 'help' for help): ")
		password, err := readLine()
		if err != nil {
			return userCancelledError("Interrupt received!")
		}
		password = strings.TrimSpace(password)
		if password == "" {
			fmt.Println("> Sorry, you must enter a password, or 'help' for help.")
			continue
		}
		if strings.ToLower(password) == "help" {
			fmt.Println("> Enter the password you were shown when you created your profile.")
			fmt.Printf("> If your email is not %s, enter 'change' to change it.\n", p.ProfileEmail)
			fmt.Println("> If you forgot your password, enter 'email' to have it mailed to you.")
			fmt.Println("> If you don't want to continue now, enter 'quit' to quit.")
			continue
		}
		if strings.ToLower(password) == "quit" || strings.ToLower(password) == "exit" {
			return newUserCancelledError("Quit!")
		}
		if strings.ToLower(password) == "change" {
			p.ProfileEmail = ""
			return nil
		}
		if strings.ToLower(password) == "email" {
			sent, err := p.requestEmail()
			if err != nil {
				return err
			}
			if sent {
				fmt.Println("> Your password has been emailed to you. It should arrive in a few minutes.")
			} else {
				fmt.Println("> You should not have been asked for a password. Please report a bug.")
				fmt.Println("> A new profile for you will be created, and you'll be told your password.")
				p.ProfileSecret = ""
				return nil
			}
			continue
		}
		_, err = uuid.Parse(password)
		if err != nil {
			fmt.Println("> Sorry, that doesn't look like a valid password. Please try again.")
			continue
		}
		p.ProfileSecret = password
		return nil
	}
}

func NewPrefs() *Prefs {
	return &Prefs{ClientId: uuid.NewString()}
}

func LoadPrefs() (*Prefs, error) {
	var prefs *Prefs
	msg, err := os.ReadFile(prefsPath)
	if err == nil {
		err = json.Unmarshal(msg, &prefs)
		if err != nil {
			return nil, newInternalError(err)
		}
	} else {
		prefs = NewPrefs()
	}
	created, err := prefs.validate()
	if err != nil {
		return nil, err
	}
	if err = prefs.save(); err != nil {
		return nil, err
	}
	if created {
		fmt.Println("> A profile has been created for your email: " + prefs.ProfileEmail)
		fmt.Println("> Your profile password is: " + prefs.ProfileSecret)
		fmt.Println("> Please write it down if you want to use your profile on other devices.")
	}
	return prefs, nil
}

func readLine() (string, error) {
	input := bufio.NewReader(os.Stdin)
	val, err := input.ReadString('\n')
	if err != nil {
		return "", userCancelledError(err.Error())
	}
	return strings.TrimSuffix(val, "\n"), nil
}
