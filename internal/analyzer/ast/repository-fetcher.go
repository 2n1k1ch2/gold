package ast

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"os"
)

type Service struct {
	Name    string
	Version string
	url     string
	root    string
}

func NewService(name, version, url string) *Service {
	return &Service{
		Name:    name,
		Version: version,
		url:     url,
		root:    "",
	}
}
func FetchServiceInfo(services ...*Service) error {
	if len(services) == 0 {
		return errors.New("no services provided")
	}
	for _, service := range services {
		tempDir, err := os.MkdirTemp("", "repo-*")

		if err != nil {
			logger.Error(fmt.Errorf("error creating temp dir: %v", err).Error())
			return err
		}
		_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
			URL:      service.url,
			Depth:    1,
			Progress: os.Stdout,
		})
		if err != nil {
			_ = os.RemoveAll(tempDir)
			return err
		}
		service.root = tempDir
	}
	return nil
}
