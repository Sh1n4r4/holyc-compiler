// Test HolyC Program
// Simple test file for holyc-go compiler

U0 TestFunction()
{
    I64 i;
    i = 42;
    return;
}

I64 Add(I64 a, I64 b)
{
    return a + b;
}

I64 Factorial(I64 n)
{
    if (n <= 1) {
        return 1;
    }
    return n * Factorial(n - 1);
}

I64 Fibonacci(I64 n)
{
    I64 a, b, temp, i;
    
    if (n <= 0) return 0;
    if (n == 1) return 1;
    
    a = 0;
    b = 1;
    
    for (i = 2; i < n + 1; i = i + 1) {
        temp = a + b;
        a = b;
        b = temp;
    }
    
    return b;
}

U0 TestLoops()
{
    I64 i;
    
    // While loop
    i = 0;
    while (i < 10) {
        i = i + 1;
    }
    
    // For loop
    for (i = 0; i < 10; i = i + 1) {
        // Empty
    }
    
    // Do-while
    i = 0;
    do {
        i = i + 1;
    } while (i < 5);
    
    return;
}

I64 Max(I64 a, I64 b)
{
    if (a > b) {
        return a;
    } else {
        return b;
    }
}

I64 Abs(I64 x)
{
    if (x < 0) {
        return 0 - x;
    }
    return x;
}

// Class example
class Point
{
    I64 x;
    I64 y;
};

I64 Main()
{
    I64 result;
    
    result = Add(10, 20);
    result = Factorial(5);
    result = Fibonacci(10);
    result = Max(100, 200);
    result = Abs(0 - 42);
    
    TestLoops();
    
    return 0;
}

// Extern functions (TempleOS)
extern U0 Print(U8 *fmt, ...);
extern I64 StrLen(U8 *s);
