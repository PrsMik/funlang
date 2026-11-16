let myTail := comptime fn(T: comptime type) -> fn([]T) -> []T {
	let step := fn(arr: []T, acc: []T, i: int) -> []T {
		return if i < 1 {
			return acc;
		} else {
			return step(arr, push(acc, arr[i]), i - 1);
		};
	};
	return fn(arr: []T) -> []T {
		return step(arr, []T{}, len(arr) - 1);
	};
};

let myTailInt := comptime myTail(int);
let nums := []int{1, 2, 3, 4, 5};

let t := myTailInt(nums);
let o := puts(t); // [2, 3, 4, 5]

let t2 := myTailInt([]int{10, 20, 30});
let o2 := puts(t2); // [20, 30]
