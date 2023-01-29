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
  test2::int8_t _res;
  _res = test2::foo(a);
  auto* out = new test2::int8_t(_res);
  Napi::External<test2::int8_t> _external_out = Externalizeint8_t(env, out);
  return _external_out;
}

// NAPI exports

Napi::Object Init(Napi::Env env, Napi::Object exports) {
  exports.Set(Napi::String::New(env, "_foo"), Napi::Function::New(env, _foo));
  return exports;
}

NODE_API_MODULE(addon, Init)
