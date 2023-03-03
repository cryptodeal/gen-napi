// Code generated by gen-napi. DO NOT EDIT.
/* eslint-disable */

const {
  _qux,
  _quux,
  _test,
  _test2,
  _test3,
  _foo,
  _bar,
  _baz,
  _test4,
  _test5
} = import.meta.require("../../../../build/Release/test.node")

export const qux = (a: Array<number | bigint> | BigInt64Array, b: number): number => {
  return _qux(a instanceof BigInt64Array ? a : new BigInt64Array(a.map((v) => typeof v === 'number' ? BigInt(v) : v)), b);
}

export const quux = (a: boolean, b: boolean): Uint8Array => {
  return _quux(a, b);
}

export const test = (a: string): string => {
  return _test(a);
}

export const test2 = (a: number[] | Float64Array): Int8Array => {
  return _test2(a instanceof Float64Array ? a : new Float64Array(a));
}

export const test3 = (a: number[] | Float64Array): [number, number] => {
  return _test3(a instanceof Float64Array ? a : new Float64Array(a));
}

export const foo = (a: number) => {
  return _foo(a);
}

export const bar = (a: number[] | Float64Array, b: number): Float64Array => {
  return _bar(a instanceof Float64Array ? a : new Float64Array(a), b);
}

export const baz = (a: number[] | Float32Array, b: number): Float32Array => {
  return _baz(a instanceof Float32Array ? a : new Float32Array(a), b);
}

export const test4 = (a: string, b: string): string => {
  return _test4(a, b);
}

export const test5 = (a: number[] | Float64Array): BigInt64Array => {
  return _test5(a instanceof Float64Array ? a : new Float64Array(a));
}

