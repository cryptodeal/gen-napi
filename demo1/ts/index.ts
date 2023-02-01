// Code generated by gen-napi. DO NOT EDIT.
/* eslint-disable */
import { Tensor } from './tensor'
const {
	_toFloat32Array,
	_toFloat64Array,
	_toBoolInt8Array,
	_toInt16Array,
	_toInt32Array,
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
	_eval: __eval,
	_dispose,
	_tensorFromFloat32Buffer,
	_tensorFromFloat64Buffer,
	_tensorFromBoolInt8Buffer,
	_tensorFromInt16Buffer,
	_tensorFromInt32Buffer,
	_tensorFromInt64Buffer,
	_tensorFromUint8Buffer,
	_tensorFromUint16Buffer,
	_tensorFromUint32Buffer,
	_tensorFromUint64Buffer,
	_save,
	_arange,
	_erf,
	_absolute,
	_sqrt,
	_tanh,
	_sign,
	_minimum,
	_argmin,
	_var: __var,
	_exp,
	_cos,
	_floor,
	_sigmoid,
	_roll,
	_transpose,
	_concatenate,
	_std,
	_isInvalidArray,
	_iota,
	_negative,
	_isnan,
	_matmul,
	_amax,
	_log1p,
	_sort,
	_maximum,
	_cumsum,
	_all,
	_identity,
	_log,
	_power,
	_amin,
	_mean,
	_median,
	_sin,
	_flip,
	_sum,
	_clip,
	_any,
	_full,
	_nonzero,
	_ceil,
	_rint,
	_tril,
	_where,
	_reshape,
	_norm,
	_logicalNot,
	_isinf,
	_triu,
	_argmax,
	_countNonzero,
	_tile,
	_sub,
	_logicalOr,
	_lessThanEqual,
	_add,
	_greaterThanEqual,
	_rShift,
	_eq,
	_greaterThan,
	_mod,
	_mul,
	_lShift,
	_logicalAnd,
	_div,
	_neq,
	_lessThan,
	_bitwiseAnd,
	_bitwiseXor,
	_bitwiseOr,
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
} = import.meta.require("../../../../build/Release/shumai_bindings.node")

export enum TensorBackendType {
	Stub = 0,
	Tracer = 1,
	ArrayFire = 2,
	OneDnn = 3,
	Jit = 4
}

export enum Location {
	Host = 0,
	Device = 1
}

export enum StorageType {
	Dense = 0,
	CSR = 1,
	CSC = 2,
	COO = 3
}

export enum PadType {
	Constant = 0,
	Edge = 1,
	Symmetric = 2
}

export enum SortMode {
	Descending = 0,
	Ascending = 1
}

export enum MatrixProperty {
	None = 0,
	Transpose = 1
}

export const erf = (tensor: Tensor): Tensor => {
	return new Tensor(_erf(tensor._native_self));
}

export const arange = (start: number, end: number, step: number): Tensor => {
	return new Tensor(_arange(start, end, step));
}

export const sqrt = (tensor: Tensor): Tensor => {
	return new Tensor(_sqrt(tensor._native_self));
}

export const absolute = (tensor: Tensor): Tensor => {
	return new Tensor(_absolute(tensor._native_self));
}

export const sign = (tensor: Tensor): Tensor => {
	return new Tensor(_sign(tensor._native_self));
}

export const minimum = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_minimum(lhs._native_self, rhs._native_self));
}

export const argmin = (input: Tensor, axis: number, keepDims: boolean): Tensor => {
	return new Tensor(_argmin(input._native_self, axis, keepDims));
}

export const tanh = (tensor: Tensor): Tensor => {
	return new Tensor(_tanh(tensor._native_self));
}

export const cos = (tensor: Tensor): Tensor => {
	return new Tensor(_cos(tensor._native_self));
}

export const floor = (tensor: Tensor): Tensor => {
	return new Tensor(_floor(tensor._native_self));
}

export const sigmoid = (tensor: Tensor): Tensor => {
	return new Tensor(_sigmoid(tensor._native_self));
}

export const roll = (tensor: Tensor, shift: number, axis: number): Tensor => {
	return new Tensor(_roll(tensor._native_self, shift, axis));
}

export const _var = (input: Tensor, axes: number[], bias: boolean, keepDims: boolean): Tensor => {
	return new Tensor(__var(input._native_self, axes, bias, keepDims));
}

export const exp = (tensor: Tensor): Tensor => {
	return new Tensor(_exp(tensor._native_self));
}

export const transpose = (tensor: Tensor, axes: number[]): Tensor => {
	return new Tensor(_transpose(tensor._native_self, axes));
}

