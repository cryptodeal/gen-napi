#include <napi.h>
#include <atomic>
#include <iostream>
#include <string>
#include "flashlight/fl/autograd/Functions.h"
#include "flashlight/fl/autograd/tensor/AutogradExtension.h"
#include "flashlight/fl/autograd/tensor/AutogradOps.h"
#include "flashlight/fl/common/DynamicBenchmark.h"
#include "flashlight/fl/nn/Init.h"
#include "flashlight/fl/runtime/Device.h"
#include "flashlight/fl/runtime/Stream.h"
#include "flashlight/fl/tensor/Compute.h"
#include "flashlight/fl/tensor/Index.h"
#include "flashlight/fl/tensor/Init.h"
#include "flashlight/fl/tensor/Random.h"
#include "flashlight/fl/tensor/TensorAdapter.h"

namespace global_vars {
static std::atomic<size_t> g_bytes_used = 0;
static std::atomic<bool> g_row_major = true;
}  // namespace global_vars

namespace exported_global_methods {
static void init(const Napi::CallbackInfo& /*info*/) {
  fl::init();
}

static Napi::Value bytesUsed(const Napi::CallbackInfo& info) {
  return Napi::BigInt::New(info.Env(), static_cast<int64_t>(g_bytes_used));
}

static void setRowMajor(const Napi::CallbackInfo& /*info*/) {
  g_row_major = true;
}

static void setColMajor(const Napi::CallbackInfo& /*info*/) {
  g_row_major = true;
}

static Napi::Value isRowMajor(const Napi::CallbackInfo& info) {
  return Napi::Boolean::New(info.Env(), g_row_major);
}

static Napi::Value isColMajor(const Napi::CallbackInfo& info) {
  return Napi::Boolean::New(info.Env(), !g_row_major);
}

static Napi::Value dtypeFloat32(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::f32));
}

static Napi::Value dtypeFloat64(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::f64));
}

static Napi::Value dtypeBoolInt8(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::b8));
}

static Napi::Value dtypeInt16(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::s16));
}

static Napi::Value dtypeInt32(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::s32));
}

static Napi::Value dtypeInt64(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::s64));
}

static Napi::Value dtypeUint8(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u8));
}

static Napi::Value dtypeUint16(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u16));
}

static Napi::Value dtypeUint32(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u32));
}

static Napi::Value dtypeUint64(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u64));
}

static Napi::Value rand(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 1 || !info[0].IsArray()) {
    Napi::TypeError::New(env,
                         "`rand` expects exactly 1 arg; (typeof `number[]`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  std::vector<long long> shape =
      jsArrayArg<long long>(info[0].As<Napi::Array>(), g_row_major, false, env);
  fl::Tensor t;
  t = fl::rand(fl::Shape(shape));
  auto _out_bytes_used = t.bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  auto* tensor = new fl::Tensor(t);
  Napi::External<fl::Tensor> wrapped = ExternalizeTensor(env, tensor);
  return wrapped;
}

static Napi::Value randn(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 1 || !info[0].IsArray()) {
    Napi::TypeError::New(env,
                         "`randn` expects exactly 1 arg; (typeof `number[]`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  std::vector<long long> shape =
      jsArrayArg<long long>(info[0].As<Napi::Array>(), g_row_major, false, env);
  fl::Tensor t;
  t = fl::randn(fl::Shape(shape));
  auto _out_bytes_used = t.bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  auto* tensor = new fl::Tensor(t);
  Napi::External<fl::Tensor> wrapped = ExternalizeTensor(env, tensor);
  return wrapped;
}
}  // namespace exported_global_methods

namespace private_helpers {
template <typename T>
static inline std::vector<T> arrayArg(const void* ptr,
                                      int len,
                                      bool reverse,
                                      int invert) {
  std::vector<T> out;
  out.reserve(len);
  for (auto i = 0; i < len; ++i) {
    const auto idx = reverse ? len - i - 1 : i;
    auto v = reinterpret_cast<const int64_t*>(ptr)[idx];
    if (invert && v < 0) {
      v = -v - 1;
    } else if (invert) {
      v = invert - v - 1;
    }
    out.emplace_back(v);
  }
  return out;
}

template <typename T>
static inline std::vector<T> jsArrayArg(Napi::Array arr,
                                        bool reverse,
                                        int invert,
                                        Napi::Env env) {
  std::vector<T> out;
  const size_t len = static_cast<size_t>(arr.Length());
  out.reserve(len);
  for (size_t i = 0; i < len; ++i) {
    const auto idx = reverse ? len - i - 1 : i;
    Napi::Value val = arr[idx];
    if (!val.IsNumber()) {
      Napi::TypeError::New(env, "jsArrayArg requires `number[]`")
          .ThrowAsJavaScriptException();
      return out;
    } else {
      int64_t v = val.As<Napi::Number>().Int64Value();
      if (invert && v < 0) {
        v = -v - 1;
      } else if (invert) {
        v = invert - v - 1;
      }
      out.emplace_back(v);
    }
  }
  return out;
}

template <typename T>
static inline std::vector<T> jsTensorArrayArg(Napi::Array arr, Napi::Env env) {
  std::vector<T> out;
  const size_t len = static_cast<size_t>(arr.Length());
  out.reserve(len);
  for (size_t i = 0; i < len; ++i) {
    Napi::Value temp = arr[i];
    if (temp.IsExternal()) {
      fl::Tensor* tensor = UnExternalize<fl::Tensor>(temp);
      out.emplace_back(*(tensor));
    } else {
      Napi::TypeError::New(env, "jsTensorArrayArg requires `Tensor[]`")
          .ThrowAsJavaScriptException();
      return out;
    }
  }
  return out;
}

static inline uint32_t axisArg(int32_t axis, bool reverse, int ndim) {
  if (!reverse) {
    return static_cast<uint32_t>(axis);
  }
  if (axis >= 0) {
    return static_cast<uint32_t>(ndim - axis - 1);
  } else {
    return static_cast<uint32_t>(-axis - 1);
  }
}

template <typename T>
static inline std::vector<T> ptrArrayArg(const void* ptr, int len) {
  std::vector<T> out;
  out.reserve(len);
  for (auto i = 0; i < len; ++i) {
    auto ptrAsInt = reinterpret_cast<const int64_t*>(ptr)[i];
    auto ptr = reinterpret_cast<T*>(ptrAsInt);
    out.emplace_back(*ptr);
  }
  return out;
}
}