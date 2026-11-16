let loop := fn(initial: int, condition: fn(int) -> bool, update: fn(int) -> int, body: fn(int) -> int) -> bool {
	let iter := fn(current: int) -> bool {
		return if (condition(current) == true) {
			let _ := body(current);
			return iter(update(current));
		} else {
			return false;
		};
	};
	return iter(initial);
};

let body := fn(i: int) -> int {
	let _: int = puts(i);
	return 0;
};

let res := loop(
	0, fn(i) {
		return i <= 100;
	}, fn(i) {
		return i + 1;
	}, body
);
