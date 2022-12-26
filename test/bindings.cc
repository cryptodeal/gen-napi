// Code generated by gen-napi. DO NOT EDIT.
#include "bindings.h"
#include <atomic>
#include <iostream>
#include <string>
using namespace Napi;

//////////
// source: TensorBase.h

// globally scoped variables
static std::atomic<size_t> g_bytes_used = 0;
static std::atomic<bool> g_row_major = true;

// non-exported helpers
template <typename T>
std::vector<T> jsTensorArrayArg(Napi::Array arr, Napi::Env env) {
  std::vector<T> out;
  const size_t len = static_cast<size_t>(arr.Length());
  out.reserve(len);
  for (size_t i = 0; i < len; ++i) {
    Napi::Value temp = arr[i];
    if (temp.IsObject()) {
      Napi::Object tensor_obj = temp.As<Napi::Object>();
      if (tensor_obj.InstanceOf(Tensor::constructor->Value())) {
        Tensor* tensor = Napi::ObjectWrap<Tensor>::Unwrap(tensor_obj);
        out.emplace_back(*(tensor->_tensor));
      } else {
        Napi::TypeError::New(env, "jsTensorArrayArg requires `Tensor[]`")
            .ThrowAsJavaScriptException();
        return out;
      }
    } else {
      Napi::TypeError::New(env, "jsTensorArrayArg requires `Tensor[]`")
          .ThrowAsJavaScriptException();
      return out;
    }
  }
  return out;
}

