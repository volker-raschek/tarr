package config

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"git.cryptic.systems/volker.raschek/tarr/pkg/domain"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type XMLConfig struct {
	XMLName              xml.Name `xml:"Config"`
	APIToken             string   `xml:"ApiKey,omitempty"`
	AuthenticationMethod string   `xml:"AuthenticationMethod,omitempty"`
	BindAddress          string   `xml:"BindAddress,omitempty"`
	Branch               string   `xml:"Branch,omitempty"`
	EnableSSL            string   `xml:"EnableSsl,omitempty"`
	InstanceName         string   `xml:"InstanceName,omitempty"`
	LaunchBrowser        string   `xml:"LaunchBrowser,omitempty"`
	LogLevel             string   `xml:"LogLevel,omitempty"`
	Port                 string   `xml:"Port,omitempty"`
	SSLCertPassword      string   `xml:"SSLCertPassword,omitempty"`
	SSLCertPath          string   `xml:"SSLCertPath,omitempty"`
	SSLPort              string   `xml:"SslPort,omitempty"`
	UpdateMechanism      string   `xml:"UpdateMechanism,omitempty"`
	URLBase              string   `xml:"UrlBase,omitempty"`
}

type YAMLConfigAuth struct {
	APIToken string `yaml:"apikey,omitempty"`
	Password string `yaml:"password,omitempty"`
	Type     string `yaml:"type,omitempty"`
	Username string `yaml:"username,omitempty"`
}

type YAMLConfig struct {
	Auth YAMLConfigAuth `yaml:"auth,omitempty"`
}

// Read reads the config struct from a file. The decoding format will be determined by the file extension like
// `xml` or `yaml`.
func ReadConfig(name string) (*domain.Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch {
	case strings.HasSuffix(name, "xml"):
		return readXMLConfig(f)
	case strings.HasSuffix(name, "yml") || strings.HasSuffix(name, "yaml"):
		return readYAMLConfig(f)
	default:
		return nil, fmt.Errorf("Unsupported file extension")
	}
}

func readXMLConfig(r io.Reader) (*domain.Config, error) {
	xmlConfig := new(XMLConfig)

	xmlDecoder := xml.NewDecoder(r)
	err := xmlDecoder.Decode(xmlConfig)
	if err != nil {
		return nil, err
	}

	return &domain.Config{
		API: &domain.API{
			Token: xmlConfig.APIToken,
		},
	}, nil
}

func readYAMLConfig(r io.Reader) (*domain.Config, error) {
	yamlConfig := new(YAMLConfig)

	yamlDecoder := yaml.NewDecoder(r)
	err := yamlDecoder.Decode(yamlConfig)
	if err != nil {
		return nil, err
	}

	return &domain.Config{
		API: &domain.API{
			Password: yamlConfig.Auth.Password,
			Token:    yamlConfig.Auth.APIToken,
			Username: yamlConfig.Auth.Username,
		},
	}, nil
}

func WatchConfig(ctx context.Context, name string) (<-chan *domain.Config, <-chan error) {
	configChannel := make(chan *domain.Config, 0)
	errorChannel := make(chan error, 0)

	go func() {
		wait := time.Second * 3
		timer := time.NewTimer(wait)
		<-timer.C

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			errorChannel <- err
			return
		}
		watcher.Add(name)

		for {
			select {
			case <-ctx.Done():
				close(configChannel)
				close(errorChannel)
				break
			case event, open := <-watcher.Events:
				if !open {
					errorChannel <- fmt.Errorf("FSWatcher closed channel: %w", err)
					break
				}

				switch event.Op {
				case fsnotify.Write:
					timer.Reset(wait)
				}
			case <-timer.C:
				config, err := ReadConfig(name)
				if err != nil {
					errorChannel <- err
					continue
				}
				configChannel <- config
			}
		}
	}()

	return configChannel, errorChannel
}

// WriteConfig writes the config struct into the file. The encoding format will be determined by the file extension like
// `xml` or `yaml`.
func WriteConfig(name string, config *domain.Config) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	switch {
	case strings.HasSuffix(name, "xml"):
		return writeXMLConfig(f, config)
	case strings.HasSuffix(name, "yml") || strings.HasSuffix(name, "yaml"):
		return writeYAMLConfig(f, config)
	default:
		return fmt.Errorf("Unsupported file extension")
	}
}

func writeXMLConfig(w io.Writer, config *domain.Config) error {
	xmlEncoder := xml.NewEncoder(w)
	defer xmlEncoder.Close()

	xmlConfig := &XMLConfig{
		APIToken: config.API.Token,
	}

	xmlEncoder.Indent("", "  ")
	err := xmlEncoder.Encode(xmlConfig)
	if err != nil {
		return err
	}

	return nil
}

func writeYAMLConfig(w io.Writer, config *domain.Config) error {
	yamlEncoder := yaml.NewEncoder(w)
	defer yamlEncoder.Close()

	yamlConfig := &YAMLConfig{
		Auth: YAMLConfigAuth{
			APIToken: config.API.Token,
			Password: config.API.Password,
			Username: config.API.Username,
		},
	}

	yamlEncoder.SetIndent(2)

	err := yamlEncoder.Encode(yamlConfig)
	if err != nil {
		return err
	}

	return nil
}
