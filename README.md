# textedit-go
Basic terminal text editor written in Go. This was implemented as a test bed
to learn about how text editor data structures are implemented. No plans on 
making this production worthy.

## Notes:
* Uses a [gap buffer](https://en.wikipedia.org/wiki/Gap_buffer) for editing.
* Makes use of [Bubble Tea](https://github.com/charmbracelet/bubbletea) for application loop and terminal control.

## Todo:
- [x] Unit test full coverage for gap buffer
- [ ] Horizontal line scrolling
- [ ] Word wrap mode
- [ ] Undo/redo
- [ ] Rope or other data structure

