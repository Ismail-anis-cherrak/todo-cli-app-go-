package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*t = append(*t, todo)
}

func (t *Todos) Complete(i int) error {
	ls := *t
	if i <= 0 || i > len(ls) {
		return errors.New("invalid index")
	}

	ls[i-1].CompletedAt = time.Now()
	ls[i-1].Done = true
	return nil

}

func (t *Todos) Delete(i int) error {
	ls := *t
	if i <= 0 || i > len(ls) {
		return errors.New("invalid index")
	}

	*t = append(ls[:i-1], ls[i:]...)

	return nil

}

func (t *Todos) Load(filename string) error {

	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)

	if err != nil {
		return err
	}

	return nil

}

func (t *Todos) Store(filename string) error {

	data, err := json.Marshal(t)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {

	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "Created At"},
			{Align: simpletable.AlignCenter, Text: "Done At"},
		},
	}
	var cells [][]*simpletable.Cell

	for i, item := range *t {
		i++
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i)},
			{Text: item.Task},
			{Text: fmt.Sprintf("%t", item.Done)},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CreatedAt.Format(time.RFC822)},
		})
	}
	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: "Your todos are here"},
	}}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}
