package main

import (
	"os"
	"bufio"
)

//a stream, for lack of a better name is significant augment to go's channels
//it changes both the way the sender, and the reciever the communicate
//for the sender messages can be sent in a delayed fashion
//this allows one to significantly reorder tokens between phases of compilation
//put() is the basic send function
//begin_buffering() causes all calls to put to be buffered until a call to end_buffering() is made
//put_buffer() causes everything buffered to be sent
//for example:
//begin_buffering()
//put(A)
//end_buffering()
//put_buffer()
//this code segment will emit A once put_buffer() is reached
//like this it isn't very useful, however put calls can be made between end_buffering() and put_buffer()
//for example:
//begin_buffering()
//put(A)
//put(B)
//end_buffering()
//put(C)
//put(D)
//put_buffer()
//this will send C, D, A, B in that order
//in this way the three buffering functions can be though of as a binary operator with 3 delimiters
//where each function acts as one delimiter, and it swaps the order of its 2 outputs
//all 3 functions should be called exactly the same number of times
//still thinking of this as a binary operation, it can be nested
//for example:
//begin_buffering()
//begin_buffering()
//put(A)
//end_buffering()
//put(B)
//put_buffer()
//end_buffering()
//put(C)
//put(D)
//put_buffer()
//this will output C, D, B, A in that order
//how do we implement this?
//first we will see a naive implementation, then sort out some details
//we will use a stack of buffers
//begin_buffering() puts an empty buffer onto the stack
//end_buffering() puts an empty buffer onto the stack
//put_buffer() pops the top 2 buffers and concatenates them with the one originally on top of the stack ending up on the left
//put() adds an element to the buffer on top of the stack
//this is mostly complete, however nothing ever gets emitted
//we can change put_buffer() to check if the resulting stack has one element, and if it does emit it
//also change put() to simply emit if the stack is empty
//now we have described a working implementation
//in fact, this solution is simpler than the one given here, because begin_buffering, and end_buffering can be the same function
//that is because this implementation will often delay emissions unnecessarily
//for example:
//begin_buffering()
//put(A)
//end_buffering()
//put(B)
//put(C)
//put(D)
//put_buffer()
//this will not emit B until put_buffer(), although it should be emitted right away
//we can solve this using a boolean can_emit
//end_buffering() should check if the stack has height 1, and if it does set can_emit to true
//the buffers are implemented as linked lists because they have O(1) concatenation
//the reciever thankfully is much simpler
//streams allow the reciever to peek at the input
//get() recieves a message from the sender
//undo() puts the last message back so that the next time get() is called it will get the same thing again
//undo() can't be called multiple times in a row

//ch is the channel that this is built atop
type Stream struct {
	ch chan byte
	buffer_stack []List
	can_emit bool
	buf byte
	buf_in_use bool
}

func new_stream() Stream {
	return Stream{make(chan byte), make([]List, 0), true, 0, false}
}

func from_file(file *File) Stream {
	s := new_stream()
	go func() {
		br := bufio.NewReader(file)
		for b, err := br.ReadByte {
			if err != nil {
				panic(err)
			}
			s.put(b)
		}
	}()
	return s
}

func put(s Stream, b byte) {
	if can_emit {
		s.ch <- b
	} else {
		s.delay_stack[len(s.delay_stack)-1].push(b)
	}
}


func begin_buffering(s Stream) {
	can_emit = false
	s.delay_stack = append(s.delay_stack, new_list())
}

func end_buffering(s Stream) {
	if len(s.delay_stack) == 0 {
		panic("nothing to end")
	}
	if len(s.delay_stack) == 1 {
		s.can_emit = true
	} else {
		s.delay_stack = append(s.delay_stack, new_list())
	}
}

func put_buffer(s Stream) {
	if len(s.delay_stack) <= 1 {
		panic("nothing to put")
	} else {
		s.buffer_stack[len(s.buffer_stack)-2].prepend(s.buffer_stack[len(s.buffer_stack)-1])
		s.buffer_stack = s.buffer_stack[:len(s.buffer_stack)-1]
		if len(s.buffer_stack) == 1 {
			//it iterates through the elements and outputs them to ch
			it := s.delay_stack[0].head
			for it != nil {
				s.ch <- it.data
				it = it.next
			}
			buffer_stack = make([]List, 0)
		}
	}
}

func get(s Stream) byte {
	if buf_in_use {
		buf_in_use = false
	} else {
		buf = <- ch
	}
	return buf
}

func undo(s Stream) {
	if buf_in_use {
		panic("multiple undo")
	}
	buf_in_use = true
}

//here are the list methods used earlier
//these should all be fairly standard
type Node struct {
	data byte
	next *Node
}

type List struct {
	head *Node
	tail *Node
}

func new_list() List {
	List{nil, nil}
}

func push(l List, b byte) {
	if tail == nil {
		tail = &Node{b, nil}
		head = tail
	} else {
		tail.next = &Node{b, nil}
		tail = tail.next
	}
}

//adds the contents of l1 to the front of l0
func prepend(l0, l1 List) {
	//if l1 is empty nothing has to happen
	if l1.tail != nil {
		l1.tail.next = l0.head
		l0.head = l1.head
		if l0.tail == nil {
			l0.tail = l1.tail
		}
	}
}
