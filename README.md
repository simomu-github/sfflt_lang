# Simple FFLT lang

Simple FFLT lang is transpiler from a simple programming langauge to the esolang [FFLT lang](https://github.com/simomu-github/fflt_lang).

## Usage

```
sfflt_lang program.sflt
```

## Building yourself

```
make build
```

## Example

```plain text:fibonacci.sflt
// This is fibonacci number
func fib(n) {
  if (n < 2) return n;

  return fib(n - 1) + fib(n - 2);
}

var input = getn;
putn fib(input);
```

Compile to FFLT lang

```
sfflt_lang fibonacci.sflt
```

<details><summary>Compiled FFLT lang</summary>

```plain text:fibonacci.fflt
FFFLLFFLLLLLLFLLFLLFLFFLLLFLLLFLLLFFFFLLLFLFLLLLFLFFTFFFLLLFLFFL
LLFLFFLLFLLFLLLLFFFFTLTLLFFFLLLFLFFLLLFLFFLLFLLFLLLLFFFFTLLLLLFF
FFLLFFLLLLLLFLLFLLFLFFLLLFLLLFLLLFFFFLLLFLFLLLLFLFFTLLLTFLLLFFLL
LLLFFLLFLLFFLLFLLFLFFLLLFFFLFTLTFLTTTTFFLLFFLLLLLFFLLFLLFFLLFLLF
LFFLLLFFFLFTFLFFFTFFFLFTLFFLTLLLLFLLFFFTFFFFTTFTLLFLLFFLTTFFLLFL
LFFFTFFFLTTFFLLFLLFFLTTLFLLFLLFFLFTFLFFFTFLTFLTTLTTFTLLFLLFFLLTT
FFLLFLLFFLFTTFFLLFLLFFLLTFLFFFTFFFLTLFFLTFLLLFFLLLLLFFLLFLLFFLLF
LLFLFFLLLFFFLFTFLFFLTFFFLFTLFFLTFLLLFFLLLLLFFLLFLLFFLLFLLFLFFLLL
FFFLFTLFFFFLTFLTTLTFFFFTFLTFLTTLT
```

</details>

execute FFLT lang

```
fflt_lang fibonacci.fflt
```

## Simple FFLT lang specification

### Expressions

#### Integer literal

`1`, `23`, `456`, ...

#### Character literal

`'a'`, `'\n'`, ...

#### String literal

`"Hello World!"`, `"One\nTwo\n"`, ...

#### Array literal

`[123, 456]`, `['a', 'b', 'c']`, ...

#### Boolean literal

`true` or `false`

#### Variable

`a`, `hoge`, `piyo123`...

#### Arithmetic operations

```
+<expression>
-<expression>
<expression> + <expression>
<expression> - <expression>
<expression> * <expression>
<expression> / <expression>
<expression> % <expression>
```

#### Logical operations

Support short-circuit evaluation.

```
!<expression>
<expression> && <expression>
<expression> || <expression>
```

#### Comparison operations

```
<expression> == <expression>
<expression> != <expression>
<expression> < <expression>
<expression> <= <expression>
<expression> > <expression>
<expression> >= <expression>
```

#### Parentheses

```
(<expression>)
```

#### Call function

```
<identifier>(<expression>, <expression>, ...)
```

#### Assignment

```
<identifier> = <expression>
```

### Statements

#### Expression statement

```
<expression>;
```

#### Variable declaration

```
var <identifier> = <expression>;
```

#### Function declaration

Function parameters are localized

```
func <identifier>(<identifier>, <identifier>, ...) {
  <statement>

  // Support return statement
  return <expression>;
}
```

#### If statement

```
if (<expression>) {
  <statement>
} else if (<expression>) {
  <statement>
} else {
  <statement>
}
```

#### While statement

```
while (<expression>) {
  <statement>

  // Support break statement
  break;
}
```

#### For statement

```
for (<declaration>;<expression>;<expression>) {
  <statement>

  // Support break statement
  break;
}
```

#### include statement

```
include "./path/to/other_file.sflt";
include "strings" // Load standard library.
```

### I/O

#### stdout

Write integer.

```
putn 123; 
```

Write character.

```
putc 'c'; 
```

#### stdin

Read integer.

```
var integer = getn;
```

Read character.

```
var character = getc;
```

### Build in functions

#### len

Get the length of an array.

```
len([1, 2, 3]);
```

#### append

Appending elements to an array.

```
var array = [1, 2, 3];
append(array, 4);
```

### Standard library


#### `strings`

- `println`

Write string with new line.

```
include "strings";

println("Hello World!");
```


### Comment

```
// This is comment
```


## TODO

- Static typing
- Memory management
...
