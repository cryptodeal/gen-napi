const {
  _add,
  _tensorFromFloat32Array,
  _mean,
  _toFloat32Scalar
} = require('../../build/Release/shumai_bindings.node')

class Tensor {
  constructor(t) {
    if (t instanceof Float32Array || t.constructor === Float32Array) {
      this._native_self = _tensorFromFloat32Array(t)
    } else {
      this._native_self = t
    }
  }

  add(other) {
    return new Tensor(_add(this._native_self, other._native_self))
  }

  mean(axes = [], keepDims = false) {
    return new Tensor(_mean(this._native_self, axes, keepDims))
  }

  toFloat32Scalar() {
    return _toFloat32Scalar(this._native_self)
  }
}

module.exports = { Tensor }
