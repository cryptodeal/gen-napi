// Code generated by gen-napi. DO NOT EDIT.
#pragma once
#include <napi.h>
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

//////////
// source: TensorBase.h

class Tensor : public Napi::ObjectWrap<Tensor> {
 public:
  Tensor(const Napi::CallbackInfo&);
  static Napi::FunctionReference* constructor;
  fl::Tensor* _tensor;
  static Napi::Function GetClass(Napi::Env);

  // methods defined in src, wrapped as class methods
  Napi::Value copy(const Napi::CallbackInfo&);
  Napi::Value shape(const Napi::CallbackInfo&);
  Napi::Value elements(const Napi::CallbackInfo&);
  Napi::Value ndim(const Napi::CallbackInfo&);
  Napi::Value isEmpty(const Napi::CallbackInfo&);
  Napi::Value bytes(const Napi::CallbackInfo&);
  Napi::Value type(const Napi::CallbackInfo&);
  Napi::Value isSparse(const Napi::CallbackInfo&);
  Napi::Value strides(const Napi::CallbackInfo&);
  Napi::Value astype(const Napi::CallbackInfo&);
  Napi::Value flatten(const Napi::CallbackInfo&);
  Napi::Value asContiguousTensor(const Napi::CallbackInfo&);
  Napi::Value isContiguous(const Napi::CallbackInfo&);
  Napi::Value tril(const Napi::CallbackInfo&);
  Napi::Value power(const Napi::CallbackInfo&);
  Napi::Value amin(const Napi::CallbackInfo&);
  Napi::Value reshape(const Napi::CallbackInfo&);
  Napi::Value erf(const Napi::CallbackInfo&);
  Napi::Value sign(const Napi::CallbackInfo&);
  Napi::Value argmin(const Napi::CallbackInfo&);
  Napi::Value median(const Napi::CallbackInfo&);
  Napi::Value countNonzero(const Napi::CallbackInfo&);
  Napi::Value cos(const Napi::CallbackInfo&);
  Napi::Value clip(const Napi::CallbackInfo&);
  Napi::Value sort(const Napi::CallbackInfo&);
  Napi::Value cumsum(const Napi::CallbackInfo&);
  Napi::Value std(const Napi::CallbackInfo&);
  Napi::Value logicalNot(const Napi::CallbackInfo&);
  Napi::Value log(const Napi::CallbackInfo&);
  Napi::Value sin(const Napi::CallbackInfo&);
  Napi::Value norm(const Napi::CallbackInfo&);
  Napi::Value all(const Napi::CallbackInfo&);
  Napi::Value sigmoid(const Napi::CallbackInfo&);
  Napi::Value roll(const Napi::CallbackInfo&);
  Napi::Value mean(const Napi::CallbackInfo&);
  Napi::Value amax(const Napi::CallbackInfo&);
  Napi::Value floor(const Napi::CallbackInfo&);
  Napi::Value maximum(const Napi::CallbackInfo&);
  Napi::Value log1p(const Napi::CallbackInfo&);
  Napi::Value isinf(const Napi::CallbackInfo&);
  Napi::Value matmul(const Napi::CallbackInfo&);
  Napi::Value sum(const Napi::CallbackInfo&);
  Napi::Value tile(const Napi::CallbackInfo&);
  Napi::Value exp(const Napi::CallbackInfo&);
  Napi::Value ceil(const Napi::CallbackInfo&);
  Napi::Value absolute(const Napi::CallbackInfo&);
  Napi::Value isnan(const Napi::CallbackInfo&);
  Napi::Value var(const Napi::CallbackInfo&);
  Napi::Value any(const Napi::CallbackInfo&);
  Napi::Value sqrt(const Napi::CallbackInfo&);
  Napi::Value tanh(const Napi::CallbackInfo&);
  Napi::Value rint(const Napi::CallbackInfo&);
  Napi::Value flip(const Napi::CallbackInfo&);
  Napi::Value triu(const Napi::CallbackInfo&);
  Napi::Value where(const Napi::CallbackInfo&);
  Napi::Value transpose(const Napi::CallbackInfo&);
  Napi::Value nonzero(const Napi::CallbackInfo&);
  Napi::Value negative(const Napi::CallbackInfo&);
  Napi::Value lessThan(const Napi::CallbackInfo&);
  Napi::Value add(const Napi::CallbackInfo&);
  Napi::Value lessThanEqual(const Napi::CallbackInfo&);
  Napi::Value bitwiseOr(const Napi::CallbackInfo&);
  Napi::Value neq(const Napi::CallbackInfo&);
  Napi::Value greaterThan(const Napi::CallbackInfo&);
  Napi::Value eq(const Napi::CallbackInfo&);
  Napi::Value lShift(const Napi::CallbackInfo&);
  Napi::Value mod(const Napi::CallbackInfo&);
  Napi::Value div(const Napi::CallbackInfo&);
  Napi::Value rShift(const Napi::CallbackInfo&);
  Napi::Value logicalOr(const Napi::CallbackInfo&);
  Napi::Value greaterThanEqual(const Napi::CallbackInfo&);
  Napi::Value logicalAnd(const Napi::CallbackInfo&);
  Napi::Value bitwiseAnd(const Napi::CallbackInfo&);
  Napi::Value bitwiseXor(const Napi::CallbackInfo&);
  Napi::Value mul(const Napi::CallbackInfo&);
  Napi::Value sub(const Napi::CallbackInfo&);
  Napi::Value toFloat32Array(const Napi::CallbackInfo&);
  Napi::Value toFloat64Array(const Napi::CallbackInfo&);
  Napi::Value toBoolInt8Array(const Napi::CallbackInfo&);
  Napi::Value toInt16Array(const Napi::CallbackInfo&);
  Napi::Value toInt32Array(const Napi::CallbackInfo&);
  Napi::Value toInt64Array(const Napi::CallbackInfo&);
  Napi::Value toUint8Array(const Napi::CallbackInfo&);
  Napi::Value toUint16Array(const Napi::CallbackInfo&);
  Napi::Value toUint32Array(const Napi::CallbackInfo&);
  Napi::Value toUint64Array(const Napi::CallbackInfo&);
  Napi::Value toFloat32Scalar(const Napi::CallbackInfo&);
  Napi::Value toFloat64Scalar(const Napi::CallbackInfo&);
  Napi::Value toBoolInt8Scalar(const Napi::CallbackInfo&);
  Napi::Value toInt16Scalar(const Napi::CallbackInfo&);
  Napi::Value toInt32Scalar(const Napi::CallbackInfo&);
  Napi::Value toInt64Scalar(const Napi::CallbackInfo&);
  Napi::Value toUint8Scalar(const Napi::CallbackInfo&);
  Napi::Value toUint16Scalar(const Napi::CallbackInfo&);
  Napi::Value toUint32Scalar(const Napi::CallbackInfo&);
  Napi::Value toUint64Scalar(const Napi::CallbackInfo&);
  void eval(const Napi::CallbackInfo&);
  void dispose(const Napi::CallbackInfo&);

 private:
};