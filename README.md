# README get-patches.go #

List of patches should be there in patch-list.txt before running. Patches should be separated by newlines. Last patch should be followed by a new line.

Script download and copy patches to zip folder

Make sure you have following,

$ go get golang.org/x/crypto/ssh/terminal
$ go get golang.org/x/net/html

Then do

$ go run get-patches.go
