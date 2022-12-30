const { _Tensor } = require('./build/Release/flashlight_napi_bindings.node')

class Tensor {
  constructor(t) {
    this.addon = new _Tensor(t)
  }

  get ndim() {
    return this.addon.ndim()
  }

  copy() {
    return new Tensor(this.addon.copy())
  }

  strides() {
    return this.addon.strides()
  }

  get elements() {
    return this.addon.elements()
  }
}

const test = new Tensor(new Float32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]))
console.log(test.elements)
console.log(test.copy().ndim)
console.log(test.strides())
