package response

import "fmt"

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("cannot write body in state %s ", w.state)
	}
	return w.writer.Write(body)
}

func (w *Writer) WriteChunkedBody(body []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("cannot write body in state %s ", w.state)
	}

	body = append(body, []byte("\r\n")...)

	return w.writer.Write(body)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.state != WriterStateBody {
		return 0, fmt.Errorf("cannot write body in state %s ", w.state)
	}

	_, err := w.writer.Write([]byte(fmt.Sprintf("%X\r\n", 0)))
	if err != nil {
		return 0, err
	}

	return w.writer.Write([]byte("\r\n"))
}
