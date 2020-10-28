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

## Copyright and License Information

Unless otherwise specified, all content, including all source code files and documentation files in this repository are:

**Copyright (c) 2020 Vodacom Mozambique**

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
