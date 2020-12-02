package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Next  *ListItem
	Prev  *ListItem
	Value interface{}
}

// Извлечь элемент из цепочки.
func (thisItem *ListItem) Extract() {
	if thisItem.Next != nil {
		thisItem.Next.Prev = thisItem.Prev
	}

	if thisItem.Prev != nil {
		thisItem.Prev.Next = thisItem.Next
	}

	thisItem.Next = nil
	thisItem.Prev = nil
}

// Установить следующий со взаимной привязкой.
func (thisItem *ListItem) SetNext(newNext *ListItem) {
	if newNext == nil {
		thisItem.Next = newNext
		return
	}

	newNext.Next = thisItem.Next
	thisItem.Next = newNext
	newNext.Prev = thisItem
}

// Установить предыдущий со взаимной привязкой.
func (thisItem *ListItem) SetPrev(newPrev *ListItem) {
	if newPrev == nil {
		thisItem.Prev = newPrev
		return
	}

	newPrev.Prev = thisItem.Prev
	thisItem.Prev = newPrev
	newPrev.Next = thisItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int // длина списка
}

func (lst *list) Len() int {
	return lst.len
}

// первый элемент списка.
func (lst *list) Front() *ListItem {
	return lst.front
}

// последний элемент списка.
func (lst *list) Back() *ListItem {
	return lst.back
}

// Добавить элемент в начало.
func (lst *list) PushFront(v interface{}) *ListItem {
	lst.len++

	var li ListItem
	li.Value = v

	// nil i! <- (prev) front <-> ... <-> elem <-> ... <-> back (next) -> nil
	if lst.front != nil {
		lst.front.SetPrev(&li)
	}

	lst.front = &li

	if lst.back == nil {
		lst.back = &li
	}

	return lst.Front()
}

// добавить значение в конец.
func (lst *list) PushBack(v interface{}) *ListItem {
	lst.len++

	var li ListItem
	li.Value = v

	// nil <- (prev) front <-> ... <-> elem <-> ... <-> back (next) -> i! nil
	if lst.back != nil {
		lst.back.SetNext(&li)
	}

	lst.back = &li

	if lst.front == nil {
		lst.front = &li
	}

	return lst.Back()
}

// удалить элемент.
func (lst *list) Remove(i *ListItem) {
	lst.len--

	// nil <- (prev) front i! <-> ... <-> elem <-> ... <-> back (next) -> nil
	if i == lst.front {
		lst.front = i.Next
	}

	// nil <- (prev) front <-> ... <-> elem <-> ... <-> i! back (next) -> nil
	if i == lst.back {
		lst.back = i.Prev
	}

	i.Extract()
}

// переместить элемент в начало.
func (lst *list) MoveToFront(i *ListItem) {
	if lst.front == i {
		return
	}

	if lst.back == i {
		lst.back = i.Prev
	}

	i.Extract()

	// i! <- (prev) front <-> ... <-> elem <-> ... <-> back (next) -> nil
	if lst.front != nil {
		lst.front.SetPrev(i)
	}

	lst.front = i
}

func NewList() List {
	return &list{}
}
