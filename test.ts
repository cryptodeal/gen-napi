import { Tensor } from './test/js/index.ts'

const test = new Tensor(new Float64Array([1, 2, 3, 4, 5, 6]))
const test2 = new Tensor(new Float64Array([1, 2, 3, 4, 5, 6])).add(test)

console.log(test.toFloat64Array())
console.log(test2.toFloat64Array())
