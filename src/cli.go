package main

import (
	"fmt"
	"github.com/thatisuday/commando"
)

const (
	ENDPOINT  string = "server endpoint"
	TIMEOUT   string = "timeout for the HTTP Request"
	FILE      string = "file or directory (whose files are) intend to be uploaded"
	USERNAME  string = "username for basic authentication (in case the servers supports it)"
	PASSWORD  string = "password for basic authentication (in case the servers supports it)"
	VERBOSE   string = "whether the command should run in verbose mode or not, true for verbose"
	ATTRIBUTE string = "form attribute to which the files are mapped into on the multipart/form-data HTTP Request"
	COMPACT   string = "determines if the files are to be sent on a single HTTP Request (must be used in combination with --attribute option)"
	STRATEGY  string = "the post upload validation criteria to determine whether the operation was successful or not. EVERY_FILE (only successful when every file is correctly uploaded); AT_LEAST_ONE(successful whenever at least one (1) file is correctly uploaded)"
)

func main() {

	engine := GoEngine{}
	engine.SetSender(&PublisherImpl{})
	engine.SetDiscoveryClient(&BasicDiscoveryClient{})

	var appName = "gofile"
	var appVersion = "1.0.0-RELEASE"
	var description = "Uploads a specific file or directory"

	commando.
		SetExecutableName(appName).SetVersion(appVersion).
		SetDescription("This cli uploads a file or files(within a directory recursively)")

	commando.Register(nil).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			fmt.Println(appName + " v" + appVersion)
		})

	commando.Register("upload").
		AddArgument("url", ENDPOINT, "").
		AddArgument("file", FILE, "./").
		AddFlag("compact,c", COMPACT, commando.Bool, false).
		AddFlag("verbose,v", VERBOSE, commando.Bool, false).
		AddFlag("username,u", USERNAME, commando.String, ".").
		AddFlag("password,p", PASSWORD, commando.String, ".").
		AddFlag("strategy,s", STRATEGY, commando.String, AtLeastOne).
		AddFlag("attribute,attr", ATTRIBUTE, commando.String, ".").
		AddFlag("timeout,t", TIMEOUT, commando.Int, 10).
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

			var file = args["file"].Value
			var endpoint = args["url"].Value
			var timeout = flags["timeout"].Value.(int)
			var isCompact = flags["compact"].Value.(bool)
			var isVerbose = flags["verbose"].Value.(bool)
			var username = flags["username"].Value.(string)
			var password = flags["password"].Value.(string)
			var strategy = flags["strategy"].Value.(string)
			var attribute = flags["attribute"].Value.(string)

			if timeout < 0 {
				timeout = 10
			}

			if username == "." && password == "." {
				username = ""
				password = ""
			}

			if attribute == "." {
				attribute = ""
			}

			engine.Run(Configuration{
				File:      file,
				Timeout:   timeout,
				Username:  username,
				Password:  password,
				Strategy:  strategy,
				Endpoint:  endpoint,
				Compact:   isCompact,
				Attribute: attribute,
				Verbose:   isVerbose,
			})

		}).SetShortDescription(description).
		SetDescription(description+"into a specific HTTP Server Endpoint")
	// parse command-line arguments
	commando.Parse(nil)
}
