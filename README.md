---
layout: docs
permalink: /documentation/
---
# Oya

Oya is a command line tool aiming to help you bootstrap and manage deployable
projects.


## Quick start

Install oya and its dependencies:

    $ curl https://oya.sh/get | bash

Initialize a project:

    $ oya init

Add an example task to the bottom of the generated `Oyafile`:

    build: |
      echo "Hello, world"

> A task name is any valid camel-case identifier starting with a lower case
> letter. Identifiers starting with caps are reserved by Oya.

List available tasks:

    $ oya tasks

Run the task:

    $ oya run build
    Hello, world

> If you're familiar with Makefiles, you may have noticed some similarity here.
> The main difference is because we're using standard YAML files, is the pipe
> character after task name. An added bonus you don't have to use tabs. :>

If all Oya offered was poor-man's Makefiles, you'd better find [something
better](http://www.dougmcinnes.com/html-5-asteroids) to do. Fortunately, Oya has
much more to offer so keep reading.


## Key concepts

- **Oyafile -** is a YAML file containing Oya configuration and task
  definitions.
- **Oya project -** is a directory and any number of subdirectories containing
  `Oyafiles`; the top-level `Oyafile` must contain a `Project:` directive.
- **Oya task -** a named bash script you can run using `oya run <task
  name>`.
- **Oya pack -** an installable Oya project containing reusable tasks you can
  easily use in other projects.


## Installation

To install the latest version of Oya run the following command:

    $ curl https://oya.sh/get | bash

> You can also specify which version should be installed
>
> ```
> $ curl https://oya.sh/get | bash -s v0.0.7`.
> oya --version
> ```


## Initializing a project

To get started using Oya in an existing project you need to initialize it by
running the following command in its top-level directory:


    $ oya init

All the command does is generate a file named `Oyafile` that looks like this.

    Project: project

You may want to change project name. We'll change it to `OyaExample`:

    Project: OyaExample


## Creating your first task

Oya task is a named bash script defined in an `Oyafile`. Let's pretend our
project is a Golang HTTP server and we need tasks for building and running the
server. Edit the generated `Oyafile` so it looks like this:

    Project: OyaExample

    build: |
      go build .

    start: |
      go run .

> Notice the pipe characters after task names. This is YAML and the pipe is
> required for multi-line script definitions.

Here's how you can list available tasks:

    $ oya tasks
    # in ./Oyafile
    oya run build
    oya run start

To make it work, let's create a simple server implementation:

```
cat << EOT >> app.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func main() {
	host := flag.String("host", "0.0.0.0", "host name to bind to")
	port := flag.Int("port", 8080, "port number")
	flag.Parse()
	http.HandleFunc("/", handler)
	bind := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Printf("Starting web server on %s\n", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
EOT
```

> Well, to really make it work, you also need the Go language tools
> [installed](https://golang.org/doc/install).


Let's start the server

    $ oya run start

Ok, but does our server work?

    $ curl localhost:8000
    Hello, world!

Success!


## Parametrizing tasks

Alongside the `Oyafile` you can create any number of YAML files with an `.oya`
extension containing constants you can use in your tasks (and in generated
boilerplate as you'll find out later).

Let's put the default port number and host name into an '.oya' file. Create file
named `values.oya` with the following contents:

    port: 4000
    host: localhost

Let's now modify our task definitions so we pass port and host name explicitly
when starting the server:

    Project: OyaExample

    build: |
      go build .

    start: |
      go run . --port ${Oya[port]} --host ${Oya[host]}

> The `${Oya[...]}` syntax is how you access bash associative arrays. Oya comes
> with its own shell implementation aiming to be compatible with Bash 4.

After restarting the server (using `oya run start`) it's reachable on a
different port:

    $ curl localhost:4000
    Hello, world!

## Passing arguments

TODO: Override port/host via flags.


## Storing confidential information

You can also store confidential data right in your projects. Oya uses
[SOPS](<https://github.com/mozilla/sops>) to store them in an encrypted format.

Imagine you want to protect our oh so very secret HTTP endpoint using a password
you need to supply as a parameter, example:

    $ curl localhost:4000?password=badpassword
    Unauthorized

First, configure SOPS for encryption method, check
<https://github.com/mozilla/sops/blob/master/README.rst#usage>.

For our example we can use a sample PGP key:

    $ export SOPS_PGP_FP="317D 6971 DD80 4501 A6B8  65B9 0F1F D46E 2E8C 7202"

Oya secrets commands:

    $ oya secrets --help
    ...
      edit        Edit secrets file
      encrypt     Encrypt secrets file
      view        View secrets
    ...

Let's first slightly modify our HTTP server so it checks if the provided
password matches one in an environment variable. Change `app.go` so it looks
like this:

    package main

    import (
        "flag"
        "fmt"
        "log"
        "net/http"
        "os"
    )

    func handler(w http.ResponseWriter, r *http.Request) {
        requiredPassword := os.Getenv("PASSWORD")
        password, ok := r.URL.Query()["password"]
        if !ok || len(password) != 1 || password[0] != requiredPassword {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        fmt.Fprintf(w, "Hello, world!")
    }

    func main() {
        host := flag.String("host", "0.0.0.0", "host name to bind to")
        port := flag.Int("port", 8080, "port number")
        flag.Parse()
        http.HandleFunc("/", handler)
        bind := fmt.Sprintf("%s:%d", *host, *port)
        fmt.Printf("Starting web server on %s\n", bind)
        log.Fatal(http.ListenAndServe(bind, nil))
    }

Long story short, the server now requires the `PASSWORD` environment variable to
be present. Let's modify our `Oyafile` so it sets that variable:

    Project: OyaExample

    build: |
      go build .

    start: |
      PASSWORD=${password} go run . --port ${Oya[port]} --host ${Oya[host]}

Because we don't want to store the password in the plain, we'll encrypt it.

First you need to create `secrets.oya` file and encrypt it:

    $ cat << EOT >> secrets.oya
    password: hokuspokus
    EOT
    $ oya secrets encrypt secrets.oya

> There's nothing special about the name of the file. You can encrypt any YAML
> file with .oya extension possibly grouping your secrets in larger projects.

Now our precious secret is safe!

    $ cat secrets.oya
    {
            "data": "ENC[AES256_GCM,data:[...]=,tag:[...]==,type:str]",
            "sops": {
                    ...
                    "pgp": [...],
                    ...
            }
    }%

Only SOPS metadata is out in the plain, the password itself is encrypted and safe.

Restart the HTTP server and test it:

    $ curl localhost:4000?password=badpassword
    Unauthorized
    $ curl localhost:4000?password=hokuspokus
    Hello, world!

To view or edit an encrypted file later:


    $ oya secrets view secrets.oya
    password: hokuspokus
    $ oya secrets edit secrets.oya

> You can use your favorite editor by setting the `EDITOR` environment variable.
> The default is vim but you should be able to make it even with [GUI
> editors](https://github.com/mozilla/sops/issues/127).


## Using Oya packs


Pack is an installable Oya project containing reusable tasks you can easily use
in other projects. Oya installs pack in your home `~/.oya` directory.

In this tutorial, let's use the `docker` pack to generate a Dockerfile for the
application:


    $ oya import github.com/tooploox/oya-packs/docker

Import will automatically resolve dependencies with newest versions and add them
to the `Require:` directive. It will make the pack's tasks available under the
`docker` alias by default. You can change the alias to something more suitable
by editing the Oyafile:

    $ cat Oyafile
    Project: OyaExample
    Require:
      github.com/tooploox/oya-packs/docker: v0.0.6
    Import:
      docker: github.com/tooploox/oya-packs/docker
    [...]


Imported pack added bunch of new commands into our project:


    $ oya tasks
    # in ./Oyafile
    oya run build
    oya run docker.build
    oya run docker.generate
    oya run docker.run
    oya run docker.stop
    oya run docker.version
    oya run start

Here's how you can generate the Dockerfile:

    $ oya run docker.generate

This is what what the Dockerfile will look like:

    FROM golang

    COPY . /go/src/app
    WORKDIR /go/src/app

    RUN go get
    RUN go build -o app

    CMD [ "app" ]

You can build the image and start the server in a container:

    $ oya run docker.build
    [...]
    $ oya run docker.run
    Starting web server on 0.0.0.0:8080


## Generating boilerplate

Oya can also render files and even entire directories from templates. Oya uses
[Plush templating engine](https://github.com/gobuffalo/plush).

In an earlier section you were asked to copy & paste a simple web server. For
the sake of illustration imagine that you want to make creating HTTP web servers
easier by generating the boilerplate from a template.

Let's create `app.go` in `templates/` directory:

    mkdir templates
    cat << EOT >> templates/app.go
    package main

    import (
        "flag"
        "fmt"
        "log"
        "net/http"
    )

    func handler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, world!")
    }

    func main() {
        host := flag.String("host", "0.0.0.0", "host name to bind to")
        port := flag.Int("port", 8080, "port number")
        flag.Parse()
        http.HandleFunc("/", handler)
        bind := fmt.Sprintf("%s:%d", *host, *port)
        fmt.Printf("Starting web server on %s\n", bind)
        log.Fatal(http.ListenAndServe(bind, nil))
    }
    EOT

This is how you can render it to the current directory (the command will
override `app.go` so ****be careful!**):

    $ oya render templates/app.go


## Parametrizing boilerplate

TODO: I was sure we had --set implemented for render so we can override values.
Now Flags doesn't seem to work either, I probably forgot something. But setting
values just to generate boilerplate is idiotic as a basic example.

## Reusing your scripts

Ok, all is good and fine but how to make the ^ code available when creating new
projects? Easy, turn it into a pack!

Technically, all you need to do to turn the project we created into an Oya pack
is push it to Github and tag it with a version number by adding & pushing a git
tag in the right format (e.g. `v0.1.0`) along with a few small changes.

Rather than doing that, let's create it step-by-step in a fresh new git repo.

First, let's add a task you'll use to generate boilerplate to the original
Oyafile so it looks like this:

    Project: project

    build: |
      go build .

    start: |
      go run . --port ${Oya[port]} --host ${Oya[host]}

    generate: |
      oya render ${Oya[BasePath]}/templates/app.go

The new `generate` task will generate files into the current directory based on
the contents of the templates directory.

> BasePath is the base directory of the path so the `oya render` command knows
> where to take templates from.

Because the script needs `port` and `host`, let's also create `values.oya`
containing pack defaults:

    port: 8080
    host: localhost

To share your pack all you need is push the project to a Github repository and
tag it with a version number, roughly:

    $ git push origin
    $ git tag v0.1.0
    $ git push --tags

That's it!

> Currently, only Github is supported as a way of sharing packs but if you want
> to help with adding support for Bitbucket and others, do get in touch!

So how do you use the pack? Easy. First, create a new empty project and
initialize it:

    $ oya init

Then import the pack. Here, I'm assuming it's under
github.com/tooploox/oya-gohttp:

    $ oya import github.com/tooploox/oya-gohttp

> TODO: There must be a way to set alias.

That's it! Let's see what tasks we have available:

    $ oya tasks
    # in ./Oyafile
    oya run oya-gohttp.build
    oya run oya-gohttp.generate
    oya run oya-gohttp.start

> TODO: Error: Internal error: values.oya file not found while loading

Let's generate the server source:

    oya run oya-gohttp.generate

This will generate `app.go` file in the current directory:

    package main

    import (
    [...]

Let's start the server:

    $ oya run oya-gohttp.start


## Overriding pack values

> TODO: Overriding port.


## Pack repositories

TODO: Describe multiple packs in a single repo (i.e. how to version and use them).


## Running tasks recursively

So far we only talked about one Oyafile in the project's top-level directory.
But you can put Oyafiles in subdirectories.

This is especially useful in monorepos so we'll use that as an example. Imagine
we have a web application consisting of a REST API back-end server and front-end
SPA application.

For the sake of illustration let's create a mono-repository containing both
back-end and front-end in separate directories.


    $ tree
    .
    ├── Oyafile
    ├── backend
    │   ├── Oyafile
    │   ├── server.go
    └── frontend
        ├── Oyafile
        └── main.ts

    2 directories, 3 files

In addition to the top-level `Oyafile` both front-end and back-end have their
own `Oyafiles`, let's have a look at each.

## Top-level Oyafile

The top level Oyafile contains a `Project:` directive to mark the top-level
project directory and a task named `build` preparing the output directory.

    $ cat ./Oyafile
    Project: myproject

    build: |
      echo "Preparing build/ folder"
      rm -rf build/ && mkdir -p build/public

### Backend

In this example, let's assume that back-end is an HTTP API written in Go. Here
is how the back-end `Oyafile` looks like:


    $ cat ./backend/Oyafile
    build: |
      echo "Compiling server"
      go build -o ../build/server .

Notice that it contains just one task, named `build` for compiling the back-end
server.

### Frontend

Front-end is a TypeScript SPA application. The `build` task compiles TypeScript
source to JavaScript:


    $ cat ./frontend/Oyafile
    build: |
      echo "Compiling front-end"
      tsc main.ts --outFile ../build/public/main.js

### Recursive run

Now let’s list the available tasks. To do it for the whole project including
subdirectories we need to use `-r or --recurse` flag.

    $ oya tasks -r
    # in ./Oyafile
    oya run build

    # in ./backend/Oyafile
    oya run build

    # in ./frontend/Oyafile
    oya run build

As you can see we have three `build` tasks one per Oyafile. We can now run them
all.


    $ oya run -r build
    Preparing build/ folder
    Compiling server
    Compiling front-end

The result is the compiled server as well as the JavaScript it serves.

TODO: For a working example, clone ...


## Contributing

1.  Install go 1.11 (goenv is recommended, example: `goenv install 1.11.4`).
2.  Checkout oya outside GOHOME.
3.  Install godog: `go get -u github.com/DATA-DOG/godog/cmd/godog`.
4.  Run acceptance tests: `godog`.

> For all test to pass you need to import the PGP key used to encrypt secrets.
>    $ gpg --import testutil/pgp/private.rsa

5.  Run tests: `go test ./...`.
6.  Run Oya: `go run oya.go`.
