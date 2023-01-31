// Code generated by gen-napi. DO NOT EDIT.
const {
  _foo,
  _bar,
  _baz,
  _qux,
  _quux,
  _test
} = require("../../build/Release/test.node")

export const test = (a: string): string => {
  return _test(a);
}

export const foo = (a: number): number => {
  return _foo(a);
}

export const bar = (a: Float64Array, b: number): Float64Array => {
  return _bar(a, b);
}

export const baz = (a: Float32Array, b: number): Float32Array => {
  return _baz(a, b);
}

export const qux = (a: BigInt64Array, b: number): BigInt64Array => {
  return _qux(a, b);
}

export const quux = (a: boolean, b: boolean): Uint8Array => {
  return _quux(a, b);
}

