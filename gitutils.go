package mgmirr

import (
	"fmt"
	"gopkg.in/libgit2/git2go.v27"
	"log"
)

type RemoteConfig struct {
	Name string
	URL  string
}

// For an existing Git repo and an RPM (e.g. cowsay) Setup the remotes.
//
// This is a best effort procedure.  Not all remotes will be available
// (fedora might not have package x).  As long as at least one remote
// works it is a success.
func SetupRpmRemotes(repo *git.Repository, rcs []RemoteConfig) error {

	var one_worked bool = false

	for _, rc := range rcs {

		// try to set up the remote, continue if it doesn't work
		err := setupRpmRemote(repo, &rc)
		if err != nil {
			log.Println(err)
		} else {
			one_worked = true
		}
	}

	if one_worked {
		return nil
	} else {
		return fmt.Errorf("unable to setup any remotes")
	}
}

func setupRpmRemote(repo *git.Repository, cfg *RemoteConfig) error {
	_, err := repo.Remotes.Create(cfg.Name, cfg.URL)
	if err != nil {
		return fmt.Errorf("git add remote for '%v' failed: %v", cfg.Name, err)
	}

	return nil
}

func FetchAll(repo *git.Repository) error {
	var one_worked bool = false

	remotes, err := repo.Remotes.List()
	if err != nil {
		return fmt.Errorf("unable to list remotes: %v", err)
	}

	for _, remote := range remotes {
		r, err := repo.Remotes.Lookup(remote) // get Remote obj
		if err != nil {
			log.Printf("unable to find remote '%v': %v", remote, err)
			continue
		}

		err = r.Fetch(nil, nil, "")
		if err != nil {
			log.Printf("git fetch remote '%v' failed: %v", remote, err)
		} else {
			one_worked = true
		}
	}

	if one_worked {
		return nil
	} else {
		return fmt.Errorf("unable to fetch any remotes")
	}
}
