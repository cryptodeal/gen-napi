// Code generated by gen-napi. DO NOT EDIT.
const {
,
  _baz,
  _qux,
  _quux,
  _foo,
  _bar
} = require("../../build/Release/test.node")

export const foo = (a: number): number => {
  return _foo(a);
}

export const bar = (a: number, b: number): number => {
  return _bar(a, b);
}

export const baz = (a: number, b: number): number => {
  return _baz(a, b);
}

export const qux = (a: number, b: number): number => {
  return _qux(a, b);
}

export const quux = (a: boolean, b: boolean): boolean => {
  return _quux(a, b);
}

