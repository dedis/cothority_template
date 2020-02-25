# How to run this example

In this example, you are going to run a webpage that lets you interact with the
"keyvalue" contract defined in `cothority_template/byzcoin/keyvalue.go`.

Note: you should have go installed before, see https://golang.org/doc/install.

Clone this repo:

```bash
git clone https://github.com/dedis/cothority_template.git
```

Create a local distributed network of nodes:

```bash
# go to the conode folder
cd cothority_template/conode
# build the conode executable
go build
# run the nodes
./run_nodes.sh -d tmp -v 3
# now the nodes are runing in your terminal, open a new terminal to continue
```

Build the javascript app

```bash
cd cothority_template/external/js/
# install the dependencies
npm install
# bundle the javascript (can take a while)
npm run bundle
```

Now you can open in your browser the `cothority_template/external/js/index.html`
file.

If everything went well, you should be able to perform step `1.` by loading the
file in `cothority_template/conode/tmp/public.toml`. Then you can display the
roster info by doing step `2.`. To do the rest of the example you must create a
new BycCoin with the `bcadmin`
[utiliy](https://github.com/dedis/cothority/tree/master/byzcoin/bcadmin).

# How to use the cothority JS library

In this example, we are using `npm` as the package manager, `typescript` as the
scripting language, `webpack` as the bundler, and `babel` as the ES6 transpiler.

## Set up a new package

In this step we initialize a new node package and add the needed dependencies.

```bash
# Create the package.json
$ npm init
> set up the desired package configurations
# Add dependencies
# Note: -D is used to specify a dev dependency
$ npm i @dedis/cothority
$ npm i -D typescript
$ npm i -D webpack webpack-cli ts-loader
$ npm i -D babel-loader @babel/core @babel/preset-env

$ npm i -D @types/long  
$ npm i -D @types/jasmine 
$ npm i -D jasmine  
$ npm i -D dockerode   
$ npm i -D @types/dockerode  
$ npm i -D nyc  
$ npm i -D ts-node
$ npm i -D tslint
$ npm i -D jasmine-console-reporter
```

## Add the webpack configuration

Copy the existing `webpack.config.js` at the root of this project.

This file is telling webpack to start from `src/index.ts` and look for
typescript and js files. Typescript files are transpiled to ES6, while ES6 files
are transpiled to ES5 with babel. Then everything is bundled in the
`dis/bundle.min.js`.

To be able to easily bundle your app, we add a new command `bundle` accessible
from the `npm run COMMAND` cli. In this case, we will bundle our app with the
command `npm run bundle`.

Add the following in `package.json`:

```json
"scripts": {
    "test": "echo \"Error: no test specified\" && exit 1",
    "bundle": "webpack"
},
```

## Add the tsconfig

Copy the `tsconfig.json` file from the root of this project.

## Add your index.html

The `index.html` will only load the bundled js file in `dis/bundle.min.js`.

The basic scaffold is the following:

```html
<!-- index.html -->
<html>
  <head>
    <meta charset="UTF-8">
    <script src="dist/bundle.min.js" type="text/javascript"></script>
  </head>
  <body>
    <h1>Hi</h1>
  </body>
</html>
```

you will want to look directly in the `index.html` to see a more complete
example.

## Add your index.ts

`src/index.tx` is the entry point of your javascript, it is responsible for loading
all the necessary js/ts files. Here is the minimum code to load the cothority
library.

```ts
// src/index.ts
import * as Cothority from "@dedis/cothority";

export {
    Cothority
};

export function sayHi() {
    console.log("hi from ts");
}
```

you will want to look directly in the `src/index.ts` to see a more complete
example.

## Run your app

To use your app, you must first install all the dependencies, then bundle the
javascript. Here is how:

```bash
npm install
npm run bundle
```

You should then be able to open `index.html` without errors on the javascript
console.

## Running the test

First build the docker image, from the root of this repo run

```bash
make docker
```

Then, from the current location, run

```bash
npm run test
```

## Add custom messages

In our keyValue example, we used the protobuf definition of the KeyValueData
stored in `external/proto/keyvalue.proto`. To use it, you must first compile it
with the custom command `npm run protobuf`, which will generate the protobuf
definition in `external/js/src/protobuf/models.json`. Then, in the custom
message class that we can find in `external/js/src/keyval.ts`, we must extend
the current models of cothority with our newly generated one with the following:

```js
import models from "./protobuf/models.json";

addJSON(models)
```

This allows us to register custom message:

```js
KeyValue.register()
KeyValueData.register()
```
