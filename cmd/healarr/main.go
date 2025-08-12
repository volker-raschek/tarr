package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.cryptic.systems/volker.raschek/tarr/pkg/config"
	"git.cryptic.systems/volker.raschek/tarr/pkg/domain"
	"git.cryptic.systems/volker.raschek/tarr/pkg/health"
	"github.com/spf13/cobra"
)

var version string

func main() {
	bazarrCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `healarr bazarr https://bazarr.example.com:8443 --config /etc/bazarr/config.yaml
healarr bazarr https://bazarr.example.com:8443 --api-token my-token`,
		RunE:  runBazarrE,
		Short: "Check if a bazarr instance is healthy",
		Use:   `bazarr`,
	}

	lidarrCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `healarr lidarr https://lidarr.example.com:8443 --config /etc/lidarr/config.xml
healarr lidarr https://lidarr.example.com:8443 --api-token my-token`,
		RunE:  runLidarrE,
		Short: "Check if a lidarr instance is healthy",
		Use:   `lidarr`,
	}

	prowlarrCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `healarr prowlarr https://prowlarr.example.com:8443 --config /etc/prowlarr/config.xml
healarr prowlarr https://prowlarr.example.com:8443 --api-token my-token`,
		RunE:  runProwlarrE,
		Short: "Check if a prowlarr instance is healthy",
		Use:   `prowlarr`,
	}

	radarrCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `healarr radarr https://radarr.example.com:8443 --config /etc/radarr/config.xml
healarr radarr https://radarr.example.com:8443 --api-token my-token`,
		RunE:  runRadarrE,
		Short: "Check if a radarr instance is healthy",
		Use:   `radarr`,
	}

	readarrCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `healarr readarr https://readarr.example.com:8443 --config /etc/readarr/config.xml
healarr readarr https://readarr.example.com:8443 --api-token my-token`,
		RunE:  runReadarrrE,
		Short: "Check if a readarr instance is healthy",
		Use:   `readarr`,
	}

	sabnzbdCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `healarr sabnzbd https://sabnzbd.example.com:8443 --config /etc/sabnzbd/config.xml
healarr sabnzbd https://sabnzbd.example.com:8443 --api-token my-token`,
		RunE:  runSabNZBdE,
		Short: "Check if a sabnzbd instance is healthy",
		Use:   `sabnzbd`,
	}

	sonarrCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `healarr sonarr https://sonarr.example.com:8443 --config /etc/sonarr/config.xml
healarr sonarr https://sonarr.example.com:8443 --api-token my-token`,
		RunE:  runSonarr,
		Short: "Check if a sonarr instance is healthy",
		Use:   `sonarr`,
	}

	rootCmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Long:    `healarr exits with a non zero exit code, when the *arr application is not healthy`,
		Version: version,
		Use:     "healarr",
	}
	rootCmd.AddCommand(bazarrCmd)
	rootCmd.AddCommand(lidarrCmd)
	rootCmd.AddCommand(prowlarrCmd)
	rootCmd.AddCommand(radarrCmd)
	rootCmd.AddCommand(readarrCmd)
	rootCmd.AddCommand(sabnzbdCmd)
	rootCmd.AddCommand(sonarrCmd)
	rootCmd.PersistentFlags().String("api-token", "", "Token to access the *arr application")
	rootCmd.PersistentFlags().String("config", "", "Path to an XML or YAML config file")
	rootCmd.PersistentFlags().Bool("insecure", false, "Trust insecure TLS certificates")
	rootCmd.PersistentFlags().Duration("timeout", time.Minute, "Timeout")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func runBazarrE(cmd *cobra.Command, args []string) error {
	return runE(cmd, args, domain.BazarrAPIQueryKeyAPIToken)
}

func runLidarrE(cmd *cobra.Command, args []string) error {
	return runE(cmd, args, domain.LidarrAPIQueryKeyAPIToken)
}

func runProwlarrE(cmd *cobra.Command, args []string) error {
	return runE(cmd, args, domain.ProwlarrAPIQueryKeyAPIToken)
}

func runRadarrE(cmd *cobra.Command, args []string) error {
	return runE(cmd, args, domain.RadarrAPIQueryKeyAPIToken)
}

func runReadarrrE(cmd *cobra.Command, args []string) error {
	return runE(cmd, args, domain.ReadarrAPIQueryKeyAPIToken)
}

func runSabNZBdE(cmd *cobra.Command, args []string) error {
	return runE(cmd, args, domain.SabNZBdAPIQueryKeyAPIToken)
}

func runSonarr(cmd *cobra.Command, args []string) error {
	return runE(cmd, args, domain.SonarrAPIQueryKeyAPIToken)
}

func runE(cmd *cobra.Command, args []string, queryKey string) error {
	apiToken, err := cmd.Flags().GetString("api-token")
	if err != nil {
		return err
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	insecure, err := cmd.Flags().GetBool("insecure")
	if err != nil {
		return err
	}

	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		return err
	}

	readinessProbeCtx, cancel := context.WithTimeout(cmd.Context(), timeout)
	defer func() { cancel() }()

	switch {
	case len(apiToken) <= 0 && len(configPath) <= 0:
		return fmt.Errorf("at least --api-token oder --config must be defined")
	case len(apiToken) > 0 && len(configPath) <= 0:
		err = health.NewReadinessProbe(args[0]).
			QueryAdd(queryKey, apiToken).
			Insecure(insecure).
			Run(readinessProbeCtx)
		if err != nil {
			return err
		}
	case len(apiToken) <= 0 && len(configPath) > 0:
		config, err := config.ReadConfig(configPath)
		if err != nil {
			return err
		}

		err = health.NewReadinessProbe(args[0]).
			QueryAdd(queryKey, config.API.Token).
			Insecure(insecure).
			Run(readinessProbeCtx)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("neither --api-token nor --config can be used at the same time")
	}

	return nil
}
