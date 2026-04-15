package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"shotgun.dev/odek/openai"
)

type RuneExpansionInfo struct {
	FullPath            string
	Depth               int
	ParentDecomposition *AutoDecomposition
}

type AutoDecomposition struct {
	Path       string
	Depth      int
	Response   *DecompositionResponse
	ParentPath string
	ChildPaths []string
}

type expansionResult struct {
	runeInfo RuneExpansionInfo
	resp     *DecompositionResponse
	err      error
}

// expandRecursively drains the expansion queue level-by-level, decomposing each
// rune up to cfg.MaxDepth and stitching child decompositions back into the tree.
// Each level dispatches all expansions in parallel since each expansion is
// independent — it only needs the initial decomposition context, not the
// results of sibling expansions.
func expandRecursively(ctx context.Context, baseMessages []openai.ChatMessage, root *AutoDecomposition, queue []RuneExpansionInfo, cfg RunConfig) {
	for i := range queue {
		queue[i].ParentDecomposition = root
	}

	allDecompositions := []*AutoDecomposition{root}
	totalRunesCount := countTotalRunes(root.Response)
	visitedRunePaths := map[string]bool{"root": true}

	fmt.Printf("\n🔄 Starting auto-recursion (depth 0: %d runes)\n", len(queue))

	currentLevel := queue
	for len(currentLevel) > 0 {
		if totalRunesCount >= cfg.RuneCap {
			fmt.Printf("\n⚠️  Max total runes (%d) reached. Stopping expansion.\n", cfg.RuneCap)
			break
		}

		var toExpand []RuneExpansionInfo
		for _, ri := range currentLevel {
			if visitedRunePaths[ri.FullPath] || ri.Depth >= cfg.MaxDepth {
				continue
			}
			visitedRunePaths[ri.FullPath] = true
			toExpand = append(toExpand, ri)
		}
		if len(toExpand) == 0 {
			break
		}

		fmt.Printf("\n📤 Dispatching %d expansions...\n", len(toExpand))

		results := make([]expansionResult, len(toExpand))
		var wg sync.WaitGroup
		var totalReqNanos int64
		levelStart := time.Now()

		for i, ri := range toExpand {
			wg.Add(1)
			go func(i int, ri RuneExpansionInfo) {
				defer wg.Done()
				results[i] = expandOne(ctx, ri, baseMessages, &totalReqNanos)
			}(i, ri)
		}

		wg.Wait()

		levelDur := time.Since(levelStart)
		sumDur := time.Duration(atomic.LoadInt64(&totalReqNanos))
		factor := float64(sumDur) / float64(levelDur)
		fmt.Printf("   ⏱️  level wall-clock: %s, sum of %d requests: %s (parallelism factor: %.1fx)\n",
			levelDur.Round(time.Millisecond),
			len(toExpand),
			sumDur.Round(time.Millisecond),
			factor,
		)

		var nextLevel []RuneExpansionInfo
		for _, r := range results {
			if r.resp == nil {
				continue
			}

			newRunes := collectRunesForExpansion(r.resp)
			if len(newRunes) == 0 {
				continue
			}

			childDecomposition := &AutoDecomposition{
				Path:       r.runeInfo.FullPath,
				Depth:      r.runeInfo.Depth + 1,
				Response:   r.resp,
				ParentPath: "",
				ChildPaths: make([]string, 0),
			}
			allDecompositions = append(allDecompositions, childDecomposition)

			if r.runeInfo.ParentDecomposition != nil {
				r.runeInfo.ParentDecomposition.ChildPaths = append(r.runeInfo.ParentDecomposition.ChildPaths, r.runeInfo.FullPath)
				childDecomposition.ParentPath = r.runeInfo.ParentDecomposition.Path
			}

			for j := range newRunes {
				newRunes[j].Depth = r.runeInfo.Depth + 1
				newRunes[j].ParentDecomposition = childDecomposition
			}
			nextLevel = append(nextLevel, newRunes...)
			totalRunesCount += countTotalRunes(r.resp)
		}

		currentLevel = nextLevel
	}

	fmt.Printf("\n")
	printCompleteTree(allDecompositions, "root", 0, true)

	separator := strings.Repeat("=", 70)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("📊 SUMMARY: %d decompositions, %d runes discovered (max depth %d)\n", len(allDecompositions), totalRunesCount, cfg.MaxDepth)
	fmt.Printf("%s\n", separator)
}

