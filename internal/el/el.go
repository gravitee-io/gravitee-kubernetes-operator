// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package el

import "fmt"

type Expression string

var False = Expression("false")

func Empty() Expression {
	return Expression("")
}

func (c Expression) Parenthesized() Expression {
	return "(" + c + ")"
}

func (c Expression) Closed() Expression {
	return "{" + c + "}"
}

func (c Expression) String() string {
	return string(c)
}

func (c Expression) And(other Expression) Expression {
	if c.IsEmpty() {
		return other
	}
	return c + " and " + other
}

func (c Expression) Or(other Expression) Expression {
	if c.IsEmpty() {
		return other
	}
	return c + " or " + other
}

func (c Expression) Format(args ...interface{}) Expression {
	return Expression(fmt.Sprintf(string(c), args...))
}

func (c Expression) Negated() Expression {
	return c.Equals(False).Parenthesized()
}

func (c Expression) Equals(other Expression) Expression {
	return c + " eq " + other
}

func (c Expression) Matches(other Expression) Expression {
	return c + " matches " + other
}

func (c Expression) IsEmpty() bool {
	return c == ""
}
