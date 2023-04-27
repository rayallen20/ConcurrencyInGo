package main

func main() {
}

func tee(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			// Tips: channel是引用类型 因此写入到局部变量out1的操作会影响到goroutine外定义的out1
			out1, out2 := out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					// Tips: 还是因为channel是引用类型,因此声明一个channel类型的变量,实际上是在栈上分配了一个
					// Tips: 指向底层数据结构的指针.这个指针指向在堆上分配的实际数据结构
					// Tips: 因此此处out1 = nil实际上修改的是循环中的局部变量out1的值,而非是指向channel底层数据结构的指针
					// Tips: 在下一次迭代时,循环中的局部变量out1将会重新被分配为第7行的out1
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func orDone(done <-chan interface{}, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}

				select {
				case valueStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valueStream
}
