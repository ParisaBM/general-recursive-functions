package main

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// a stream, for lack of a better name is significant augment to go's channels
// it changes both the way the sender, and the reciever the communicate
// for the sender messages can be sent in a delayed fashion
// this allows one to significantly reorder tokens between phases of compilation
// put() is the basic send function
// beginBuffering() causes all calls to put to be buffered until a call to endBuffering() is made
// putBuffer() causes everything buffered to be sent
// for example:
// beginBuffering()
// put(A)
// endBuffering()
// putBuffer()
// this code segment will emit A once putBuffer() is reached
// like this it isn't very useful, however put calls can be made between endBuffering() and putBuffer()
// for example:
// beginBuffering()
// put(A)
// put(B)
// endBuffering()
// put(C)
// put(D)
// putBuffer()
// this will send C, D, A, B in that order
// in this way the three buffering functions can be though of as a binary operator with 3 delimiters
// where each function acts as one delimiter, and it swaps the order of its 2 outputs
// all 3 functions should be called exactly the same number of times
// still thinking of this as a binary operation, it can be nested
// for example:
// beginBuffering()
// beginBuffering()
// put(A)
// endBuffering()
// put(B)
// putBuffer()
// endBuffering()
// put(C)
// put(D)
// putBuffer()
// this will output C, D, B, A in that order
// how do we implement this?
// we will use a stack of buffers
// beginBuffering() puts an empty buffer onto the stack
// endBuffering() puts an empty buffer onto the stack
// putBuffer() pops the top, and 2nd from top and adds them to the 3rd from top element of the stack
// if there aren't 3 elements on the stack they get emitted instead
// put() adds an element to the buffer on top of the stack, unless its empty in which case it emits it
// note that beginBuffering() and endBuffering() are actually the same thing
// in our implementation we just have 1 function delimitBuffering()
// it is accompanied by a comment explaining which one it everywhere it is used
// the reciever thankfully is much simpler
// streams allow the reciever to peek at the input
// get() recieves a message from the sender
// undo() puts the last message back so that the next time get() is called it will get the same thing again
// undo() can't be called multiple times in a row

// ch is the channel that this is built atop
type Stream struct {
	ch          chan byte
	bufferStack []List
	buf         byte
	bufInUse    bool
}

func newStream() Stream {
	return Stream{make(chan byte), make([]List, 0), 0, false}
}

// this function wasn't described in the long sequence of comments above
// it takes the contents of a file, makes a new stream, and sends the file down the stream
// it also closes the file
func fromFile(file *os.File) Stream {
	s := newStream()
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
	if len(s.bufferStack) == 0 {
		s.ch <- b
	} else {
		s.bufferStack[len(s.bufferStack)-1].push(b)
	}
}

func (s *Stream) delimitBuffering() {
	s.bufferStack = append(s.bufferStack, newList())
}

func (s *Stream) putBuffer() {
	s.bufferStack[len(s.bufferStack)-2].prepend(s.bufferStack[len(s.bufferStack)-1])
	if len(s.bufferStack) == 2 {
		// it iterates through the elements and outputs them to ch
		it := s.bufferStack[0].head
		for it != nil {
			s.ch <- it.data
			it = it.next
		}
	} else {
		s.bufferStack[len(s.bufferStack)-3].append(s.bufferStack[len(s.bufferStack)-2])
	}
	s.bufferStack = s.bufferStack[:len(s.bufferStack)-2]
}

func (s *Stream) get() byte {
	if s.bufInUse {
		s.bufInUse = false
	} else {
		s.buf = <-s.ch
	}
	return s.buf
}

func (s *Stream) undo() {
	if s.bufInUse {
		panic("multiple undo")
	}
	s.bufInUse = true
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

func newList() List {
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
	// if l1 is empty nothing has to happen
	if l1.tail != nil {
		l1.tail.next = l0.head
		l0.head = l1.head
		if l0.tail == nil {
			l0.tail = l1.tail
		}
	}
}

// adds the contents of l1 after l0
func (l0 *List) append(l1 List) {
	l1.prepend(*l0)
	*l0 = l1
}
