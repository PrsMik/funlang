let OkType: fn(T: type) -> type = fn(T: type) {
    return interface {
        let value: type = T;
    };
};

let ErrType: fn(E: type) -> type = fn(E: type) {
    return interface {
        let error: type = E;
    };
};

let Result: fn(T: type, E: type) -> type = fn(T: type, E: type) {
    return OkType(T) | ErrType(E); 
};

let Ok: fn(T: type, val: T) -> OkType(T) = fn(T: type, val: T) {
    return impl (OkType(T)) {
        let value: T = val;
    };
};

let Err: fn(E: type, err_val: E) -> ErrType(E) = fn(E: type, err_val: E) {
    return impl (ErrType(E)) {
        let error: E = err_val;
    };
};

let divide: fn(int, int) -> Result(int, string) = fn(a, b) {
    if (b == 0) {
        return Err(string, "Division by zero"); 
    } else {
        return Ok(int, a / b); 
    }
};

let res: Result(int, string) = divide(10, 0);

let error_text: string = switch type(res) {
    case ErrType(string) {
        return res.error; 
    }
    case OkType(int) {
        return "Ошибки нет, результат: " + intToString(res.value); // псевдовызов
    }
};

puts(error_text);