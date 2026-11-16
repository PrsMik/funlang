let reduce := comptime fn(T: comptime type, U: comptime type) -> fn([]T, U, fn(U, T) -> U) -> U {
	let step := fn(arr: []T, acc: U, f: fn(U, T) -> U, i: int) -> U {
		return if i >= len(arr) {
			return acc;
		} else {
			return step(arr, f(acc, arr[i]), f, i + 1);
		};
	};
	return fn(arr: []T, init: U, f: fn(U, T) -> U) -> U {
		return step(arr, init, f, 0);
	};
};

let nums := []int{1, 2, 3, 4, 5};

// сумма
let reduceIntInt := comptime reduce(int, int);
let sum := reduceIntInt(
	nums, 0, fn(acc, x) {
		return acc + x;
	}
);
let o1 := puts(sum); // 15

// произведение
let prod := reduceIntInt(
	nums, 1, fn(acc, x) {
		return acc * x;
	}
);
let o2 := puts(prod); // 120
