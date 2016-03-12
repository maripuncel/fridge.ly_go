# fridge.ly

## What we're using

* Migrations: https://github.com/mattes/migrate
* npm: dependency management
* Webpack/Babel: bundling
* React


## Client

The main function is declared in `index.js`.

The client code will be compiled into `static/`. Use this command to compile:
```
./node_modules/.bin/webpack -d
```

## Server

To compile the server:
```
cd cmd/fridge.ly/
go install
```
