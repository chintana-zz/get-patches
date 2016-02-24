package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
)

// This programs connects to patch SVN, fetch all patches in
// patch-list.txt and put them into a folder named zip.
//
// Why do this?
//
// A given product will have the set of patches in repository/components/patches.
// Patches in the above folder just contain the patch directory and set of JAR files
// README + related configurations are in the original patch file. This script will
// get those original patch archives

const (
	// PATCHES_DIR should contain the correct SVN URL containing patches folder
	PATCHES_DIR = ".../carbon/turing/patches/"
	PATCH_LIST  = "patch-list.txt"
)

var (
	username = ""
	password = ""
)

func main() {
	username, password = getCredentials()
	patches := getPatchList()
	c := make(chan string)

	os.Mkdir("zip", os.ModePerm)
	for _, v := range patches {
		go downloadPatch(v, c)
	}

	for _ = range patches {
		x := <-c
		fmt.Println(x)
	}
}

// Get the full name of the patch ending in .zip
func getPatchName(patchId string) string {
	transport := &http.Transport{}
	client := http.Client{
		Transport: transport,
	}
	defer transport.CloseIdleConnections()

	req, _ := http.NewRequest("GET", PATCHES_DIR+patchId+"/", nil)
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			panic(z.Err())
		case tt == html.StartTagToken:
			t := z.Token()

			if isLink := t.Data == "a"; isLink {
				for _, a := range t.Attr {
					if a.Key == "href" && strings.HasSuffix(a.Val, ".zip") {
						return a.Val
					}
				}
			}
		}
	}
}

func downloadPatch(patchId string, c chan string) {
	fmt.Println("Downloading patch", patchId)
	p := getPatchName(patchId)

	out, err := os.Create("zip/" + p)
	defer out.Close()
	if err != nil {
		panic(err)
	}

	transport := &http.Transport{}
	client := http.Client{
		Transport: transport,
	}
	defer transport.CloseIdleConnections()

	req, _ := http.NewRequest("GET", PATCHES_DIR+patchId+"/"+p, nil)
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	c <- p
}

// Read the patch list file and return what's in the file in a slice
func getPatchList() []string {
	f, err := os.Open(PATCH_LIST)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	patches := make([]string, 0)
	reader := bufio.NewReader(f)
	for {
		patch, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		patch = strings.TrimSpace(patch)
		if strings.Contains(patch, "patch") {
			patches = append(patches, patch)
		}
	}
	return patches
}

// Read username password for connecting to private SVN containing patches
func getCredentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Password: ")
	bytePwd, _ := terminal.ReadPassword(0)
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(string(bytePwd))
}
