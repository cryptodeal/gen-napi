/* eslint-disable @typescript-eslint/no-var-requires */
const { Tensor } = require('./demo1/js/tensor.cjs');
const sm = require('./demo1/js/index.cjs');

sm.init();

function genRand() {
  const out = new Float32Array(128);
  for (let i = 0; i < 128; ++i) {
    out[i] = Math.random();
  }
  return out;
}

const t0 = performance.now() / 1e3;
let m = 0;
for (let i = 0; i < 1000000; ++i) {
  // console.log('bytes: ', Number(sm.bytesUsed()))
  const a = sm.rand([128]);
  const b = new Tensor(genRand());
  m += a.add(b).mean([], false).toFloat32Scalar();
  // console.log('bytes: ', Number(sm.bytesUsed()))
}
const t1 = performance.now() / 1e3;
console.log(t1 - t0, 'seconds to calculate', m);
m = null;
console.log('bytes: ', Number(sm.bytesUsed()));
