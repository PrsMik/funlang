let OkType := comptime fn(T: type) -> type {
	return interface {
		let value: type = T;
	};
};

let ErrType := comptime fn(E: type) -> type {
	return interface {
		let error: type = E;
	};
};

let Result := comptime fn(T: type, E: type) -> type {
	return OkType(T) | ErrType(E);
};

let MakeOkFor := comptime fn(T: comptime type) -> fn(T) -> comptime OkType(T) {
	return fn(val: comptime T) -> comptime OkType(T) {
		let ok: type = impl (comptime OkType(T)){
			let value: T = val;
		};
		return ok.{value:val};
	};
};

let MakeErrFor := comptime fn(E: comptime type) -> fn(E) -> comptime ErrType(E) {
	return fn(err_val: comptime E) -> comptime ErrType(E) {
		let err: type = impl (comptime ErrType(E)){
			let error: E = err_val;
		};
		return err.{error:err_val};
	};
};

let resTp := comptime Result(int, string);
let divide := fn(a: int, b: int) -> resTp {
	return if (b == 0) {
		let res: resTp = (comptime MakeErrFor(string))("Division by zero");
		return res;
	} else {
		let res: resTp = (comptime MakeOkFor(int))(a / b);
		return res;
	};
};

let res := divide(10, 2);

let error_text := switch res.(type) {case comptime ErrType(string){let o:= puts(res.error);return res.error;}case comptime OkType(int){let o:= puts(res.value);return "";}default {return "";}};
