// Code generated by gen-napi. DO NOT EDIT.
const { Tensor } = require('./tensor.cjs')
const {
  _toFloat32Array,
  _toFloat64Array,
  _toBoolInt8Array,
  _toInt16Array,
  _toInt32Array,
  _save,
  _toInt64Array,
  _toUint8Array,
  _toUint16Array,
  _toUint32Array,
  _toUint64Array,
  _toFloat32Scalar,
  _toFloat64Scalar,
  _toBoolInt8Scalar,
  _toInt16Scalar,
  _toInt32Scalar,
  _toInt64Scalar,
  _toUint8Scalar,
  _toUint16Scalar,
  _toUint32Scalar,
  _toUint64Scalar,
  _eval,
  _dispose,
  _tensorFromFloat32Array,
  _tensorFromFloat64Array,
  _tensorFromBoolInt8Array,
  _tensorFromInt16Array,
  _tensorFromInt32Array,
  _tensorFromInt64Array,
  _tensorFromUint8Array,
  _tensorFromUint16Array,
  _tensorFromUint32Array,
  _tensorFromUint64Array,
  _mean,
  _iota,
  _log1p,
  _sort,
  _countNonzero,
  _maximum,
  _power,
  _argmax,
  _tanh,
  _where,
  _minimum,
  _matmul,
  _amin,
  _identity,
  _transpose,
  _logicalNot,
  _median,
  _floor,
  _absolute,
  _sum,
  _var: __var,
  _any,
  _arange,
  _reshape,
  _nonzero,
  _roll,
  _sign,
  _triu,
  _tile,
  _log,
  _cos,
  _tril,
  _cumsum,
  _std,
  _all,
  _flip,
  _clip,
  _isinf,
  _full,
  _negative,
  _rint,
  _sqrt,
  _ceil,
  _sigmoid,
  _erf,
  _isnan,
  _concatenate,
  _exp,
  _sin,
  _amax,
  _argmin,
  _norm,
  _mul,
  _logicalOr,
  _bitwiseAnd,
  _bitwiseXor,
  _lShift,
  _add,
  _lessThan,
  _greaterThan,
  _logicalAnd,
  _greaterThanEqual,
  _eq,
  _mod,
  _lessThanEqual,
  _neq,
  _div,
  _sub,
  _bitwiseOr,
  _rShift,
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
} = require('../../build/Release/shumai_bindings.node')

const clip = (tensor, low, high) => {
  return new Tensor(_clip(tensor._native_self, low._native_self, high._native_self))
}

const isinf = (tensor) => {
  return new Tensor(_isinf(tensor._native_self))
}

const tril = (tensor) => {
  return new Tensor(_tril(tensor._native_self))
}

const cumsum = (input, axis) => {
  return new Tensor(_cumsum(input._native_self, axis))
}

const std = (input, axes, keepDims) => {
  return new Tensor(_std(input._native_self, axes, keepDims))
}

const all = (input, axes, keepDims) => {
  return new Tensor(_all(input._native_self, axes, keepDims))
}

const flip = (tensor, dim) => {
  return new Tensor(_flip(tensor._native_self, dim))
}

const negative = (tensor) => {
  return new Tensor(_negative(tensor._native_self))
}

const rint = (tensor) => {
  return new Tensor(_rint(tensor._native_self))
}

const full = (dims, val) => {
  return new Tensor(_full(dims, val))
}

const exp = (tensor) => {
  return new Tensor(_exp(tensor._native_self))
}

const sin = (tensor) => {
  return new Tensor(_sin(tensor._native_self))
}

const sqrt = (tensor) => {
  return new Tensor(_sqrt(tensor._native_self))
}

const ceil = (tensor) => {
  return new Tensor(_ceil(tensor._native_self))
}

const sigmoid = (tensor) => {
  return new Tensor(_sigmoid(tensor._native_self))
}

const erf = (tensor) => {
  return new Tensor(_erf(tensor._native_self))
}

const isnan = (tensor) => {
  return new Tensor(_isnan(tensor._native_self))
}

const concatenate = (tensors, axis) => {
  return new Tensor(_concatenate(tensors, axis))
}

const argmin = (input, axis, keepDims) => {
  return new Tensor(_argmin(input._native_self, axis, keepDims))
}

const norm = (input, axes, p, keepDims) => {
  return new Tensor(_norm(input._native_self, axes, p, keepDims))
}

const amax = (input, axes, keepDims) => {
  return new Tensor(_amax(input._native_self, axes, keepDims))
}

const log1p = (tensor) => {
  return new Tensor(_log1p(tensor._native_self))
}

const sort = (input, axis) => {
  return new Tensor(_sort(input._native_self, axis))
}

const mean = (input, axes, keepDims) => {
  return new Tensor(_mean(input._native_self, axes, keepDims))
}

const iota = (dims, tileDims) => {
  return new Tensor(_iota(dims, tileDims))
}

const power = (lhs, rhs) => {
  return new Tensor(_power(lhs._native_self, rhs._native_self))
}

const argmax = (input, axis, keepDims) => {
  return new Tensor(_argmax(input._native_self, axis, keepDims))
}

const countNonzero = (input, axes, keepDims) => {
  return new Tensor(_countNonzero(input._native_self, axes, keepDims))
}

const maximum = (lhs, rhs) => {
  return new Tensor(_maximum(lhs._native_self, rhs._native_self))
}

const transpose = (tensor, axes) => {
  return new Tensor(_transpose(tensor._native_self, axes))
}

const logicalNot = (tensor) => {
  return new Tensor(_logicalNot(tensor._native_self))
}

