let map := comptime fn(T: comptime type, U: comptime type) -> fn([]T, fn(T) -> U) -> []U {
	let step := fn(arr: []T, f: fn(T) -> U, acc: []U, i: int) -> []U {
		return if (i < 0) {
			return acc;
		} else {
			return step(arr, f, push(acc, f(arr[i])), i - 1);
		};
	};
	return fn(arr: []T, aplliedFn: fn(T) -> U) -> []U {
		return step(arr, aplliedFn, []T{}, len(arr) - 1);
	};
};

let nums := []int{1, 2, 3, 4, 5, 6};

let evenFn := fn(x: int) -> bool {
	let div: int = x - (x / 2) * 2;
	return if (div == 0) {
		return true;
	} else {
		return false;
	};
};

let mapIntToInt := comptime map(int, int);
let doubled := mapIntToInt(
	nums, fn(x) {
		return x * 2;
	}
);
let o1 := puts(doubled[2]); // 6

let mapIntToBool := comptime map(int, bool);
let flags := mapIntToBool(nums, evenFn);
let o2 := puts(flags);
