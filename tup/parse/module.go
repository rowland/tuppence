package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// module = { top_level_item } .

func Module(source *source.Source, module *ast.Module) (mod *ast.Module, err error) {
	var tokens []tok.Token
	if tokens, err = tok.Tokenize(source.Contents, source.Filename); err != nil {
		return nil, err
	}
	remainder := tokens
	for {
		var item ast.TopLevelItem
		if item, remainder, err = TopLevelItem(remainder); err != nil {
			return nil, err
		} else if item == nil {
			break
		}
		module.AddTopLevelItem(item)
	}
	return module, nil
}
