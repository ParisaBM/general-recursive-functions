@fstr = global [4 x i8] c"%d\0A\00"

declare i32 @printf(i8* %0, ...)

declare i32 @atoi(...)

define void @main(i32 %argc, i8** %argv) {
0:
	%1 = getelementptr i8*, i8** %argv, i32 1
	%2 = load i8*, i8** %1
	%3 = call i32 (...) @atoi(i8* %2)
	%4 = add i32 %3, 1
	%5 = getelementptr [4 x i8], [4 x i8]* @fstr, i32 0, i32 0
	%6 = call i32 (i8*, ...) @printf(i8* %5, i32 %4)
	ret void
}
