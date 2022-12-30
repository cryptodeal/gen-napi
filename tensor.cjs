const { _Tensor } = require('./build/Release/flashlight_napi_bindings.node')

class Tensor {
  constructor(t) {
    this.addon = new _Tensor(t)
  }
}

const test = new Tensor(new Float32Array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10]))
