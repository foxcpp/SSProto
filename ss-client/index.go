// index.go - client files hashing ("indexing")
// Copyright (c) 2018  Hexawolf
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
package main

import (
	"path/filepath"
	"regexp"
	"os"
	"io/ioutil"
	"crypto/sha256"
)

// excludedGlob is a collection of snowflakes ❄️
// This is a list of files and dirs that should not be hashed. That is, their existence is ignored by updater.
var excludedGlob = []string{
	"/?ignored_*",
	"assets",
	"screenshots",
	"saves",
	"library",
}

func shouldExclude(path string) bool {
	for _, pattern := range excludedGlob {
		if match, _ := regexp.MatchString(pattern, filepath.ToSlash(path)); match {
			return true
		}
	}
	return false
}

func collectRecurse(root string) ([]string, error) {
	var res []string
	walkfn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldExclude(path) {
				return filepath.SkipDir
			}
			return nil
		}
		if shouldExclude(path) {
			return nil
		}

		res = append(res, path)
		return nil
	}
	err := filepath.Walk(root, walkfn)
	return res, err
}

func collectHashList() (map[string][]byte, error) {
	res := make(map[string][]byte)

	list, err := collectRecurse(".")
	if err != nil {
		return nil, err
	}

	// A very special snowflake for Hexamine ❄️
	authlib := "libraries/com/mojang/authlib/1.5.25/authlib-1.5.25.jar"
	if fileExists(authlib) {
		list = append(list, filepath.ToSlash(authlib))
	}

	for _, path := range list {
		blob, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(blob)
		res[path] = sum[:]
	}
	return res, nil
}
