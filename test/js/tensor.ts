const {
  _add,
  _tensorFromFloat32Array,
  _mean,
  _toFloat32Scalar
} = require('../../build/Release/shumai_bindings.node')

export class Tensor {
  public _native_self: any
  constructor(t) {
    if (t instanceof Float32Array || t.constructor === Float32Array) {
      this._native_self = _tensorFromFloat32Array(t)
    } else {
      this._native_self = t
    }
  }

  add(other: Tensor): Tensor {
    return new Tensor(_add(this._native_self, other._native_self))
  }

  mean(axes = [], keepDims = false): Tensor {
    return new Tensor(_mean(this._native_self, axes, keepDims))
  }

  toFloat32Scalar(): number {
    return _toFloat32Scalar(this._native_self)
  }
}
