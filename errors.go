/*
 * Copyright 2024 Daniel C. Brotsky. All rights reserved.
 * All the copyrighted work in this repository is licensed under the
 * GNU Affero General Public License v3, reproduced in the LICENSE file.
 */

package main

import (
	"fmt"
	"io"
)

type userCancelledError string

func (e userCancelledError) Error() string {
	return string(e)
}

func newUserCancelledError(s string) error {
	return userCancelledError(s)
}

type serverError struct {
	code    int
	message string
}

func (e *serverError) Error() string {
	return fmt.Sprintf("server error: %d %s", e.code, e.message)
}

func newServerError(code int, body io.Reader) error {
	message, _ := io.ReadAll(body)
	return &serverError{code, string(message)}
}

type networkError struct {
	err error
}

func (e *networkError) Error() string {
	return fmt.Sprintf("network error: %s", e.err.Error())
}

func newNetworkError(err error) error {
	return &networkError{err}
}

type jsonError struct {
	err error
}

func (e *jsonError) Error() string {
	return fmt.Sprintf("json error: %s", e.err.Error())
}

func newJsonError(err error) error {
	return &jsonError{err}
}

type internalError struct {
	err error
}

func (e *internalError) Error() string {
	return fmt.Sprintf("internal error (report a bug!): %s", e.err.Error())
}

func newInternalError(err error) error {
	return &internalError{err}
}