const tanh = (tensor) => {
  return new Tensor(_tanh(tensor._native_self))
}

const where = (condition, x, y) => {
  return new Tensor(_where(condition._native_self, x._native_self, y._native_self))
}

const minimum = (lhs, rhs) => {
  return new Tensor(_minimum(lhs._native_self, rhs._native_self))
}

const matmul = (lhs, rhs) => {
  return new Tensor(_matmul(lhs._native_self, rhs._native_self))
}

const amin = (input, axes, keepDims) => {
  return new Tensor(_amin(input._native_self, axes, keepDims))
}

const identity = (dim) => {
  return new Tensor(_identity(dim))
}

const median = (input, axes, keepDims) => {
  return new Tensor(_median(input._native_self, axes, keepDims))
}

const reshape = (tensor, shape) => {
  return new Tensor(_reshape(tensor._native_self, shape))
}

const nonzero = (tensor) => {
  return new Tensor(_nonzero(tensor._native_self))
}

const floor = (tensor) => {
  return new Tensor(_floor(tensor._native_self))
}

const absolute = (tensor) => {
  return new Tensor(_absolute(tensor._native_self))
}

const sum = (input, axes, keepDims) => {
  return new Tensor(_sum(input._native_self, axes, keepDims))
}

const _var = (input, axes, bias, keepDims) => {
  return new Tensor(__var(input._native_self, axes, bias, keepDims))
}

const any = (input, axes, keepDims) => {
  return new Tensor(_any(input._native_self, axes, keepDims))
}

const arange = (start, end, step) => {
  return new Tensor(_arange(start, end, step))
}

const log = (tensor) => {
  return new Tensor(_log(tensor._native_self))
}

const cos = (tensor) => {
  return new Tensor(_cos(tensor._native_self))
}

const roll = (tensor, shift, axis) => {
  return new Tensor(_roll(tensor._native_self, shift, axis))
}

const sign = (tensor) => {
  return new Tensor(_sign(tensor._native_self))
}

const triu = (tensor) => {
  return new Tensor(_triu(tensor._native_self))
}

const tile = (tensor, shape) => {
  return new Tensor(_tile(tensor._native_self, shape))
}

const div = (lhs, rhs) => {
  return new Tensor(_div(lhs._native_self, rhs._native_self))
}

const sub = (lhs, rhs) => {
  return new Tensor(_sub(lhs._native_self, rhs._native_self))
}

const lessThanEqual = (lhs, rhs) => {
  return new Tensor(_lessThanEqual(lhs._native_self, rhs._native_self))
}

const neq = (lhs, rhs) => {
  return new Tensor(_neq(lhs._native_self, rhs._native_self))
}

const bitwiseOr = (lhs, rhs) => {
  return new Tensor(_bitwiseOr(lhs._native_self, rhs._native_self))
}

const rShift = (lhs, rhs) => {
  return new Tensor(_rShift(lhs._native_self, rhs._native_self))
}

const add = (lhs, rhs) => {
  return new Tensor(_add(lhs._native_self, rhs._native_self))
}

const lessThan = (lhs, rhs) => {
  return new Tensor(_lessThan(lhs._native_self, rhs._native_self))
}

const mul = (lhs, rhs) => {
  return new Tensor(_mul(lhs._native_self, rhs._native_self))
}

const logicalOr = (lhs, rhs) => {
  return new Tensor(_logicalOr(lhs._native_self, rhs._native_self))
}

const bitwiseAnd = (lhs, rhs) => {
  return new Tensor(_bitwiseAnd(lhs._native_self, rhs._native_self))
}

const bitwiseXor = (lhs, rhs) => {
  return new Tensor(_bitwiseXor(lhs._native_self, rhs._native_self))
}

const lShift = (lhs, rhs) => {
  return new Tensor(_lShift(lhs._native_self, rhs._native_self))
}

const eq = (lhs, rhs) => {
  return new Tensor(_eq(lhs._native_self, rhs._native_self))
}

const mod = (lhs, rhs) => {
  return new Tensor(_mod(lhs._native_self, rhs._native_self))
}

const greaterThan = (lhs, rhs) => {
  return new Tensor(_greaterThan(lhs._native_self, rhs._native_self))
}

const logicalAnd = (lhs, rhs) => {
  return new Tensor(_logicalAnd(lhs._native_self, rhs._native_self))
}

const greaterThanEqual = (lhs, rhs) => {
  return new Tensor(_greaterThanEqual(lhs._native_self, rhs._native_self))
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
  iota,
  log1p,
  sort,
  mean,
  maximum,
  power,
  argmax,
  countNonzero,
  matmul,
  amin,
  identity,
  transpose,
  logicalNot,
  tanh,
  where,
  minimum,
  median,
  _var,
  any,
  arange,
  reshape,
  nonzero,
  floor,
  absolute,
  sum,
  tile,
  log,
  cos,
  roll,
  sign,
  triu,
  all,
  flip,
  clip,
  isinf,
  tril,
  cumsum,
  std,
  full,
  negative,
  rint,
  erf,
  isnan,
  concatenate,
  exp,
  sin,
  sqrt,
  ceil,
  sigmoid,
  amax,
  argmin,
  norm,
  eq,
  mod,
  greaterThan,
  logicalAnd,
  greaterThanEqual,
  div,
  sub,
  lessThanEqual,
  neq,
  bitwiseOr,
  rShift,
  lShift,
  add,
  lessThan,
  mul,
  logicalOr,
  bitwiseAnd,
  bitwiseXor,
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
