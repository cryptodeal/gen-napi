const { Tensor } = require('./test/js/index.cjs')
const test = new Tensor(new Float64Array([1, 2, 3, 4, 5, 6]))
console.log(test.toFloat64Array())
