Smalldoc
=========
Markdown documentation manager written in Go (and React)

**How to install?**

- Install [Go](https://golang.org/doc/install) and setup your go [workspace](https://golang.org/doc/code.html)
- Install [bazaar](http://wiki.bazaar.canonical.com/Download) as some of go packages use `bzr`
- Install [MongoDB](http://docs.mongodb.org/manual/installation/) and start mongo daemon
- [Node](https://github.com/joyent/node/wiki/installation) and npm installtion
- Install `bower` and `react-tools` using `npm`
```bash
npm install -g bower
npm install -g react-tools
```
- Finally, build and start
```bash
cd $GO_WORKSPACE/src/github.com/loconsolutions/smalldocs
make install
make build
make start
```
