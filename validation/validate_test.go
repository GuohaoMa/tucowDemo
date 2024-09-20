package validation

import (
	"os"
	"testing"
)

func createTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile
}

func TestValidate_ValidXML(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<id>1</id>
		<name>Test Graph</name>
		<nodes>
			<node>
				<id>1</id>
				<name>Node1</name>
			</node>
			<node>
				<id>2</id>
				<name>Node2</name>
			</node>
		</nodes>
		<edges>
			<node>
				<from>1</from>
				<to>2</to>
				<cost>10</cost>
			</node>
		</edges>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidate_MissingNodes(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<id>1</id>
		<name>Test Graph</name>
		<nodes></nodes>
		<edges>
			<node>
				<from>1</from>
				<to>2</to>
				<cost>10</cost>
			</node>
		</edges>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err == nil || err.Error() != "There must be at least one <node> in the <nodes> group" {
		t.Errorf("Expected error about missing nodes, got %v", err)
	}
}

func TestValidate_MissingGraphID(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<name>Test Graph</name>
		<nodes>
			<node>
				<id>1</id>
				<name>Node1</name>
			</node>
		</nodes>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err == nil || err.Error() != "There must be an <id> in the <graph>" {
		t.Errorf("Expected error about missing graph ID, got %v", err)
	}
}

func TestValidate_MissingGraphName(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<id>1</id>
		<nodes>
			<node>
				<id>1</id>
				<name>Node1</name>
			</node>
		</nodes>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err == nil || err.Error() != "There must be an <name> in the <graph>" {
		t.Errorf("Expected error about missing graph name, got %v", err)
	}
}

func TestValidate_DuplicateNodeID(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<id>1</id>
		<name>Test Graph</name>
		<nodes>
			<node>
				<id>1</id>
				<name>Node1</name>
			</node>
			<node>
				<id>1</id>
				<name>Node2</name>
			</node>
		</nodes>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err == nil || err.Error() != "All nodes must have different <id> tags." {
		t.Errorf("Expected error about duplicate node IDs, got %v", err)
	}
}

func TestValidate_NegativeEdgeCost(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<id>1</id>
		<name>Test Graph</name>
		<nodes>
			<node>
				<id>1</id>
				<name>Node1</name>
			</node>
			<node>
				<id>2</id>
				<name>Node2</name>
			</node>
		</nodes>
		<edges>
			<node>
				<from>1</from>
				<to>2</to>
				<cost>-10</cost>
			</node>
		</edges>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err == nil || err.Error() != "Cost of an edge must be non-negative." {
		t.Errorf("Expected error about negative edge cost, got %v", err)
	}
}

func TestValidate_MissingEdgeFrom(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<id>1</id>
		<name>Test Graph</name>
		<nodes>
			<node>
				<id>1</id>
				<name>Node1</name>
			</node>
			<node>
				<id>2</id>
				<name>Node2</name>
			</node>
		</nodes>
		<edges>
			<node>
				<to>2</to>
				<cost>10</cost>
			</node>
		</edges>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err == nil || err.Error() != "For every <edge>, there must be a single <from> tag" {
		t.Errorf("Expected error about missing <from> tag, got %v", err)
	}
}

func TestValidate_MissingEdgeTo(t *testing.T) {
	t.Parallel()
	xmlContent := `
	<graph>
		<id>1</id>
		<name>Test Graph</name>
		<nodes>
			<node>
				<id>1</id>
				<name>Node1</name>
			</node>
			<node>
				<id>2</id>
				<name>Node2</name>
			</node>
		</nodes>
		<edges>
			<node>
				<from>1</from>
				<cost>10</cost>
			</node>
		</edges>
	</graph>`

	tmpfile := createTempFile(t, xmlContent)
	defer os.Remove(tmpfile.Name())

	err := Validate(tmpfile.Name())
	if err == nil || err.Error() != "For every <edge>, there must be a single <to> tag" {
		t.Errorf("Expected error about missing <to> tag, got %v", err)
	}
}
