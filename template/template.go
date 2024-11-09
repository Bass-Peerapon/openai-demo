package template

import _ "embed"

//go:embed order/order.tmpl
var Order []byte

//go:embed metaprompt.tmpl
var Metaprompt []byte

//go:embed question.tmpl
var Question []byte
