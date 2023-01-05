// Code generated by gen-napi. DO NOT EDIT.
const {
	_Tensor,
	_log,
	_sin,
	_clip,
	_isinf,
	_sign,
	_triu,
	_ceil,
	_rint,
	_erf,
	_norm,
	_countNonzero,
	_argmin,
	_median,
	_full,
	_negative,
	_logicalNot,
	_sigmoid,
	_matmul,
	_amax,
	_concatenate,
	_cos,
	_roll,
	_cumsum,
	_var: __var,
	_all,
	_transpose,
	_tanh,
	_floor,
	_std,
	_any,
	_tril,
	_where,
	_arange,
	_iota,
	_tile,
	_log1p,
	_flip,
	_isnan,
	_sort,
	_maximum,
	_power,
	_mean,
	_sum,
	_identity,
	_reshape,
	_nonzero,
	_sqrt,
	_absolute,
	_minimum,
	_exp,
	_amin,
	_argmax,
	_mul,
	_lessThan,
	_bitwiseXor,
	_eq,
	_greaterThan,
	_bitwiseAnd,
	_bitwiseOr,
	_lessThanEqual,
	_sub,
	_greaterThanEqual,
	_add,
	_mod,
	_neq,
	_rShift,
	_lShift,
	_logicalOr,
	_div,
	_logicalAnd,
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
} = require("../../build/Release/flashlight_napi_bindings.node")

export class Tensor {
	#_native_self: any;

	constructor(t) {
		this.#_native_self = new _Tensor(t);
	}

	get _native_self() {
		return this.#_native_self;
	}

