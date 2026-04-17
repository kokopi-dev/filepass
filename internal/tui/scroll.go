package tui

// visibleRange returns [start, end) indices of rows that should be rendered.
func visibleRange(total, selected, scrollOff, viewHeight int) (int, int) {
	if total <= 0 {
		return 0, 0
	}

	maxVis := viewHeight
	if maxVis <= 0 {
		maxVis = total // unconstrained
	}
	if maxVis >= total {
		return 0, total
	}

	// Reserve lines for scroll indicators.
	hasAbove := scrollOff > 0
	hasBelow := scrollOff+maxVis < total
	if hasAbove {
		maxVis--
	}
	if hasBelow {
		maxVis--
	}
	if maxVis < 1 {
		maxVis = 1
	}

	start := scrollOff
	end := start + maxVis

	if selected < start {
		start = selected
		end = start + maxVis
	}
	if selected >= end {
		end = selected + 1
		start = end - maxVis
	}
	if start < 0 {
		start = 0
		end = min(maxVis, total)
	}
	if end > total {
		end = total
		start = max(0, end-maxVis)
	}

	return start, end
}

// syncScroll adjusts scrollOff so selected stays inside the viewport.
func syncScroll(selected, total, viewHeight, scrollOff int) int {
	if total <= 0 {
		return 0
	}
	maxVis := viewHeight
	if maxVis <= 0 || maxVis >= total {
		return 0
	}

	if selected < 0 {
		selected = 0
	}
	if selected >= total {
		selected = total - 1
	}

	if scrollOff < 0 {
		scrollOff = 0
	}
	if scrollOff > total-1 {
		scrollOff = total - 1
	}

	if selected < scrollOff {
		scrollOff = selected
	}
	if selected >= scrollOff+maxVis {
		scrollOff = selected - maxVis + 1
	}

	if scrollOff < 0 {
		scrollOff = 0
	}
	return scrollOff
}