export const concatenate = (tensors: Tensor[], axis: number): Tensor => {
	return new Tensor(_concatenate(tensors, axis));
}

export const negative = (tensor: Tensor): Tensor => {
	return new Tensor(_negative(tensor._native_self));
}

export const isnan = (tensor: Tensor): Tensor => {
	return new Tensor(_isnan(tensor._native_self));
}

export const matmul = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_matmul(lhs._native_self, rhs._native_self));
}

export const amax = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_amax(input._native_self, axes, keepDims));
}

export const std = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_std(input._native_self, axes, keepDims));
}

export const isInvalidArray = (tensor: Tensor): boolean => {
	return _isInvalidArray(tensor._native_self);
}

export const iota = (dims: number[], tileDims: number[]): Tensor => {
	return new Tensor(_iota(dims, tileDims));
}

export const sort = (input: Tensor, axis: number): Tensor => {
	return new Tensor(_sort(input._native_self, axis));
}

export const maximum = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_maximum(lhs._native_self, rhs._native_self));
}

export const log1p = (tensor: Tensor): Tensor => {
	return new Tensor(_log1p(tensor._native_self));
}

export const log = (tensor: Tensor): Tensor => {
	return new Tensor(_log(tensor._native_self));
}

export const power = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_power(lhs._native_self, rhs._native_self));
}

export const amin = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_amin(input._native_self, axes, keepDims));
}

export const cumsum = (input: Tensor, axis: number): Tensor => {
	return new Tensor(_cumsum(input._native_self, axis));
}

export const all = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_all(input._native_self, axes, keepDims));
}

export const identity = (dim: number): Tensor => {
	return new Tensor(_identity(dim));
}

export const sin = (tensor: Tensor): Tensor => {
	return new Tensor(_sin(tensor._native_self));
}

export const flip = (tensor: Tensor, dim: number): Tensor => {
	return new Tensor(_flip(tensor._native_self, dim));
}

export const sum = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_sum(input._native_self, axes, keepDims));
}

export const mean = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_mean(input._native_self, axes, keepDims));
}

export const median = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_median(input._native_self, axes, keepDims));
}

export const nonzero = (tensor: Tensor): Tensor => {
	return new Tensor(_nonzero(tensor._native_self));
}

export const ceil = (tensor: Tensor): Tensor => {
	return new Tensor(_ceil(tensor._native_self));
}

export const rint = (tensor: Tensor): Tensor => {
	return new Tensor(_rint(tensor._native_self));
}

export const clip = (tensor: Tensor, low: Tensor, high: Tensor): Tensor => {
	return new Tensor(_clip(tensor._native_self, low._native_self, high._native_self));
}

export const any = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_any(input._native_self, axes, keepDims));
}

export const full = (dims: number[], val: number): Tensor => {
	return new Tensor(_full(dims, val));
}

export const tril = (tensor: Tensor): Tensor => {
	return new Tensor(_tril(tensor._native_self));
}

export const where = (condition: Tensor, x: Tensor, y: Tensor): Tensor => {
	return new Tensor(_where(condition._native_self, x._native_self, y._native_self));
}

export const norm = (input: Tensor, axes: number[], p: number, keepDims: boolean): Tensor => {
	return new Tensor(_norm(input._native_self, axes, p, keepDims));
}

export const reshape = (tensor: Tensor, shape: number[]): Tensor => {
	return new Tensor(_reshape(tensor._native_self, shape));
}

export const isinf = (tensor: Tensor): Tensor => {
	return new Tensor(_isinf(tensor._native_self));
}

export const triu = (tensor: Tensor): Tensor => {
	return new Tensor(_triu(tensor._native_self));
}

export const argmax = (input: Tensor, axis: number, keepDims: boolean): Tensor => {
	return new Tensor(_argmax(input._native_self, axis, keepDims));
}

export const countNonzero = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_countNonzero(input._native_self, axes, keepDims));
}

export const logicalNot = (tensor: Tensor): Tensor => {
	return new Tensor(_logicalNot(tensor._native_self));
}

export const tile = (tensor: Tensor, shape: number[]): Tensor => {
	return new Tensor(_tile(tensor._native_self, shape));
}

export const toFloat32Array = (tensor: Tensor): Float32Array => {
	return _toFloat32Array(tensor._native_self);
}

export const toFloat64Array = (tensor: Tensor): Float64Array => {
	return _toFloat64Array(tensor._native_self);
}

export const toBoolInt8Array = (tensor: Tensor): Int8Array => {
	return _toBoolInt8Array(tensor._native_self);
}

