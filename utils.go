/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

//goland:noinspection GoDfaNilDereference
func SendRequest(prefs *Prefs, path, verb string, query url.Values, body any) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s", apiRoot, path)
	if query != nil {
		uri += "?" + query.Encode()
	}
	var req *http.Request
	var err error
	switch verb {
	case "GET", "HEAD", "DELETE":
		req, err = http.NewRequest(verb, uri, nil)
	case "POST", "PUT", "PATCH":
		payload, err := json.Marshal(body)
		if err == nil {
			req, err = http.NewRequest(verb, uri, bytes.NewReader(payload))
		}
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	default:
		err = fmt.Errorf("unknown HTTP verb: %s", verb)
	}
	if err != nil {
		return nil, newInternalError(err)
	}
	req.Header.Set("X-Client-Id", prefs.ClientId)
	if prefs.ProfileId != "" && prefs.ProfileSecret != "" {
		if token, err := makeJwt(prefs); err != nil {
			return nil, newInternalError(err)
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}
	}
	return http.DefaultClient.Do(req)
}

func makeJwt(prefs *Prefs) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   prefs.ClientId,
		Subject:  prefs.ProfileId,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	})
	key, err := uuid.Parse(prefs.ProfileSecret)
	if err != nil {
		return "", err
	}
	keyBytes, err := key.MarshalBinary()
	if err != nil {
		return "", err
	}
	signedToken, err := token.SignedString(keyBytes)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// from https://stackoverflow.com/a/10701951/558006
func makeSha1(s string) string {
	hashFn := sha1.New()
	hashFn.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hashFn.Sum(nil))
}
