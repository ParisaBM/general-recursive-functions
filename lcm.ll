declare i32 @printf(i8* %0, ...)

declare i32 @atoi(...)

define i32 @pred(i32 %p0) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = alloca i32
	store i32 0, i32* %2
	br label %3

3:
	%4 = load i32, i32* %1
	%5 = icmp ult i32 %4, %p0
	br i1 %5, label %6, label %10

6:
	%7 = load i32, i32* %2
	%8 = add i32 %4, 0
	store i32 %8, i32* %2
	%9 = add i32 %4, 1
	store i32 %9, i32* %1
	br label %3

10:
	ret i32 %8
}

define i32 @add(i32 %p0, i32 %p1) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = alloca i32
	store i32 %p1, i32* %2
	br label %3

3:
	%4 = load i32, i32* %1
	%5 = icmp ult i32 %4, %p0
	br i1 %5, label %6, label %10

6:
	%7 = load i32, i32* %2
	%8 = add i32 %7, 1
	store i32 %8, i32* %2
	%9 = add i32 %4, 1
	store i32 %9, i32* %1
	br label %3

10:
	ret i32 %8
}

define i32 @mul(i32 %p0, i32 %p1) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = alloca i32
	store i32 0, i32* %2
	br label %3

3:
	%4 = load i32, i32* %1
	%5 = icmp ult i32 %4, %p0
	br i1 %5, label %6, label %10

6:
	%7 = load i32, i32* %2
	%8 = call i32 @add(i32 %7, i32 %p1)
	store i32 %8, i32* %2
	%9 = add i32 %4, 1
	store i32 %9, i32* %1
	br label %3

10:
	ret i32 %8
}

define i32 @sub(i32 %p0, i32 %p1) {
0:
	%1 = alloca i32
	store i32 0, i32* %1
	%2 = alloca i32
	store i32 %p1, i32* %2
	br label %3

3:
	%4 = load i32, i32* %1
	%5 = icmp ult i32 %4, %p0
	br i1 %5, label %6, label %10

6:
	%7 = load i32, i32* %2
	%8 = call i32 @pred(i32 %7)
	store i32 %8, i32* %2
	%9 = add i32 %4, 1
	store i32 %9, i32* %1
	br label %3

10:
	ret i32 %8
}

define i32 @div(i32 %p0, i32 %p1) {
0:
	%1 = alloca i32
	store i32 -1, i32* %1
	br label %2

2:
	%3 = load i32, i32* %1
	%4 = add i32 %3, 1
	store i32 %4, i32* %1
	%5 = call i32 @mul(i32 %4, i32 %p1)
	%6 = call i32 @sub(i32 %5, i32 %p0)
	%7 = icmp eq i32 %6, 0
	br i1 %7, label %8, label %2

8:
	ret i32 %4
}

define i32 @factors(i32 %p0, i32 %p1) {
0:
	%1 = call i32 @div(i32 %p1, i32 %p0)
	%2 = call i32 @mul(i32 %p0, i32 %1)
	%3 = call i32 @sub(i32 %p1, i32 %2)
	ret i32 %3
}

define i32 @main(i32 %argc, i8** %argv) {
0:
	%1 = getelementptr i8*, i8** %argv, i32 0
	%2 = load i8*, i8** %1
	%3 = call i32 (...) @atoi(i8* %2)
	%4 = getelementptr i8*, i8** %argv, i32 1
	%5 = load i8*, i8** %4
	%6 = call i32 (...) @atoi(i8* %5)
	%7 = alloca i32
	store i32 -1, i32* %7
	br label %8

8:
	%9 = load i32, i32* %7
	%10 = add i32 %9, 1
	store i32 %10, i32* %7
	%11 = add i32 %10, 1
	%12 = call i32 @factors(i32 %11, i32 %3)
	%13 = add i32 %10, 1
	%14 = call i32 @factors(i32 %13, i32 %6)
	%15 = call i32 @add(i32 %12, i32 %14)
	%16 = icmp eq i32 %15, 0
	br i1 %16, label %17, label %8

17:
	%18 = call i32 @pred(i32 %10)
	ret i32 %18
}
