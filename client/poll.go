package client

import (
	"context"
	"time"
)

// pollConfig captures the resolved polling parameters for a single
// PollOperation / PollResource call.
type pollConfig struct {
	interval time.Duration
	maxWait  time.Duration
	backoff  float64
}

// PollOption tunes a polling call. All options are optional with sane
// defaults (2s interval, 10min max-wait, constant interval).
type PollOption func(*pollConfig)

// WithPollInterval sets the wait between status checks. Default 2 seconds.
// Caller's ctx deadline always wins — a 30s ctx with a 5s interval gives
// 6 attempts max.
func WithPollInterval(d time.Duration) PollOption {
	return func(p *pollConfig) {
		if d > 0 {
			p.interval = d
		}
	}
}

// WithPollMaxWait caps the total time spent polling regardless of ctx.
// Default 10 minutes. Pass 0 to disable (rely purely on ctx).
//
// Useful when callers want a hard SLA on a Wait helper without setting a
// ctx deadline — e.g. `client.Apps().CreateAndWait(ctx, req,
// client.WithPollMaxWait(5*time.Minute))` to abort even if ctx has no
// deadline.
func WithPollMaxWait(d time.Duration) PollOption {
	return func(p *pollConfig) { p.maxWait = d }
}

// WithPollBackoff multiplies the interval after each poll. Default 1.0
// (constant interval). Set to 1.5 for gentle backoff that's friendly to
// long-running operations (apps that take 2-3 minutes).
//
// Backoff is capped by maxWait so an aggressive factor won't blow past
// the SLA.
func WithPollBackoff(factor float64) PollOption {
	return func(p *pollConfig) {
		if factor >= 1.0 {
			p.backoff = factor
		}
	}
}

func defaultPollConfig() *pollConfig {
	return &pollConfig{
		interval: 2 * time.Second,
		maxWait:  10 * time.Minute,
		backoff:  1.0,
	}
}

// PollResource polls a resource until a terminal state is reached.
//
// fetch is called once per attempt to retrieve the current state.
// terminal classifies the state: return (true, nil) on success,
// (true, error) on terminal failure (e.g. resource.Status=="failed"),
// (false, nil) to keep polling.
//
// Returns the last fetched value alongside any error. On ctx cancellation
// or maxWait expiry, returns the last fetched value (if any) and ctx.Err()
// / context.DeadlineExceeded.
//
// This generic helper covers VPS (ActionStatus on the instance) and
// Volume (Status on the volume). For Apps use AppsService.PollOperation
// which knows the operation-id polling shape.
//
// Example:
//
//	vol, err := client.PollResource(ctx,
//	    func(ctx context.Context) (*types.VolumeResponse, error) {
//	        return c.Volumes().Get(ctx, volumeID)
//	    },
//	    func(v *types.VolumeResponse) (bool, error) {
//	        switch v.Status {
//	        case "ready", "detached":
//	            return true, nil
//	        case "failed":
//	            return true, fmt.Errorf("volume failed: %s", *v.ErrorMessage)
//	        default:
//	            return false, nil
//	        }
//	    },
//	)
func PollResource[T any](
	ctx context.Context,
	fetch func(context.Context) (T, error),
	terminal func(T) (done bool, err error),
	opts ...PollOption,
) (T, error) {
	cfg := defaultPollConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// Apply maxWait as a derived ctx deadline so both bounds collapse to
	// "whichever expires first wins."
	if cfg.maxWait > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.maxWait)
		defer cancel()
	}

	var zero T
	interval := cfg.interval

	for {
		// ctx check before each fetch — cheap insurance.
		if err := ctx.Err(); err != nil {
			return zero, err
		}

		current, err := fetch(ctx)
		if err != nil {
			// Fetch failures during polling are surfaced immediately —
			// the SDK's HTTP retry policy already handled transient
			// network/5xx blips, so anything reaching here is genuine.
			return current, err
		}
		done, termErr := terminal(current)
		if done {
			return current, termErr
		}

		if err := sleepCtx(ctx, interval); err != nil {
			return current, err
		}
		// Apply backoff factor for next iteration. Cap at maxWait so an
		// aggressive factor doesn't sleep past the deadline.
		if cfg.backoff > 1.0 {
			next := time.Duration(float64(interval) * cfg.backoff)
			if cfg.maxWait > 0 && next > cfg.maxWait {
				next = cfg.maxWait
			}
			interval = next
		}
	}
}
