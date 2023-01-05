// Code generated by gen-napi. DO NOT EDIT.
const {
  _Tensor,
  _full,
  _arange,
  _concatenate,
  _floor,
  _flip,
  _triu,
  _power,
  _cumsum,
  _negative,
  _sqrt,
  _rint,
  _where,
  _minimum,
  _maximum,
  _argmax,
  _iota,
  _logicalNot,
  _sin,
  _roll,
  _sign,
  _sum,
  _median,
  _tile,
  _cos,
  _absolute,
  _sigmoid,
  _sort,
  _countNonzero,
  _transpose,
  _exp,
  _log1p,
  _clip,
  _tril,
  _amin,
  _amax,
  _var: __var,
  _std,
  _all,
  _nonzero,
  _ceil,
  _erf,
  _mean,
  _norm,
  _any,
  _isnan,
  _identity,
  _reshape,
  _log,
  _tanh,
  _isinf,
  _matmul,
  _argmin,
  _bitwiseOr,
  _greaterThan,
  _logicalAnd,
  _rShift,
  _neq,
  _lessThan,
  _lessThanEqual,
  _greaterThanEqual,
  _logicalOr,
  _bitwiseAnd,
  _div,
  _eq,
  _bitwiseXor,
  _add,
  _mod,
  _sub,
  _mul,
  _lShift,
  _init,
  _bytesUsed,
  _setRowMajor,
  _setColMajor,
  _isRowMajor,
  _isColMajor,
  _dtypeFloat32,
  _dtypeFloat64,
  _dtypeBoolInt8,
  _dtypeInt16,
  _dtypeInt32,
  _dtypeInt64,
  _dtypeUint8,
  _dtypeUint16,
  _dtypeUint32,
  _dtypeUint64,
  _rand,
  _randn
} = require('../../build/Release/flashlight_napi_bindings.node')

class Tensor {
  #_native_self

  constructor(t) {
    this.#_native_self = new _Tensor(t)
  }

  get _native_self() {
    return this.#_native_self
  }

