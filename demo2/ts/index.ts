// Code generated by gen-napi. DO NOT EDIT.
/* eslint-disable */
const {
  _test,
  _baz
} = import.meta.require("../../../../build/Release/test.node")

export enum dtype {
  f16 = 0,
  f32 = 1,
  f64 = 2,
  b8 = 3,
  s16 = 4,
  s32 = 5,
  s64 = 6,
  u8 = 7,
  u16 = 8,
  u32 = 9,
  u64 = 10
}

export const test = (a: string): string => {
  return _test(a);
}

export const baz = (a: Float32Array, b: number): Float32Array => {
  return _baz(a, b);
}