// expandOne runs a single rune expansion in its own goroutine. The work was
// previously inline inside expandRecursively; pulling it out drops the nesting
// by two levels and makes the error handling easier to read.
func expandOne(ctx context.Context, ri RuneExpansionInfo, baseMessages []openai.ChatMessage, totalReqNanos *int64) expansionResult {
	extendedReq := fmt.Sprintf(`Forget the prior decomposition. Imagine you are seeing "%s" for the first time, in isolation, as a black-box function you have to implement.

Question: what 0–3 PRIVATE helper functions would you write inside "%s"'s body to do its job? Helpers that no other function would ever call. Implementation details only.

Call the decompose tool. The runes map keys must be of the form "%s.<new_helper_name>". Example, for a different rune: if you were expanding "image.compress", reasonable helpers would be "image.compress.detect_format", "image.compress.choose_quality", "image.compress.encode_bytes". Each is a verb-phrase describing one internal step.

If "%s" is a single primitive operation (like an arithmetic op or a single syscall) and would have no private helpers in its body, return an empty runes map ({}). That is the correct answer.

Hard rules:
- Reply ONLY by calling the decompose tool.
- Never include sibling-level functions, never repeat existing names, never include "%s" itself.
- At most 3 helpers.`, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath, ri.FullPath)

	localMsgs := make([]openai.ChatMessage, 0, len(baseMessages)+1)
	localMsgs = append(localMsgs, baseMessages...)
	localMsgs = append(localMsgs, openai.ChatMessage{
		Role:    openai.RoleUser,
		Content: extendedReq,
	})

	reqStart := time.Now()
	response, _, err := client.Decompose(ctx, localMsgs)
	reqDur := time.Since(reqStart)
	atomic.AddInt64(totalReqNanos, int64(reqDur))
	dur := reqDur.Round(time.Millisecond)

	if err != nil {
		stdoutMu.Lock()
		fmt.Printf("   ⚠️  %s: %v (%s)\n", ri.FullPath, err, dur)
		stdoutMu.Unlock()
		return expansionResult{ri, nil, err}
	}

	if clar, isClar := response.(ClarificationRequest); isClar {
		stdoutMu.Lock()
		fmt.Printf("   ⚠️  %s: model returned text instead of tool call (%s): %q\n", ri.FullPath, dur, clar.Message)
		stdoutMu.Unlock()
		return expansionResult{ri, nil, fmt.Errorf("unexpected response type %T", response)}
	}

	respVal, ok := response.(DecompositionResponse)
	if !ok {
		stdoutMu.Lock()
		fmt.Printf("   ⚠️  %s: unexpected response type %T (%s): %+v\n", ri.FullPath, response, dur, response)
		stdoutMu.Unlock()
		return expansionResult{ri, nil, fmt.Errorf("unexpected response type %T", response)}
	}
	if respVal.ProjectPackage.Name == "" {
		stdoutMu.Lock()
		fmt.Printf("   ⚠️  %s: tool call had empty project_package.name (%s)\n      parsed response: %+v\n", ri.FullPath, dur, respVal)
		stdoutMu.Unlock()
		return expansionResult{ri, nil, fmt.Errorf("empty project_package.name")}
	}

	newRunes := collectRunesForExpansion(&respVal)
	stdoutMu.Lock()
	if len(newRunes) == 0 {
		fmt.Printf("   ✓ %s: leaf (%s)\n", ri.FullPath, dur)
	} else {
		fmt.Printf("   ➜ %s: %d sub-runes (%s)\n", ri.FullPath, len(newRunes), dur)
	}
	stdoutMu.Unlock()

	return expansionResult{ri, &respVal, nil}
}

func collectRunesForExpansion(resp *DecompositionResponse) []RuneExpansionInfo {
	var runes []RuneExpansionInfo

	if resp == nil || resp.ProjectPackage.Name == "" {
		return runes
	}

	if len(resp.ProjectPackage.Runes) > 0 {
		for name := range resp.ProjectPackage.Runes {
			path := name
			if !strings.HasPrefix(name, resp.ProjectPackage.Name+".") {
				path = fmt.Sprintf("%s.%s", resp.ProjectPackage.Name, name)
			}
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		for name := range resp.StdPackage.Runes {
			path := name
			if !strings.HasPrefix(name, resp.StdPackage.Name+".") {
				path = fmt.Sprintf("%s.%s", resp.StdPackage.Name, name)
			}
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	return runes
}

func countTotalRunes(resp *DecompositionResponse) int {
	if resp == nil {
		return 0
	}
	count := 0
	if len(resp.ProjectPackage.Runes) > 0 {
		count += len(resp.ProjectPackage.Runes)
	}
	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		count += len(resp.StdPackage.Runes)
	}
	return count
}
