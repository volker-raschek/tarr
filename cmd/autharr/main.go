package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git.cryptic.systems/volker.raschek/tarr/pkg/config"
	"git.cryptic.systems/volker.raschek/tarr/pkg/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version string

func main() {
	rootCmd := &cobra.Command{
		Args: cobra.RangeArgs(1, 2),
		Long: `autharr reads the XML or YAML configuration file and prints the API token on stdout`,
		Example: `autharr /etc/bazarr/config.yaml
autharr /etc/lidarr/config.xml`,
		RunE:    runE,
		Version: version,
		Use:     "autharr",
	}
	rootCmd.Flags().Bool("watch", false, "Listens for changes to the configuration and writes the token continuously to the output")

	rootCmd.Execute()
}

func runE(cmd *cobra.Command, args []string) error {
	watchCfg, err := cmd.Flags().GetBool("watch")
	if err != nil {
		return err
	}

	switch watchCfg {
	case true:
		return runWatch(cmd, args)
	case false:
		return runSingle(cmd, args)
	}

	return nil
}

func runSingle(_ *cobra.Command, args []string) error {
	cfg, err := config.ReadConfig(args[0])
	if err != nil {
		return err
	}

	var dest string
	if len(args) == 2 {
		dest = args[1]
	}

	return writeConfig(cfg, dest)
}

func runWatch(cmd *cobra.Command, args []string) error {
	// Initial output
	err := runSingle(cmd, args)
	if err != nil {
		return err
	}

	// Watcher output
	configChannel, errorChannel := config.WatchConfig(cmd.Context(), args[0])

	var dest string
	if len(args) == 2 {
		dest = args[1]
	}

	waitFor := time.Millisecond * 100
	timer := time.NewTimer(waitFor)
	<-timer.C

	var cachedConfig *domain.Config = nil

	for {
		select {
		case <-cmd.Context().Done():
			return nil
		case err := <-errorChannel:
			logrus.WithError(err).Errorln("Received from config watcher")
		case <-timer.C:
			writeConfig(cachedConfig, dest)
		case config := <-configChannel:
			cachedConfig = config
			timer.Reset(waitFor)
		}
	}
}

func writeConfig(config *domain.Config, dest string) error {
	switch {
	case len(dest) <= 0:
		_, err := fmt.Fprintf(os.Stdout, "%s", config.API.Token)
		if err != nil {
			return err
		}
	case len(dest) > 0:
		dirname := filepath.Dir(dest)

		err := os.MkdirAll(dirname, 0755)
		if err != nil {
			return err
		}

		f, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		f.WriteString(config.API.Token)
	}

	return nil
}
