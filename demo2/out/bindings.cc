// Code generated by gen-napi. DO NOT EDIT.
#include <napi.h>
#include <atomic>
#include <string>

//////////
// source: BasicLogic.h

// non-exported helpers
template <typename T>
static inline std::vector<T> jsArrayToVector(Napi::Array arr,
                                             bool reverse,
                                             int invert) {
  std::vector<T> out;
  size_t len = arr.Length();
  out.reserve(len);
  for (size_t i = 0; i < len; ++i) {
    const auto idx = reverse ? len - i - 1 : i;
    Napi::Value val = arr[idx];
    auto v = static_cast<const T>(val.As<Napi::Number>().Int64Value());
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
static inline void DeleteArrayBuffer(Napi::Env env,
                                     void* /*data*/,
                                     std::vector<T>* hint) {
  size_t bytes = hint->size() * sizeof(T);
  std::unique_ptr<std::vector<T>> vectorPtrToDelete(hint);
  Napi::MemoryManagement::AdjustExternalMemory(env, -bytes);
}

// exported functions

static Napi::Value _qux(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 2) {
    Napi::TypeError::New(env, "`qux` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`qux` expects args[0] to be typeof `BigInt64Array`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  long long* a = reinterpret_cast<long long*>(
      info[0].As<Napi::TypedArrayOf<int64_t>>().Data());
  if (!info[1].IsNumber()) {
    Napi::TypeError::New(env, "`qux` expects args[1] to be typeof `number`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  int b = static_cast<int>(info[1].As<Napi::Number>().Int64Value());
  int64_t* _res;
  _res = reinterpret_cast<int64_t*>(demo2::qux(a, b));
  size_t _res_byte_len = sizeof(_res);
  size_t _res_elem_len = _res_byte_len / sizeof(*_res);
  std::unique_ptr<std::vector<int64_t>> _res_native_array =
      std::make_unique<std::vector<int64_t>>(_res, _res + _res_elem_len);
  Napi::ArrayBuffer _res_arraybuffer = Napi::ArrayBuffer::New(
      env, _res_native_array->data(), _res_byte_len, DeleteArrayBuffer<int64_t>,
      _res_native_array.get());
  _res_native_array.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, _res_byte_len);
  return Napi::TypedArrayOf<int64_t>::New(env, _res_elem_len, _res_arraybuffer,
                                          0, napi_bigint64_array);
}

static Napi::Value _quux(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 2) {
    Napi::TypeError::New(env, "`quux` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsBoolean()) {
    Napi::TypeError::New(env, "`quux` expects args[0] to be typeof `boolean`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  bool a = info[0].As<Napi::Boolean>().Value();
  if (!info[1].IsBoolean()) {
    Napi::TypeError::New(env, "`quux` expects args[1] to be typeof `boolean`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  bool b = info[1].As<Napi::Boolean>().Value();
  uint8_t* _res;
  _res = reinterpret_cast<uint8_t*>(demo2::quux(a, b));
  size_t _res_byte_len = sizeof(_res);
  size_t _res_elem_len = _res_byte_len / sizeof(*_res);
  std::unique_ptr<std::vector<uint8_t>> _res_native_array =
      std::make_unique<std::vector<uint8_t>>(_res, _res + _res_elem_len);
  Napi::ArrayBuffer _res_arraybuffer = Napi::ArrayBuffer::New(
      env, _res_native_array->data(), _res_byte_len, DeleteArrayBuffer<uint8_t>,
      _res_native_array.get());
  _res_native_array.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, _res_byte_len);
  return Napi::TypedArrayOf<uint8_t>::New(env, _res_elem_len, _res_arraybuffer,
                                          0, napi_uint8_array);
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

static Napi::Value _foo(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 1) {
    Napi::TypeError::New(env, "`foo` expects exactly 1 arg")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsNumber()) {
    Napi::TypeError::New(env, "`foo` expects args[0] to be typeof `number`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  int8_t a = static_cast<int8_t>(info[0].As<Napi::Number>().Int32Value());
  int8_t _res;
  _res = demo2::foo(a);
  return Napi::Number::New(env, _res);
}

static Napi::Value _bar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  if (info.Length() != 2) {
    Napi::TypeError::New(env, "`bar` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`bar` expects args[0] to be typeof `Float64Array`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  double* a = info[0].As<Napi::TypedArrayOf<double>>().Data();
  if (!info[1].IsNumber()) {
    Napi::TypeError::New(env, "`bar` expects args[1] to be typeof `number`)")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  int32_t b = static_cast<int32_t>(info[1].As<Napi::Number>().Int32Value());
  double* _res;
  _res = demo2::bar(a, b);
  size_t _res_byte_len = sizeof(_res);
  size_t _res_elem_len = _res_byte_len / sizeof(*_res);
  std::unique_ptr<std::vector<double>> _res_native_array =
      std::make_unique<std::vector<double>>(_res, _res + _res_elem_len);
  Napi::ArrayBuffer _res_arraybuffer = Napi::ArrayBuffer::New(
      env, _res_native_array->data(), _res_byte_len, DeleteArrayBuffer<double>,
      _res_native_array.get());
  _res_native_array.release();
  Napi::MemoryManagement::AdjustExternalMemory(env, _res_byte_len);
  return Napi::TypedArrayOf<double>::New(env, _res_elem_len, _res_arraybuffer,
                                         0, napi_float64_array);
}

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

// NAPI exports

Napi::Object Init(Napi::Env env, Napi::Object exports) {
  exports.Set(Napi::String::New(env, "_quux"), Napi::Function::New(env, _quux));
  exports.Set(Napi::String::New(env, "_test"), Napi::Function::New(env, _test));
  exports.Set(Napi::String::New(env, "_foo"), Napi::Function::New(env, _foo));
  exports.Set(Napi::String::New(env, "_bar"), Napi::Function::New(env, _bar));
  exports.Set(Napi::String::New(env, "_baz"), Napi::Function::New(env, _baz));
  exports.Set(Napi::String::New(env, "_qux"), Napi::Function::New(env, _qux));
  return exports;
}

NODE_API_MODULE(addon, Init)
