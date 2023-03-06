// Code generated by gen-napi. DO NOT EDIT.
#include <napi.h>
#include <atomic>
#include <string>
// non-exported helpers
template <typename T>
static inline void DeleteArrayBufferFromVector(Napi::Env env,
                                               void* /*data*/,
                                               std::vector<T>* hint) {
  size_t bytes = hint->size() * sizeof(T);
  std::unique_ptr<std::vector<T>> vectorPtrToDelete(hint);
  Napi::MemoryManagement::AdjustExternalMemory(env, -bytes);
}

//////////
// source: BasicLogic.h

// exported functions

static Napi::Value _bar(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 2) {
    Napi::TypeError::New(env, "`bar` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`bar` expects `a` (args[0]) to be instanceof "
                         "`number[] | Float64Array`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[1].IsNumber()) {
    Napi::TypeError::New(env,
                         "`bar` expects `b` (args[1]) to be typeof `number`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  Napi::TypedArrayOf<double> _gen_tmp_a =
      info[0].As<Napi::TypedArrayOf<double>>();
  auto* a = _gen_tmp_a.Data();
  auto b = info[1].As<Napi::Number>().Int32Value();
  auto _gen_tmp_res_value = demo2::bar(a, b);
  std::unique_ptr<std::vector<double>> _gen_tmp_res_value_vec;
  _gen_tmp_res_value_vec.reset(&_gen_tmp_res_value);
  auto _gen_tmp_res_value_vec_elem_len = _gen_tmp_res_value_vec->size();
  auto _gen_tmp_res_value_vec_byte_len =
      _gen_tmp_res_value_vec_elem_len * sizeof(double);
  Napi::ArrayBuffer _gen_tmp_res_value_js_tform = Napi::ArrayBuffer::New(
      env, _gen_tmp_res_value_vec->data(), _gen_tmp_res_value_vec_byte_len,
      DeleteArrayBufferFromVector<double>, _gen_tmp_res_value_vec.get());
  _gen_tmp_res_value_vec.release();
  Napi::MemoryManagement::AdjustExternalMemory(env,
                                               _gen_tmp_res_value_vec_byte_len);
  return Napi::TypedArrayOf<double>::New(env, _gen_tmp_res_value_vec_elem_len,
                                         _gen_tmp_res_value_js_tform, 0);
}

static Napi::Value _quux(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 2) {
    Napi::TypeError::New(env, "`quux` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsBoolean()) {
    Napi::TypeError::New(env,
                         "`quux` expects `a` (args[0]) to be typeof `boolean`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[1].IsBoolean()) {
    Napi::TypeError::New(env,
                         "`quux` expects `b` (args[1]) to be typeof `boolean`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  bool a = info[0].As<Napi::Boolean>().Value();
  bool b = info[1].As<Napi::Boolean>().Value();
  auto _gen_tmp_res_value = demo2::quux(a, b);
  std::unique_ptr<std::vector<uint8_t>> _gen_tmp_res_value_cast;
  auto _gen_tmp_res_value_cast_elem_len = _gen_tmp_res_value_cast->size();
  _gen_tmp_res_value_cast->reserve(_gen_tmp_res_value_cast_elem_len);
  for (size_t i = 0; i < _gen_tmp_res_value_cast_elem_len; ++i) {
    (*_gen_tmp_res_value_cast)[i] = static_cast<uint8_t>(_gen_tmp_res_value[i]);
  }
  auto _gen_tmp_res_value_cast_byte_len =
      _gen_tmp_res_value_cast_elem_len * sizeof(uint8_t);
  Napi::ArrayBuffer _gen_tmp_res_value_js_tform = Napi::ArrayBuffer::New(
      env, _gen_tmp_res_value_cast->data(), _gen_tmp_res_value_cast_byte_len,
      DeleteArrayBufferFromVector<uint8_t>, _gen_tmp_res_value_cast.get());
  _gen_tmp_res_value_cast.release();
  Napi::MemoryManagement::AdjustExternalMemory(
      env, _gen_tmp_res_value_cast_byte_len);
  return Napi::TypedArrayOf<uint8_t>::New(env, _gen_tmp_res_value_cast_elem_len,
                                          _gen_tmp_res_value_js_tform, 0);
}

static Napi::Value _test3(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 1) {
    Napi::TypeError::New(env, "`test3` expects exactly 1 arg")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`test3` expects `a` (args[0]) to be instanceof "
                         "`number[] | Float64Array`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  Napi::TypedArrayOf<double> _gen_tmp_a =
      info[0].As<Napi::TypedArrayOf<double>>();
  auto* a = _gen_tmp_a.Data();
  auto _gen_tmp_res_value = demo2::test3(a);
  Napi::Array _gen_tmp_res_value_js_tform = Napi::Array::New(env, 2);
  _gen_tmp_res_value_js_tform[0] =
      Napi::Number::New(env, static_cast<double>(_gen_tmp_res_value.first));
  _gen_tmp_res_value_js_tform[1] =
      Napi::Number::New(env, static_cast<double>(_gen_tmp_res_value.second));
  return _gen_tmp_res_value_js_tform;
}

static Napi::Value _test5(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 1) {
    Napi::TypeError::New(env, "`test5` expects exactly 1 arg")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`test5` expects `a` (args[0]) to be instanceof "
                         "`number[] | Float64Array`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  Napi::TypedArrayOf<double> _gen_tmp_a =
      info[0].As<Napi::TypedArrayOf<double>>();
  auto* _gen_tmp_a_ptr = _gen_tmp_a.Data();
  size_t _gen_tmp_a_len = _gen_tmp_a.ElementLength();
  std::vector<double> a(_gen_tmp_a_ptr, _gen_tmp_a_ptr + _gen_tmp_a_len);
  auto _gen_tmp_res_value = demo2::test5(a);
  std::unique_ptr<std::vector<int64_t>> _gen_tmp_res_value_cast;
  auto _gen_tmp_res_value_cast_elem_len = _gen_tmp_res_value_cast->size();
  _gen_tmp_res_value_cast->reserve(_gen_tmp_res_value_cast_elem_len);
  for (size_t i = 0; i < _gen_tmp_res_value_cast_elem_len; ++i) {
    (*_gen_tmp_res_value_cast)[i] = static_cast<int64_t>(_gen_tmp_res_value[i]);
  }
  auto _gen_tmp_res_value_cast_byte_len =
      _gen_tmp_res_value_cast_elem_len * sizeof(int64_t);
  Napi::ArrayBuffer _gen_tmp_res_value_js_tform = Napi::ArrayBuffer::New(
      env, _gen_tmp_res_value_cast->data(), _gen_tmp_res_value_cast_byte_len,
      DeleteArrayBufferFromVector<int64_t>, _gen_tmp_res_value_cast.get());
  _gen_tmp_res_value_cast.release();
  Napi::MemoryManagement::AdjustExternalMemory(
      env, _gen_tmp_res_value_cast_byte_len);
  return Napi::TypedArrayOf<int64_t>::New(env, _gen_tmp_res_value_cast_elem_len,
                                          _gen_tmp_res_value_js_tform, 0);
}

static Napi::Value _test2(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 1) {
    Napi::TypeError::New(env, "`test2` expects exactly 1 arg")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`test2` expects `a` (args[0]) to be instanceof "
                         "`number[] | Float64Array`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  Napi::TypedArrayOf<double> _gen_tmp_a =
      info[0].As<Napi::TypedArrayOf<double>>();
  auto* _gen_tmp_a_ptr = _gen_tmp_a.Data();
  size_t _gen_tmp_a_len = _gen_tmp_a.ElementLength();
  std::vector<double> a(_gen_tmp_a_ptr, _gen_tmp_a_ptr + _gen_tmp_a_len);
  auto _gen_tmp_res_value = demo2::test2(a);
  std::unique_ptr<std::array<int8_t, 20>> _gen_tmp_res_value_std_array;
  _gen_tmp_res_value_std_array.reset(&_gen_tmp_res_value);
  auto _gen_tmp_res_value_std_array_byte_len = 20 * sizeof(int8_t);
  Napi::ArrayBuffer _gen_tmp_res_value_js_tform = Napi::ArrayBuffer::New(
      env, 20, _gen_tmp_res_value_std_array_byte_len,
      [](Napi::Env env, void* /*data*/, std::array<int8_t, 20>* hint) {
        std::unique_ptr<std::array<int8_t, 20>> arrayPtrToDelete(hint);
        Napi::MemoryManagement::AdjustExternalMemory(env,
                                                     -(20 * sizeof(int8_t)));
      },
      _gen_tmp_res_value_std_array.get());
  _gen_tmp_res_value_std_array.release();
  Napi::MemoryManagement::AdjustExternalMemory(
      env, _gen_tmp_res_value_std_array_byte_len);
  return Napi::TypedArrayOf<int8_t>::New(env, 20, _gen_tmp_res_value_js_tform,
                                         0);
}

static Napi::Value _test4(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 2) {
    Napi::TypeError::New(env, "`test4` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsString()) {
    Napi::TypeError::New(env,
                         "`test4` expects `a` (args[0]) to be typeof `string`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[1].IsString()) {
    Napi::TypeError::New(env,
                         "`test4` expects `b` (args[1]) to be typeof `string`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  auto _gen_tmp_std_str_a = info[0].As<Napi::String>().Utf8Value();
  auto* a = _gen_tmp_std_str_a.c_str();
  auto _gen_tmp_std_str_b = info[1].As<Napi::String>().Utf16Value();
  auto* b = _gen_tmp_std_str_b.c_str();
  auto* _gen_tmp_res_value = demo2::test4(a, b);
  return Napi::String::New(env, _gen_tmp_res_value);
}

static void _foo(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 1) {
    Napi::TypeError::New(env, "`foo` expects exactly 1 arg")
        .ThrowAsJavaScriptException();
    return;
  }
  if (!info[0].IsNumber()) {
    Napi::TypeError::New(env,
                         "`foo` expects `a` (args[0]) to be typeof `number`")
        .ThrowAsJavaScriptException();
    return;
  }
  auto a = static_cast<int8_t>(info[0].As<Napi::Number>().Int32Value());
  demo2::foo(a);
}

static Napi::Value _baz(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 2) {
    Napi::TypeError::New(env, "`baz` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`baz` expects `a` (args[0]) to be instanceof "
                         "`number[] | Float32Array`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[1].IsNumber()) {
    Napi::TypeError::New(env,
                         "`baz` expects `b` (args[1]) to be typeof `number`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  Napi::TypedArrayOf<float> _gen_tmp_a =
      info[0].As<Napi::TypedArrayOf<float>>();
  auto* a = _gen_tmp_a.Data();
  auto b = static_cast<int>(info[1].As<Napi::Number>().Int32Value());
  auto _gen_tmp_res_value = demo2::baz(a, b);
  std::unique_ptr<std::vector<float>> _gen_tmp_res_value_vec;
  _gen_tmp_res_value_vec.reset(&_gen_tmp_res_value);
  auto _gen_tmp_res_value_vec_elem_len = _gen_tmp_res_value_vec->size();
  auto _gen_tmp_res_value_vec_byte_len =
      _gen_tmp_res_value_vec_elem_len * sizeof(float);
  Napi::ArrayBuffer _gen_tmp_res_value_js_tform = Napi::ArrayBuffer::New(
      env, _gen_tmp_res_value_vec->data(), _gen_tmp_res_value_vec_byte_len,
      DeleteArrayBufferFromVector<float>, _gen_tmp_res_value_vec.get());
  _gen_tmp_res_value_vec.release();
  Napi::MemoryManagement::AdjustExternalMemory(env,
                                               _gen_tmp_res_value_vec_byte_len);
  return Napi::TypedArrayOf<float>::New(env, _gen_tmp_res_value_vec_elem_len,
                                        _gen_tmp_res_value_js_tform, 0);
}

static Napi::Value _qux(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 2) {
    Napi::TypeError::New(env, "`qux` expects exactly 2 args")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsTypedArray()) {
    Napi::TypeError::New(env,
                         "`qux` expects `a` (args[0]) to be instanceof "
                         "`Array<number | bigint> | BigInt64Array`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[1].IsNumber()) {
    Napi::TypeError::New(env,
                         "`qux` expects `b` (args[1]) to be typeof `number`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  Napi::TypedArrayOf<int64_t> _gen_tmp_a =
      info[0].As<Napi::TypedArrayOf<int64_t>>();
  auto* a = static_cast<long long*>(_gen_tmp_a.Data());
  auto b = static_cast<int>(info[1].As<Napi::Number>().Int32Value());
  auto _gen_tmp_res_value = demo2::qux(a, b);
  return Napi::Number::New(env, _gen_tmp_res_value);
}

static Napi::Value _test(const Napi::CallbackInfo& info) {
  Napi::Env env = info.Env();
  const auto _arg_count = info.Length();
  if (_arg_count != 1) {
    Napi::TypeError::New(env, "`test` expects exactly 1 arg")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  if (!info[0].IsString()) {
    Napi::TypeError::New(env,
                         "`test` expects `a` (args[0]) to be typeof `string`")
        .ThrowAsJavaScriptException();
    return env.Undefined();
  }
  auto a = info[0].As<Napi::String>().Utf8Value();
  auto _gen_tmp_res_value = demo2::test(a);
  return Napi::String::New(env, _gen_tmp_res_value);
}

// NAPI exports

Napi::Object Init(Napi::Env env, Napi::Object exports) {
  exports.Set(Napi::String::New(env, "_bar"), Napi::Function::New(env, _bar));
  exports.Set(Napi::String::New(env, "_quux"), Napi::Function::New(env, _quux));
  exports.Set(Napi::String::New(env, "_test3"),
              Napi::Function::New(env, _test3));
  exports.Set(Napi::String::New(env, "_test5"),
              Napi::Function::New(env, _test5));
  exports.Set(Napi::String::New(env, "_foo"), Napi::Function::New(env, _foo));
  exports.Set(Napi::String::New(env, "_baz"), Napi::Function::New(env, _baz));
  exports.Set(Napi::String::New(env, "_qux"), Napi::Function::New(env, _qux));
  exports.Set(Napi::String::New(env, "_test"), Napi::Function::New(env, _test));
  exports.Set(Napi::String::New(env, "_test2"),
              Napi::Function::New(env, _test2));
  exports.Set(Napi::String::New(env, "_test4"),
              Napi::Function::New(env, _test4));
  return exports;
}

NODE_API_MODULE(addon, Init)
