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

test('should compress and decompress data correctly', () => {
	const encryptedData = {
		iv: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12],
		data: [13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26]
	};

	const compressedData = Sym.Compress(encryptedData);
	const decompressedData = Sym.Decompress(compressedData);
	expect(decompressedData.iv).toEqual(encryptedData.iv);
	expect(decompressedData.data).toEqual(encryptedData.data);
});

test('should handle empty data arrays', () => {
	const encryptedData = {
		iv: [],
		data: []
	};

	const compressedData = Sym.Compress(encryptedData);
	const decompressedData = Sym.Decompress(compressedData);
	expect(decompressedData.iv).toEqual(encryptedData.iv);
	expect(decompressedData.data).toEqual(encryptedData.data);
});
