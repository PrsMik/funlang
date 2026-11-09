let reduce: fn(T: comptime type, U: comptime type) -> fn([T], U, fn(U, T) -> U) -> U = comptime fn(T, U) {
	let step: fn([T], U, fn(U, T) -> U, int) -> U = fn(arr, acc, f, i) {
		return if i >= len(arr) {
			return acc;
		} else {
			return step(arr, f(acc, arr[i]), f, i + 1);
		};
	};
	return fn(arr: [T], init: U, f: fn(U, T) -> U) -> U {
		return step(arr, init, f, 0);
	};
};

let nums: [int] = [1, 2, 3, 4, 5];

// сумма
let reduceIntInt: fn(arr: [int], init: int, f: fn(int, int) -> int) -> int = comptime reduce(int, int);
let sum: int = reduceIntInt(nums, 0, fn(acc, x) { return acc + x; });
let o1: int = puts(sum);   // 15

// произведение
let prod: int = reduceIntInt(nums, 1, fn(acc, x) { return acc * x; });
let o2: int = puts(prod);  // 120