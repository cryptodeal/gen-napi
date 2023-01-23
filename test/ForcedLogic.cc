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
