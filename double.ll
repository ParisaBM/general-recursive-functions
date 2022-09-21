@fstr = global [4 x i8] c"%d\0A\00"

declare i32 @printf(i8* %0, ...)

declare i32 @atoi(...)

define i32 @add2(i32 %p0) {
0:
	%1 = add i32 %p0, 1
	%2 = add i32 %1, 1
	ret i32 %2
}

define void @main(i32 %argc, i8** %argv) {
0:
	%1 = getelementptr i8*, i8** %argv, i32 1
	%2 = load i8*, i8** %1
	%3 = call i32 (...) @atoi(i8* %2)
	%4 = alloca i32
	store i32 0, i32* %4
	%5 = alloca i32
	store i32 0, i32* %5
	br label %6

6:
	%7 = load i32, i32* %4
	%8 = icmp ult i32 %7, %3
	br i1 %8, label %9, label %13

9:
	%10 = load i32, i32* %5
	%11 = call i32 @add2(i32 %10)
	store i32 %11, i32* %5
	%12 = add i32 %7, 1
	store i32 %12, i32* %4
	br label %6

13:
	%14 = getelementptr [4 x i8], [4 x i8]* @fstr, i32 0, i32 0
	%15 = call i32 (i8*, ...) @printf(i8* %14, i32 %11)
	ret void
}
