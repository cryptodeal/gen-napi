// Code generated by gen-napi. DO NOT EDIT.
import { Tensor } from './tensor'
const {
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
	_randn,
	_any,
	_identity,
	_absolute,
	_sigmoid,
	_isnan,
	_isinf,
	_norm,
	_minimum,
	_maximum,
	_arange,
	_reshape,
	_log,
	_tanh,
	_sign,
	_tril,
	_amin,
	_mean,
	_full,
	_concatenate,
	_negative,
	_exp,
	_sum,
	_var: __var,
	_matmul,
	_amax,
	_clip,
	_power,
	_all,
	_argmin,
	_cumsum,
	_tile,
	_cos,
	_floor,
	_rint,
	_roll,
	_where,
	_std,
	_sort,
	_argmax,
	_iota,
	_transpose,
	_nonzero,
	_logicalNot,
	_log1p,
	_erf,
	_median,
	_sin,
	_sqrt,
	_ceil,
	_flip,
	_triu,
	_countNonzero,
	_mul,
	_logicalAnd,
	_rShift,
	_sub,
	_neq,
	_greaterThanEqual,
	_bitwiseAnd,
	_logicalOr,
	_mod,
	_lessThan,
	_lessThanEqual,
	_div,
	_eq,
	_bitwiseOr,
	_lShift,
	_greaterThan,
	_add,
	_bitwiseXor,
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
} = require("../../build/Release/shumai_bindings.node")

export const reshape = (tensor: Tensor, shape: number[]): Tensor => {
	return new Tensor(_reshape(tensor._native_self, shape));
}

export const log = (tensor: Tensor): Tensor => {
	return new Tensor(_log(tensor._native_self));
}

export const tanh = (tensor: Tensor): Tensor => {
	return new Tensor(_tanh(tensor._native_self));
}

export const sign = (tensor: Tensor): Tensor => {
	return new Tensor(_sign(tensor._native_self));
}

export const tril = (tensor: Tensor): Tensor => {
	return new Tensor(_tril(tensor._native_self));
}

export const minimum = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_minimum(lhs._native_self, rhs._native_self));
}

export const maximum = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_maximum(lhs._native_self, rhs._native_self));
}

export const arange = (start: number, end: number, step: number): Tensor => {
	return new Tensor(_arange(start, end, step));
}

export const mean = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_mean(input._native_self, axes, keepDims));
}

export const amin = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_amin(input._native_self, axes, keepDims));
}

export const concatenate = (tensors: Tensor[], axis: number): Tensor => {
	return new Tensor(_concatenate(tensors, axis));
}

export const negative = (tensor: Tensor): Tensor => {
	return new Tensor(_negative(tensor._native_self));
}

export const exp = (tensor: Tensor): Tensor => {
	return new Tensor(_exp(tensor._native_self));
}

export const sum = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_sum(input._native_self, axes, keepDims));
}

export const _var = (input: Tensor, axes: number[], bias: boolean, keepDims: boolean): Tensor => {
	return new Tensor(__var(input._native_self, axes, bias, keepDims));
}

export const full = (dims: number[], val: number): Tensor => {
	return new Tensor(_full(dims, val));
}

export const amax = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_amax(input._native_self, axes, keepDims));
}

export const matmul = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_matmul(lhs._native_self, rhs._native_self));
}

export const power = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_power(lhs._native_self, rhs._native_self));
}

export const all = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_all(input._native_self, axes, keepDims));
}

export const clip = (tensor: Tensor, low: Tensor, high: Tensor): Tensor => {
	return new Tensor(_clip(tensor._native_self, low._native_self, high._native_self));
}

export const cos = (tensor: Tensor): Tensor => {
	return new Tensor(_cos(tensor._native_self));
}

export const floor = (tensor: Tensor): Tensor => {
	return new Tensor(_floor(tensor._native_self));
}

