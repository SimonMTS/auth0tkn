# Auth0tkn

## Usage

```shell
$ curl -H "$(auth0tkn)" example.com
$ curl -H "$(auth0tkn -p henk)" example.com
$ curl -H "Authorization: Bearer $(auth0tkn -raw)" example.com
```

## Config

```shell
$ cat ~/.config/auth0tkn/profiles
tenant  example https://company.eu.auth0.com https://example.com client-id client-secret
profile default example freek@example.com pa$$word123
profile henk    example henk@example.com  321abc
```

## Install

```shell
go install s14.nl/auth0tkn@0.0.2
```
