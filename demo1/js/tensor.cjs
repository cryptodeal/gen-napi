/* eslint-disable @typescript-eslint/no-var-requires */
const {
  _add,
  _tensorFromUint8Buffer,
  _tensorFromFloat32Buffer,
  _tensorFromFloat64Buffer,
  _mean,
  _toFloat32Scalar,
  _toFloat64Array,
  _toUint8Array,
  _toFloat32Array
} = require('../../build/Release/shumai_bindings.node');

class Tensor {
  #native_self;

  constructor(t) {
    if (t instanceof Uint8Array || t.constructor === Uint8Array) {
      this.#native_self = _tensorFromUint8Buffer(t.buffer);
    } else if (t instanceof Float32Array || t.constructor === Float32Array) {
      this.#native_self = _tensorFromFloat32Buffer(t.buffer);
    } else if (t instanceof Float64Array || t.constructor === Float64Array) {
      this.#native_self = _tensorFromFloat64Buffer(t.buffer);
    } else {
      this.#native_self = t;
    }
  }

  get _native_self() {
    return this.#native_self;
  }

  add(other) {
    return new Tensor(_add(this._native_self, other._native_self));
  }

  mean(axes = [], keepDims = false) {
    return new Tensor(_mean(this._native_self, axes, keepDims));
  }

  toFloat32Scalar() {
    return _toFloat32Scalar(this._native_self);
  }

  toFloat32Array() {
    return _toFloat32Array(this._native_self);
  }

  toFloat64Array() {
    return _toFloat64Array(this._native_self);
  }

  toUint8Array() {
    return _toUint8Array(this._native_self);
  }
}

module.exports = { Tensor };
