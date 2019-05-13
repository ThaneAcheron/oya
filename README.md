---
layout: docs
permalink: /documentation/
---
# Oya

## Usage

Install oya and its dependencies:

    $ curl https://oya.sh/get | bash

Initialize project.

    $ oya init

Define a task you can run:

    # ./Oyafile
    build: echo "Hello, world"
    
View tasks

    $ oya tasks

Run the task:

    $ oya run build
    Hello, world

The task in the above example is called "build" but there's nothing special about the name. In fact, a task name is any camel case identifier as long as it starts with a lower-case letter. You can have as many different tasks as you like.

# Concept

- **Oyafile -** is a yaml formatted fileontaining Oya configuration and task definitions.
- **Oya project -** is a directory with Oyafile inside and `Project: name` defined. Project is a set of tasks and files.
- **Oya task -** tasks are bash scripts defined in Oyafiles under name like `task: |`.

# Install Oya


    $ curl https://oya.sh/get | bash

This will install latest version of oya. It’s also possible to specify which version should be installed 

`$ curl https://oya.sh/get | bash -s v0.0.7`.
`$ oya` `--``version`


# Oya Project

You can create Oyafile by hand or with init command.


    $ oya init
    $ cat Oyafile
    Project: project

Don’t forget to change project name.


# Tasks

Oya task is a bash script defined as a Oyafile key.
Tasks in Oyafile are defined as a yaml keys, with pipe at the beginning line and bash code in following.


    $ cat Oyafile
    Project: OyaExample
    
    build: |
      go build app.go
    
    start: |
      ./server

To list available tasks use tasks command:


    $ oya tasks
    # in ./Oyafile
    oya run build
    oya run start

Execute task

    oya run build
    oya run start
    
# *.oya files

Inside a project you can have many files with named `*.oya` they will be read as a Value files, expected syntax is a pair of  `key: value`. 


    $ cat values.oya
    fruit: banana


    $ cat Oyafile
    Project: OyaExample
    
    eat: |
      echo ${Oya[fruit]}

`$ oya run eat`


# Values

You can also store all values inside Oyafile.


    $ cat Oyafile
    Project: OyaExample
    
    Values:
      fruit: banana
    
    eat: |
      echo ${Oya[fruit]}


# Secrets

You can also store confidential data using Oya secrets.
Oya uses SOPS (https://github.com/mozilla/sops) to help with secrets management. 
First you need to configure SOPS for encryption method, check https://github.com/mozilla/sops/blob/master/README.rst#usage.
For our example we can use PGP key.

    $ export SOPS_PGP_FP="317D 6971 DD80 4501 A6B8  65B9 0F1F D46E 2E8C 7202"

Oya secrets commands:

    $ oya secrets --help
    ...
      edit        Edit secrets file
      encrypt     Encrypt secrets file
      view        View secrets
    ...

**First run:**

First you need to create `secrets.oya` file. with `key: value` in each line. And encrypt it. (You can also go straight to edit with `$ oya secrets edit secrets.oya`).


    $ cat secrets.oya
    magical_spell: hokus pokus czary mary
    $ oya secrets encrypt secrets.oya

Done your secrets are safe. Check how secrets.oya looks like.


    $ cat secrets.oya
    {
            "data": "ENC[AES256_GCM,data:XXXX=,tag:XXXX==,type:str]",
            "sops": {
                    ...
                    "pgp": [...],
                    ...
            }
    }%

We can see only sops metadata, our data are safe and encrypted.

To view or edit use:


    $ oya secrets view secrets.oya
    magical_spell: hokus pokus czary mary
    $ oya secrets edit secrets.oya

You can access secret value with  `${Oya[magical_spell]}` from task.


    $ cat Oyafile
    Project: OyaExample
    ...
    spell: |
      echo ${Oya[magical_spell]}
    
    $ oya run spell
    hokus pokus czary mary


# Render

Oya can also render a template files or even whole directory. Oya uses Plush templating system You can find out more here https://github.com/gobuffalo/plush.


    $ cat template/index.html
    <h1><%= title %></h1>

Let’s add our title under  `Value:`


    $ cat Oyafile
    Project: OyaExample
    
    Values:
      title: Hello from Oya
    ...

Render output into `public/` dir so our server can see it.


    $ oya render template/index.html
    $ cat index.html
    <h1>Hello from Oya</h1>


# Recursive oya

It’s possible to organize oya project into directories with separated logic. 

We can separate it into `backend/` and `frontend/`. 


    $ tree
    .
    ├── Oyafile
    ├── backend
    │   ├── Oyafile
    └── frontend
        ├── Oyafile
    
    2 directories, 3 files

Each of them will have own `Oyafile`, and thanks to recursive tasks each file can have task with the same name in our case `build`.

## Backend

Here is how our backend Oyafile looks like (note that there is no `Project:` for subdirectories):


    $ cat ./backend/Oyafile
    build: |
      echo "Compiling server"
      go build -o ../build/server app.go

## Frontend


    $ cat ./frontend/Oyafile
    Values:
      title: Hello from Oya
    
    build: |
      echo "Rednering template"
      oya render template/index.html -o ../build/public

Our frontend holds only template file and Values necessary to render it.

## Project Oyafile


    $ cat Oyafile
    Project: OyaExample
    
    build: |
      echo "Preparing build/ folder"
      rm -rf build/ && mkdir -p build/public

## Recursive run

Now let’s see what tasks we have. To do it for whole project including subdirectories we need to use `-r or --recurse`  flag.


    $ oya tasks -r
    # in ./Oyafile
    oya run build
    
    # in ./backend/Oyafile
    oya run build
    
    # in ./frontend/Oyafile
    oya run build

As you can see we have three `build` tasks one per Oyafile. We can now run them all.


    $ oya run -r build
    Preparing build/ folder
    Compiling server
    Rendering template.html

And now we can start the app with `$ oya run start`.

# Packs

Pack is a Oya project with general purpose tasks which can be easily shared and used inside other projects. Oya installs pack in your home `~/.oya` directory. Each time you rune oya command dependencies will be resolved and installed.


    $ oya import github.com/tooploox/oya-packs/docker
    $ cat Oyafile
    Project: OyaExample
    Require:
      github.com/tooploox/oya-packs/docker: v0.0.6
    Import:
      docker: github.com/tooploox/oya-packs/docker
    ...

Import will add importing pack under `Import:`, key of imported pack is his alias and can be accessed by this name, (you can change it if needed). 

## Packs versioning

Import will automatically resolve dependencies with newest versions and add them under `Require:`.

Imported pack added bunch of new commands into our project


    $ oya tasks
    # in ./Oyafile
    oya run build
    oya run docker.build
    oya run docker.generate
    oya run docker.run
    oya run docker.stop
    oya run docker.version
    oya run start

We can easily generate Dockerfile, update it, and build our project.


# Pack development - sharing oya’s

Each oya project is a pack itself, all you need to do is push it to git and tag it version `name/v0.0.1`. Import it as before with `$ oya import github.com/tooploox/oya-packs/name` oya should automatically resolve newest version and add `Require` to your Oyafile.


# Tests PGP keys

To have all tests passing successfull it's require to have our pgp key for secrets

    $ gpg --import testutil/pgp/private.rsa
    
# Contributing

1. Install go 1.11 (goenv is recommended, example: `goenv install 1.11.4`).
2. Checkout oya outside GOHOME.
3. Install godog: `go get -u github.com/DATA-DOG/godog/cmd/godog`.
4. Run acceptance tests: `godog`.
5. Run tests: `go test ./...`.
6. Run Oya: `go run oya.go`.

