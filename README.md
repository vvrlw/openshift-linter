OpenShift Linter
================

This is a utility for OpenShift users/admins who want to know if certain (very basic) rules have been followed. You can also specify naming conventions for namespaces (i.e. projects), names, containers and environment variables.

As it's very early days the focus is on `DeploymentConfig` objects.

<img src="screenshots/openshift-linter.png" width="512" alt="Screenshot of the OpenShift Linter GUI"/>

**Fig. 1** OpenShift Linter GUI

Usage
-----
```
Usage: ./openshift-linter [<JSON file> [<JSON file>]]
  -container string
    	pattern for containers (default "^[a-z0-9_-]+$")
  -env string
    	pattern for environment variables (default "^[A-Z0-9_-]+$")
  -n string
    	hostname (default "localhost")
  -name string
    	pattern for names (default "^[a-z0-9_-]+$")
  -namespace string
    	pattern for namespaces/projects (default "^[a-z0-9_-]*$")
  -p int
    	listen on port (default 8000)
Commands:
  list	Print list of available checks
```

The two main use cases are:

* You already have a bunch of configuration files (the output of `oc export dc --all-namespaces`, say, assuming you're lucky enough to be `cluster-admin`)
```
$ ./openshift-linter i-contain-multitudes.json
```

* Run `./openshift-linter` and open the GUI at `http://localhost:8000/openshift-linter/report` (configure hostname and port using the -n and -p switches, respectively)

When setting naming conventions for namespaces, names, containers and environment variables, be sure to use anchors to describe the string as a whole.

Listing
-------
To print a list of the available linter items, enter:
```
$ ./openshift-linter list
health
image pull policy
invalid key
invalid name
limits
security
similar key
```

Build
-----
Install Go using one of the installers available from `https://golang.org/dl/` and set up your `$GOPATH` and `$GOBIN` as you see fit (exporting `GOPATH=~/golang` and `GOBIN=$GOPATH/bin` in your `.bash_profile` will do).

Then clone `github.com/gerald1248/openshift-linter`. The folder structure below `$GOPATH` should look roughly as follows:
```
src
└── github.com
    └── gerald1248
        └── openshift-linter
            ├── LICENSE
            ├── README.md
            ├── bindata.go
            ├── bower.json
            ├── bower_components
            ├── contributors.txt
            ├── data
            ├── gulpfile.js
            ├── item-env.go
            ├── item-health.go
            ├── item-image-pull-policy.go
            ├── item-limits.go
            ├── item-pattern.go
            ├── item-security.go
            ├── items.go
            ├── makelist.go
            ├── openshift-linter.go
            ├── package
            ├── package.json
            ├── preflight.go
            ├── preflight_test.go
            ├── preprocess.go
            ├── process.go
            ├── screenshots
            ├── server.go
            ├── src
            ├── static
            ├── summary.go
            ├── types.go
            └── types_test.go
```

Next, install Node.js with npm using your package manager. `cd` into the working directory `openshift-linter` and enter:

```
$ sudo npm install -g gulp-cli
$ npm install
```

Note for Ubuntu users: as `gulp-cli` currently expects `node`, but Ubuntu installs `nodejs`, `gulp` has to be triggered as follows:

```
$ nodejs node_modules/gulp/bin/gulp.js
```

In other words, it's very nearly the invocation to use when installing `gulp-cli` globally is not possible or desirable:

```
$ node node_modules/gulp/bin/gulp.js
```

Before running `gulp` (which builds and tests the program), fetch and install the dependencies (`go get` also runs at build time, albeit without the -u switch):

```
$ go get -u github.com/jteeuwen/go-bindata/...
$ go get -u
```

With that, the workspace is ready. The default task (triggered by `gulp`) compiles `openshift-linter` from source, runs (minimal for now) tests, checks the source format, generates a binary in `package` and writes out a distributable zip for your operating system.

You can also run `gulp build`, `gulp test`, `gulp watch`, etc. individually if you wish.

How do I create my own checks?
------------------------------
Add types that conform to the `LinterItem` interface, then register them in `items.go`.

Cross-compile for Windows
-------------------------
To cross-compile Windows binaries on Linux or Mac, enter:
```
$ GOOS=windows GOARCH=amd64 go install
$ gulp build-win32
```
