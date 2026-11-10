let myTail: fn(T: comptime type) -> fn([]T) -> []T = comptime fn(T) {
	let step: fn([]T, []T, int) -> []T = fn(arr, acc, i) {
		return if i < 1 {
			return acc;
		} else {
			return step(arr, push(acc, arr[i]), i - 1);
		};
	};
	return fn(arr: []T) -> ([]T) {
		return step(arr, []T{}, len(arr) - 1);
	};
};

let myTailInt: fn(arr: []int) -> []int = comptime myTail(int);
let nums: []int = []int{1, 2, 3, 4, 5};

let t: []int = myTailInt(nums);
let o: int = puts(t); // [2, 3, 4, 5]

let t2: []int = myTailInt([]int{10, 20, 30});
let o2: int = puts(t2); // [20, 30]
