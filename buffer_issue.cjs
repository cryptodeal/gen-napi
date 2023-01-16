const { Tensor } = require('./test/js/tensor.cjs')
const sm = require('./test/js/index.cjs')

const fillArray = (arr) => {
  const len = arr.length
  for (let i = 0; i < len; i++) {
    arr[i] = Math.random()
  }
  return arr
}

const t0 = performance.now() / 1e3

for (let i = 0; i < 100000; i++) {
  const backingArray = fillArray(new Float32Array(10000))
  const a = new Tensor(backingArray)
  const out = a.toFloat32Array()
}

const t1 = performance.now() / 1e3

const time = t1 - t0
console.log(time, 'seconds to init 100,000 tensors of 10,000 elements (dtype float64)')
