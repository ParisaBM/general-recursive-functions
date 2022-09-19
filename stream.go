package main

import (
	"bufio"
	"errors"
	"io"
	"os"
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
//we will use a stack of buffers
//begin_buffering() puts an empty buffer onto the stack
//end_buffering() puts an empty buffer onto the stack
//put_buffer() pops the top, and 2nd from top and adds them to the 3rd from top element of the stack
//if there aren't 3 elements on the stack they get emitted instead
//put() adds an element to the buffer on top of the stack, unless its empty in which case it emits it
//note that begin_buffering() and end_buffering() are actually the same thing
//in our implementation we just have 1 function delimit_buffering()
//it is accompanied by a comment explaining which one it everywhere it is used
//the reciever thankfully is much simpler
//streams allow the reciever to peek at the input
//get() recieves a message from the sender
//undo() puts the last message back so that the next time get() is called it will get the same thing again
//undo() can't be called multiple times in a row

// ch is the channel that this is built atop
type Stream struct {
	ch           chan byte
	buffer_stack []List
	buf          byte
	buf_in_use   bool
}

func new_stream() Stream {
	return Stream{make(chan byte), make([]List, 0), 0, false}
}

// this function wasn't described in the long sequence of comments above
// it takes the contents of a file, makes a new stream, and sends the file down the stream
// it also closes the file
func from_file(file *os.File) Stream {
	s := new_stream()
	go func() {
		br := bufio.NewReader(file)
		for {
			b, err := br.ReadByte()
			if err != nil {
				if errors.Is(err, io.EOF) {
					s.put('\x00')
					file.Close()
				} else {
					panic(err)
				}
			}
			s.put(b)
		}
	}()
	return s
}

func (s *Stream) put(b byte) {
	if len(s.buffer_stack) == 0 {
		s.ch <- b
	} else {
		s.buffer_stack[len(s.buffer_stack)-1].push(b)
	}
}

func (s *Stream) delimit_buffering() {
	s.buffer_stack = append(s.buffer_stack, new_list())
}

func (s *Stream) put_buffer() {
	s.buffer_stack[len(s.buffer_stack)-2].prepend(s.buffer_stack[len(s.buffer_stack)-1])
	if len(s.buffer_stack) == 2 {
		//it iterates through the elements and outputs them to ch
		it := s.buffer_stack[0].head
		for it != nil {
			s.ch <- it.data
			it = it.next
		}
	} else {
		s.buffer_stack[len(s.buffer_stack)-3].list_append(s.buffer_stack[len(s.buffer_stack)-2])
	}
	s.buffer_stack = s.buffer_stack[:len(s.buffer_stack)-2]
}

func (s *Stream) get() byte {
	if s.buf_in_use {
		s.buf_in_use = false
	} else {
		s.buf = <-s.ch
	}
	return s.buf
}

func (s *Stream) undo() {
	if s.buf_in_use {
		panic("multiple undo")
	}
	s.buf_in_use = true
}

// here are the list methods used earlier
// these should all be fairly standard
type Node struct {
	data byte
	next *Node
}

type List struct {
	head *Node
	tail *Node
}

func new_list() List {
	return List{nil, nil}
}

func (l *List) push(b byte) {
	if l.tail == nil {
		l.tail = &Node{b, nil}
		l.head = l.tail
	} else {
		l.tail.next = &Node{b, nil}
		l.tail = l.tail.next
	}
}

// adds the contents of l1 to the front of l0
func (l0 *List) prepend(l1 List) {
	//if l1 is empty nothing has to happen
	if l1.tail != nil {
		l1.tail.next = l0.head
		l0.head = l1.head
		if l0.tail == nil {
			l0.tail = l1.tail
		}
	}
}

// adds the contents of l1 after l0
// the name append is taken
func (l0 *List) list_append(l1 List) {
	l1.prepend(*l0)
	*l0 = l1
}
