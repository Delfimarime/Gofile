package main

func main() {
	engine := GoEngine{}
	engine.SetSender(&PublisherImpl{})
	engine.SetDiscoveryClient(&BasicDiscoveryClient{})

	engine.Run(Configuration{
		File:      "/Users/delfimarime/Desktop/github/test",
		Strategy:  AtLeastOne,
		Endpoint:  "http://localhost:3000",
		Verbose: true,
		Attribute: "file",
		Compact: true,
	})

}
