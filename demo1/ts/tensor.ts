const { _add, _tensorFromFloat32Buffer, _mean, _toFloat32Scalar, _toUint8Array } =
  import.meta.require('../../build/Release/shumai_bindings.node');

export class Tensor {
  #native_self: any;
  #underlying;

  constructor(t) {
    if (t instanceof Float32Array || t.constructor === Float32Array) {
      this.#underlying = t;
      this.#native_self = _tensorFromFloat32Buffer(t.buffer);
    } else {
      this.#native_self = t;
    }
  }

  get _native_self() {
    return this.#native_self;
  }

  add(other: Tensor) {
    return new Tensor(_add(this._native_self, other._native_self));
  }

  mean(axes: number[] = [], keepDims = false) {
    return new Tensor(_mean(this._native_self, axes, keepDims));
  }

  toFloat32Scalar(): number {
    return _toFloat32Scalar(this._native_self);
  }

  toUint8Array(): Uint8Array {
    return _toUint8Array(this._native_self);
  }
}
