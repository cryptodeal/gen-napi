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
}  // namespace global_vars

namespace exported_global_methods {
/*
  @gen-napi-`ts_return_type`: void
*/
static void init(const Napi::CallbackInfo& /*info*/) {
  fl::init();
}

/*
  @gen-napi-`ts_return_type`: bigint
*/
static Napi::Value bytesUsed(const Napi::CallbackInfo& info) {
  return Napi::BigInt::New(info.Env(), static_cast<int64_t>(g_bytes_used));
}

/*
  @gen-napi-`ts_return_type`: void
*/
static void setRowMajor(const Napi::CallbackInfo& /*info*/) {
  g_row_major = true;
}

/*
  @gen-napi-`ts_return_type`: void
*/
static void setColMajor(const Napi::CallbackInfo& /*info*/) {
  g_row_major = false;
}

/*
  @gen-napi-`ts_return_type`: boolean
*/
static Napi::Value isRowMajor(const Napi::CallbackInfo& info) {
  return Napi::Boolean::New(info.Env(), g_row_major);
}

/*
  @gen-napi-`ts_return_type`: boolean
*/
static Napi::Value isColMajor(const Napi::CallbackInfo& info) {
  return Napi::Boolean::New(info.Env(), !g_row_major);
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeFloat32(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::f32));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeFloat64(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::f64));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeBoolInt8(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::b8));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeInt16(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::s16));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeInt32(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::s32));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeInt64(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::s64));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeUint8(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u8));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeUint16(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u16));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeUint32(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u32));
}

/*
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value dtypeUint64(const Napi::CallbackInfo& info) {
  return Napi::Number::New(info.Env(), static_cast<double>(fl::dtype::u64));
}

/*
  @gen-napi-`ts_args`: (shape: number[])
  @gen-napi-`ts_return_type`: Tensor
*/
static Napi::Value rand(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 1 || !info[0].IsArray()) {
    Napi::TypeError::New(env,
                         "`rand` expects exactly 1 arg; (typeof `number[]`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  std::vector<long long> shape = jsArrayToVector<long long>(
      env, info[0].As<Napi::Array>(), g_row_major, false);
  fl::Tensor t;
  t = fl::rand(fl::Shape(shape));
  auto _out_bytes_used = t.bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  auto* tensor = new fl::Tensor(t);
  Napi::External<fl::Tensor> wrapped = ExternalizeTensor(env, tensor);
  return wrapped;
}

