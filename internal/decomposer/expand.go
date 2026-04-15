package decomposer

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	openai "shotgun.dev/odek/openai"
)

// ExpandStreaming runs level-by-level BFS expansion of the session's root
// decomposition up to cfg.MaxDepth, emitting events on a buffered channel.
// The channel is closed after EventDone (or EventCancelled + EventDone on
// ctx cancellation, or EventCapReached + EventDone on the rune-cap path).
// The caller must range over the channel until it closes.
//
// Invariant: the decomposer updates session state in-line via sess.Apply,
// so consumers can read sess.Snapshot() after receiving each event and see
// consistent state without needing to Apply themselves.
func (d *Decomposer) ExpandStreaming(ctx context.Context, sess *Session, cfg Config) <-chan ExpansionEvent {
	ch := make(chan ExpansionEvent, 16)
	ctx, cancel := context.WithCancel(ctx)
	sess.setEvents(ch, cancel)

	go func() {
		defer close(ch)
		defer sess.clearEvents()
		d.runExpansion(ctx, sess, cfg, ch)
	}()

	return ch
}

// runExpansion drives the level-by-level BFS. The producer goroutine owns
// all writes to the session's state: emit calls sess.Apply before pushing
// the event onto ch, so consumers can treat events as pure notifications
// and read sess.Snapshot() without worrying about who mutates.
func (d *Decomposer) runExpansion(ctx context.Context, sess *Session, cfg Config, ch chan<- ExpansionEvent) {
	emit := func(evt ExpansionEvent) {
		sess.Apply(evt)
		select {
		case ch <- evt:
		case <-ctx.Done():
		}
	}
	var emitMu sync.Mutex
	safeEmit := func(evt ExpansionEvent) {
		emitMu.Lock()
		defer emitMu.Unlock()
		emit(evt)
	}

	rootResp := sess.Root.Response
	initial := collectRunesForExpansion(rootResp)
	for i := range initial {
		initial[i].Depth = 1
		initial[i].ParentDecomposition = sess.Root
	}

	totalRunes := countTotalRunes(rootResp)
	totalDecomps := 1
	maxDepthSeen := 0
	visited := map[string]bool{"root": true}

	finish := func(terminal ExpansionEvent) {
		if terminal != nil {
			safeEmit(terminal)
		}
		safeEmit(EventDone{
			TotalDecompositions: totalDecomps,
			TotalRunes:          totalRunes,
			MaxDepth:            maxDepthSeen,
		})
	}

	currentLevel := initial
	for len(currentLevel) > 0 {
		if ctx.Err() != nil {
			finish(EventCancelled{})
			return
		}
		if cfg.RuneCap > 0 && totalRunes >= cfg.RuneCap {
			finish(EventCapReached{TotalRunes: totalRunes, Cap: cfg.RuneCap})
			return
		}

		var toExpand []RuneExpansionInfo
		for _, ri := range currentLevel {
			if visited[ri.FullPath] {
				continue
			}
			if cfg.MaxDepth > 0 && ri.Depth > cfg.MaxDepth {
				continue
			}
			visited[ri.FullPath] = true
			toExpand = append(toExpand, ri)
		}
		if len(toExpand) == 0 {
			break
		}

		levelDepth := toExpand[0].Depth
		safeEmit(EventLevelStarted{Depth: levelDepth, Count: len(toExpand)})

		results := make([]expansionResult, len(toExpand))
		var wg sync.WaitGroup
		var totalReqNanos int64
		levelStart := time.Now()

		for i, ri := range toExpand {
			wg.Add(1)
			go func(i int, ri RuneExpansionInfo) {
				defer wg.Done()
				safeEmit(EventRuneStarted{Path: ri.FullPath, Depth: ri.Depth})
				results[i] = d.expandOne(ctx, ri, sess.BaseMessages, &totalReqNanos, safeEmit)
			}(i, ri)
		}
		wg.Wait()

		levelDur := time.Since(levelStart)
		sumDur := time.Duration(atomic.LoadInt64(&totalReqNanos))

		var nextLevel []RuneExpansionInfo
		for i := range results {
			r := results[i]
			ri := toExpand[i]
			parentPath := "root"
			if ri.ParentDecomposition != nil {
				parentPath = ri.ParentDecomposition.Path
			}

			if r.err != nil {
				safeEmit(EventRuneError{
					Path:      ri.FullPath,
					Depth:     ri.Depth,
					Err:       r.err.Error(),
					ElapsedMs: r.elapsedMs,
				})
				continue
			}
			if r.resp == nil {
				continue
			}

			childCount := countTotalRunes(r.resp)
			safeEmit(EventRuneExpanded{
				Path:       ri.FullPath,
				ParentPath: parentPath,
				Depth:      ri.Depth,
				Response:   r.resp,
				ElapsedMs:  r.elapsedMs,
				ChildCount: childCount,
			})

			childDecomposition := &AutoDecomposition{
				Path:       ri.FullPath,
				Depth:      ri.Depth,
				Response:   r.resp,
				ParentPath: parentPath,
				ChildPaths: make([]string, 0),
			}
			totalDecomps++
			totalRunes += childCount
			if ri.Depth > maxDepthSeen {
				maxDepthSeen = ri.Depth
			}

			newRunes := collectRunesForExpansion(r.resp)
			for j := range newRunes {
				newRunes[j].Depth = ri.Depth + 1
				newRunes[j].ParentDecomposition = childDecomposition
			}
			nextLevel = append(nextLevel, newRunes...)
		}

		safeEmit(EventLevelCompleted{
			Depth:        levelDepth,
			WallClockMs:  levelDur.Milliseconds(),
			SumRequestMs: sumDur.Milliseconds(),
		})

		currentLevel = nextLevel
	}

	finish(nil)
}

