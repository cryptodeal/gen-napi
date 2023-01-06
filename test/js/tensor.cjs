const {
  _add,
  _tensorFromFloat32Array,
  _mean,
  _toFloat32Scalar
} = require('../../build/Release/shumai_bindings.node')

class Tensor {
  #native_self

  constructor(t) {
    if (t instanceof Float32Array || t.constructor === Float32Array) {
      this.#native_self = _tensorFromFloat32Array(t)
    } else {
      this.#native_self = t
    }
  }

  get _native_self() {
    return this.#native_self
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
