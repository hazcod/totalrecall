# totalrecall-go

A Go SDK and commandline utility to abuse the latest Windows Copilot+ Recall feature.<br/>
This was inspired by [Kevin Beaumonts excellent blog article](https://doublepulsar.com/recall-stealing-everything-youve-ever-typed-or-viewed-on-your-own-windows-pc-is-now-possible-da3e12e9465e).

This will extract any Recall extracts which contains the following information:
- Timestamp of the extract
- Window title
- Window token
- Screenshot contents
- Web traffic information (domain, URL)

This does not need SYSTEM privileges on Windows, as it's just reading data in %APPDATA%.
Current user privileges are sufficient, however the SDK supports other user accounts if you are indeed SYSTEM.

## Usage

First make sure you have a Windows CoPilot+ device on ARM64 or enable it using [Amperage](https://github.com/thebookisclosed/AmperageKit).

Either use the CLI utility:

```shell
./totalrecall -log=info
```

Or use the SDK:
```go
package main

import (
	"log"
	"os"
	recallPkg "github.com/hazcod/totalrecall/pkg/recall"
)

func main() {
	recall, err := recallPkg.New(nil) // or set a Logrus.Logger
	if err != nil { log.Fatal(err) }

	// find any OCRd images
	extracts, err := recallPkg.ExtractImagesForCurrentUser()
	if err != nil {
		log.Printf("could not extract Recall Images: %w", err)
		os.Exit(1)
	}

	for i, extract := range extracts {
		log.Printf("%d - %s - %s - %s", i+1, extract.Timestamp, extract.WindowTitle, extract.WindowToken)
	}
}
```

## Documentation

See the SDK documentation on [pkg.go.dev](https://pkg.go.dev/github.com/hazcod/totalrecall).