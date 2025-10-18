package field

import (
	"testing"
)

func TestParseDirective_Text(t *testing.T) {
	line := `# @text label="Username" required`
	info := ParseDirective(line)
	
	if info == nil {
		t.Fatal("Expected info, got nil")
	}
	if info.Type != TextField {
		t.Errorf("Expected TextField, got %v", info.Type)
	}
	if info.Label != "Username" {
		t.Errorf("Expected label 'Username', got %s", info.Label)
	}
	if !info.Required {
		t.Error("Expected required to be true")
	}
}

func TestParseDirective_Hidden(t *testing.T) {
	line := `# @hidden note="Internal value"`
	info := ParseDirective(line)
	
	if info == nil {
		t.Fatal("Expected info, got nil")
	}
	if info.Type != HiddenField {
		t.Errorf("Expected HiddenField, got %v", info.Type)
	}
	if info.Note != "Internal value" {
		t.Errorf("Expected note 'Internal value', got %s", info.Note)
	}
}

func TestParseDirective_Number(t *testing.T) {
	line := `# @number label="Port" required readonly`
	info := ParseDirective(line)
	
	if info == nil {
		t.Fatal("Expected info, got nil")
	}
	if info.Type != NumberField {
		t.Errorf("Expected NumberField, got %v", info.Type)
	}
	if !info.Required || !info.Readonly {
		t.Error("Expected required and readonly to be true")
	}
}

func TestParseDirective_Select(t *testing.T) {
	line := `# @select label="Environment" options=dev,staging,prod`
	info := ParseDirective(line)
	
	if info == nil {
		t.Fatal("Expected info, got nil")
	}
	if info.Type != SelectField {
		t.Errorf("Expected SelectField, got %v", info.Type)
	}
	if len(info.Options) != 3 {
		t.Errorf("Expected 3 options, got %d", len(info.Options))
	}
}

func TestParseDirective_Script(t *testing.T) {
	line := `# @text script="echo hello"`
	info := ParseDirective(line)
	
	if info == nil {
		t.Fatal("Expected info, got nil")
	}
	if info.Script != "echo hello" {
		t.Errorf("Expected script 'echo hello', got %s", info.Script)
	}
}

func TestParseDirective_NoDirective(t *testing.T) {
	line := `# Just a comment`
	info := ParseDirective(line)
	
	if info != nil {
		t.Errorf("Expected nil for non-directive comment, got %v", info)
	}
}
