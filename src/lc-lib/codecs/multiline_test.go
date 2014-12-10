package codecs

import (
  "lc-lib/core"
  "sync"
  "testing"
  "time"
)

var multiline_t *testing.T
var multiline_lines int
var multiline_lock sync.Mutex

func createMultilineCodec(unused map[string]interface{}, callback core.CodecCallbackFunc, t *testing.T) core.Codec {
  config := core.NewConfig()
  config.General.MaxLineBytes = 1048576
  config.General.SpoolMaxBytes = 10485760

  factory, err := NewMultilineCodecFactory(config, "", unused, "multiline")
  if err != nil {
    t.Logf("Failed to create multiline codec: %s", err)
    t.FailNow()
  }

  return factory.NewCodec(callback, 0)
}

func checkMultiline(start_offset int64, end_offset int64, text string) {
  multiline_lock.Lock()
  defer multiline_lock.Unlock()
  multiline_lines++

  if multiline_lines == 1 {
    if text != "DEBUG First line\nNEXT line\nANOTHER line" {
      multiline_t.Logf("Event data incorrect [% X]", text)
      multiline_t.FailNow()
    }

    if start_offset != 0 {
      multiline_t.Logf("Event start offset is incorrect [%d]", start_offset)
      multiline_t.FailNow()
    }

    if end_offset != 5 {
      multiline_t.Logf("Event end offset is incorrect [%d]", end_offset)
      multiline_t.FailNow()
    }

    return
  }

  if text != "DEBUG Next line" {
    multiline_t.Logf("Event data incorrect [% X]", text)
    multiline_t.FailNow()
  }

  if start_offset != 6 {
    multiline_t.Logf("Event start offset is incorrect [%d]", start_offset)
    multiline_t.FailNow()
  }

  if end_offset != 7 {
    multiline_t.Logf("Event end offset is incorrect [%d]", end_offset)
    multiline_t.FailNow()
  }
}

func TestMultilinePrevious(t *testing.T) {
  multiline_t = t
  multiline_lines = 0

  codec := createMultilineCodec(map[string]interface{}{
    "pattern": "^(ANOTHER|NEXT) ",
    "what": "previous",
    "negate": false,
  }, checkMultiline, t)

  // Send some data
  codec.Event(0, 1, "DEBUG First line")
  codec.Event(2, 3, "NEXT line")
  codec.Event(4, 5, "ANOTHER line")
  codec.Event(6, 7, "DEBUG Next line")

  if multiline_lines != 1 {
    t.Logf("Wrong line count received")
    t.FailNow()
  }
}

func TestMultilinePreviousNegate(t *testing.T) {
  multiline_t = t
  multiline_lines = 0

  codec := createMultilineCodec(map[string]interface{}{
    "pattern": "^DEBUG ",
    "what": "previous",
    "negate": true,
  }, checkMultiline, t)

  // Send some data
  codec.Event(0, 1, "DEBUG First line")
  codec.Event(2, 3, "NEXT line")
  codec.Event(4, 5, "ANOTHER line")
  codec.Event(6, 7, "DEBUG Next line")

  if multiline_lines != 1 {
    t.Logf("Wrong line count received")
    t.FailNow()
  }
}

func TestMultilinePreviousTimeout(t *testing.T) {
  multiline_t = t
  multiline_lines = 0

  codec := createMultilineCodec(map[string]interface{}{
    "pattern": "^(ANOTHER|NEXT) ",
    "what": "previous",
    "negate": false,
    "previous timeout": "5s",
  }, checkMultiline, t)

  // Send some data
  codec.Event(0, 1, "DEBUG First line")
  codec.Event(2, 3, "NEXT line")
  codec.Event(4, 5, "ANOTHER line")
  codec.Event(6, 7, "DEBUG Next line")

  // Allow 3 seconds
  time.Sleep(3 * time.Second)

  multiline_lock.Lock()
  if multiline_lines != 1 {
    t.Logf("Timeout triggered too early")
    t.FailNow()
  }
  multiline_lock.Unlock()

  // Allow 7 seconds
  time.Sleep(7 * time.Second)

  multiline_lock.Lock()
  if multiline_lines != 2 {
    t.Logf("Wrong line count received")
    t.FailNow()
  }
  multiline_lock.Unlock()

  codec.Teardown()
}

func TestMultilineNext(t *testing.T) {
  multiline_t = t
  multiline_lines = 0

  codec := createMultilineCodec(map[string]interface{}{
    "pattern": "^(DEBUG|NEXT) ",
    "what": "next",
    "negate": false,
  }, checkMultiline, t)

  // Send some data
  codec.Event(0, 1, "DEBUG First line")
  codec.Event(2, 3, "NEXT line")
  codec.Event(4, 5, "ANOTHER line")
  codec.Event(6, 7, "DEBUG Next line")

  if multiline_lines != 1 {
    t.Logf("Wrong line count received")
    t.FailNow()
  }
}

func TestMultilineNextNegate(t *testing.T) {
  multiline_t = t
  multiline_lines = 0

  codec := createMultilineCodec(map[string]interface{}{
    "pattern": "^ANOTHER ",
    "what": "next",
    "negate": true,
  }, checkMultiline, t)

  // Send some data
  codec.Event(0, 1, "DEBUG First line")
  codec.Event(2, 3, "NEXT line")
  codec.Event(4, 5, "ANOTHER line")
  codec.Event(6, 7, "DEBUG Next line")

  if multiline_lines != 1 {
    t.Logf("Wrong line count received")
    t.FailNow()
  }
}

func checkMultilineMaxBytes(start_offset int64, end_offset int64, text string) {
  multiline_lines++

  if multiline_lines == 1 {
    if text != "DEBUG First line\nsecond line\nthi" {
      multiline_t.Logf("Event data incorrect [% X]", text)
      multiline_t.FailNow()
    }

    return
  }

  if text != "rd line" {
    multiline_t.Logf("Second event data incorrect [% X]", text)
    multiline_t.FailNow()
  }
}

func TestMultilineMaxBytes(t *testing.T) {
  multiline_t = t
  multiline_lines = 0

  codec := createMultilineCodec(map[string]interface{}{
    "max multiline bytes": int64(32),
    "pattern": "^DEBUG ",
    "negate": true,
  }, checkMultilineMaxBytes, t)

  // Send some data
  codec.Event(0, 1, "DEBUG First line")
  codec.Event(2, 3, "second line")
  codec.Event(4, 5, "third line")
  codec.Event(6, 7, "DEBUG Next line")

  if multiline_lines != 2 {
    t.Logf("Wrong line count received")
    t.FailNow()
  }
}
