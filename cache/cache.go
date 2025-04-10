package cache

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
	"time"

	"s14.nl/auth0tkn/profile"
	"s14.nl/auth0tkn/token"
)

func Check(p profile.Profile) (t token.Token, hit bool, err error) {
	cache, err := read()
	if err != nil {
		return t, false, err
	}

	line, hit := cache[hash(p)]
	t = token.Token{
		Prefix:     line.prefix,
		Token:      line.token,
		ValidUntil: line.expiry,
	}
	return t, hit, nil
}

func Update(p profile.Profile, t token.Token) error {
	line := cacheLine{
		hash:   hash(p),
		expiry: t.ValidUntil,
		prefix: t.Prefix,
		token:  t.Token,
	}

	cache, err := read()
	if err != nil {
		return err
	}

	cache[line.hash] = line

	return write(cache)
}

type cacheLine struct {
	hash, prefix, token string
	expiry              int
}

const sep = "\t"

func read() (map[string]cacheLine, error) {
	cache := make(map[string]cacheLine)
	now := int(time.Now().Unix())

	f, err := os.Open(cachePath())
	if err != nil {
		return cache, nil
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		fields := strings.Split(s.Text(), sep)
		if len(fields) != 4 {
			return cache, fmt.Errorf("bad cache entry: expected 4 fields, found %d", len(fields))
		}
		line := cacheLine{
			hash:   fields[0],
			prefix: fields[2],
			token:  fields[3],
		}
		expiry, err := strconv.Atoi(fields[1])
		if err != nil {
			return cache, fmt.Errorf("bad cache entry: %w", err)
		}
		line.expiry = expiry

		if expiry > now {
			cache[line.hash] = line
		}
	}

	if s.Err() != nil {
		return cache, s.Err()
	}

	return cache, nil
}

func write(cache map[string]cacheLine) error {
	err := os.MkdirAll(cacheDir(), 0700)
	if err != nil {
		return err
	}

	f, err := os.Create(cachePath())
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range cache {
		_, err := f.WriteString(strings.Join([]string{
			line.hash,
			strconv.Itoa(line.expiry),
			line.prefix,
			line.token,
		}, sep) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func hash(p profile.Profile) string {
	h := fnv.New64()
	h.Write([]byte(p.String()))
	return hex.EncodeToString(h.Sum(nil))
}

func cacheDir() string {
	dir, ok := os.LookupEnv("XDG_CACHE_HOME")
	if !ok {
		dir = "~/.cache"
	}
	return dir + "/auth0tkn"
}

func cachePath() string {
	return cacheDir() + "/tokens"
}
