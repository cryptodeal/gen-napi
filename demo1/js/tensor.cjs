/* eslint-disable @typescript-eslint/no-var-requires */
const { _init } = require('../../build/Release/shumai_bindings.node');
const { _Base_Tensor, gen_Tensor_ops_shim } = require('./gen_Tensor_methods_shim.cjs');

_init();

class Tensor extends _Base_Tensor {
  #underlying;

  constructor(t) {
    super();

    if (t instanceof Float32Array || t.constructor === Float32Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromFloat32Array(t);
      return;
    }

    if (t instanceof Float64Array || t.constructor === Float64Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromFloat64Array(t);
      return;
    }

    if (t instanceof Int8Array || t.constructor === Int8Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromInt8Array(t);
      return;
    }

    if (t instanceof Int16Array || t.constructor === Int16Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromInt16Array(t);
      return;
    }

    if (t instanceof Int32Array || t.constructor === Int32Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromInt32Array(t);
      return;
    }

    if (t instanceof BigInt64Array || t.constructor === BigInt64Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromBigInt64Array(t);
      return;
    }

    if (t instanceof Uint8Array || t.constructor === Uint8Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromUint8Array(t);
      return;
    }

    if (t instanceof Uint16Array || t.constructor === Uint16Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromUint16Array(t);
      return;
    }

    if (t instanceof Uint32Array || t.constructor === Uint32Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromUint32Array(t);
      return;
    }

    if (t instanceof BigUint64Array || t.constructor === BigUint64Array) {
      this.#underlying = t;
      this._native_Tensor = this.tensorFromBigUint64Array(t);
      return;
    }

    this._native_Tensor = t;
  }
}

for (const [method, closure] of Object.entries(gen_Tensor_ops_shim(Tensor))) {
  Tensor.prototype[method] = closure;
}

module.exports = { Tensor };
