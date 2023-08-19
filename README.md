# async
go async library


## Fiber

类似PHP8的Fiber

纤程（Fiber）表示一组有完整栈、可中断的功能。 纤程可以在调用堆栈中的任何位置被挂起，在纤程内暂停执行，直到稍后恢复。

纤程可以暂停整个执行堆栈，所以该函数的直接调用者不需要改变调用这个函数的方式。

你可以在调用堆栈的任意地方使用 suspend() 中断执行（也就是说，suspend() 的调用位置可以在一个深度嵌套的函数中，甚至可以不存在）。

纤程一旦被暂停，可以使用 Fiber::resume() 传递任意值.

```golang
	fib := fiber.New(func(suspend SuspendFunc[string]) string {
		value := suspend("fiber")
		fmt.Println("Value used to resume fiber:", value)
		return "return val"
	})

	value, _ := fib.Start()
	fmt.Println("Value from fiber suspending:", value)
	fib.Resume("test")
	ret, _ := fib.GetReturn()
	fmt.Println("Value from fiber return:", ret)
```

以上示例会输出：

```
Value from fiber suspending: fiber
Value used to resume fiber: test
Value from fiber return: return val
```