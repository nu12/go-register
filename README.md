# go-register

A little package to help you avoid the Bart-Simpson-Syndrome:

![""](https://awalterschulze.github.io/blog/monads-for-goprogrammers/bartiferr.png)

## tl;dr

Replace something like this:
```go
func myFunc() error{
    var err error

    err = doSomething()
    if err != nil{
        myPrintError("Print error:", err)
        return err
    }

    err = doSomethingElse()
    if err != nil{
        myPrintError("Print error:", err)
        return err
    }

    err = doLastThing()
    if err != nil{
        myPrintError("Print error:", err)
        return err
    }

    return nil
}
```
With something like this:
```go
func myFunc() error{
    re := register.New()

    re.Run(
        func(){re.Err = doSomething()},
        func(){re.Err = doSomethingElse()},
        func(){re.Err = doLastThing()},
    ).IfError(
        func(){myPrintError("Print error:", re.Error())},
    )
    
    return re.Error()
}
```
## How it works

Instead of writing the infamous `if err != nil` multiple times, pass a series of anonymous functions to the `Run` function, assigning any error to the `re.Err` variable. The next function will only execute if error was nil in the previous step.

This means you can chain several steps without worrying about verifying errors as the register will... well, register... them for you. When the first one appears, subsequent steps will be skipped, avoiding nil pointers and other issues. 

At the end, you can get the error contained in the register with `re.Error()` and execute any error handling task with the `IfError` function, that will only execute the code block if error is not nil.

## Usage

All together:

```go
package main

import "github.com/nu12/go-register"

func main(){
    var b bool
    re := register.New()

    re.Run(
        func(){re.Err = executeThis("This runs anyway as it's the first step")},
        func(){re.Err = executeThat("This runs if previous step didn't returned an error")},
        func(){
            executeMultiple("We can have a multiline function, no problem")    
            b, re.Err = executeThat("Let's get that boolean value")
        },
    ).If(b,
        func(){executeIf("This runs if b is true and if no error was registered")},
    ).IfError(func(){errorHandle("This executes only if an error was registered, otherwise it'll be skipped")})

    print("Number of successful steps: " + re.Steps)
    print("You can also get the error: " + re.Error()) // or re.Err
}
```

## Acknowledgments

This approach does reduce the amount of `if err != nil`, but it's not necessarily pretty. If you have a solution to avoid repeating `func(){ ... }` all the time, please consider contributing to this project =]

If you like the idea and want to implement it in bigger projects, you can also build your own "register", possibly with your custom types built-in for easy of use.