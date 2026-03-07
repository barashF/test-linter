package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

func isLowerCase(msgLit *ast.BasicLit, pass *analysis.Pass, cfg *Config) {
	if slices.Contains(cfg.DisableRules, "lowercase") {
		return
	}

	msg := strings.Trim(msgLit.Value, `"`)
	firstRune, size := utf8.DecodeRuneInString(msg)
	if !unicode.IsUpper(firstRune) {
		return
	}

	lowerMsg := string(unicode.ToLower(firstRune)) + msg[size:]

	start := msgLit.Pos() + 1
	end := msgLit.End() - 1

	pass.Report(analysis.Diagnostic{
		Pos:     msgLit.Pos(),
		End:     msgLit.End(),
		Message: "log message must start with lowercase letter",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Convert first letter to lowercase",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     start,
						End:     end,
						NewText: []byte(lowerMsg),
					},
				},
			},
		},
	})
}

func hasCyrillic(msg string, pos token.Pos, pass *analysis.Pass, cfg *Config) {
	if slices.Contains(cfg.DisableRules, "cyrillic") {
		return
	}
	for _, r := range msg {
		if unicode.Is(unicode.Cyrillic, r) {
			pass.Reportf(pos, "log message must be in English (found Cyrillic: %q)", r)
			return
		}
	}
}

func checkSpecialChars(msgLit *ast.BasicLit, msg string, pass *analysis.Pass, cfg *Config) {
	if slices.Contains(cfg.DisableRules, "special-chars") {
		return
	}

	valueStart := msgLit.ValuePos

	for i, r := range msg {
		if isEmojiOrSymbol(r) {
			runeLen := utf8.RuneLen(r)
			charStart := valueStart + token.Pos(i) + 1
			charEnd := charStart + token.Pos(runeLen) - 2

			pass.Report(analysis.Diagnostic{
				Pos:     charStart,
				End:     charEnd,
				Message: fmt.Sprintf("log message should not contain emoji or symbol: %q", r),
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Remove emoji or symbol",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     charStart,
								End:     charEnd,
								NewText: []byte(""),
							},
						},
					},
				},
			})
			return
		}

		if r < 128 && strings.ContainsRune("!@#$%^&*?\\|~`", r) {
			charStart := valueStart + token.Pos(i) + 1
			charEnd := charStart + 1

			pass.Report(analysis.Diagnostic{
				Pos:     charStart,
				End:     charEnd,
				Message: fmt.Sprintf("log message should not contain special character: %q", r),
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Remove special character",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     charStart,
								End:     charEnd,
								NewText: []byte(""),
							},
						},
					},
				},
			})
			return
		}
	}

	if idx := strings.Index(msg, "..."); idx != -1 {
		ellipsisStart := valueStart + token.Pos(idx)
		ellipsisEnd := ellipsisStart + 4

		pass.Report(analysis.Diagnostic{
			Pos:     ellipsisStart,
			End:     ellipsisEnd,
			Message: "log message should not contain ellipsis '...'",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Remove ellipsis",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     ellipsisStart,
							End:     ellipsisEnd,
							NewText: []byte(""),
						},
					},
				},
			},
		})
	}
}

var defaultSensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"token", "api_key", "apikey", "secret",
	"credential", "private_key", "auth_token",
}

func checkSensitiveData(call *ast.CallExpr, pass *analysis.Pass, cfg *Config) {
	if slices.Contains(cfg.DisableRules, "sensitive-keywords") {
		return
	}
	if len(call.Args) == 0 {
		return
	}

	arg := call.Args[0]

	keywords := defaultSensitiveKeywords
	if len(cfg.SensitiveKeywords) > 0 {
		keywords = cfg.SensitiveKeywords
	}
	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg := strings.Trim(lit.Value, `"`)
		checkKeyword(msg, lit.Pos(), pass, keywords)
		return
	}

	if bin, ok := arg.(*ast.BinaryExpr); ok && bin.Op == token.ADD {
		if leftLit, ok := bin.X.(*ast.BasicLit); ok && leftLit.Kind == token.STRING {
			leftMsg := strings.Trim(leftLit.Value, `"`)
			lowerLeft := strings.ToLower(leftMsg)

			for _, kw := range keywords {
				if strings.Contains(lowerLeft, kw) &&
					(strings.HasSuffix(leftMsg, "=") || strings.HasSuffix(leftMsg, ":") || strings.HasSuffix(leftMsg, ": ")) {
					pass.Reportf(leftLit.ValuePos, "log message may contain sensitive keyword %q", kw)
					return
				}
			}
		}
	}
}

func checkKeyword(msg string, pos token.Pos, pass *analysis.Pass, keywords []string) {
	lowerMsg := strings.ToLower(msg)
	for _, kw := range keywords {
		if strings.Contains(lowerMsg, kw) {
			pass.Reportf(pos, "log message may contain sensitive keyword %q", kw)
			return
		}
	}
}

func isEmojiOrSymbol(r rune) bool {
	return unicode.Is(unicode.So, r) ||
		unicode.Is(unicode.Sk, r) ||
		unicode.Is(unicode.Sc, r)
}
