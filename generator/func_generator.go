package generator

import (
	"github.com/moznion/gowrtr/internal/errmsg"
)

// FuncGenerator represents a code generator for the func.
type FuncGenerator struct {
	FuncReceiver  *FuncReceiverGenerator
	FuncSignature *FuncSignatureGenerator
	Statements    []StatementGenerator
}

// NewFuncGenerator returns a new `FuncGenerator`.
func NewFuncGenerator(receiver *FuncReceiverGenerator, signature *FuncSignatureGenerator, statements ...StatementGenerator) *FuncGenerator {
	return &FuncGenerator{
		FuncReceiver:  receiver,
		FuncSignature: signature,
		Statements:    statements,
	}
}

// AddStatements adds statements for the func to `FuncGenerator`.
// This method returns a *new* `FuncGenerator`; it means this method acts as immutable.
func (fg *FuncGenerator) AddStatements(statements ...StatementGenerator) *FuncGenerator {
	return &FuncGenerator{
		FuncReceiver:  fg.FuncReceiver,
		FuncSignature: fg.FuncSignature,
		Statements:    append(fg.Statements, statements...),
	}
}

// Generate generates a func block as golang code.
func (fg *FuncGenerator) Generate(indentLevel int) (string, error) {
	indent := buildIndent(indentLevel)

	stmt := indent + "func "

	receiver := ""
	if fg.FuncReceiver != nil {
		var err error
		receiver, err = fg.FuncReceiver.Generate(0)
		if err != nil {
			return "", err
		}
	}
	if receiver != "" {
		stmt += receiver + " "
	}

	if fg.FuncSignature == nil {
		return "", errmsg.FuncSignatureIsNilError()
	}
	sig, err := fg.FuncSignature.Generate(0)
	if err != nil {
		return "", err
	}
	stmt += sig + " {\n"

	nextIndentLevel := indentLevel + 1
	for _, c := range fg.Statements {
		gen, err := c.Generate(nextIndentLevel)
		if err != nil {
			return "", err
		}
		stmt += gen
	}

	stmt += indent + "}\n"

	return stmt, nil
}
