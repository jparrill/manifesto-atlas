package preview

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Chrome's debugging pipe exchanges NUL-delimited CDP JSON on inherited file
// descriptors. Device emulation needs no generated JavaScript or WebSocket SDK.
type browserPipe struct {
	in      io.Writer
	out     *bufio.Reader
	next    int
	session string
	events  map[string]bool
}

func (p *browserPipe) read() (map[string]json.RawMessage, error) {
	raw, err := p.out.ReadBytes(0)
	if err != nil {
		return nil, err
	}
	var message map[string]json.RawMessage
	if err := json.Unmarshal(raw[:len(raw)-1], &message); err != nil {
		return nil, err
	}
	if method, ok := message["method"]; ok {
		var name string
		_ = json.Unmarshal(method, &name)
		p.events[name] = true
	}
	return message, nil
}

func (p *browserPipe) call(method string, params any, result any) error {
	p.next++
	request := map[string]any{"id": p.next, "method": method, "params": params}
	if p.session != "" {
		request["sessionId"] = p.session
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return err
	}
	if _, err := p.in.Write(append(raw, 0)); err != nil {
		return err
	}
	for {
		message, err := p.read()
		if err != nil {
			return err
		}
		var id int
		_ = json.Unmarshal(message["id"], &id)
		if id != p.next {
			continue
		}
		if detail, ok := message["error"]; ok {
			return fmt.Errorf("%s: %s", method, detail)
		}
		if result != nil {
			return json.Unmarshal(message["result"], result)
		}
		return nil
	}
}

func (p *browserPipe) loaded() error {
	for !p.events["Page.loadEventFired"] {
		if _, err := p.read(); err != nil {
			return err
		}
	}
	return nil
}

func (p *browserPipe) page(target string, width, height int, mobile bool) (string, string, error) {
	var created struct {
		TargetID string `json:"targetId"`
	}
	if err := p.call("Target.createTarget", map[string]any{"url": "about:blank"}, &created); err != nil {
		return "", "", err
	}
	var attached struct {
		SessionID string `json:"sessionId"`
	}
	if err := p.call("Target.attachToTarget", map[string]any{"targetId": created.TargetID, "flatten": true}, &attached); err != nil {
		return "", "", err
	}
	p.session = attached.SessionID
	p.events = map[string]bool{}
	if err := p.call("Emulation.setDeviceMetricsOverride", map[string]any{"width": width, "height": height, "deviceScaleFactor": 1, "mobile": mobile}, nil); err != nil {
		return "", "", err
	}
	if err := p.call("Page.enable", map[string]any{}, nil); err != nil {
		return "", "", err
	}
	if err := p.call("Page.navigate", map[string]any{"url": target}, nil); err != nil {
		return "", "", err
	}
	if err := p.loaded(); err != nil {
		return "", "", err
	}
	var root struct {
		Root struct {
			NodeID int `json:"nodeId"`
		} `json:"root"`
	}
	if err := p.call("DOM.getDocument", map[string]any{"depth": 0}, &root); err != nil {
		return "", "", err
	}
	var document struct {
		OuterHTML string `json:"outerHTML"`
	}
	if err := p.call("DOM.getOuterHTML", map[string]any{"nodeId": root.Root.NodeID}, &document); err != nil {
		return "", "", err
	}
	var metrics struct {
		Content struct {
			Width float64 `json:"width"`
		} `json:"cssContentSize"`
	}
	if err := p.call("Page.getLayoutMetrics", map[string]any{}, &metrics); err != nil {
		return "", "", err
	}
	if metrics.Content.Width > float64(width)+1 {
		return "", "", errors.New("actual viewport has horizontal overflow")
	}
	var image struct {
		Data string `json:"data"`
	}
	if err := p.call("Page.captureScreenshot", map[string]any{"format": "png", "captureBeyondViewport": false}, &image); err != nil {
		return "", "", err
	}
	return document.OuterHTML, image.Data, nil
}
