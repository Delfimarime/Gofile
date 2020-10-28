# gofile/cli 
> A standalone cli toolkit that enables file upload 

## Build

To build **gofile** you need Golang runtime version 1.15.1+

Clone the repository to your local machine 

```
git clone https://github.com/VodacomMZ/Gofile.git
```

Once on the source directory , get the dependencies to you machine 
```
$ go get -u "github.com/thatisuday/commando"
```

Build the binary
```
$ go build -o ./bin/gofile ./src
```

## Available Commands

Enact supports several commands, For help on individual commands, add `--help` following the command name. The commands are:

### `gofile version`
Prints gofile version

### `gofile`
Prints gofile runtime information

### `gofile upload` 
Uploads a file to an HTTP Server
```
Usage:
   gofile <url> <file> {flags}

Arguments: 
   url                           server endpoint
   file                          file or directory (whose files are) intend to be uploaded (default: ./)

Flags: 
   -a, --attribute               form attribute to which the files are mapped into on the multipart/form-data HTTP Request (default: .)
   -c, --compact                 determines if the files are to be sent on a single HTTP Request (must be used in combination with --attribute option) (default: false)
   -h, --help                    displays usage information of the application or a command (default: false)
   -p, --password                password for basic authentication (in case the servers supports it)
   -s, --strategy                the post upload validation criteria to determine whether the operation was successful or not. EVERY_FILE (only successful when every file is correctly uploaded); AT_LEAST_ONE(successful whenever at least one (1) file is correctly uploaded) (default: AT_LEAST_ONE)
   -t, --timeout                 timeout for the HTTP Request (default: 10)
   -u, --username                username for basic authentication (in case the servers supports it)
   -V, --verbose                 whether the command should run in verbose mode or not, true for verbose (default: false)
 ```

## Copyright and License Information

Unless otherwise specified, all content, including all source code files and documentation files in this repository are:

**Copyright (c) 2020 Vodacom Mozambique**

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
