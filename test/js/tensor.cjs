const {
  _add,
  _tensorFromFloat32Buffer,
  _mean,
  _toFloat32Scalar,
  _toFloat32Array
} = require('../../build/Release/shumai_bindings.node')

class Tensor {
  #native_self

  constructor(t) {
    if (t instanceof Float32Array || t.constructor === Float32Array) {
      this.#native_self = _tensorFromFloat32Buffer(t.buffer)
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

  toFloat32Array() {
    return _toFloat32Array(this._native_self)
  }
}

module.exports = { Tensor }
