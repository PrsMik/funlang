let OkType: fn(T: type) -> type = fn(T) {
    return interface {
        let value: type = T;
    };
};

let ErrType: fn(E: type) -> type = fn(E) {
    return interface {
        let error: type = E;
    };
};

let Result: fn(T: type, E: type) -> type = fn(T, E) {
    return OkType(T) | ErrType(E); 
};

let Ok: fn(T: type, val: T) -> OkType(T) = fn(T, val) {
    let ok: type = impl (OkType(T)) {
        let value: T = val;
    };
    return ok.{value : val};
};

let Err: fn(E: type, err_val: E) -> ErrType(E) = fn(E, err_val) {
    let err: type = impl (ErrType(E)) {
        let error: E = err_val;
    };
    return err.{error : err_val};
};

let divide: fn(int, int) -> Result(int, string) = fn(a, b) {
    return if (b == 0) {
        return Err(string, "Division by zero"); 
    } else {
        return Ok(int, a / b); 
    };
};

let res: Result(int, string) = divide(10, 0);

let error_text: string = switch res.(type) {
    case ErrType(string) {
        return res.error; 
    }
    case OkType(int) {
        return "Ошибки нет, результат: "; // псевдовызов
    }
    default {
        return "";
    }
};

return puts(error_text);