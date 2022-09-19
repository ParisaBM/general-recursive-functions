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

define i32 @main(i32 %p0, i32 %p1) {
0:
	%1 = alloca i32
	store i32 -1, i32* %1
	br label %2

2:
	%3 = load i32, i32* %1
	%4 = add i32 %3, 1
	store i32 %4, i32* %1
	%5 = add i32 %4, 1
	%6 = call i32 @factors(i32 %5, i32 %p0)
	%7 = add i32 %4, 1
	%8 = call i32 @factors(i32 %7, i32 %p1)
	%9 = call i32 @add(i32 %6, i32 %8)
	%10 = icmp eq i32 %9, 0
	br i1 %10, label %11, label %2

11:
	%12 = call i32 @pred(i32 %4)
	ret i32 %12
}
