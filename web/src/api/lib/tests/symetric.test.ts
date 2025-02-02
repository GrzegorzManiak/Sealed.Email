import * as Sym from '../symetric';
import { expect, test } from "bun:test";
import {BytesToBigInt} from "gowl-client-lib";

test("New Key", () => {
	const key = Sym.NewKey();
	expect(key).toBeTruthy();
});

test("Key Randomness", () => {
	const key1 = Sym.NewKey();
	const key2 = Sym.NewKey();
	expect(BytesToBigInt(key1)).not.toBe(BytesToBigInt(key2));
});

test("Key Length", () => {
	const key = Sym.NewKey();
	expect(key.length).toBe(32);

	const key2 = Sym.NewKey(16);
	expect(key2.length).toBe(16);
});


test("Encrypt", async () => {
	const key = Sym.NewKey();
	const data = "Hello, world!";
	const encrypted = await Sym.Encrypt(data, key);
	expect(encrypted).toBeTruthy();
});

test("Decrypt", async () => {
	const key = Sym.NewKey();
	const data = "Hello, world!";
	const encrypted = await Sym.Encrypt(data, key);
	const decrypted = await Sym.Decrypt(encrypted, key);
	expect(decrypted).toBe(data);
});

test("Compress", async () => {
	const key = Sym.NewKey();
	const data = "Hello, world!";
	const encrypted = await Sym.Encrypt(data, key);
	const compressed = Sym.Compress(encrypted);
	expect(compressed).toBeTruthy();
});

test("Decompress", async () => {
	const key = Sym.NewKey();
	const data = "Hello, world!";
	const encrypted = await Sym.Encrypt(data, key);
	const compressed = Sym.Compress(encrypted);
	const decompressed = Sym.Decompress(compressed);
	expect(decompressed).toBeTruthy();
});
