export class StringBuffer {
	#buffer: Buffer;
	constructor(v: string) {
		if (!(v.slice(-2) == '\0')) v += '\0';
		this.#buffer = Buffer.from(v, 'utf8');
	}

	get buffer() {
		return this.#buffer;
	}

	toString(): string {
		return this.#buffer.slice(0, -2).toString();
	}

	toUint8Array(): Uint8Array {
		return new Uint8Array(this.#buffer);
	}
}
