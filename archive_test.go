package main

import (
	"testing"
	"time"
)

func TestArchive(t *testing.T) {
	a := NewArchive()
	t.Logf("Archive: '%v'", a)

	p1 := NewPage()
	p2 := NewPage()
	p3 := NewPage()

	p1.PublishDate = time.Now().Add(time.Hour * -2)
	p2.PublishDate = time.Now()
	p3.PublishDate = time.Now().Add(time.Hour * -1)

	a.AddPage(p1)
	a.AddPage(p2)
	a.AddPage(p3)

	t.Logf("Pages: '%v'", a.Pages)

	a.SortByDate()

	t.Logf("Pages: '%v'", a.Pages)

	p1.PublishDate = time.Now().Add(time.Hour * 1)

	a.SortByDate()

	t.Logf("Pages: '%v'", a.Pages)

	t.Logf("Time now: '%v", time.Now().Nanosecond())

}
