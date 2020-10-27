package main

import (
	"fmt"
	"github.com/thatisuday/commando"
)

const (
	ENDPOINT  string = "server url"
	COMPACT   string = "whether the upload is to be single request or not"
	FILE      string = "file or directory (whose files are) intend for upload"
	ATTRIBUTE string = "form attribute to which the files are to be provided in"
	USERNAME  string = "username for basic authentication in case the servers supports it"
	PASSWORD  string = "password for basic authentication in case the servers supports it"
	VERBOSE   string = "whether the command should run in verbose mode or not, true for verbose"
	STRATEGY  string = "strategy to determine whether the operation was successful.\nEVERY_FILE (only successful when every file is correctly uploaded)\nAT_LEAST_ONE(successful whenever at least one (1) file is correctly uploaded) "
)

func main() {

	engine := GoEngine{}
	engine.SetSender(&PublisherImpl{})
	engine.SetTransformer(&DefaultTransformer{})
	engine.SetDiscoveryClient(&BasicDiscoveryClient{})

	var appName = "gofile"
	var appVersion = "1.0.0-RELEASE"

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
		AddFlag("attribute,attr", ATTRIBUTE, commando.String, "file").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

			var file = args["file"].Value
			var endpoint = args["url"].Value
			var isCompact = flags["compact"].Value.(bool)
			var isVerbose = flags["verbose"].Value.(bool)
			var username = flags["username"].Value.(string)
			var password = flags["password"].Value.(string)
			var strategy = flags["strategy"].Value.(string)
			var attribute = flags["attribute"].Value.(string)

			if username == "." && password == "." {
				username = ""
				password = ""
			}

			engine.Run(Configuration{
				File:      file,
				Username:  username,
				Password:  password,
				Strategy:  strategy,
				Endpoint:  endpoint,
				Verbose:   isVerbose,
				Compact:   isCompact,
				Attribute: attribute,
			})

		})

	// parse command-line arguments
	commando.Parse(nil)
}
