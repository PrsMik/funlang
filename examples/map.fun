let map: fn(T: comptime type, U: comptime type) -> fn([T], fn(T) -> U) -> [U] = comptime fn(T, U) {
    let step: fn([T], fn(T) -> U, [U], int) -> [U] = fn(arr, f, acc, i) {
        return if (i < 0) {
            return acc;
        } else {
            return step(arr, f, push(acc, f(arr[i])), i - 1);
        };
    };
    return fn(arr: [T], aplliedFn: fn(T) -> U) -> [U] {
        return step(arr, aplliedFn, [], len(arr) - 1);
    };
};

let nums: [int] = [1, 2, 3, 4, 5, 6];

let evenFn: fn(int) -> bool = fn(x) {
	let div: int = x - (x / 2) * 2;
	return if (div == 0) {
		return true;
	} else {
		return false;
	};
};

let mapIntToInt: fn(arr: [int], aplliedFn: fn(int) -> int) -> [int] = comptime map(int, int);
let doubled: [int] = mapIntToInt(nums, fn(x) { return x * 2; });
let o1: int = puts(doubled[2]);              // 6

let mapIntToBool: fn(arr: [int], aplliedFn: fn(int) -> bool) -> [bool] = comptime map(int, bool);
let flags: [bool] = mapIntToBool(nums, evenFn);
let o2: int = puts(flags);
