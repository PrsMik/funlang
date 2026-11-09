let foldRange: fn(U: comptime type) -> fn(int, int, int, U, fn(U, int) -> U) -> U = comptime fn(U) {
	let step: fn(int, int, int, U, fn(U, int) -> U) -> U = fn(from, to, stepBy, acc, body) {
		return if from >= to {
			return acc;
		} else {
			return step(from + stepBy, to, stepBy, body(acc, from), body);
		};
	};
	return fn(from: int, to: int, stepBy: int, init: U, body: fn(U, int) -> U) -> U {
		return step(from, to, stepBy, init, body);
	};
};

let reverse: fn(T: comptime type) -> fn([T]) -> [T] = comptime fn(T) {
	let step: fn([T], [T], int, int) -> [T] = fn(arr, acc, n, i) {
		return if i >= n {
			return acc;
		} else {
			return step(arr, push(acc, arr[i]), n, i + 1);
		};
	};
	return fn(arr: [T]) -> [T] {
		return step(arr, [], len(arr), 0);
	};
};

// reduce: [T] × U → U, порядок слева направо
let reduce: fn(T: comptime type, U: comptime type) -> fn([T], U, fn(U, T) -> U) -> U = comptime fn(T, U) {
	let f: fn(int, int, int, U, fn(U, int) -> U) -> U = comptime foldRange(U);
	return fn(arr: [T], init: U, comb: fn(U, T) -> U) -> U {
		return f(
			0, len(arr), 1, init, fn(acc, i) {
				return comb(acc, arr[i]);
			}
		);
	};
};

// map: [T] × (T → U) → [U]
let map: fn(T: comptime type, U: comptime type) -> fn([T], fn(T) -> U) -> [U] = comptime fn(T, U) {
	let foldType: type = comptime [U];
	let f: fn(int, int, int, [U], fn([U], int) -> [U]) -> [U] = comptime foldRange(foldType);
	return fn(arr: [T], g: fn(T) -> U) -> [U] {
		return f(
			0, len(arr), 1, [], fn(acc, i) {
				return push(acc, g(arr[i]));
			}
		);
	};
};

// gen(n, f): n элементов f(0..n-1)
let gen: fn(T: comptime type) -> fn(int, fn(int) -> T) -> [T] = comptime fn(T) {
	let foldType: type = comptime [T];
	let f: fn(int, int, int, [T], fn([T], int) -> [T]) -> [T] = comptime foldRange(foldType);
	return fn(n: int, g: fn(int) -> T) -> [T] {
		return f(
			0, n, 1, [], fn(acc, i) {
				return push(acc, g(i));
			}
		);
	};
};

// genRange(from, to, stepBy, f)
let genRange: fn(T: comptime type) -> fn(int, int, int, fn(int) -> T) -> [T] = comptime fn(T) {
	let foldType: type = comptime [T];
	let f: fn(int, int, int, [T], fn([T], int) -> [T]) -> [T] = comptime foldRange(foldType);
	return fn(from: int, to: int, stepBy: int, g: fn(int) -> T) -> [T] {
		return f(
			from, to, stepBy, [], fn(acc, i) {
				return push(acc, g(i));
			}
		);
	};
};

// genIf(from, to, pred, f)
let genIf: fn(T: comptime type) -> fn(int, int, fn(int) -> bool, fn(int) -> T) -> [T] = comptime fn(T) {
	let foldType: type = comptime [T];
	let fold: fn(int, int, int, [T], fn([T], int) -> [T]) -> [T] = comptime foldRange(foldType);
	return fn(from: int, to: int, predicate: fn(int) -> bool, genFunc: fn(int) -> T) -> [T] {
		return fold(
			from, to, 1, [], fn(acc, i) {
				return if predicate(i) == true {
					return push(acc, genFunc(i));
				} else {
					return acc;
				};
			}
		);
	};
};

let nums: [int] = [1, 2, 3, 4, 5, 6];

let reverseInt: fn(arr: [int]) -> [int] = comptime reverse(int);
let reverseBool: fn(arr: [bool]) -> [bool] = comptime reverse(bool);

let mapIntInt: fn(arr: [int], g: fn(int) -> int) -> [int] = comptime map(int, int);
let mapIntBool: fn(arr: [int], g: fn(int) -> bool) -> [bool] = comptime map(int, bool);
let reduceIntInt: fn(arr: [int], init: int, comb: fn(int, int) -> int) -> int = comptime reduce(int, int);
let genInt: fn(n: int, g: fn(int) -> int) -> [int] = comptime gen(int);
let genRangeInt: fn(from: int, to: int, stepBy: int, g: fn(int) -> int) -> [int] = comptime genRange(int);
let genIfInt: fn(from: int, to: int, pred: fn(int) -> bool, g: fn(int) -> int) -> [int] = comptime genIf(int);

let doubled: [int] = mapIntInt(
	nums, fn(x) {
		return x * 2;
	}
);

let flags: [bool] = mapIntBool(
	nums, fn(x) {
		return x / 2 * 2 == x;
	}
);

let sum: int = reduceIntInt(
	nums, 0, fn(a, x) {
		return a + x;
	}
);

let sq: [int] = genInt(
	10, fn(i) {
		return i * i;
	}
);

let evens: [int] = genRangeInt(
	0, 10, 2, fn(i) {
		return i;
	}
);

let sqEven: [int] = genIfInt(
	0, 20, fn(i) {
		return i / 2 * 2 == i;
	}, fn(i) {
		return i * i;
	}
);

let o1: int = puts(reverseInt(doubled));
let o2: int = puts(reverseBool(flags));
let o3: int = puts(sum);
let o5: int = puts(reverseInt(sq));
