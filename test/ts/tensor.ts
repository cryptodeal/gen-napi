import { add, mean, tensorFromFloat32Array, toFloat32Scalar } from '.'

export class Tensor {
  #native_self: any
  constructor(t) {
    if (t instanceof Float32Array || t.constructor === Float32Array) {
      this.#native_self = tensorFromFloat32Array(t)
    } else {
      this.#native_self = t
    }
  }
  get _native_self() {
    return this.#native_self
  }

  add(other) {
    return new Tensor(add(this, other))
  }

  mean(axes = [], keepDims = false) {
    return new Tensor(mean(this, axes, keepDims))
  }

  toFloat32Scalar() {
    return toFloat32Scalar(this)
  }
}
