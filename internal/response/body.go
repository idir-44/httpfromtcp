package response

import "fmt"

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("cannot write body in state %s ", w.state)
	}
	return w.writer.Write(body)
}