type expansionResult struct {
	resp      *DecompositionResponse
	err       error
	elapsedMs int64
}

// expandOne runs a single rune expansion. It stitches the base conversation
// with a rune-specific extension prompt and invokes the tool loop. On
// success, returns the parsed DecompositionResponse; on failure, returns an
// error. Read_example events emitted during the call are pushed through
// emit.
func (d *Decomposer) expandOne(ctx context.Context, ri RuneExpansionInfo, baseMessages []openai.ChatMessage, totalReqNanos *int64, emit func(ExpansionEvent)) expansionResult {
	extendedReq := fmt.Sprintf(`Forget the prior decomposition. Imagine you are seeing "%s" for the first time, in isolation, as a function you have to implement.

The user is browsing this decomposition as an interactive hierarchy: each rune is a column in a Miller-column (macOS-Finder-style) view, and the user will drill from parent to child to child. Your job is to continue that hierarchical breakdown by one more level beneath "%s".

Question: what 0–3 child units make up "%s"'s implementation? Each child should be a self-contained step the user would naturally drill into — a private helper, a distinct pipeline stage, or an internal subsystem, depending on the parent's granularity. They will appear as the next column to the right of "%s".

Call the decompose tool. The runes map keys must be of the form "%s.<child_name>". Example, for a different rune: if you were expanding "image.compress", reasonable children would be "image.compress.detect_format", "image.compress.choose_quality", "image.compress.encode_bytes". Each is a verb-phrase describing one internal step.

If "%s" is a single primitive operation (like an arithmetic op or a single syscall) and has no meaningful children, return an empty runes map ({}). That is the correct answer for leaves.

Hard rules:
- Reply ONLY by calling the decompose tool.
- Children exist only to serve "%s"; never include sibling-level functions, never repeat existing names, never include "%s" itself.
- A good child is one the user would click on to see its own next-level breakdown. Prefer 0–3 meaningful children over padding.
- At most 3 children.`, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath)

	localMsgs := make([]openai.ChatMessage, 0, len(baseMessages)+1)
	localMsgs = append(localMsgs, baseMessages...)
	localMsgs = append(localMsgs, openai.ChatMessage{
		Role:    openai.RoleUser,
		Content: extendedReq,
	})

	reqStart := time.Now()
	response, _, err := d.Decompose(ctx, localMsgs, emit)
	reqDur := time.Since(reqStart)
	atomic.AddInt64(totalReqNanos, int64(reqDur))
	elapsed := reqDur.Milliseconds()

	if err != nil {
		return expansionResult{nil, err, elapsed}
	}
	if clar, isClar := response.(ClarificationRequest); isClar {
		return expansionResult{nil, fmt.Errorf("model returned clarification: %s", clar.Message), elapsed}
	}
	respVal, ok := response.(DecompositionResponse)
	if !ok {
		return expansionResult{nil, fmt.Errorf("unexpected response type %T", response), elapsed}
	}
	if respVal.ProjectPackage.Name == "" {
		return expansionResult{nil, fmt.Errorf("empty project_package.name"), elapsed}
	}

	return expansionResult{&respVal, nil, elapsed}
}

// collectRunesForExpansion flattens the project+std rune maps of a response
// into a queue of RuneExpansionInfo, with fully-qualified paths.
func collectRunesForExpansion(resp *DecompositionResponse) []RuneExpansionInfo {
	var runes []RuneExpansionInfo
	if resp == nil {
		return runes
	}
	if resp.ProjectPackage.Name != "" && len(resp.ProjectPackage.Runes) > 0 {
		for name := range resp.ProjectPackage.Runes {
			runes = append(runes, RuneExpansionInfo{FullPath: qualify(resp.ProjectPackage.Name, name)})
		}
	}
	if resp.StdPackage != nil && resp.StdPackage.Name != "" && len(resp.StdPackage.Runes) > 0 {
		for name := range resp.StdPackage.Runes {
			runes = append(runes, RuneExpansionInfo{FullPath: qualify(resp.StdPackage.Name, name)})
		}
	}
	return runes
}

// countTotalRunes returns the total rune count across project + std packages
// of a single DecompositionResponse (not recursive).
func countTotalRunes(resp *DecompositionResponse) int {
	if resp == nil {
		return 0
	}
	count := 0
	count += len(resp.ProjectPackage.Runes)
	if resp.StdPackage != nil {
		count += len(resp.StdPackage.Runes)
	}
	return count
}
