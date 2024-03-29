//go:build linux

// The XDG Base Directory Specification is not perfectly followed.

package main

import (
	"auth0tkn/cache"
	"auth0tkn/profile"
	"auth0tkn/token"
	"flag"
	"fmt"
	"os"
)

var (
	selectedProfile = flag.String("p", "default",
		"Which profile to use")
	raw = flag.Bool("raw", false,
		"Print the raw token, without \"Authorization: Bearer \"")
	printProfile = flag.Bool("print", false,
		"Print data in selected profile instead of getting a token")
)

func main() {
	flag.Parse()

	p, err := profile.Load(*selectedProfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *printProfile {
		fmt.Println(p)
		return
	}

	tkn, hit, err := cache.Check(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	if !hit {
		tkn, err = token.Get(p)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(3)
		}

		err = cache.Update(p, tkn)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(4)
		}
	}

	if *raw {
		fmt.Printf("%s", tkn.Token)
	} else {
		fmt.Printf("Authorization: %s %s", tkn.Prefix, tkn.Token)
	}
}