uint32_t axisArg(int32_t axis, bool reverse, int ndim) {
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
std::vector<T> ptrArrayArg(const void* ptr, int len) {
  std::vector<T> out;
  out.reserve(len);
  for (auto i = 0; i < len; ++i) {
    auto ptrAsInt = reinterpret_cast<const int64_t*>(ptr)[i];
    auto ptr = reinterpret_cast<T*>(ptrAsInt);
    out.emplace_back(*ptr);
  }
  return out;
}

fl::Tensor* load(std::string filename, Napi::Env env) {
  try {
    fl::Tensor tensor;
    fl::load(filename, tensor);
    auto* t = new fl::Tensor(tensor);
    g_bytes_used += t->bytes();
    return t;
  } catch (std::exception const& e) {
    Napi::TypeError::New(env, e.what()).ThrowAsJavaScriptException();
  }
}

template <typename T>
std::vector<T> arrayArg(const void* ptr, int len, bool reverse, int invert) {
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
std::vector<T> jsArrayArg(Napi::Array arr, bool reverse, int invert, Napi::Env env) {
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

// exported functions
static Napi::Value amin(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`amin` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`amin` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`amin` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`amin` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value norm(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 4) {
		Napi::TypeError::New(info.Env(), "`norm` expects exactly 4 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`norm` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`norm` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsNumber()) {
		Napi::TypeError::New(info.Env(), "`norm` expects args[2] to be typeof `number`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[3].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`norm` expects args[3] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value any(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`any` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`any` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`any` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`any` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value all(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`all` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`all` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`all` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`all` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value std(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`std` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`std` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`std` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`std` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value exp(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`exp` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`exp` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value tanh(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`tanh` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`tanh` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value scalar(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 0) {
		Napi::TypeError::New(info.Env(), "`scalar` expects exactly 0 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value arange(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	return env.Null();
}

static Napi::Value maximum(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	return env.Null();
}

static Napi::Value transpose(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`transpose` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`transpose` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value var(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 4) {
		Napi::TypeError::New(info.Env(), "`var` expects exactly 4 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`var` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`var` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`var` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[3].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`var` expects args[3] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value full(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`full` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value concatenate(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`concatenate` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsArray()) {
		Napi::TypeError::New(info.Env(), "`concatenate` expects args[0] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value pad(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`pad` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`pad` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`pad` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value sigmoid(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`sigmoid` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`sigmoid` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value median(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`median` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`median` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`median` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`median` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value identity(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`identity` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value reshape(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`reshape` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`reshape` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value log(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`log` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`log` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value flip(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`flip` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`flip` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value log1p(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`log1p` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`log1p` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value sqrt(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`sqrt` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`sqrt` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value argmax(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`argmax` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`argmax` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`argmax` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value logicalNot(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`logicalNot` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`logicalNot` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value tril(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`tril` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`tril` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value cumsum(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`cumsum` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`cumsum` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value sign(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`sign` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`sign` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value triu(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`triu` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`triu` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value argsort(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`argsort` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`argsort` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value mean(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`mean` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`mean` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`mean` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`mean` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value sum(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`sum` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`sum` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`sum` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`sum` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value clip(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	return env.Null();
}

static Napi::Value roll(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`roll` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`roll` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsNumber()) {
		Napi::TypeError::New(info.Env(), "`roll` expects args[1] to be typeof `number`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value isnan(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`isnan` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`isnan` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value argmin(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`argmin` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`argmin` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`argmin` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value amax(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`amax` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`amax` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`amax` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`amax` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value countNonzero(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`countNonzero` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`countNonzero` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsArray()) {
		Napi::TypeError::New(info.Env(), "`countNonzero` expects args[1] to be typeof `number[]`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[2].IsBoolean()) {
		Napi::TypeError::New(info.Env(), "`countNonzero` expects args[2] to be typeof `boolean`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value tile(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`tile` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`tile` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value rint(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`rint` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`rint` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value absolute(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`absolute` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`absolute` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value erf(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`erf` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`erf` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value where(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	return env.Null();
}

static Napi::Value sort(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`sort` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`sort` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value sin(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`sin` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`sin` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value cos(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`cos` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`cos` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value power(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	return env.Null();
}

static Napi::Value fromScalar(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 2) {
		Napi::TypeError::New(info.Env(), "`fromScalar` expects exactly 2 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value iota(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 3) {
		Napi::TypeError::New(info.Env(), "`iota` expects exactly 3 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value floor(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`floor` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`floor` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value ceil(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`ceil` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`ceil` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value isinf(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`isinf` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`isinf` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value nonzero(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`nonzero` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`nonzero` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value negative(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 1) {
		Napi::TypeError::New(info.Env(), "`negative` expects exactly 1 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`negative` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value matmul(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	if (info.Length() != 4) {
		Napi::TypeError::New(info.Env(), "`matmul` expects exactly 4 args").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[0].IsObject()) {
		Napi::TypeError::New(info.Env(), "`matmul` expects args[0] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	if (!info[1].IsObject()) {
		Napi::TypeError::New(info.Env(), "`matmul` expects args[1] to be instanceof `Tensor`").ThrowAsJavaScriptException();
		return env.Null();
	}
	return env.Null();
}

static Napi::Value minimum(const Napi::CallbackInfo& info) {
	Napi::Env env = info.Env();
	return env.Null();
}

// NAPI exports
Napi::Object Init(Napi::Env env, Napi::Object exports) {
	exports.Set(Napi::String::New(env, "_log1p"), Napi::Function::New(env, log1p));
	exports.Set(Napi::String::New(env, "_sqrt"), Napi::Function::New(env, sqrt));
	exports.Set(Napi::String::New(env, "_argmax"), Napi::Function::New(env, argmax));
	exports.Set(Napi::String::New(env, "_logicalNot"), Napi::Function::New(env, logicalNot));
	exports.Set(Napi::String::New(env, "_tril"), Napi::Function::New(env, tril));
	exports.Set(Napi::String::New(env, "_cumsum"), Napi::Function::New(env, cumsum));
	exports.Set(Napi::String::New(env, "_sign"), Napi::Function::New(env, sign));
	exports.Set(Napi::String::New(env, "_triu"), Napi::Function::New(env, triu));
	exports.Set(Napi::String::New(env, "_argsort"), Napi::Function::New(env, argsort));
	exports.Set(Napi::String::New(env, "_mean"), Napi::Function::New(env, mean));
	exports.Set(Napi::String::New(env, "_argmin"), Napi::Function::New(env, argmin));
	exports.Set(Napi::String::New(env, "_sum"), Napi::Function::New(env, sum));
	exports.Set(Napi::String::New(env, "_clip"), Napi::Function::New(env, clip));
	exports.Set(Napi::String::New(env, "_roll"), Napi::Function::New(env, roll));
	exports.Set(Napi::String::New(env, "_isnan"), Napi::Function::New(env, isnan));
	exports.Set(Napi::String::New(env, "_sort"), Napi::Function::New(env, sort));
	exports.Set(Napi::String::New(env, "_amax"), Napi::Function::New(env, amax));
	exports.Set(Napi::String::New(env, "_countNonzero"), Napi::Function::New(env, countNonzero));
	exports.Set(Napi::String::New(env, "_tile"), Napi::Function::New(env, tile));
	exports.Set(Napi::String::New(env, "_rint"), Napi::Function::New(env, rint));
	exports.Set(Napi::String::New(env, "_absolute"), Napi::Function::New(env, absolute));
	exports.Set(Napi::String::New(env, "_erf"), Napi::Function::New(env, erf));
	exports.Set(Napi::String::New(env, "_where"), Napi::Function::New(env, where));
	exports.Set(Napi::String::New(env, "_sin"), Napi::Function::New(env, sin));
	exports.Set(Napi::String::New(env, "_cos"), Napi::Function::New(env, cos));
	exports.Set(Napi::String::New(env, "_power"), Napi::Function::New(env, power));
	exports.Set(Napi::String::New(env, "_fromScalar"), Napi::Function::New(env, fromScalar));
	exports.Set(Napi::String::New(env, "_iota"), Napi::Function::New(env, iota));
	exports.Set(Napi::String::New(env, "_floor"), Napi::Function::New(env, floor));
	exports.Set(Napi::String::New(env, "_ceil"), Napi::Function::New(env, ceil));
	exports.Set(Napi::String::New(env, "_isinf"), Napi::Function::New(env, isinf));
	exports.Set(Napi::String::New(env, "_nonzero"), Napi::Function::New(env, nonzero));
	exports.Set(Napi::String::New(env, "_negative"), Napi::Function::New(env, negative));
	exports.Set(Napi::String::New(env, "_matmul"), Napi::Function::New(env, matmul));
	exports.Set(Napi::String::New(env, "_minimum"), Napi::Function::New(env, minimum));
	exports.Set(Napi::String::New(env, "_amin"), Napi::Function::New(env, amin));
	exports.Set(Napi::String::New(env, "_norm"), Napi::Function::New(env, norm));
	exports.Set(Napi::String::New(env, "_any"), Napi::Function::New(env, any));
	exports.Set(Napi::String::New(env, "_all"), Napi::Function::New(env, all));
	exports.Set(Napi::String::New(env, "_std"), Napi::Function::New(env, std));
	exports.Set(Napi::String::New(env, "_exp"), Napi::Function::New(env, exp));
	exports.Set(Napi::String::New(env, "_tanh"), Napi::Function::New(env, tanh));
	exports.Set(Napi::String::New(env, "_scalar"), Napi::Function::New(env, scalar));
	exports.Set(Napi::String::New(env, "_arange"), Napi::Function::New(env, arange));
	exports.Set(Napi::String::New(env, "_maximum"), Napi::Function::New(env, maximum));
	exports.Set(Napi::String::New(env, "_transpose"), Napi::Function::New(env, transpose));
	exports.Set(Napi::String::New(env, "_var"), Napi::Function::New(env, var));
	exports.Set(Napi::String::New(env, "_full"), Napi::Function::New(env, full));
	exports.Set(Napi::String::New(env, "_concatenate"), Napi::Function::New(env, concatenate));
	exports.Set(Napi::String::New(env, "_pad"), Napi::Function::New(env, pad));
	exports.Set(Napi::String::New(env, "_sigmoid"), Napi::Function::New(env, sigmoid));
	exports.Set(Napi::String::New(env, "_median"), Napi::Function::New(env, median));
	exports.Set(Napi::String::New(env, "_identity"), Napi::Function::New(env, identity));
	exports.Set(Napi::String::New(env, "_reshape"), Napi::Function::New(env, reshape));
	exports.Set(Napi::String::New(env, "_log"), Napi::Function::New(env, log));
	exports.Set(Napi::String::New(env, "_flip"), Napi::Function::New(env, flip));
}

NODE_API_MODULE(addon, Init)
