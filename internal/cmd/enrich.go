package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dedene/raindrop-cli/internal/api"
	"github.com/dedene/raindrop-cli/internal/enrich"
	"github.com/dedene/raindrop-cli/internal/errfmt"
	"github.com/dedene/raindrop-cli/internal/output"
)

// EnrichCmd scaffolds enrichment records for links.
type EnrichCmd struct {
	Collection string `help:"Collection name/ID (default: all)" default:"0" short:"c"`
	Limit      int    `help:"Maximum number of links to enrich" default:"50"`
	DryRun     bool   `help:"Preview enrichment records only (no writeback)" default:"true" name:"dry-run"`
}

func (c *EnrichCmd) Run(flags *RootFlags) error {
	if c.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if !c.DryRun {
		return fmt.Errorf("write mode is not implemented yet; rerun with --dry-run")
	}

	client, ctx, cancel, err := getClientWithContext()
	if err != nil {
		return errfmt.Format(err)
	}
	defer cancel()

	collectionID, err := client.ResolveCollection(ctx, c.Collection)
	if err != nil {
		return errfmt.Format(err)
	}

	items, err := c.fetchItems(ctx, client, collectionID)
	if err != nil {
		return errfmt.Format(err)
	}

	now := time.Now().UTC()
	records := make([]enrich.Record, 0, len(items))
	for i := range items {
		records = append(records, enrich.FromRaindrop(&items[i], now))
	}

	if flags.JSON {
		return output.WriteJSON(os.Stdout, records)
	}

	fmt.Fprintln(os.Stdout, "Enrich scaffold preview (dry-run)")
	tw := output.NewTableWriter(os.Stdout, "ID", "DOMAIN", "TYPE", "TOPIC", "CANONICAL_URL")
	for _, rec := range records {
		tw.AddRow(
			strconv.Itoa(rec.ID),
			output.SanitizeInline(rec.SourceDomain),
			output.SanitizeInline(rec.ContentType),
			output.SanitizeInline(rec.Topic),
			output.TruncateURL(output.SanitizeInline(rec.CanonicalURL), 56),
		)
	}
	tw.Render()
	fmt.Fprintf(os.Stdout, "\nPrepared %d enrichment record(s). No changes were written.\n", len(records))
	fmt.Fprintln(os.Stdout, "Tip: run with --json to inspect full scaffold output.")

	return nil
}

func (c *EnrichCmd) fetchItems(ctx context.Context, client *api.Client, collectionID int) ([]api.Raindrop, error) {
	opts := api.ListOptions{
		Sort:    "-created",
		PerPage: min(c.Limit, 50),
	}

	items := make([]api.Raindrop, 0, c.Limit)
	page := 0

	for len(items) < c.Limit {
		opts.Page = page

		resp, err := client.ListRaindrops(ctx, collectionID, opts)
		if err != nil {
			return nil, err
		}

		if len(resp.Items) == 0 {
			break
		}

		items = append(items, resp.Items...)
		if len(items) >= c.Limit || len(items) >= resp.Count {
			break
		}

		page++
	}

	if len(items) > c.Limit {
		items = items[:c.Limit]
	}

	return items, nil
}
