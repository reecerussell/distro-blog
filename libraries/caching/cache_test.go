package caching

import (
	"encoding/json"
	"os"
	"testing"
)

var host = os.Getenv("CACHE_HOST")

func TestNew(t *testing.T) {
	_, err := New(host)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	t.Run("Invalid Host", func(t *testing.T) {
		_, err := New("invalid hostname")
		if err == nil {
			t.Errorf("expected to fail")
		}
	})
}

func TestClient_Set(t *testing.T) {
	c, _ := New(host)
	v := map[string]string{
		"Hello": "World",
	}
	bytes, _ := json.Marshal(v)

	err := c.Set("ClientSet_Object", bytes)
	if err != nil {
		t.Errorf("unexpected fail: %v", err)
	}

	t.Run("Uninitiated Client", func(t *testing.T) {
		var c *client
		err := c.Set("UnsetClient", nil)
		if err == nil {
			t.Errorf("expected to fail")
		}
	})
}

func TestClient_Get(t *testing.T) {
	c, _ := New(host)
	key := "ClientGet_IteM"
	value := map[string]string{
		"Hello": "World",
	}
	bytes, _ := json.Marshal(value)

	err := c.Set(key, bytes)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	bytes, ok := c.Get(key)
	if !ok {
		t.Errorf("expected to be ok")
		return
	}

	var out map[string]interface{}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	for k, v := range value {
		if ov, ok := out[k]; !ok {
			t.Errorf("expected output to contain key '%s'", k)
		} else if ov != v {
			t.Errorf("output[%s] expected to be '%s' but was '%s'", k, v, ov)
		}
	}

	t.Run("Un-stored Item", func(t *testing.T) {
		_, ok := c.Get("UnStoredItem")
		if ok {
			t.Errorf("expected to fail")
		}
	})

	t.Run("Uninitiated Client", func(t *testing.T) {
		var c *client
		_, ok := c.Get("UnsetClient")
		if ok {
			t.Errorf("expected to fail")
		}
	})
}