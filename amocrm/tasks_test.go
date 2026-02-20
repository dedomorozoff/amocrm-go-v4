package amocrm

import (
	"encoding/json"
	"testing"
)

func TestTaskResult_UnmarshalJSON_Object(t *testing.T) {
	// Test unmarshaling from object format
	jsonData := `{"text": "Task completed successfully"}`

	var result TaskResult
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal object format: %v", err)
	}

	expected := "Task completed successfully"
	if result.Text != expected {
		t.Errorf("Expected text '%s', got '%s'", expected, result.Text)
	}
}

func TestTaskResult_UnmarshalJSON_Array(t *testing.T) {
	// Test unmarshaling from array format
	jsonData := `[{"text": "Task completed successfully"}]`

	var result TaskResult
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal array format: %v", err)
	}

	expected := "Task completed successfully"
	if result.Text != expected {
		t.Errorf("Expected text '%s', got '%s'", expected, result.Text)
	}
}

func TestTaskResult_UnmarshalJSON_EmptyArray(t *testing.T) {
	// Test unmarshaling from empty array (API sometimes returns this)
	jsonData := `[]`

	var result TaskResult
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal empty array: %v", err)
	}

	if result.Text != "" {
		t.Errorf("Expected empty text, got '%s'", result.Text)
	}
}

func TestTask_UnmarshalWithEmptyArrayResult(t *testing.T) {
	// Test full Task unmarshaling with empty array result (real-world scenario)
	jsonData := `{
		"id": 789,
		"task_type_id": 3,
		"text": "Send email",
		"is_completed": false,
		"result": []
	}`

	var task Task
	err := json.Unmarshal([]byte(jsonData), &task)
	if err != nil {
		t.Errorf("Failed to unmarshal Task with empty array result: %v", err)
	}

	if task.Result.Text != "" {
		t.Errorf("Expected empty result text, got '%s'", task.Result.Text)
	}
}

func TestTaskResult_UnmarshalJSON_EmptyObject(t *testing.T) {
	// Test unmarshaling from empty object
	jsonData := `{}`

	var result TaskResult
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal empty object: %v", err)
	}

	if result.Text != "" {
		t.Errorf("Expected empty text, got '%s'", result.Text)
	}
}

func TestTaskResult_MarshalJSON(t *testing.T) {
	// Test marshaling to object format
	result := TaskResult{Text: "Test result"}

	data, err := json.Marshal(result)
	if err != nil {
		t.Errorf("Failed to marshal TaskResult: %v", err)
	}

	expected := `{"text":"Test result"}`
	if string(data) != expected {
		t.Errorf("Expected JSON '%s', got '%s'", expected, string(data))
	}
}

func TestTask_UnmarshalWithArrayResult(t *testing.T) {
	// Test full Task unmarshaling with array result format (real-world scenario)
	jsonData := `{
		"id": 123,
		"task_type_id": 1,
		"text": "Call the client",
		"is_completed": true,
		"result": [{"text": "Client called, will call back tomorrow"}]
	}`

	var task Task
	err := json.Unmarshal([]byte(jsonData), &task)
	if err != nil {
		t.Errorf("Failed to unmarshal Task with array result: %v", err)
	}

	expectedText := "Client called, will call back tomorrow"
	if task.Result.Text != expectedText {
		t.Errorf("Expected result text '%s', got '%s'", expectedText, task.Result.Text)
	}
}

func TestTask_UnmarshalWithObjectResult(t *testing.T) {
	// Test full Task unmarshaling with object result format
	jsonData := `{
		"id": 456,
		"task_type_id": 2,
		"text": "Meeting scheduled",
		"is_completed": true,
		"result": {"text": "Meeting completed successfully"}
	}`

	var task Task
	err := json.Unmarshal([]byte(jsonData), &task)
	if err != nil {
		t.Errorf("Failed to unmarshal Task with object result: %v", err)
	}

	expectedText := "Meeting completed successfully"
	if task.Result.Text != expectedText {
		t.Errorf("Expected result text '%s', got '%s'", expectedText, task.Result.Text)
	}
}
