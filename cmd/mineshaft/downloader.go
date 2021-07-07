package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadFile(ctx context.Context, url, destinationPath string) error {
	Logger.Printf(`Download -> %s`, url)

	headResp, err := http.Head(url)
	if err != nil {
		return err
	}
	if headResp.StatusCode < 200 || headResp.StatusCode >= 300 {
		return fmt.Errorf("http error response %s", headResp.Status)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	out, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}
