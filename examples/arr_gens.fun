// простой генератор из функции от i
let gen := comptime fn(T: comptime type) -> fn(int, fn(int) -> T) -> []T {
	let step: fn(int, fn(int) -> T, []T, int) -> []T = fn(n, genFn, acc, i) {
		return if (i >= n) {
			return acc;
		} else {
			return step(n, genFn, push(acc, genFn(i)), i + 1);
		};
	};
	return fn(n: int, genFn: fn(int) -> T) -> []T {
		return step(n, genFn, []T{}, 0);
	};
};

// пример: массив квадратов от 0 до 9
let genInt := comptime gen(int);
let xs := genInt(
	5, fn(i) {
		return i * 10;
	}
);
let _ := puts(xs); // [0, 10, 20, 30, 40]

// генератор с диапазоном (from / to / step)
let genRange: fn(T: comptime type) -> fn(int, int, int, fn(int) -> T) -> []T = comptime fn(T) {
	let step: fn(int, int, int, fn(int) -> T, []T, int) -> []T = fn(from, to, stepBy, genFn, acc, i) {
		return if (i >= to) {
			return acc;
		} else {
			return step(from, to, stepBy, genFn, push(acc, genFn(i)), i + stepBy);
		};
	};
	return fn(from: int, to: int, stepBy: int, genFn: fn(int) -> T) -> []T {
		return step(from, to, stepBy, genFn, []T{}, from);
	};
};

let genRangeInt := comptime genRange(int);
let evens := genRangeInt(
	0, 10, 2, fn(i) {
		return i;
	}
);
let _ := puts(evens);

// генератор с фильтром [f(i) for i in range if p(i)]
let genIf: fn(T: comptime type) -> fn(int, int, fn(int) -> bool, fn(int) -> T) -> []T = comptime fn(T) {
	let step: fn(int, int, fn(int) -> bool, fn(int) -> T, []T, int) -> []T = fn(from, to, predicate, genFun, acc, i) {
		return if (i >= to) {
			return acc;
		} else {
			return if (predicate(i) == true) {
				return step(from, to, predicate, genFun, push(acc, genFun(i)), i + 1);
			} else {
				return step(from, to, predicate, genFun, acc, i + 1);
			};
		};
	};
	return fn(from: int, to: int, predicate: fn(int) -> bool, genFn: fn(int) -> T) -> []T {
		return step(from, to, predicate, genFn, []T{}, from);
	};
};

// чётные квадраты от 0 до 20
let genIfInt := comptime genIf(int);
let r: []int = genIfInt(
	0, 20,
	fn(i) {
		return ((i / 2) * 2) == i;
	},
	fn(i) {
		return i * i;
	}
);
let _ := puts(r);
