const {
	_tensorFromFloat32Buffer,
	_toBoolInt8Array,
	_toUint8Array,
	_toInt16Array,
	_toUint16Array,
	_toInt32Array,
	_toUint32Array,
	_toUint64Array,
	_toInt64Array,
	_toBoolInt8Scalar,
	_toUint8Scalar,
	_toInt16Scalar,
	_toUint16Scalar,
	_toUint32Scalar,
	_toInt32Scalar,
	_toUint64Scalar,
	_toInt64Scalar,
	_toFloat32Array,
	_toFloat64Array,
	_toFloat32Scalar,
	_toFloat64Scalar
} = import.meta.require('../../build/Release/shumai_bindings.node');
import { _Base_Tensor, gen_Tensor_ops_shim } from './gen_Tensor_methods_shim';

export class Tensor extends _Base_Tensor {
	#underlying;

	constructor(t) {
		super();
		if (t instanceof Float32Array || t.constructor === Float32Array) {
			this.#underlying = t;
			this._native_Tensor = _tensorFromFloat32Buffer(t.buffer);
		} else {
			this._native_Tensor = t;
		}
	}

	toInt8Scalar() {
		return _toBoolInt8Scalar(this._native_ref);
	}

	toInt8Array() {
		return _toBoolInt8Array(this._native_ref);
	}

	toUint8Scalar() {
		return _toUint8Scalar(this._native_ref);
	}

	toUint8Array() {
		return _toUint8Array(this._native_ref);
	}

	toInt16Scalar() {
		return _toInt16Scalar(this._native_ref);
	}

	toInt16Array() {
		return _toInt16Array(this._native_ref);
	}

	toUint16Scalar() {
		return _toUint16Scalar(this._native_ref);
	}

	toUint16Array() {
		return _toUint16Array(this._native_ref);
	}

	toInt32Scalar() {
		return _toInt32Scalar(this._native_ref);
	}

	toInt32Array() {
		return _toInt32Array(this._native_ref);
	}

	toUint32Scalar() {
		return _toUint32Scalar(this._native_ref);
	}

	toUint32Array() {
		return _toUint32Array(this._native_ref);
	}

	toBigInt64Scalar() {
		return _toInt64Scalar(this._native_ref);
	}

	toBigInt64Array() {
		return _toInt64Array(this._native_ref);
	}

	toBigUint64Scalar() {
		return _toUint64Scalar(this._native_ref);
	}

	toBigUint64Array() {
		return _toUint64Array(this._native_ref);
	}

	toFloat32Scalar() {
		return _toFloat32Scalar(this._native_ref);
	}

	toFloat32Array() {
		return _toFloat32Array(this._native_ref);
	}

	toFloat64Scalar() {
		return _toFloat64Scalar(this._native_ref);
	}

	toFloat64Array() {
		return _toFloat64Array(this._native_ref);
	}
}

export interface Tensor extends ReturnType<typeof gen_Tensor_ops_shim> {} // eslint-disable-line

for (const [method, closure] of Object.entries(gen_Tensor_ops_shim(Tensor))) {
	Tensor.prototype[method] = closure;
}
