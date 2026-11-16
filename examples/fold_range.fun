let foldRange := comptime fn(U: comptime type) -> fn(int, int, int, U, fn(U, int) -> U) -> U {
	let step := fn(from: int, to: int, stepBy: int, acc: U, body: fn(U, int) -> U) -> U {
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

let reverse := comptime fn(T: comptime type) -> fn([]T) -> []T {
	let step := fn(arr: []T, acc: []T, n: int, i: int) -> []T {
		return if i >= n {
			return acc;
		} else {
			return step(arr, push(acc, arr[i]), n, i + 1);
		};
	};
	return fn(arr: []T) -> []T {
		return step(arr, []T{}, len(arr), 0);
	};
};

// reduce: [T] × U → U, порядок слева направо
let reduce := comptime fn(T: comptime type, U: comptime type) -> fn([]T, U, fn(U, T) -> U) -> U {
	let foldFunc := comptime foldRange(U);
	return fn(arr: []T, init: U, comb: fn(U, T) -> U) -> U {
		return foldFunc(
			0, len(arr), 1, init, fn(acc, i) {
				return comb(acc, arr[i]);
			}
		);
	};
};

// map: [T] × (T → U) → [U]
let map := comptime fn(T: comptime type, U: comptime type) -> fn([]T, fn(T) -> U) -> []U {
	let foldFunc := comptime foldRange([]U);
	return fn(arr: []T, g: fn(T) -> U) -> []U {
		return foldFunc(
			0, len(arr), 1, []U{}, fn(acc, i) {
				return push(acc, g(arr[i]));
			}
		);
	};
};

// gen(n, f): n элементов f(0..n-1)
let gen := comptime fn(T: comptime type) -> fn(int, fn(int) -> T) -> []T {
	let foldFunc := comptime foldRange([]T);
	return fn(n: int, g: fn(int) -> T) -> []T {
		return foldFunc(
			0, n, 1, []T{}, fn(acc, i) {
				return push(acc, g(i));
			}
		);
	};
};

// genRange(from, to, stepBy, f)
let genRange := comptime fn(T: comptime type) -> fn(int, int, int, fn(int) -> T) -> []T {
	let foldFunc := comptime foldRange([]T);
	return fn(from: int, to: int, stepBy: int, g: fn(int) -> T) -> []T {
		return foldFunc(
			from, to, stepBy, []T{}, fn(acc, i) {
				return push(acc, g(i));
			}
		);
	};
};

// genIf(from, to, pred, f)
let genIf := comptime fn(T: comptime type) -> fn(int, int, fn(int) -> bool, fn(int) -> T) -> []T {
	let fold := comptime foldRange([]T);
	return fn(from: int, to: int, predicate: fn(int) -> bool, genFunc: fn(int) -> T) -> []T {
		return fold(
			from, to, 1, []T{}, fn(acc, i) {
				return if predicate(i) == true {
					return push(acc, genFunc(i));
				} else {
					return acc;
				};
			}
		);
	};
};

let nums := []int{1, 2, 3, 4, 5, 6};

let reverseInt := comptime reverse(int);
let reverseBool := comptime reverse(bool);

let mapIntInt := comptime map(int, int);
let mapIntBool := comptime map(int, bool);
let reduceIntInt := comptime reduce(int, int);
let genInt := comptime gen(int);
let genRangeInt := comptime genRange(int);
let genIfInt := comptime genIf(int);

let doubled := mapIntInt(
	nums, fn(x) {
		return x * 2;
	}
);

let flags := mapIntBool(
	nums, fn(x) {
		return x / 2 * 2 == x;
	}
);

let sum := reduceIntInt(
	nums, 0, fn(a, x) {
		return a + x;
	}
);

let sq := genInt(
	10, fn(i) {
		return i * i;
	}
);

let evens := genRangeInt(
	0, 10, 2, fn(i) {
		return i;
	}
);

let sqEven := genIfInt(
	0, 20, fn(i) {
		return i / 2 * 2 == i;
	}, fn(i) {
		return i * i;
	}
);

let o1 := puts(reverseInt(doubled));
let o2 := puts(reverseBool(flags));
let o3 := puts(sum);
let o5 := puts(reverseInt(sq));
