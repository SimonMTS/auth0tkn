package profile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type (
	Profile struct {
		Key,
		Username,
		Password string
		Tenant Tenant
	}

	Tenant struct {
		Key,
		Url,
		Audience,
		ClientId,
		ClientSecret string
	}
)

func Load(key string) (Profile, error) {
	p := Profile{}

	f, err := os.Open(configPath())
	if err != nil {
		return p, err
	}
	defer f.Close()

	profiles := make(map[string]Profile)
	tenants := make(map[string]Tenant)

	s := bufio.NewScanner(f)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) == 0 || fields[0][0] == '#' {
			continue
		}

		switch fields[0] {
		case "profile":
			if len(fields) != 5 {
				return p, fmt.Errorf("bad profile line: expected 5 fields, found %d", len(fields))
			}
			if _, ok := profiles[fields[1]]; ok {
				return p, fmt.Errorf("bad profile line: repeat name '%s'", fields[1])
			}

			profiles[fields[1]] = Profile{
				Key:      fields[1],
				Tenant:   Tenant{Key: fields[2]},
				Username: fields[3],
				Password: fields[4],
			}

		case "tenant":
			if len(fields) != 6 {
				return p, fmt.Errorf("bad tenant line: expected 6 fields, found %d", len(fields))
			}
			if _, ok := tenants[fields[1]]; ok {
				return p, fmt.Errorf("bad tenant line: repeat name '%s'", fields[1])
			}

			tenants[fields[1]] = Tenant{
				Key:          fields[1],
				Url:          fields[2],
				Audience:     fields[3],
				ClientId:     fields[4],
				ClientSecret: fields[5],
			}

		default:
			return p, fmt.Errorf("unexpected text '%s', expected 'tenant' or 'profile'", fields[0])
		}
	}

	if s.Err() != nil {
		return p, s.Err()
	}

	p, ok := profiles[key]
	if !ok {
		return p, fmt.Errorf("no profile found for key: %s", key)
	}

	t, ok := tenants[p.Tenant.Key]
	if !ok {
		return p, fmt.Errorf("no tenant found for key: %s", p.Tenant.Key)
	}
	p.Tenant = t

	return p, nil
}

func configPath() string {
	dir, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if !ok {
		dir = "~/.config"
	}
	return dir + "/auth0tkn/profiles"
}

func (p Profile) String() string {
	return fmt.Sprintf(
		"Key:      %s\n"+
			"Username: %s\n"+
			"Password: %s\n"+
			"Tenant:\n"+
			"    Key:          %s\n"+
			"    Url:          %s\n"+
			"    Audience:     %s\n"+
			"    ClientId:     %s\n"+
			"    ClientSecret: %s",
		p.Key,
		p.Username,
		p.Password,
		p.Tenant.Key,
		p.Tenant.Url,
		p.Tenant.Audience,
		p.Tenant.ClientId,
		p.Tenant.ClientSecret,
	)
}
