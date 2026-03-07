package analyzer

import (
	"flag"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type Config struct {
	SensitiveKeywords []string
	DisableRules      []string
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	if value == "" {
		*s = []string{}
		return nil
	}
	parts := strings.Split(value, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	*s = parts
	return nil
}

func New(cfg *Config) *analysis.Analyzer {
	if cfg == nil {
		cfg = &Config{}
	}

	var flags flag.FlagSet
	flags.Var((*stringSlice)(&cfg.SensitiveKeywords), "sensitive-keywords",
		"Comma-separated list of sensitive keywords")
	flags.Var((*stringSlice)(&cfg.DisableRules), "disable-rules",
		"Comma-separated list of rules to disable")
	return &analysis.Analyzer{
		Name:             "loglint",
		Doc:              "checking logs",
		Flags:            flags,
		RunDespiteErrors: true,
		Run:              run(cfg),
	}
}

var Analyzer = New(nil)

func run(cfg *Config) func(*analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		for _, file := range pass.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)

				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				if !isLogMethod(sel.Sel.Name) {
					return true
				}

				if !isSlogCall(sel, pass.TypesInfo) && !isZapLogger(sel, pass.TypesInfo) {
					return true
				}

				if len(call.Args) == 0 {
					return true
				}

				checkSensitiveData(call, pass, cfg)
				msgLit, ok := call.Args[0].(*ast.BasicLit)
				if !ok || msgLit.Kind != token.STRING {
					return true
				}
				msg := strings.Trim(msgLit.Value, `"`)
				isLowerCase(msgLit, pass, cfg)
				hasCyrillic(msg, call.Pos(), pass, cfg)
				checkSpecialChars(msgLit, msg, pass, cfg)

				return true
			})
		}
		return nil, nil
	}
}
