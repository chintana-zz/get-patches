# get-patches.go #

List of patches should be there in patch-list.txt before running. Patches should be separated by newlines. Last patch should be followed by a new line.

Script download and copy patches to zip folder

Make sure you have following,

<code>$ go get golang.org/x/crypto/ssh/terminal</code><br/>
<code>$ go get golang.org/x/net/html</code>

Then do

$ go run get-patches.go