export const rint = (tensor: Tensor): Tensor => {
	return new Tensor(_rint(tensor._native_self));
}

export const roll = (tensor: Tensor, shift: number, axis: number): Tensor => {
	return new Tensor(_roll(tensor._native_self, shift, axis));
}

export const where = (condition: Tensor, x: Tensor, y: Tensor): Tensor => {
	return new Tensor(_where(condition._native_self, x._native_self, y._native_self));
}

export const argmin = (input: Tensor, axis: number, keepDims: boolean): Tensor => {
	return new Tensor(_argmin(input._native_self, axis, keepDims));
}

export const cumsum = (input: Tensor, axis: number): Tensor => {
	return new Tensor(_cumsum(input._native_self, axis));
}

export const tile = (tensor: Tensor, shape: number[]): Tensor => {
	return new Tensor(_tile(tensor._native_self, shape));
}

export const std = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_std(input._native_self, axes, keepDims));
}

export const transpose = (tensor: Tensor, axes: number[]): Tensor => {
	return new Tensor(_transpose(tensor._native_self, axes));
}

export const nonzero = (tensor: Tensor): Tensor => {
	return new Tensor(_nonzero(tensor._native_self));
}

export const logicalNot = (tensor: Tensor): Tensor => {
	return new Tensor(_logicalNot(tensor._native_self));
}

export const log1p = (tensor: Tensor): Tensor => {
	return new Tensor(_log1p(tensor._native_self));
}

export const erf = (tensor: Tensor): Tensor => {
	return new Tensor(_erf(tensor._native_self));
}

export const sort = (input: Tensor, axis: number): Tensor => {
	return new Tensor(_sort(input._native_self, axis));
}

export const argmax = (input: Tensor, axis: number, keepDims: boolean): Tensor => {
	return new Tensor(_argmax(input._native_self, axis, keepDims));
}

export const iota = (dims: number[], tileDims: number[]): Tensor => {
	return new Tensor(_iota(dims, tileDims));
}

export const median = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_median(input._native_self, axes, keepDims));
}

export const sqrt = (tensor: Tensor): Tensor => {
	return new Tensor(_sqrt(tensor._native_self));
}

export const ceil = (tensor: Tensor): Tensor => {
	return new Tensor(_ceil(tensor._native_self));
}

export const flip = (tensor: Tensor, dim: number): Tensor => {
	return new Tensor(_flip(tensor._native_self, dim));
}

export const triu = (tensor: Tensor): Tensor => {
	return new Tensor(_triu(tensor._native_self));
}

export const countNonzero = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_countNonzero(input._native_self, axes, keepDims));
}

export const sin = (tensor: Tensor): Tensor => {
	return new Tensor(_sin(tensor._native_self));
}

export const absolute = (tensor: Tensor): Tensor => {
	return new Tensor(_absolute(tensor._native_self));
}

export const sigmoid = (tensor: Tensor): Tensor => {
	return new Tensor(_sigmoid(tensor._native_self));
}

export const isnan = (tensor: Tensor): Tensor => {
	return new Tensor(_isnan(tensor._native_self));
}

export const isinf = (tensor: Tensor): Tensor => {
	return new Tensor(_isinf(tensor._native_self));
}

export const norm = (input: Tensor, axes: number[], p: number, keepDims: boolean): Tensor => {
	return new Tensor(_norm(input._native_self, axes, p, keepDims));
}

export const any = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_any(input._native_self, axes, keepDims));
}

export const identity = (dim: number): Tensor => {
	return new Tensor(_identity(dim));
}

export const init = (): void => {
	return _init();
}

export const bytesUsed = (): bigint => {
	return _bytesUsed();
}

export const setRowMajor = (): void => {
	return _setRowMajor();
}

export const setColMajor = (): void => {
	return _setColMajor();
}

export const isRowMajor = (): boolean => {
	return _isRowMajor();
}

export const isColMajor = (): boolean => {
	return _isColMajor();
}

