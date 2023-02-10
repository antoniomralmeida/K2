package queue

import (
	"errors"
)

type Queue struct {
	Element any
	Next    *Queue
	Last    *Queue
}

func (q *Queue) Enqueue(elem any) {
	nq := Queue{Element: elem}
	if q == nil {
		q = &nq
	} else {
		q.Last.Next = &nq
	}
	q.Last = &nq
}

func (q *Queue) Dequeue() (r any) {
	if q == nil {
		r = nil
	} else {
		r = q.Element
		if q.Next != nil {
			q.Next.Last = q.Last
		}
		q = q.Next
	}
	return
}

func (q *Queue) GetLength() (r int) {
	if q == nil {
		return
	}
	p := *q
	for {
		r++
		if p.Next == nil {
			return
		}
		p = *p.Next
	}
}

func (q *Queue) IsEmpty() bool {
	return q == nil
}

func (q *Queue) Peek() (any, error) {
	if q.IsEmpty() {
		return nil, errors.New("empty queue")
	}
	return q.Element, nil
}
