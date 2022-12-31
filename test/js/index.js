// Code generated by gen-napi. DO NOT EDIT.
const {
  _Tensor,
  _clip,
  _isinf,
  _sort,
  _amin,
  _amax,
  _var,
  _concatenate,
  _erf,
  _std,
  _absolute,
  _tril,
  _triu,
  _minimum,
  _matmul,
  _sin,
  _ceil,
  _flip,
  _sign,
  _maximum,
  _power,
  _identity,
  _nonzero,
  _mean,
  _countNonzero,
  _tile,
  _argmin,
  _cos,
  _sqrt,
  _median,
  _iota,
  _exp,
  _log1p,
  _floor,
  _sigmoid,
  _roll,
  _isnan,
  _sum,
  _transpose,
  _negative,
  _any,
  _all,
  _tanh,
  _rint,
  _argmax,
  _cumsum,
  _norm,
  _arange,
  _logicalNot,
  _log,
  _where,
  _full,
  _reshape,
  _mul,
  _bitwiseOr,
  _rShift,
  _mod,
  _logicalOr,
  _greaterThan,
  _eq,
  _logicalAnd,
  _div,
  _lessThanEqual,
  _greaterThanEqual,
  _sub,
  _lessThan,
  _bitwiseXor,
  _lShift,
  _add,
  _neq,
  _bitwiseAnd
} = require('../../build/Release/flashlight_napi_bindings.node')

class Tensor {
  #_native_self

  constructor(t) {
    this.#_native_self = new _Tensor(t)
  }

  sigmoid() {
    return this.#_native_self.sigmoid(this.#_native_self)
  }

  roll(shift, axis) {
    return this.#_native_self.roll(this.#_native_self, shift, axis)
  }

  isnan() {
    return this.#_native_self.isnan(this.#_native_self)
  }

  sum(axes, keepDims) {
    return this.#_native_self.sum(this.#_native_self, axes, keepDims)
  }

  transpose(axes) {
    return this.#_native_self.transpose(this.#_native_self, axes)
  }

  negative() {
    return this.#_native_self.negative(this.#_native_self)
  }

  log1p() {
    return this.#_native_self.log1p(this.#_native_self)
  }

  floor() {
    return this.#_native_self.floor(this.#_native_self)
  }

  any(axes, keepDims) {
    return this.#_native_self.any(this.#_native_self, axes, keepDims)
  }

  all(axes, keepDims) {
    return this.#_native_self.all(this.#_native_self, axes, keepDims)
  }

  cumsum(axis) {
    return this.#_native_self.cumsum(this.#_native_self, axis)
  }

  norm(axes, p, keepDims) {
    return this.#_native_self.norm(this.#_native_self, axes, p, keepDims)
  }

  logicalNot() {
    return this.#_native_self.logicalNot(this.#_native_self)
  }

  tanh() {
    return this.#_native_self.tanh(this.#_native_self)
  }

  rint() {
    return this.#_native_self.rint(this.#_native_self)
  }

  reshape(shape) {
    return this.#_native_self.reshape(this.#_native_self, shape)
  }

  log() {
    return this.#_native_self.log(this.#_native_self)
  }

  where(x, y) {
    return this.#_native_self.where(this.#_native_self, x, y)
  }

  sort(axis, sortMode) {
    return this.#_native_self.sort(this.#_native_self, axis, sortMode)
  }

  amin(axes, keepDims) {
    return this.#_native_self.amin(this.#_native_self, axes, keepDims)
  }

  amax(axes, keepDims) {
    return this.#_native_self.amax(this.#_native_self, axes, keepDims)
  }

  var(axes, bias, keepDims) {
    return this.#_native_self.var(this.#_native_self, axes, bias, keepDims)
  }

  erf() {
    return this.#_native_self.erf(this.#_native_self)
  }

  clip(low, high) {
    return this.#_native_self.clip(this.#_native_self, low, high)
  }

  isinf() {
    return this.#_native_self.isinf(this.#_native_self)
  }

  std(axes, keepDims) {
    return this.#_native_self.std(this.#_native_self, axes, keepDims)
  }

  triu() {
    return this.#_native_self.triu(this.#_native_self)
  }

  matmul(rhs, lhsProp, rhsProp) {
    return this.#_native_self.matmul(this.#_native_self, rhs, lhsProp, rhsProp)
  }

  sin() {
    return this.#_native_self.sin(this.#_native_self)
  }

  ceil() {
    return this.#_native_self.ceil(this.#_native_self)
  }

  absolute() {
    return this.#_native_self.absolute(this.#_native_self)
  }

  tril() {
    return this.#_native_self.tril(this.#_native_self)
  }

  maximum(rhs) {
    return this.#_native_self.maximum(this.#_native_self, rhs)
  }

  power(rhs) {
    return this.#_native_self.power(this.#_native_self, rhs)
  }

  nonzero() {
    return this.#_native_self.nonzero(this.#_native_self)
  }

  flip(dim) {
    return this.#_native_self.flip(this.#_native_self, dim)
  }

  sign() {
    return this.#_native_self.sign(this.#_native_self)
  }

  tile(shape) {
    return this.#_native_self.tile(this.#_native_self, shape)
  }

  argmin(axis, keepDims) {
    return this.#_native_self.argmin(this.#_native_self, axis, keepDims)
  }

  mean(axes, keepDims) {
    return this.#_native_self.mean(this.#_native_self, axes, keepDims)
  }

  countNonzero(axes, keepDims) {
    return this.#_native_self.countNonzero(this.#_native_self, axes, keepDims)
  }

  median(axes, keepDims) {
    return this.#_native_self.median(this.#_native_self, axes, keepDims)
  }

  exp() {
    return this.#_native_self.exp(this.#_native_self)
  }

  cos() {
    return this.#_native_self.cos(this.#_native_self)
  }

  sqrt() {
    return this.#_native_self.sqrt(this.#_native_self)
  }
}