/*
  @gen-napi-`ts_args`: (shape: number[])
  @gen-napi-`ts_return_type`: Tensor
*/
static Napi::Value randn(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 1 || !info[0].IsArray()) {
    Napi::TypeError::New(env,
                         "`randn` expects exactly 1 arg; (typeof `number[]`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  std::vector<long long> shape = jsArrayToVector<long long>(
      env, info[0].As<Napi::Array>(), g_row_major, false);
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
static inline std::vector<T> jsTensorArrayArg(Napi::Array arr, Napi::Env env) {
  std::vector<T> out;
  size_t len = arr.Length();
  out.reserve(len);
  for (size_t i = 0; i < len; ++i) {
    Napi::Value temp = arr[i];
    fl::Tensor* tensor = UnExternalize<fl::Tensor>(temp);
    out.emplace_back(*(tensor));
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

}  // namespace private_helpers

namespace Tensor_forced_methods {
/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Float32Array
*/
static Napi::Value toFloat32Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(float);
  float* ptr;
  if (contig_tensor->type() == fl::dtype::f32) {
    ptr = contig_tensor->host<float>();
  } else {
    ptr = contig_tensor->astype(fl::dtype::f32).host<float>();
  }
  std::unique_ptr<std::vector<float>> nativeArray =
      std::make_unique<std::vector<float>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<float>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<float>::New(env, elemLen, buff, 0,
                                        napi_float32_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Float64Array
*/
static Napi::Value toFloat64Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(double);
  double* ptr;
  if (contig_tensor->type() == fl::dtype::f64) {
    ptr = contig_tensor->host<double>();
  } else {
    ptr = contig_tensor->astype(fl::dtype::f64).host<double>();
  }
  std::unique_ptr<std::vector<double>> nativeArray =
      std::make_unique<std::vector<double>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<double>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<double>::New(env, elemLen, buff, 0,
                                         napi_float64_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Int8Array
*/
static Napi::Value toBoolInt8Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(int8_t);
  int8_t* ptr;
  if (contig_tensor->type() == fl::dtype::b8) {
    ptr = reinterpret_cast<int8_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<int8_t*>(
        contig_tensor->astype(fl::dtype::b8).host<int>());
  }
  std::unique_ptr<std::vector<int8_t>> nativeArray =
      std::make_unique<std::vector<int8_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<int8_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<int8_t>::New(env, elemLen, buff, 0,
                                         napi_int8_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Int16Array
*/
static Napi::Value toInt16Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(int16_t);
  int16_t* ptr;
  if (contig_tensor->type() == fl::dtype::s16) {
    ptr = reinterpret_cast<int16_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<int16_t*>(
        contig_tensor->astype(fl::dtype::s16).host<int>());
  }
  std::unique_ptr<std::vector<int16_t>> nativeArray =
      std::make_unique<std::vector<int16_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<int16_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<int16_t>::New(env, elemLen, buff, 0,
                                          napi_int16_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Int32Array
*/
static Napi::Value toInt32Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(int32_t);
  int32_t* ptr;
  if (contig_tensor->type() == fl::dtype::s32) {
    ptr = reinterpret_cast<int32_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<int32_t*>(
        contig_tensor->astype(fl::dtype::s32).host<int>());
  }
  std::unique_ptr<std::vector<int32_t>> nativeArray =
      std::make_unique<std::vector<int32_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<int32_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<int32_t>::New(env, elemLen, buff, 0,
                                          napi_int32_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: BigInt64Array
*/
static Napi::Value toInt64Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(int64_t);
  int64_t* ptr;
  if (contig_tensor->type() == fl::dtype::s64) {
    ptr = reinterpret_cast<int64_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<int64_t*>(
        contig_tensor->astype(fl::dtype::s64).host<int>());
  }
  std::unique_ptr<std::vector<int64_t>> nativeArray =
      std::make_unique<std::vector<int64_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<int64_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<int64_t>::New(env, elemLen, buff, 0,
                                          napi_bigint64_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Uint8Array
*/
static Napi::Value toUint8Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(uint8_t);
  uint8_t* ptr;
  if (contig_tensor->type() == fl::dtype::u8) {
    ptr = reinterpret_cast<uint8_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<uint8_t*>(
        contig_tensor->astype(fl::dtype::u8).host<int>());
  }
  std::unique_ptr<std::vector<uint8_t>> nativeArray =
      std::make_unique<std::vector<uint8_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<uint8_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<uint8_t>::New(env, elemLen, buff, 0,
                                          napi_uint8_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Uint16Array
*/
static Napi::Value toUint16Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(uint16_t);
  uint16_t* ptr;
  if (contig_tensor->type() == fl::dtype::u16) {
    ptr = reinterpret_cast<uint16_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<uint16_t*>(
        contig_tensor->astype(fl::dtype::u16).host<int>());
  }
  std::unique_ptr<std::vector<uint16_t>> nativeArray =
      std::make_unique<std::vector<uint16_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<uint16_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<uint16_t>::New(env, elemLen, buff, 0,
                                           napi_uint16_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: Uint32Array
*/
static Napi::Value toUint32Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(uint32_t);
  uint32_t* ptr;
  if (contig_tensor->type() == fl::dtype::u32) {
    ptr = reinterpret_cast<uint32_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<uint32_t*>(
        contig_tensor->astype(fl::dtype::u32).host<int>());
  }
  std::unique_ptr<std::vector<uint32_t>> nativeArray =
      std::make_unique<std::vector<uint32_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<uint32_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<uint32_t>::New(env, elemLen, buff, 0,
                                           napi_uint32_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: BigUint64Array
*/
static Napi::Value toUint64Array(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::Tensor* contig_tensor = t;
  bool isContiguous = t->isContiguous();
  if (!isContiguous) {
    contig_tensor = new fl::Tensor(t->asContiguousTensor());
  }
  size_t elemLen = contig_tensor->elements();
  size_t byteLen = elemLen * sizeof(uint64_t);
  uint64_t* ptr;
  if (contig_tensor->type() == fl::dtype::u32) {
    ptr = reinterpret_cast<uint64_t*>(contig_tensor->host<int>());
  } else {
    ptr = reinterpret_cast<uint64_t*>(
        contig_tensor->astype(fl::dtype::u64).host<int>());
  }
  std::unique_ptr<std::vector<uint64_t>> nativeArray =
      std::make_unique<std::vector<uint64_t>>(ptr, ptr + elemLen);
  if (!isContiguous) {
    delete contig_tensor;
  }
  Napi::ArrayBuffer buff =
      Napi::ArrayBuffer::New(env, nativeArray->data(), byteLen,
                             DeleteArrayBuffer<uint64_t>, nativeArray.get());
  nativeArray.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, byteLen);
  return Napi::TypedArrayOf<uint64_t>::New(env, elemLen, buff, 0,
                                           napi_biguint64_array);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toFloat32Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toFloat32Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<float>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toFloat64Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toFloat64Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<double>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toBoolInt8Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toBoolInt8Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<char>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toInt16Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toInt16Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<int16_t>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toInt32Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toInt32Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<int32_t>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: bigint
*/
static Napi::Value toInt64Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toInt64Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::BigInt::New(env, t->asScalar<int64_t>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toUint8Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toUint8Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<uint8_t>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toUint16Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toUint16Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<uint16_t>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: number
*/
static Napi::Value toUint32Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toUint32Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::Number::New(env, t->asScalar<uint32_t>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: bigint
*/
static Napi::Value toUint64Scalar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`toUint64Scalar` expects args[0] to be native "
                         "`Tensor` (typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  return Napi::BigInt::New(env, t->asScalar<uint64_t>());
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: void
*/
static void eval(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`eval` expects args[0] to be native `Tensor` (typeof "
                         "`Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return;
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  fl::eval(*(t));
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor)
  @gen-napi-`ts_return_type`: void
*/
static void dispose(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`dispose` expects args[0] to be native `Tensor` "
                         "(typeof `Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return;
  }
  fl::Tensor& t = *UnExternalize<fl::Tensor>(info[0]);
  auto byte_count = static_cast<int64_t>(t.bytes());
  g_bytes_used -= byte_count;
  Napi::MemoryManagement::AdjustExternalMemory(env, -byte_count);
  fl::detail::releaseAdapterUnsafe(t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromFloat32Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(env,
                     "`tensorFromFloat32Buffer` epects args[0] to be "
                     "instanceof `ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(float));
  float* ptr = reinterpret_cast<float*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromFloat64Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(env,
                     "`tensorFromFloat64Buffer` epects args[0] to be "
                     "instanceof `ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(double));
  double* ptr = reinterpret_cast<double*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromBoolInt8Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(env,
                     "`tensorFromBoolInt8Buffer` epects args[0] to be "
                     "instanceof `ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(int8_t));
  char* ptr = reinterpret_cast<char*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromInt16Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(
        env,
        "`tensorFromInt16Buffer` epects args[0] to be instanceof `ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(int16_t));
  int16_t* ptr = reinterpret_cast<int16_t*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromInt32Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(
        env,
        "`tensorFromInt32Buffer` epects args[0] to be instanceof `ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(int32_t));
  int32_t* ptr = reinterpret_cast<int32_t*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromInt64Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(env,
                     "`tensorFtensorFromInt64BufferromFloat32Buffer` epects "
                     "args[0] to be instanceof `ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(int64_t));
  int64_t* ptr = reinterpret_cast<int64_t*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromUint8Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(
        env,
        "`tensorFromUint8Buffer` epects args[0] to be instanceof `ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(uint8_t));
  uint8_t* ptr = reinterpret_cast<uint8_t*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromUint16Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(env,
                     "`tensorFromUint16Buffer` epects args[0] to be instanceof "
                     "`ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(uint16_t));
  uint16_t* ptr = reinterpret_cast<uint16_t*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromUint32Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(env,
                     "`tensorFromUint32Buffer` epects args[0] to be instanceof "
                     "`ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(uint32_t));
  uint32_t* ptr = reinterpret_cast<uint32_t*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (buffer: ArrayBuffer)
  @gen-napi-`ts_return_type`: any
*/
static Napi::Value tensorFromUint64Buffer(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (!info[0].IsArrayBuffer()) {
    Napi::Error::New(env,
                     "`tensorFromUint64Buffer` epects args[0] to be instanceof "
                     "`ArrayBuffer`")
        .ThrowAsJavaScriptException();
    return env.Null();
  }
  Napi::ArrayBuffer buf = info[0].As<Napi::ArrayBuffer>();
  int64_t length = static_cast<int64_t>(buf.ByteLength() / sizeof(uint64_t));
  uint64_t* ptr = reinterpret_cast<uint64_t*>(buf.Data());
  auto* t = new fl::Tensor(
      fl::Tensor::fromBuffer({length}, ptr, fl::MemoryLocation::Host));
  auto _out_bytes_used = t->bytes();
  g_bytes_used += _out_bytes_used;
  Napi::MemoryManagement::AdjustExternalMemory(env, _out_bytes_used);
  return ExternalizeTensor(env, t);
}

/*
  @gen-napi-`ts_args`: (tensor: Tensor, path: string)
  @gen-napi-`ts_return_type`: void
*/
static void save(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 2) {
    Napi::TypeError::New(env, "`save` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return;
  }
  if (!info[0].IsExternal()) {
    Napi::TypeError::New(env,
                         "`save` expects args[0] to be native `Tensor` (typeof "
                         "`Napi::External<fl::Tensor>`)")
        .ThrowAsJavaScriptException();
    return;
  }
  fl::Tensor* t = UnExternalize<fl::Tensor>(info[0]);
  if (!info[1].IsString()) {
    Napi::TypeError::New(env, "`save` expects args[1] to be typeof `string`")
        .ThrowAsJavaScriptException();
    return;
  }
  Napi::String str = info[1].As<Napi::String>();
  std::string filename = str.Utf8Value();
  fl::save(filename, *(t));
}
}  // namespace Tensor_forced_methods