export const toInt16Array = (tensor: Tensor): Int16Array => {
	return _toInt16Array(tensor._native_self);
}

export const toInt32Array = (tensor: Tensor): Int32Array => {
	return _toInt32Array(tensor._native_self);
}

export const toInt64Array = (tensor: Tensor): BigInt64Array => {
	return _toInt64Array(tensor._native_self);
}

export const toUint8Array = (tensor: Tensor): Uint8Array => {
	return _toUint8Array(tensor._native_self);
}

export const toUint16Array = (tensor: Tensor): Uint16Array => {
	return _toUint16Array(tensor._native_self);
}

export const toUint32Array = (tensor: Tensor): Uint32Array => {
	return _toUint32Array(tensor._native_self);
}

export const toUint64Array = (tensor: Tensor): BigUint64Array => {
	return _toUint64Array(tensor._native_self);
}

export const toFloat32Scalar = (tensor: Tensor): number => {
	return _toFloat32Scalar(tensor._native_self);
}

export const toFloat64Scalar = (tensor: Tensor): number => {
	return _toFloat64Scalar(tensor._native_self);
}

export const toBoolInt8Scalar = (tensor: Tensor): number => {
	return _toBoolInt8Scalar(tensor._native_self);
}

export const toInt16Scalar = (tensor: Tensor): number => {
	return _toInt16Scalar(tensor._native_self);
}

export const toInt32Scalar = (tensor: Tensor): number => {
	return _toInt32Scalar(tensor._native_self);
}

export const toInt64Scalar = (tensor: Tensor): bigint => {
	return _toInt64Scalar(tensor._native_self);
}

export const toUint8Scalar = (tensor: Tensor): number => {
	return _toUint8Scalar(tensor._native_self);
}

export const toUint16Scalar = (tensor: Tensor): number => {
	return _toUint16Scalar(tensor._native_self);
}

export const toUint32Scalar = (tensor: Tensor): number => {
	return _toUint32Scalar(tensor._native_self);
}

export const toUint64Scalar = (tensor: Tensor): bigint => {
	return _toUint64Scalar(tensor._native_self);
}

export const _eval = (tensor: Tensor): void => {
	return __eval(tensor._native_self);
}

export const dispose = (tensor: Tensor): void => {
	return _dispose(tensor._native_self);
}

export const tensorFromFloat32Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromFloat32Buffer(buffer);
}

export const tensorFromFloat64Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromFloat64Buffer(buffer);
}

export const tensorFromBoolInt8Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromBoolInt8Buffer(buffer);
}

export const tensorFromInt16Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromInt16Buffer(buffer);
}

export const tensorFromInt32Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromInt32Buffer(buffer);
}

export const tensorFromInt64Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromInt64Buffer(buffer);
}

export const tensorFromUint8Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromUint8Buffer(buffer);
}

export const tensorFromUint16Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromUint16Buffer(buffer);
}

export const tensorFromUint32Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromUint32Buffer(buffer);
}

export const tensorFromUint64Buffer = (buffer: ArrayBuffer) => {
	return _tensorFromUint64Buffer(buffer);
}

export const save = (tensor: Tensor, path: string): void => {
	return _save(tensor._native_self, path);
}

export const eq = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_eq(lhs._native_self, rhs._native_self));
}

export const greaterThan = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_greaterThan(lhs._native_self, rhs._native_self));
}

export const mod = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_mod(lhs._native_self, rhs._native_self));
}

export const lessThanEqual = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lessThanEqual(lhs._native_self, rhs._native_self));
}

export const add = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_add(lhs._native_self, rhs._native_self));
}

export const greaterThanEqual = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_greaterThanEqual(lhs._native_self, rhs._native_self));
}

export const rShift = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_rShift(lhs._native_self, rhs._native_self));
}

export const logicalAnd = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_logicalAnd(lhs._native_self, rhs._native_self));
}

export const mul = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_mul(lhs._native_self, rhs._native_self));
}

export const lShift = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lShift(lhs._native_self, rhs._native_self));
}

export const bitwiseXor = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseXor(lhs._native_self, rhs._native_self));
}

export const bitwiseOr = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseOr(lhs._native_self, rhs._native_self));
}

export const div = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_div(lhs._native_self, rhs._native_self));
}

export const neq = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_neq(lhs._native_self, rhs._native_self));
}

export const lessThan = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lessThan(lhs._native_self, rhs._native_self));
}

export const bitwiseAnd = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseAnd(lhs._native_self, rhs._native_self));
}

export const logicalOr = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_logicalOr(lhs._native_self, rhs._native_self));
}

export const sub = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_sub(lhs._native_self, rhs._native_self));
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

