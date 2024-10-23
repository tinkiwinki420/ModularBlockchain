package core

import "io"

type Transactions struct {
}

func (tx *Transactions) DecodeBinary(r io.Reader) error { return nil }

func (tx *Transactions) EncodeBinary(w io.Writer) error { return nil }