	argmin(axis: number, keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.argmin(axis, keepDims));
	}

	median(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.median(axes, keepDims));
	}

	negative(): Tensor {
		return new Tensor(this.#_native_self.negative());
	}

	logicalNot(): Tensor {
		return new Tensor(this.#_native_self.logicalNot());
	}

	sigmoid(): Tensor {
		return new Tensor(this.#_native_self.sigmoid());
	}

	matmul(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.matmul(rhs.#_native_self));
	}

	amax(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.amax(axes, keepDims));
	}

	cos(): Tensor {
		return new Tensor(this.#_native_self.cos());
	}

	roll(shift: number, axis: number): Tensor {
		return new Tensor(this.#_native_self.roll(shift, axis));
	}

	cumsum(axis: number): Tensor {
		return new Tensor(this.#_native_self.cumsum(axis));
	}

	var(axes: number[], bias: boolean, keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.var(axes, bias, keepDims));
	}

	all(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.all(axes, keepDims));
	}

	transpose(axes: number[]): Tensor {
		return new Tensor(this.#_native_self.transpose(axes));
	}

	tanh(): Tensor {
		return new Tensor(this.#_native_self.tanh());
	}

	floor(): Tensor {
		return new Tensor(this.#_native_self.floor());
	}

	std(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.std(axes, keepDims));
	}

	any(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.any(axes, keepDims));
	}

	tril(): Tensor {
		return new Tensor(this.#_native_self.tril());
	}

	where(x: Tensor, y: Tensor): Tensor {
		return new Tensor(this.#_native_self.where(x.#_native_self, y.#_native_self));
	}

	tile(shape: number[]): Tensor {
		return new Tensor(this.#_native_self.tile(shape));
	}

	log1p(): Tensor {
		return new Tensor(this.#_native_self.log1p());
	}

	flip(dim: number): Tensor {
		return new Tensor(this.#_native_self.flip(dim));
	}

	isnan(): Tensor {
		return new Tensor(this.#_native_self.isnan());
	}

	sort(axis: number): Tensor {
		return new Tensor(this.#_native_self.sort(axis));
	}

	maximum(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.maximum(rhs.#_native_self));
	}

	power(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.power(rhs.#_native_self));
	}

	mean(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.mean(axes, keepDims));
	}

	sum(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.sum(axes, keepDims));
	}

	reshape(shape: number[]): Tensor {
		return new Tensor(this.#_native_self.reshape(shape));
	}

	nonzero(): Tensor {
		return new Tensor(this.#_native_self.nonzero());
	}

	sqrt(): Tensor {
		return new Tensor(this.#_native_self.sqrt());
	}

	absolute(): Tensor {
		return new Tensor(this.#_native_self.absolute());
	}

	exp(): Tensor {
		return new Tensor(this.#_native_self.exp());
	}

	amin(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.amin(axes, keepDims));
	}

	log(): Tensor {
		return new Tensor(this.#_native_self.log());
	}

	sin(): Tensor {
		return new Tensor(this.#_native_self.sin());
	}

	clip(low: Tensor, high: Tensor): Tensor {
		return new Tensor(this.#_native_self.clip(low.#_native_self, high.#_native_self));
	}

	isinf(): Tensor {
		return new Tensor(this.#_native_self.isinf());
	}

	sign(): Tensor {
		return new Tensor(this.#_native_self.sign());
	}

	triu(): Tensor {
		return new Tensor(this.#_native_self.triu());
	}

	ceil(): Tensor {
		return new Tensor(this.#_native_self.ceil());
	}

	rint(): Tensor {
		return new Tensor(this.#_native_self.rint());
	}

	erf(): Tensor {
		return new Tensor(this.#_native_self.erf());
	}

	norm(axes: number[], p: number, keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.norm(axes, p, keepDims));
	}

	countNonzero(axes: number[], keepDims: boolean): Tensor {
		return new Tensor(this.#_native_self.countNonzero(axes, keepDims));
	}

	logicalAnd(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.logicalAnd(rhs.#_native_self));
	}

	neq(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.neq(rhs.#_native_self));
	}

	rShift(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.rShift(rhs.#_native_self));
	}

	lShift(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.lShift(rhs.#_native_self));
	}

	logicalOr(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.logicalOr(rhs.#_native_self));
	}

	div(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.div(rhs.#_native_self));
	}

	lessThan(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.lessThan(rhs.#_native_self));
	}

	mul(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.mul(rhs.#_native_self));
	}

	bitwiseOr(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.bitwiseOr(rhs.#_native_self));
	}

	bitwiseXor(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.bitwiseXor(rhs.#_native_self));
	}

	eq(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.eq(rhs.#_native_self));
	}

	greaterThan(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.greaterThan(rhs.#_native_self));
	}

	bitwiseAnd(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.bitwiseAnd(rhs.#_native_self));
	}

	mod(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.mod(rhs.#_native_self));
	}

	lessThanEqual(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.lessThanEqual(rhs.#_native_self));
	}

	sub(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.sub(rhs.#_native_self));
	}

	greaterThanEqual(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.greaterThanEqual(rhs.#_native_self));
	}

	add(rhs: Tensor): Tensor {
		return new Tensor(this.#_native_self.add(rhs.#_native_self));
	}

	copy(): Tensor {
		return new Tensor(this.#_native_self.copy());
	}

	shape(): number[] {
		return this.#_native_self.shape();
	}

	elements(): number {
		return this.#_native_self.elements();
	}

	ndim(): number {
		return this.#_native_self.ndim();
	}

	isEmpty(): boolean {
		return this.#_native_self.isEmpty();
	}

	bytes(): number {
		return this.#_native_self.bytes();
	}

	type(): number {
		return this.#_native_self.type();
	}

	isSparse(): boolean {
		return this.#_native_self.isSparse();
	}

	strides(): number[] {
		return this.#_native_self.strides();
	}

	astype(): Tensor {
		return new Tensor(this.#_native_self.astype());
	}

	flatten(): Tensor {
		return new Tensor(this.#_native_self.flatten());
	}

	asContiguousTensor(): Tensor {
		return new Tensor(this.#_native_self.asContiguousTensor());
	}

	isContiguous(): boolean {
		return this.#_native_self.isContiguous();
	}

	toFloat32Array(): Float32Array {
		return this.#_native_self.toFloat32Array();
	}

	toFloat64Array(): Float64Array {
		return this.#_native_self.toFloat64Array();
	}

	toBoolInt8Array(): Int8Array {
		return this.#_native_self.toBoolInt8Array();
	}

	toInt16Array(): Int16Array {
		return this.#_native_self.toInt16Array();
	}

	toInt32Array(): Int32Array {
		return this.#_native_self.toInt32Array();
	}

	save(filename: string) {
		return this.#_native_self.save(filename);
	}

	toInt64Array(): BigInt64Array {
		return this.#_native_self.toInt64Array();
	}

	toUint8Array(): Uint8Array {
		return this.#_native_self.toUint8Array();
	}

	toUint16Array(): Uint16Array {
		return this.#_native_self.toUint16Array();
	}

	toUint32Array(): Uint32Array {
		return this.#_native_self.toUint32Array();
	}

	toUint64Array(): BigUint64Array {
		return this.#_native_self.toUint64Array();
	}

	toFloat32Scalar(): number {
		return this.#_native_self.toFloat32Scalar();
	}

	toFloat64Scalar(): number {
		return this.#_native_self.toFloat64Scalar();
	}

	toBoolInt8Scalar(): number {
		return this.#_native_self.toBoolInt8Scalar();
	}

	toInt16Scalar(): number {
		return this.#_native_self.toInt16Scalar();
	}

	toInt32Scalar(): number {
		return this.#_native_self.toInt32Scalar();
	}

	toInt64Scalar(): bigint {
		return this.#_native_self.toInt64Scalar();
	}

	toUint8Scalar(): number {
		return this.#_native_self.toUint8Scalar();
	}

	toUint16Scalar(): number {
		return this.#_native_self.toUint16Scalar();
	}

	toUint32Scalar(): number {
		return this.#_native_self.toUint32Scalar();
	}

	toUint64Scalar(): bigint {
		return this.#_native_self.toUint64Scalar();
	}

	eval() {
		return this.#_native_self.eval();
	}

	dispose() {
		return this.#_native_self.dispose();
	}

}

export const transpose = (tensor: Tensor, axes: number[]): Tensor => {
	return new Tensor(_transpose(tensor._native_self, axes));
}

export const tanh = (tensor: Tensor): Tensor => {
	return new Tensor(_tanh(tensor._native_self));
}

export const floor = (tensor: Tensor): Tensor => {
	return new Tensor(_floor(tensor._native_self));
}

export const std = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_std(input._native_self, axes, keepDims));
}

export const any = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_any(input._native_self, axes, keepDims));
}

export const where = (condition: Tensor, x: Tensor, y: Tensor): Tensor => {
	return new Tensor(_where(condition._native_self, x._native_self, y._native_self));
}

export const arange = (start: number, end: number, step: number): Tensor => {
	return new Tensor(_arange(start, end, step));
}

export const iota = (dims: number[], tileDims: number[]): Tensor => {
	return new Tensor(_iota(dims, tileDims));
}

export const tile = (tensor: Tensor, shape: number[]): Tensor => {
	return new Tensor(_tile(tensor._native_self, shape));
}

export const log1p = (tensor: Tensor): Tensor => {
	return new Tensor(_log1p(tensor._native_self));
}

export const flip = (tensor: Tensor, dim: number): Tensor => {
	return new Tensor(_flip(tensor._native_self, dim));
}

export const isnan = (tensor: Tensor): Tensor => {
	return new Tensor(_isnan(tensor._native_self));
}

export const tril = (tensor: Tensor): Tensor => {
	return new Tensor(_tril(tensor._native_self));
}

export const sort = (input: Tensor, axis: number): Tensor => {
	return new Tensor(_sort(input._native_self, axis));
}

export const maximum = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_maximum(lhs._native_self, rhs._native_self));
}

export const power = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_power(lhs._native_self, rhs._native_self));
}

export const mean = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_mean(input._native_self, axes, keepDims));
}

export const identity = (dim: number): Tensor => {
	return new Tensor(_identity(dim));
}

export const reshape = (tensor: Tensor, shape: number[]): Tensor => {
	return new Tensor(_reshape(tensor._native_self, shape));
}

export const nonzero = (tensor: Tensor): Tensor => {
	return new Tensor(_nonzero(tensor._native_self));
}

export const sqrt = (tensor: Tensor): Tensor => {
	return new Tensor(_sqrt(tensor._native_self));
}

export const absolute = (tensor: Tensor): Tensor => {
	return new Tensor(_absolute(tensor._native_self));
}

export const minimum = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_minimum(lhs._native_self, rhs._native_self));
}

export const sum = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_sum(input._native_self, axes, keepDims));
}

export const exp = (tensor: Tensor): Tensor => {
	return new Tensor(_exp(tensor._native_self));
}

export const amin = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_amin(input._native_self, axes, keepDims));
}

export const argmax = (input: Tensor, axis: number, keepDims: boolean): Tensor => {
	return new Tensor(_argmax(input._native_self, axis, keepDims));
}

export const log = (tensor: Tensor): Tensor => {
	return new Tensor(_log(tensor._native_self));
}

export const sin = (tensor: Tensor): Tensor => {
	return new Tensor(_sin(tensor._native_self));
}

export const clip = (tensor: Tensor, low: Tensor, high: Tensor): Tensor => {
	return new Tensor(_clip(tensor._native_self, low._native_self, high._native_self));
}

export const isinf = (tensor: Tensor): Tensor => {
	return new Tensor(_isinf(tensor._native_self));
}

export const sign = (tensor: Tensor): Tensor => {
	return new Tensor(_sign(tensor._native_self));
}

export const triu = (tensor: Tensor): Tensor => {
	return new Tensor(_triu(tensor._native_self));
}

export const ceil = (tensor: Tensor): Tensor => {
	return new Tensor(_ceil(tensor._native_self));
}

export const rint = (tensor: Tensor): Tensor => {
	return new Tensor(_rint(tensor._native_self));
}

export const erf = (tensor: Tensor): Tensor => {
	return new Tensor(_erf(tensor._native_self));
}

export const norm = (input: Tensor, axes: number[], p: number, keepDims: boolean): Tensor => {
	return new Tensor(_norm(input._native_self, axes, p, keepDims));
}

export const countNonzero = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_countNonzero(input._native_self, axes, keepDims));
}

export const median = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_median(input._native_self, axes, keepDims));
}

export const full = (dims: number[], val: number): Tensor => {
	return new Tensor(_full(dims, val));
}

export const negative = (tensor: Tensor): Tensor => {
	return new Tensor(_negative(tensor._native_self));
}

export const logicalNot = (tensor: Tensor): Tensor => {
	return new Tensor(_logicalNot(tensor._native_self));
}

export const sigmoid = (tensor: Tensor): Tensor => {
	return new Tensor(_sigmoid(tensor._native_self));
}

export const matmul = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_matmul(lhs._native_self, rhs._native_self));
}

export const amax = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_amax(input._native_self, axes, keepDims));
}

export const argmin = (input: Tensor, axis: number, keepDims: boolean): Tensor => {
	return new Tensor(_argmin(input._native_self, axis, keepDims));
}

export const concatenate = (tensors: Tensor[], axis: number): Tensor => {
	return new Tensor(_concatenate(tensors, axis));
}

export const cos = (tensor: Tensor): Tensor => {
	return new Tensor(_cos(tensor._native_self));
}

export const roll = (tensor: Tensor, shift: number, axis: number): Tensor => {
	return new Tensor(_roll(tensor._native_self, shift, axis));
}

export const cumsum = (input: Tensor, axis: number): Tensor => {
	return new Tensor(_cumsum(input._native_self, axis));
}

export const _var = (input: Tensor, axes: number[], bias: boolean, keepDims: boolean): Tensor => {
	return new Tensor(__var(input._native_self, axes, bias, keepDims));
}

export const all = (input: Tensor, axes: number[], keepDims: boolean): Tensor => {
	return new Tensor(_all(input._native_self, axes, keepDims));
}

export const lessThanEqual = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lessThanEqual(lhs._native_self, rhs._native_self));
}

export const sub = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_sub(lhs._native_self, rhs._native_self));
}

export const greaterThanEqual = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_greaterThanEqual(lhs._native_self, rhs._native_self));
}

export const add = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_add(lhs._native_self, rhs._native_self));
}

export const mod = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_mod(lhs._native_self, rhs._native_self));
}

export const neq = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_neq(lhs._native_self, rhs._native_self));
}

export const rShift = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_rShift(lhs._native_self, rhs._native_self));
}

export const lShift = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lShift(lhs._native_self, rhs._native_self));
}

export const logicalOr = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_logicalOr(lhs._native_self, rhs._native_self));
}

export const div = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_div(lhs._native_self, rhs._native_self));
}

export const logicalAnd = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_logicalAnd(lhs._native_self, rhs._native_self));
}

export const mul = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_mul(lhs._native_self, rhs._native_self));
}

export const lessThan = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_lessThan(lhs._native_self, rhs._native_self));
}

export const bitwiseXor = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseXor(lhs._native_self, rhs._native_self));
}

export const eq = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_eq(lhs._native_self, rhs._native_self));
}

export const greaterThan = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_greaterThan(lhs._native_self, rhs._native_self));
}

export const bitwiseAnd = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseAnd(lhs._native_self, rhs._native_self));
}

export const bitwiseOr = (lhs: Tensor, rhs: Tensor): Tensor => {
	return new Tensor(_bitwiseOr(lhs._native_self, rhs._native_self));
}

export const init = () => {
	return _init();
}

export const bytesUsed = (): bigint => {
	return _bytesUsed();
}

export const setRowMajor = () => {
	return _setRowMajor();
}

export const setColMajor = () => {
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

