let loop: fn(int, fn(int) -> bool, fn(int) -> int, fn(int) -> int) -> bool = fn(initial, condition, update, body) {
    let iter: fn(int) -> bool = fn(current) {
        return if (condition(current) == true) {
            let _: int = body(current);
            return iter(update(current));
        } else {
            return false;
        };
    };
    return iter(initial);
};

let body: fn(int) -> int = fn(i) { 
    let _: int = puts(i);
    return 0;
};

let res: bool = loop(0, fn(i) {return i <= 100;}, fn(i) {return i + 1;}, body);

