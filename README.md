# Gen-NAPI (WIP)

**Not even remotely prod ready!**

Gen-NAPI is under development, and is not yet ready for production use. It is currently only tested against `test/TensorBase.h` and used to test NAPI performance for the [shumai](https://github.com/facebookresearch/shumai) NAPI bindings.

Longer term, `Gen-NAPI` will aim to provide an easy, opinionated, means of generating high performance NAPI bindings for existing C++ libraries.

## Installation

```shell
go install github.com/cryptodeal/gen-napi@external_main
```

## How it works

Using `go-tree-sitter`, we're able to parse a C++ Header file, query the `tree-sitter` syntax tree for relevant logic (e.g. Function, Class, Enum, etc... declarations), and extract the argument (overloads a WIP)/return types into data structures that provide the rough information from which we're able to generate `Node-API` bindings via `node-addon-api`. The library will provide a means to manually override logic when needed; although, obviously, the goal is that the logic will be opinionated to the point of automatically handling so that overrides are used for performance vs necessary to compile w `CMake`.

## Example

_`gen-napi` config (`gen_napi.yaml`)_

```yaml
packages:
  - path: demo2/BasicLogic.h
    bindings_out_path: demo2/out
    js_wrapper_opts:
      addon_path: ../../build/Release/test.node
      wrapper_out_path: demo2/ts/index.ts
```

_C++ Header input file_

```h
#include <stdint.h>
#include <string>

namespace demo2 {

enum class dtype {
  f16 = 0,  // 16-bit float
  f32 = 1,  // 32-bit float
  f64 = 2,  // 64-bit float
  b8 = 3,   // 8-bit boolean
  s16 = 4,  // 16-bit signed integer
  s32 = 5,  // 32-bit signed integer
  s64 = 6,  // 64-bit signed integer
  u8 = 7,   // 8-bit unsigned integer
  u16 = 8,  // 16-bit unsigned integer
  u32 = 9,  // 32-bit unsigned integer
  u64 = 10  // 64-bit unsigned integer
};

float* baz(float* a, int b);

std::string test(std::string a);
}
```

_Generated `Node-API` bindings (`.cc`)_

```cc
// Code generated by gen-napi. DO NOT EDIT.
#include <napi.h>
#include <atomic>
#include <string>

//////////
// source: BasicLogic.h

// non-exported helpers (omitted for brevity)

// exported functions

static Napi::Value _baz(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 2) {
    Napi::TypeError::New(env, "`baz` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`baz` expects args[0] to be typeof `Float32Array`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  float* a = info[0].As<Napi::TypedArrayOf<float>>().Data();
  if (!info[1].IsNumber()) {
    Napi::TypeError::New(env, "`baz` expects args[1] to be typeof `number`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  int b = static_cast<int>(info[1].As<Napi::Number>().Int64Value());
  float* _res;
  _res = demo2::baz(a, b);
  size_t _res_byte_len = sizeof(_res);
  size_t _res_elem_len = _res_byte_len / sizeof(*_res);
  std::unique_ptr<std::vector<float>> _res_native_array =
      std::make_unique<std::vector<float>>(_res, _res + _res_elem_len);
  Napi::ArrayBuffer _res_arraybuffer =
      Napi::ArrayBuffer::New(env, _res_native_array->data(), _res_byte_len,
                             DeleteArrayBuffer<float>, _res_native_array.get());
  _res_native_array.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, _res_byte_len);
  return Napi::TypedArrayOf<float>::New(env, _res_elem_len, _res_arraybuffer, 0,
                                        napi_float32_array);
}

static Napi::Value _test(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 1) {
    Napi::TypeError::New(env, "`test` expects exactly 1 arg")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsString()) {
    Napi::TypeError::New(env, "`test` expects args[0] to be typeof `string`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  std::string a = info[0].As<Napi::String>().Utf8Value();
  std::string _res;
  _res = demo2::test(a);
  return Napi::String::New(env, _res);
}

// NAPI exports

Napi::Object Init(Napi::Env env, Napi::Object exports) {
  exports.Set(Napi::String::New(env, "_baz"), Napi::Function::New(env, _baz));
  exports.Set(Napi::String::New(env, "_test"), Napi::Function::New(env, _test));
  return exports;
}

NODE_API_MODULE(addon, Init)
```

_Generated TS wrapper_

```ts
// Code generated by gen-napi. DO NOT EDIT.
/* eslint-disable */
const { _baz, _test } = import.meta.require('../../../../build/Release/test.node');

export enum dtype {
  f16 = 0,
  f32 = 1,
  f64 = 2,
  b8 = 3,
  s16 = 4,
  s32 = 5,
  s64 = 6,
  u8 = 7,
  u16 = 8,
  u32 = 9,
  u64 = 10
}

export const baz = (a: Float32Array, b: number): Float32Array => {
  return _baz(a, b);
};

export const test = (a: string): string => {
  return _test(a);
};
```

# TODO: Documentation

_Docs are a WIP; pending testing against other header libraries and general clean up of the logic/config setup; expect breaking changes as we work to flesh out the best opinionated setup for ease of use/quick start up time._

## Usage

### Option A: CLI (recommended)

Create a file `gen-napi.yaml` in which you specify the header file to parse/generate bindings from (as well as additional type handlers, output path, etc.). For instance:

```yaml
packages:
  - path: ../flashlight/flashlight/fl/tensor/TensorBase.h
    lib_root_dir: ../flashlight
    bindings_out_path: demo1/out
    indent: "\t"
    js_wrapper_opts:
      front_matter: |
        import { Tensor } from './tensor'
      addon_path: ../../build/Release/shumai_bindings.node
      wrapper_out_path: demo1/ts/index.ts
    type_mappings:
      Shape:
        ts: number[]
      Dim:
        ts: number
```

Then run

```shell
gen-napi generate
```

The output node native addon logic `.cc` will be generated to the specified `bindings_out_path` and the TS/JS will be generated to the specified `wrapper_out_path`.

### Option B: Library Mode (TODO: Docs)

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
