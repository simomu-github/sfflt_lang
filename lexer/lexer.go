package lexer

import "github.com/simomu-github/sfflt_lang/token"

type Lexer struct {
	source string

	start   int
	current int
	line    int
	column  int
}

func New(source string) *Lexer {
	return &Lexer{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
		column:  0,
	}
}

func (l *Lexer) ScanToken() token.Token {
	l.skipWhitespace()

	l.start = l.current
	char := l.readChar()

	switch char {
	case '+':
		return l.makeToken(token.PLUS, string(char))
	case '-':
		return l.makeToken(token.MINUS, string(char))
	case '*':
		return l.makeToken(token.ASTERISK, string(char))
	case '/':
		return l.makeToken(token.SLASH, string(char))

	case '(':
		return l.makeToken(token.LPAREN, string(char))
	case ')':
		return l.makeToken(token.RPAREN, string(char))
	case '{':
		return l.makeToken(token.LBRACE, string(char))
	case '}':
		return l.makeToken(token.RBRACE, string(char))
	case ';':
		return l.makeToken(token.SEMICOLON, string(char))

	case '=':
		if l.peekChar() == '=' {
			nextChar := l.readChar()
			return l.makeToken(token.EQ, string(char)+string(nextChar))
		} else {
			return l.makeToken(token.ASSIGN, string(char))
		}
	case '!':
		if l.peekChar() == '=' {
			nextChar := l.readChar()
			return l.makeToken(token.NOT_EQ, string(char)+string(nextChar))
		} else {
			return l.makeToken(token.BANG, string(char))
		}
	case '<':
		if l.peekChar() == '=' {
			nextChar := l.readChar()
			return l.makeToken(token.LTEQ, string(char)+string(nextChar))
		} else {
			return l.makeToken(token.LT, string(char))
		}
	case '>':
		if l.peekChar() == '=' {
			nextChar := l.readChar()
			return l.makeToken(token.GTEQ, string(char)+string(nextChar))
		} else {
			return l.makeToken(token.GT, string(char))
		}
	case '\'':
		return l.scanChar()
	case 0:
		return l.makeToken(token.EOF, string(char))
	default:
		if isDigit(char) {
			return l.scanNumber()
		} else if isLetter(char) {
			return l.scanIdentifier()
		}
	}

	return l.makeToken(token.ILLEGAL, string(char))
}

func (l *Lexer) readChar() byte {
	if l.current >= len(l.source) {
		return 0
	}

	char := l.source[l.current]
	l.current += 1
	l.column += 1
	return char
}

func (l *Lexer) peekChar() byte {
	if l.current >= len(l.source) {
		return 0
	}

	return l.source[l.current]
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.peekChar() {
		case ' ':
			l.readChar()
		case '\r':
			l.readChar()
		case '\t':
			l.readChar()
		case '\n':
			l.readChar()
			l.line += 1
			l.column = 0
		default:
			return
		}
	}
}

func (l *Lexer) scanChar() token.Token {

	char := l.readChar()

	if l.peekChar() != '\'' || l.isAtEnd() {
		return l.makeToken(token.ILLEGAL, "Unterminated char.")
	}

	l.readChar()
	return l.makeToken(token.CHAR, string(char))
}

func (l *Lexer) scanNumber() token.Token {
	for isDigit(l.peekChar()) {
		l.readChar()
	}

	return l.makeToken(token.INT, l.source[l.start:l.current])
}

func (l *Lexer) scanIdentifier() token.Token {
	for isLetter(l.peekChar()) {
		l.readChar()
	}

	identifier := l.source[l.start:l.current]
	return l.makeToken(token.LookupIdent(identifier), identifier)
}

func (l *Lexer) makeToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
