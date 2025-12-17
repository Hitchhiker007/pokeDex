package main

import (
	"os"
	"os/exec"
)

func displaySprite(url string) error {
	if url == "" {
		return nil
	}

	curl := exec.Command("curl", "-s", url)

	chafa := exec.Command(
		"chafa",
		// "--fill=block",
		// "--colors=none",
		// "96x96",
		"-",
	)

	pipe, err := curl.StdoutPipe()
	if err != nil {
		return err
	}

	chafa.Stdin = pipe
	chafa.Stdout = os.Stdout
	chafa.Stderr = os.Stderr

	if err := chafa.Start(); err != nil {
		return err
	}

	if err := curl.Start(); err != nil {
		return err
	}

	if err := curl.Wait(); err != nil {
		return err
	}

	return chafa.Wait()

}
