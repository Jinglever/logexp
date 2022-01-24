# Logical Expression Compiler

# Functions:
	- Compile(exp string)
	- Match(text string)

Usage Example:
```
exp := "(hello|hi)&world"
expression, err := logexp.Compile(exp)
if err != nil {
	// TODO
} else {
	text := "hello world"
	hit := expression.Match(text)
	fmt.Println(hit)
}
```