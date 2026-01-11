package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// module = { top_level_item } .

func Module(source *source.Source, module *ast.Module) (*ast.Module, error) {
	tokens, err := tok.Tokenize(source.Contents, source.Filename)
	if err != nil {
		return nil, err
	}
	for {
		item, remainder, err := TopLevelItem(tokens)
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		module.AddTopLevelItem(item)
		tokens = remainder
	}
	return module, nil
}
