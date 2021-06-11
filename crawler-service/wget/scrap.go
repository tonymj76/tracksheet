package wget

import (
	"bytes"
	"cache_service/crawler-service/structure"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
)

// Call shell command wget to download all files in url. The reason to use wget is that wget
// works only in Linux.
func WgetLink(url_path, depth string) []string {
	cmd := exec.Command("wget", "-U", "Mozilla/5.0", "-r", "-l", depth, url_path)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse(url_path)
	if err != nil {
		log.Fatal("error parsing url", err)
	}
	if _, err := os.Stat(u.Host); !os.IsNotExist(err) {
		//delete folder wget generates
		os.RemoveAll(u.Host)
	}
	if err != nil {
		fmt.Println("error ", err)
	}
	return structure.ParseSiteMapURL(errb.String(), u.Hostname())
}
