# Gen-NAPI (WIP)

** Not even remotely prod ready! (see below for details) **

Gen-NAPI is under development, and is not yet ready for production use. It is currently only tested against `test/TensorBase.h` and used to test NAPI performance for the [shumai](https://github.com/facebookresearch/shumai) NAPI bindings.

Longer term, `Gen-NAPI` will aim to provide an easy, opinionated, means of generating high performance NAPI bindings for existing C++ libraries.

## How it works

Using `go-tree-sitter`, we're able to parse a C++ Header file, query the `tree-sitter` syntax tree for relevant logic (e.g. Function, Class, Enum, etc... declarations), and extract the argument (overloads a WIP)/return types into data structures that provide the rough information from which we're able to generate `Node-API` bindings via `node-addon-api`. The library will provide a means to manually override logic when needed; although, obviously, the goal is that the logic will be opinionated to the point of automatically handling so that overrides are used for performance vs necessary to compile w `CMake`.

## TODO: Documentation

_Docs are pending testing against other header libraries, removing the one(?) hardcoded ref to a shumai global function (row/col major inversion), and general clean up of the logic/config setup._
