package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

func MakeProvider[ConfigT any](configFilename string) (*Provider[ConfigT], error) {
	pathPriority := []func() (string, error){
		func() (string, error) {
			return os.Getwd()
		},
		func() (string, error) {
			ex, err := os.Executable()
			if err != nil {
				return "", fmt.Errorf("could not get path to executable: %w", err)
			}
			return filepath.Dir(ex), nil
		},
		func() (string, error) {
			curUser, err := user.Current()
			if err != nil {
				return "", fmt.Errorf("error getting current OS user context: %w", err)
			}
			return curUser.HomeDir, nil
		},
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating config watcher: %w", err)
	}

	p := &Provider[ConfigT]{
		watcher: watcher,
	}

	for _, pathGetter := range pathPriority {
		var path string
		if path, err = pathGetter(); err == nil {
			p.path = filepath.Join(path, configFilename)

			if err = p.loadConfig(); err != nil {
				return nil, fmt.Errorf("error opening config file: %w", err)
			}

			if err = p.watcher.Add(p.path); err != nil {
				return nil, fmt.Errorf("error adding config file to watcher: %w", err)
			}

			p.watchForChanges()
			return p, nil
		}
	}
	return nil, fmt.Errorf("could not open config file: %w", err)
}

type Provider[ConfigT any] struct {
	watcher        *fsnotify.Watcher
	path           string
	cfg            ConfigT
	ChangeCallback func(newConfig ConfigT)
}

func (p *Provider[ConfigT]) GetConfig() ConfigT {
	return p.cfg
}

func (p *Provider[ConfigT]) loadConfig() error {
	cfgFile, err := os.Open(p.path)
	if err != nil {
		return err
	}

	var cfg ConfigT
	d := yaml.NewDecoder(cfgFile)
	if err = d.Decode(&cfg); err != nil {
		return err
	}

	p.cfg = cfg
	return nil
}

func (p *Provider[ConfigT]) watchForChanges() {
	go func() {
		for {
			select {
			case event, ok := <-p.watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) {
					if err := p.loadConfig(); err == nil {
						p.ChangeCallback(p.cfg)
					}
				}
			case _, ok := <-p.watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()
}
