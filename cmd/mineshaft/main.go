package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	ctx := context.Background()
	config := init_setup(ctx)

	signalChan := installSignalHandlers(ctx)

	exec, err := StartMinecraftExecution(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-signalChan

		tCtx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		Logger.Print(`Attempting Graceful Shutdown`)

		exec.Stop()

		select {
		case <-tCtx.Done():
			Logger.Print(`Failed graceful shutdown in time allotted`)
			exec.cmd.Process.Kill()
			os.Exit(-1)
		case <-signalChan:
			Logger.Print(`Interupt... Killing`)
			exec.cmd.Process.Kill()
			os.Exit(-1)
		}
	}()

	<-exec.Stopped()

	Logger.Print(`Server shutdown completed. Goodbye`)
}

func init_setup(ctx context.Context) Config {
	workingDir, err := os.Getwd()
	if err != nil {
		Logger.Fatal(err)
	}

	var configFilePath string

	flag.StringVar(&configFilePath, "f", "./mineshaft.json", "the config file path for mineshaft")

	flag.Parse()

	if !filepath.IsAbs(configFilePath) {
		configFilePath = filepath.Join(workingDir, configFilePath)
	}

	config, err := loadConfig(configFilePath, workingDir)
	if err != nil {
		Logger.Fatal(err)
	}

	if !fileExists(config.Machine.WorkingDir) {
		Logger.Print(`Creating working directory`)

		if err := os.MkdirAll(config.Machine.WorkingDir, 0755); err != nil {
			Logger.Fatal(err)
		}
	}

	if !fileExists(config.JarPath()) {
		Logger.Print(`Downloading Server File...`)

		if err := downloadFile(ctx, config.Server.SourceURL, config.JarPath()); err != nil {
			Logger.Fatal(err)
		}
	}

	return config
}
