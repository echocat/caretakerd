package config

import (
    "testing"
    "github.com/stretchr/testify/assert"
    . "github.com/echocat/caretakerd/values"
)

func TestEnvironment_parseCmd(t *testing.T) {
    assert.Equal(t, []String{"a","b"}, parseCmd("a b"))
    assert.Equal(t, []String{"a b"}, parseCmd("\"a b\""))
    assert.Equal(t, []String{"\"a","b"}, parseCmd("\\\"a b"))
    assert.Equal(t, []String{"\"a b"}, parseCmd("\"\\\"a b\""))
    assert.Equal(t, []String{"\\"}, parseCmd("\\\\"))
    assert.Equal(t, []String{"v1=a ","v2= b"}, parseCmd("v1=\"a \" \"v2= b\""))
}