  countNonzero(axes, keepDims) {
    return new Tensor(this.#_native_self.countNonzero(axes, keepDims))
  }

  tile(shape) {
    return new Tensor(this.#_native_self.tile(shape))
  }

  cos() {
    return new Tensor(this.#_native_self.cos())
  }

  absolute() {
    return new Tensor(this.#_native_self.absolute())
  }

  sigmoid() {
    return new Tensor(this.#_native_self.sigmoid())
  }

  sort(axis) {
    return new Tensor(this.#_native_self.sort(axis))
  }

  amin(axes, keepDims) {
    return new Tensor(this.#_native_self.amin(axes, keepDims))
  }

  amax(axes, keepDims) {
    return new Tensor(this.#_native_self.amax(axes, keepDims))
  }

  var(axes, bias, keepDims) {
    return new Tensor(this.#_native_self.var(axes, bias, keepDims))
  }

  transpose(axes) {
    return new Tensor(this.#_native_self.transpose(axes))
  }

  exp() {
    return new Tensor(this.#_native_self.exp())
  }

  log1p() {
    return new Tensor(this.#_native_self.log1p())
  }

  clip(low, high) {
    return new Tensor(this.#_native_self.clip(low.#_native_self, high.#_native_self))
  }

  tril() {
    return new Tensor(this.#_native_self.tril())
  }

  std(axes, keepDims) {
    return new Tensor(this.#_native_self.std(axes, keepDims))
  }

  all(axes, keepDims) {
    return new Tensor(this.#_native_self.all(axes, keepDims))
  }

  any(axes, keepDims) {
    return new Tensor(this.#_native_self.any(axes, keepDims))
  }

  nonzero() {
    return new Tensor(this.#_native_self.nonzero())
  }

  ceil() {
    return new Tensor(this.#_native_self.ceil())
  }

  erf() {
    return new Tensor(this.#_native_self.erf())
  }

  mean(axes, keepDims) {
    return new Tensor(this.#_native_self.mean(axes, keepDims))
  }

  norm(axes, p, keepDims) {
    return new Tensor(this.#_native_self.norm(axes, p, keepDims))
  }

  isnan() {
    return new Tensor(this.#_native_self.isnan())
  }

  matmul(rhs) {
    return new Tensor(this.#_native_self.matmul(rhs.#_native_self))
  }

  argmin(axis, keepDims) {
    return new Tensor(this.#_native_self.argmin(axis, keepDims))
  }

  reshape(shape) {
    return new Tensor(this.#_native_self.reshape(shape))
  }

  log() {
    return new Tensor(this.#_native_self.log())
  }

  tanh() {
    return new Tensor(this.#_native_self.tanh())
  }

  isinf() {
    return new Tensor(this.#_native_self.isinf())
  }

  triu() {
    return new Tensor(this.#_native_self.triu())
  }

  power(rhs) {
    return new Tensor(this.#_native_self.power(rhs.#_native_self))
  }

  cumsum(axis) {
    return new Tensor(this.#_native_self.cumsum(axis))
  }

  floor() {
    return new Tensor(this.#_native_self.floor())
  }

  flip(dim) {
    return new Tensor(this.#_native_self.flip(dim))
  }

  maximum(rhs) {
    return new Tensor(this.#_native_self.maximum(rhs.#_native_self))
  }

  negative() {
    return new Tensor(this.#_native_self.negative())
  }

  sqrt() {
    return new Tensor(this.#_native_self.sqrt())
  }

  rint() {
    return new Tensor(this.#_native_self.rint())
  }

  where(x, y) {
    return new Tensor(this.#_native_self.where(x.#_native_self, y.#_native_self))
  }

  sum(axes, keepDims) {
    return new Tensor(this.#_native_self.sum(axes, keepDims))
  }

  median(axes, keepDims) {
    return new Tensor(this.#_native_self.median(axes, keepDims))
  }

  logicalNot() {
    return new Tensor(this.#_native_self.logicalNot())
  }

  sin() {
    return new Tensor(this.#_native_self.sin())
  }

  roll(shift, axis) {
    return new Tensor(this.#_native_self.roll(shift, axis))
  }

  sign() {
    return new Tensor(this.#_native_self.sign())
  }

  greaterThan(rhs) {
    return new Tensor(this.#_native_self.greaterThan(rhs.#_native_self))
  }

  logicalAnd(rhs) {
    return new Tensor(this.#_native_self.logicalAnd(rhs.#_native_self))
  }

  rShift(rhs) {
    return new Tensor(this.#_native_self.rShift(rhs.#_native_self))
  }

  logicalOr(rhs) {
    return new Tensor(this.#_native_self.logicalOr(rhs.#_native_self))
  }

  bitwiseAnd(rhs) {
    return new Tensor(this.#_native_self.bitwiseAnd(rhs.#_native_self))
  }

  div(rhs) {
    return new Tensor(this.#_native_self.div(rhs.#_native_self))
  }

  neq(rhs) {
    return new Tensor(this.#_native_self.neq(rhs.#_native_self))
  }

  lessThan(rhs) {
    return new Tensor(this.#_native_self.lessThan(rhs.#_native_self))
  }

  lessThanEqual(rhs) {
    return new Tensor(this.#_native_self.lessThanEqual(rhs.#_native_self))
  }

  greaterThanEqual(rhs) {
    return new Tensor(this.#_native_self.greaterThanEqual(rhs.#_native_self))
  }

  sub(rhs) {
    return new Tensor(this.#_native_self.sub(rhs.#_native_self))
  }

  mul(rhs) {
    return new Tensor(this.#_native_self.mul(rhs.#_native_self))
  }

  lShift(rhs) {
    return new Tensor(this.#_native_self.lShift(rhs.#_native_self))
  }

  eq(rhs) {
    return new Tensor(this.#_native_self.eq(rhs.#_native_self))
  }

  bitwiseXor(rhs) {
    return new Tensor(this.#_native_self.bitwiseXor(rhs.#_native_self))
  }

  add(rhs) {
    return new Tensor(this.#_native_self.add(rhs.#_native_self))
  }

  mod(rhs) {
    return new Tensor(this.#_native_self.mod(rhs.#_native_self))
  }

  bitwiseOr(rhs) {
    return new Tensor(this.#_native_self.bitwiseOr(rhs.#_native_self))
  }

  copy() {
    return new Tensor(this.#_native_self.copy())
  }

  shape() {
    return this.#_native_self.shape()
  }

  elements() {
    return this.#_native_self.elements()
  }

  ndim() {
    return this.#_native_self.ndim()
  }

  isEmpty() {
    return this.#_native_self.isEmpty()
  }

  bytes() {
    return this.#_native_self.bytes()
  }

  type() {
    return this.#_native_self.type()
  }

  isSparse() {
    return this.#_native_self.isSparse()
  }

  strides() {
    return this.#_native_self.strides()
  }

  astype() {
    return new Tensor(this.#_native_self.astype())
  }

  flatten() {
    return new Tensor(this.#_native_self.flatten())
  }

  asContiguousTensor() {
    return new Tensor(this.#_native_self.asContiguousTensor())
  }

  isContiguous() {
    return this.#_native_self.isContiguous()
  }

  toFloat32Array() {
    return this.#_native_self.toFloat32Array()
  }

  toFloat64Array() {
    return this.#_native_self.toFloat64Array()
  }

  toBoolInt8Array() {
    return this.#_native_self.toBoolInt8Array()
  }

  toInt16Array() {
    return this.#_native_self.toInt16Array()
  }

  toInt32Array() {
    return this.#_native_self.toInt32Array()
  }

  save(filename) {
    return this.#_native_self.save(filename)
  }

  toInt64Array() {
    return this.#_native_self.toInt64Array()
  }

  toUint8Array() {
    return this.#_native_self.toUint8Array()
  }

  toUint16Array() {
    return this.#_native_self.toUint16Array()
  }

  toUint32Array() {
    return this.#_native_self.toUint32Array()
  }

  toUint64Array() {
    return this.#_native_self.toUint64Array()
  }

  toFloat32Scalar() {
    return this.#_native_self.toFloat32Scalar()
  }

  toFloat64Scalar() {
    return this.#_native_self.toFloat64Scalar()
  }

  toBoolInt8Scalar() {
    return this.#_native_self.toBoolInt8Scalar()
  }

  toInt16Scalar() {
    return this.#_native_self.toInt16Scalar()
  }

  toInt32Scalar() {
    return this.#_native_self.toInt32Scalar()
  }

  toInt64Scalar() {
    return this.#_native_self.toInt64Scalar()
  }

  toUint8Scalar() {
    return this.#_native_self.toUint8Scalar()
  }

  toUint16Scalar() {
    return this.#_native_self.toUint16Scalar()
  }

  toUint32Scalar() {
    return this.#_native_self.toUint32Scalar()
  }

  toUint64Scalar() {
    return this.#_native_self.toUint64Scalar()
  }

  eval() {
    return this.#_native_self.eval()
  }

  dispose() {
    return this.#_native_self.dispose()
  }
}

const isnan = (tensor) => {
  return new Tensor(_isnan(tensor._native_self))
}

const log = (tensor) => {
  return new Tensor(_log(tensor._native_self))
}

const tanh = (tensor) => {
  return new Tensor(_tanh(tensor._native_self))
}

const isinf = (tensor) => {
  return new Tensor(_isinf(tensor._native_self))
}

const matmul = (lhs, rhs) => {
  return new Tensor(_matmul(lhs._native_self, rhs._native_self))
}

const argmin = (input, axis, keepDims) => {
  return new Tensor(_argmin(input._native_self, axis, keepDims))
}

const identity = (dim) => {
  return new Tensor(_identity(dim))
}

const reshape = (tensor, shape) => {
  return new Tensor(_reshape(tensor._native_self, shape))
}

const concatenate = (tensors, axis) => {
  return new Tensor(_concatenate(tensors, axis))
}

const floor = (tensor) => {
  return new Tensor(_floor(tensor._native_self))
}

const flip = (tensor, dim) => {
  return new Tensor(_flip(tensor._native_self, dim))
}

const triu = (tensor) => {
  return new Tensor(_triu(tensor._native_self))
}

const power = (lhs, rhs) => {
  return new Tensor(_power(lhs._native_self, rhs._native_self))
}

const cumsum = (input, axis) => {
  return new Tensor(_cumsum(input._native_self, axis))
}

const full = (dims, val) => {
  return new Tensor(_full(dims, val))
}

const arange = (start, end, step) => {
  return new Tensor(_arange(start, end, step))
}

const rint = (tensor) => {
  return new Tensor(_rint(tensor._native_self))
}

const where = (condition, x, y) => {
  return new Tensor(_where(condition._native_self, x._native_self, y._native_self))
}

const minimum = (lhs, rhs) => {
  return new Tensor(_minimum(lhs._native_self, rhs._native_self))
}

const maximum = (lhs, rhs) => {
  return new Tensor(_maximum(lhs._native_self, rhs._native_self))
}

const argmax = (input, axis, keepDims) => {
  return new Tensor(_argmax(input._native_self, axis, keepDims))
}

const negative = (tensor) => {
  return new Tensor(_negative(tensor._native_self))
}

const sqrt = (tensor) => {
  return new Tensor(_sqrt(tensor._native_self))
}

const sin = (tensor) => {
  return new Tensor(_sin(tensor._native_self))
}

const roll = (tensor, shift, axis) => {
  return new Tensor(_roll(tensor._native_self, shift, axis))
}

const sign = (tensor) => {
  return new Tensor(_sign(tensor._native_self))
}

const sum = (input, axes, keepDims) => {
  return new Tensor(_sum(input._native_self, axes, keepDims))
}

const median = (input, axes, keepDims) => {
  return new Tensor(_median(input._native_self, axes, keepDims))
}

const iota = (dims, tileDims) => {
  return new Tensor(_iota(dims, tileDims))
}

const logicalNot = (tensor) => {
  return new Tensor(_logicalNot(tensor._native_self))
}

const absolute = (tensor) => {
  return new Tensor(_absolute(tensor._native_self))
}

const sigmoid = (tensor) => {
  return new Tensor(_sigmoid(tensor._native_self))
}

const sort = (input, axis) => {
  return new Tensor(_sort(input._native_self, axis))
}

const countNonzero = (input, axes, keepDims) => {
  return new Tensor(_countNonzero(input._native_self, axes, keepDims))
}

const tile = (tensor, shape) => {
  return new Tensor(_tile(tensor._native_self, shape))
}

const cos = (tensor) => {
  return new Tensor(_cos(tensor._native_self))
}

const log1p = (tensor) => {
  return new Tensor(_log1p(tensor._native_self))
}

const clip = (tensor, low, high) => {
  return new Tensor(_clip(tensor._native_self, low._native_self, high._native_self))
}

const tril = (tensor) => {
  return new Tensor(_tril(tensor._native_self))
}

const amin = (input, axes, keepDims) => {
  return new Tensor(_amin(input._native_self, axes, keepDims))
}

const amax = (input, axes, keepDims) => {
  return new Tensor(_amax(input._native_self, axes, keepDims))
}

const _var = (input, axes, bias, keepDims) => {
  return new Tensor(__var(input._native_self, axes, bias, keepDims))
}

const transpose = (tensor, axes) => {
  return new Tensor(_transpose(tensor._native_self, axes))
}

const exp = (tensor) => {
  return new Tensor(_exp(tensor._native_self))
}

const std = (input, axes, keepDims) => {
  return new Tensor(_std(input._native_self, axes, keepDims))
}

const all = (input, axes, keepDims) => {
  return new Tensor(_all(input._native_self, axes, keepDims))
}

const erf = (tensor) => {
  return new Tensor(_erf(tensor._native_self))
}

const mean = (input, axes, keepDims) => {
  return new Tensor(_mean(input._native_self, axes, keepDims))
}

const norm = (input, axes, p, keepDims) => {
  return new Tensor(_norm(input._native_self, axes, p, keepDims))
}

const any = (input, axes, keepDims) => {
  return new Tensor(_any(input._native_self, axes, keepDims))
}

const nonzero = (tensor) => {
  return new Tensor(_nonzero(tensor._native_self))
}

const ceil = (tensor) => {
  return new Tensor(_ceil(tensor._native_self))
}

const mul = (lhs, rhs) => {
  return new Tensor(_mul(lhs._native_self, rhs._native_self))
}

const lShift = (lhs, rhs) => {
  return new Tensor(_lShift(lhs._native_self, rhs._native_self))
}

const eq = (lhs, rhs) => {
  return new Tensor(_eq(lhs._native_self, rhs._native_self))
}

const bitwiseXor = (lhs, rhs) => {
  return new Tensor(_bitwiseXor(lhs._native_self, rhs._native_self))
}

const add = (lhs, rhs) => {
  return new Tensor(_add(lhs._native_self, rhs._native_self))
}

const mod = (lhs, rhs) => {
  return new Tensor(_mod(lhs._native_self, rhs._native_self))
}

const sub = (lhs, rhs) => {
  return new Tensor(_sub(lhs._native_self, rhs._native_self))
}

const bitwiseOr = (lhs, rhs) => {
  return new Tensor(_bitwiseOr(lhs._native_self, rhs._native_self))
}

const logicalAnd = (lhs, rhs) => {
  return new Tensor(_logicalAnd(lhs._native_self, rhs._native_self))
}

const rShift = (lhs, rhs) => {
  return new Tensor(_rShift(lhs._native_self, rhs._native_self))
}

const greaterThan = (lhs, rhs) => {
  return new Tensor(_greaterThan(lhs._native_self, rhs._native_self))
}

const bitwiseAnd = (lhs, rhs) => {
  return new Tensor(_bitwiseAnd(lhs._native_self, rhs._native_self))
}

const div = (lhs, rhs) => {
  return new Tensor(_div(lhs._native_self, rhs._native_self))
}

const neq = (lhs, rhs) => {
  return new Tensor(_neq(lhs._native_self, rhs._native_self))
}

const lessThan = (lhs, rhs) => {
  return new Tensor(_lessThan(lhs._native_self, rhs._native_self))
}

const lessThanEqual = (lhs, rhs) => {
  return new Tensor(_lessThanEqual(lhs._native_self, rhs._native_self))
}

const greaterThanEqual = (lhs, rhs) => {
  return new Tensor(_greaterThanEqual(lhs._native_self, rhs._native_self))
}

const logicalOr = (lhs, rhs) => {
  return new Tensor(_logicalOr(lhs._native_self, rhs._native_self))
}

const init = () => {
  return _init()
}

const bytesUsed = () => {
  return _bytesUsed()
}

const setRowMajor = () => {
  return _setRowMajor()
}

const setColMajor = () => {
  return _setColMajor()
}

const isRowMajor = () => {
  return _isRowMajor()
}

const isColMajor = () => {
  return _isColMajor()
}

const dtypeFloat32 = () => {
  return _dtypeFloat32()
}

const dtypeFloat64 = () => {
  return _dtypeFloat64()
}

const dtypeBoolInt8 = () => {
  return _dtypeBoolInt8()
}

const dtypeInt16 = () => {
  return _dtypeInt16()
}

const dtypeInt32 = () => {
  return _dtypeInt32()
}

const dtypeInt64 = () => {
  return _dtypeInt64()
}

const dtypeUint8 = () => {
  return _dtypeUint8()
}

const dtypeUint16 = () => {
  return _dtypeUint16()
}

const dtypeUint32 = () => {
  return _dtypeUint32()
}

const dtypeUint64 = () => {
  return _dtypeUint64()
}

const rand = (shape) => {
  return new Tensor(_rand(shape))
}

const randn = (shape) => {
  return new Tensor(_randn(shape))
}

module.exports = {
  Tensor,
  full,
  arange,
  concatenate,
  floor,
  flip,
  triu,
  power,
  cumsum,
  negative,
  sqrt,
  rint,
  where,
  minimum,
  maximum,
  argmax,
  iota,
  logicalNot,
  sin,
  roll,
  sign,
  sum,
  median,
  tile,
  cos,
  absolute,
  sigmoid,
  sort,
  countNonzero,
  transpose,
  exp,
  log1p,
  clip,
  tril,
  amin,
  amax,
  _var,
  std,
  all,
  nonzero,
  ceil,
  erf,
  mean,
  norm,
  any,
  isnan,
  identity,
  reshape,
  log,
  tanh,
  isinf,
  matmul,
  argmin,
  bitwiseOr,
  rShift,
  greaterThan,
  logicalAnd,
  div,
  neq,
  lessThan,
  lessThanEqual,
  greaterThanEqual,
  logicalOr,
  bitwiseAnd,
  lShift,
  eq,
  bitwiseXor,
  add,
  mod,
  sub,
  mul,
  init,
  bytesUsed,
  setRowMajor,
  setColMajor,
  isRowMajor,
  isColMajor,
  dtypeFloat32,
  dtypeFloat64,
  dtypeBoolInt8,
  dtypeInt16,
  dtypeInt32,
  dtypeInt64,
  dtypeUint8,
  dtypeUint16,
  dtypeUint32,
  dtypeUint64,
  rand,
  randn
}
