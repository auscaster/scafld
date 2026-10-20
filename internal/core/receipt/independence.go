// Package receipt owns the signed classification shared by finalization and verification.
package receipt

import (
	"fmt"
	"strings"
)

// ClassifyIndependence derives the actual reviewer separation, ignoring caller claims.
func ClassifyIndependence(reviewerProvider, hostAgent string) Independence {
	var ind Independence
	if strings.TrimSpace(reviewerProvider) == "human" {
		return Independence{Level: IndependenceLevelHumanOverride, Reason: "explicit operator override; no independent model review is claimed"}
	}
	// Fully re-derive the signed level and distinctness from the reviewer that
	// actually ran versus the detected host, never trusting the passed-in level.
	// cross_vendor is stamped only when separation is genuinely proven, so neither a
	// forged cross_vendor input nor a stale isolation_only input can survive.
	if crossVendorProven(reviewerProvider, hostAgent) {
		ind.Level = IndependenceLevelCrossVendor
		ind.Distinct = true
	} else {
		ind.Level = IndependenceLevelIsolationOnly
		ind.Distinct = false
	}
	ind.Downgraded = independenceDowngrade(ind, reviewerProvider, hostAgent)
	ind.Reason = independenceReason(ind, reviewerProvider, hostAgent)
	return ind
}

// crossVendorProven reports whether the reviewer and host are both known and from
// different vendors, the only case in which cross_vendor independence is real.
func crossVendorProven(reviewerProvider, hostAgent string) bool {
	reviewerProvider = strings.TrimSpace(reviewerProvider)
	hostAgent = strings.TrimSpace(hostAgent)
	if !modelVendor(reviewerProvider) || !modelVendor(hostAgent) {
		return false
	}
	return !strings.EqualFold(reviewerProvider, hostAgent)
}

func independenceDowngrade(ind Independence, reviewerProvider string, hostAgent string) string {
	if ind.Level == IndependenceLevelCrossVendor {
		return ""
	}
	reviewerProvider = strings.TrimSpace(reviewerProvider)
	hostAgent = strings.TrimSpace(hostAgent)
	if reviewerProvider == "" || strings.EqualFold(reviewerProvider, "unknown") {
		return IndependenceDowngradeUnknownReviewer
	}
	if hostAgent == "" || strings.EqualFold(hostAgent, "unknown") {
		return IndependenceDowngradeUnknownHost
	}
	if strings.EqualFold(reviewerProvider, hostAgent) {
		return IndependenceDowngradeSameVendor
	}
	return ""
}

func independenceReason(ind Independence, reviewerProvider string, hostAgent string) string {
	reviewerProvider = strings.TrimSpace(reviewerProvider)
	hostAgent = strings.TrimSpace(hostAgent)
	switch ind.Level {
	case IndependenceLevelCrossVendor:
		return fmt.Sprintf("cross_vendor: reviewer %q and host %q are different model vendors; this is multi-model review that reduces correlated blind spots but remains single-party local tooling", reviewerProvider, hostAgent)
	case IndependenceLevelIsolationOnly:
		if hostAgent == "" || strings.EqualFold(hostAgent, "unknown") {
			return fmt.Sprintf("isolation_only: reviewer %q ran in an isolated path, but the host vendor was not detected so cross_vendor separation is not proven", reviewerProvider)
		}
		if strings.EqualFold(reviewerProvider, hostAgent) {
			return fmt.Sprintf("isolation_only: reviewer %q matches host %q; the review is isolated but same-vendor, so correlated blind spots remain possible", reviewerProvider, hostAgent)
		}
		return fmt.Sprintf("isolation_only: reviewer %q ran separately from host %q, but cross_vendor separation is not proven", reviewerProvider, hostAgent)
	default:
		return fmt.Sprintf("unrecognized independence level %q; verify policy remains authoritative", ind.Level)
	}
}

func modelVendor(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "codex", "claude", "gemini":
		return true
	default:
		return false
	}
}
