import { describe, expect, it } from 'bun:test';
import { dtype as jsDtype } from '../demo1/js/index.cjs';
import { dtype } from '../demo1/ts/index';

describe('C++ Enums -> Generated Enum', () => {
	const dtype_keys = ['f16', 'f32', 'f64', 'b8', 's16', 's32', 's64', 'u8', 'u16', 'u32', 'u64'];
	it('pre-compiled JS Enum works', () => {
		for (let i = 0; i < 11; i++) {
			expect(typeof jsDtype[dtype_keys[i]]).toBe('number');
			expect(typeof jsDtype[i]).toBe('string');
		}
	});

	it('TS Enum works', () => {
		for (let i = 0; i < 11; i++) {
			expect(typeof dtype[dtype_keys[i]]).toBe('number');
			expect(typeof dtype[i]).toBe('string');
		}
	});
});
