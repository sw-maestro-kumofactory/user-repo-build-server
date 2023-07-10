package repodownload

import (
	"io"
	"net/http"
	"os"
)

func RepoDownload(filepath string, user string, repo string, branch string) error {
	url := "https://api.github.com/repos/" + user + "/" + repo + "/tarball/" + branch

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
