package cli

import (
	"os"

	term "github.com/appscode/go-term"
	"github.com/appscode/go/io"
	homeDir "github.com/mitchellh/go-homedir"
)

var (
	home, _   = homeDir.Dir()
	apprcPath = home + "/.appscode/apprc.json"
)

type Apprc struct {
	Context string  `json:"context"`
	Auths   []*Auth `json:"auths"`
}

/* Exits if there is any error.*/
func (rc *Apprc) GetAuth() *Auth {
	if rc.Context != "" {
		for _, a := range rc.Auths {
			if a.TeamAddr() == rc.Context {
				term.Env = a.Env
				return a
			}
		}
	}
	return nil
}

func (rc *Apprc) SetAuth(a *Auth) error {
	for i, b := range rc.Auths {
		if b.TeamAddr() == a.TeamAddr() {
			rc.Auths = append(rc.Auths[:i], rc.Auths[i+1:]...)
			break
		}
	}
	rc.Context = a.TeamAddr()
	rc.Auths = append(rc.Auths, a)
	return rc.Write()
}

func (rc *Apprc) DeleteAuth() error {
	if rc.Context != "" {
		for i, a := range rc.Auths {
			if a.TeamAddr() == rc.Context {
				rc.Auths = append(rc.Auths[:i], rc.Auths[i+1:]...)
				rc.Context = ""
				break
			}
		}
	}
	return rc.Write()
}

func (rc *Apprc) Write() error {
	err := io.WriteJson(apprcPath, rc)
	if err != nil {
		return err
	}
	os.Chmod(apprcPath, 0600)
	return nil
}

func LoadApprc() (*Apprc, error) {
	if _, err := os.Stat(apprcPath); err != nil {
		return nil, err
	}

	os.Chmod(apprcPath, 0600)
	rc := &Apprc{}
	err := io.ReadFileAs(apprcPath, rc)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

/* Exits if there is any error.*/
func GetAuthOrDie() *Auth {
	rc, err := LoadApprc()
	if err != nil {
		term.Fatalln("Command requires authentication, please run `appctl login`")
	}
	a := rc.GetAuth()
	if a == nil {
		term.Fatalln("Command requires authentication, please run `appctl login`")
	}
	return a
}

/* Exits if there is any error.*/
func GetAuthOrAnon() (*Auth, bool) {
	rc, err := LoadApprc()
	if err != nil {
		return NewAnonAUth(), false
	}
	a := rc.GetAuth()
	if a == nil {
		return NewAnonAUth(), false
	}
	return a, true
}

func SetAuth(a *Auth) error {
	rc, err := LoadApprc()
	if err != nil {
		rc = &Apprc{}
		rc.Auths = make([]*Auth, 0)
	}
	return rc.SetAuth(a)
}