export const dtypeFloat32 = (): number => {
	return _dtypeFloat32();
}

export const dtypeFloat64 = (): number => {
	return _dtypeFloat64();
}

export const dtypeBoolInt8 = (): number => {
	return _dtypeBoolInt8();
}

export const dtypeInt16 = (): number => {
	return _dtypeInt16();
}

export const dtypeInt32 = (): number => {
	return _dtypeInt32();
}

export const dtypeInt64 = (): number => {
	return _dtypeInt64();
}

export const dtypeUint8 = (): number => {
	return _dtypeUint8();
}

export const dtypeUint16 = (): number => {
	return _dtypeUint16();
}

export const dtypeUint32 = (): number => {
	return _dtypeUint32();
}

export const dtypeUint64 = (): number => {
	return _dtypeUint64();
}

export const rand = (shape: number[]): Tensor => {
	return new Tensor(_rand(shape));
}

export const randn = (shape: number[]): Tensor => {
	return new Tensor(_randn(shape));
}

export const div = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_div(lhs._native_self, rhs._native_self));
}

export const eq = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_eq(lhs._native_self, rhs._native_self));
}

export const bitwiseOr = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseOr(lhs._native_self, rhs._native_self));
}

export const lShift = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lShift(lhs._native_self, rhs._native_self));
}

export const greaterThan = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_greaterThan(lhs._native_self, rhs._native_self));
}

export const add = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_add(lhs._native_self, rhs._native_self));
}

export const bitwiseXor = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseXor(lhs._native_self, rhs._native_self));
}

export const logicalAnd = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_logicalAnd(lhs._native_self, rhs._native_self));
}

export const rShift = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_rShift(lhs._native_self, rhs._native_self));
}

export const sub = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_sub(lhs._native_self, rhs._native_self));
}

export const neq = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_neq(lhs._native_self, rhs._native_self));
}

export const greaterThanEqual = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_greaterThanEqual(lhs._native_self, rhs._native_self));
}

export const bitwiseAnd = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseAnd(lhs._native_self, rhs._native_self));
}

export const logicalOr = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_logicalOr(lhs._native_self, rhs._native_self));
}

export const mul = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_mul(lhs._native_self, rhs._native_self));
}

export const mod = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_mod(lhs._native_self, rhs._native_self));
}

export const lessThan = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lessThan(lhs._native_self, rhs._native_self));
}

export const lessThanEqual = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lessThanEqual(lhs._native_self, rhs._native_self));
}

export const init = (): void => {
	return _init();
}

export const bytesUsed = (): bigint => {
	return _bytesUsed();
}

export const setRowMajor = (): void => {
	return _setRowMajor();
}

export const setColMajor = (): void => {
	return _setColMajor();
}

export const isRowMajor = (): boolean => {
	return _isRowMajor();
}

export const isColMajor = (): boolean => {
	return _isColMajor();
}

export const dtypeFloat32 = (): number => {
	return _dtypeFloat32();
}

export const dtypeFloat64 = (): number => {
	return _dtypeFloat64();
}

export const dtypeBoolInt8 = (): number => {
	return _dtypeBoolInt8();
}

export const dtypeInt16 = (): number => {
	return _dtypeInt16();
}

export const dtypeInt32 = (): number => {
	return _dtypeInt32();
}

export const dtypeInt64 = (): number => {
	return _dtypeInt64();
}

export const dtypeUint8 = (): number => {
	return _dtypeUint8();
}

export const dtypeUint16 = (): number => {
	return _dtypeUint16();
}

export const dtypeUint32 = (): number => {
	return _dtypeUint32();
}

export const dtypeUint64 = (): number => {
	return _dtypeUint64();
}

export const rand = (shape: number[]): Tensor => {
	return new Tensor(_rand(shape));
}

export const randn = (shape: number[]): Tensor => {
	return new Tensor(_randn(shape));
}

