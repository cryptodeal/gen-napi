# Gen-NAPI (WIP)

**Not even remotely prod ready!**

Gen-NAPI is under development, and is not yet ready for production use. It is currently only tested against `test/TensorBase.h` and used to test NAPI performance for the [shumai](https://github.com/facebookresearch/shumai) NAPI bindings.

Longer term, `Gen-NAPI` will aim to provide an easy, opinionated, means of generating high performance NAPI bindings for existing C++ libraries.

## How it works

Using `go-tree-sitter`, we're able to parse a C++ Header file, query the `tree-sitter` syntax tree for relevant logic (e.g. Function, Class, Enum, etc... declarations), and extract the argument (overloads a WIP)/return types into data structures that provide the rough information from which we're able to generate `Node-API` bindings via `node-addon-api`. The library will provide a means to manually override logic when needed; although, obviously, the goal is that the logic will be opinionated to the point of automatically handling so that overrides are used for performance vs necessary to compile w `CMake`.

## TODO: Documentation

_Docs are pending testing against other header libraries, removing the one(?) hardcoded ref to a shumai global function (row/col major inversion), and general clean up of the logic/config setup._

## Demo

**Prerequisite**

- `go` (see Golang [installation instructions](https://go.dev/doc/install))
- `Node.js` / `npm` (see [Node.js installation instructions](https://nodejs.org/en/download/))
- Alternatively, you can use [Bun](https://bun.sh) (see [Bun installation instructions](https://bun.sh))

Clone the repository (use branch: `external_main` at present for latest changes):

```sh
git clone git@github.com:cryptodeal/gen-napi.git
```

Install the `gen-napi` CLI tool:

```sh
go install github.com/cryptodeal/gen-napi@external_main
```

We parse `test/TensorBase.h` (specified in `gen_napi.yaml`) for the demo as it defines the `fl::Tensor` class that shumai wraps (s/o to the amazing [Flashlight](https://github.com/flashlight/flashlight) library; highly recommend checking it out if you're interested Machine Learning in C++).

To verify that the logic is being generated from scratch, run the following (deletes the generated bindings that are pushed to the repo for demonstrative purposes):

```sh
rm -rf test/out && rm -rf test/ts/index.ts
```

Don't worry, the next thing `gen-napi` will re-generate the bindings and the Typescript file containing TS strongly typed function wrappers.

To generate the bindings & associated TS function wrappers, run the following from the root directory of local copy of `gen-napi`:

```sh
gen-napi generate
```

To build, run:

```sh
npm install
```

Or, to use [Bun](https://bun.sh):

```sh
bun install
```
